"""Generates SDK model classes from openapi.yaml.

Go, TypeScript and Python use established generators. For the other SDKs,
this script reads the response schemas of the spec and writes plain model
classes, with the same rules everywhere:

- request and response envelope schemas are skipped (each client builds
  its requests and unwraps the envelope itself);
- the four "flexible" schemas (StringOrInteger, StringOrStrings, MagnetList,
  QuotaMap) map to hand-written types of each SDK;
- enums stay plain strings or integers, so a new value never breaks decoding;
- required fields are non-optional, every other field is optional.

Usage: python scripts/gen_models.py <language> <output file>
Languages: kotlin, rust, csharp, swift, php (the output is the src directory)
"""

from __future__ import annotations

import re
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Any

import yaml

SPEC = Path(__file__).resolve().parent.parent / "openapi.yaml"
FLEX = {"StringOrInteger", "StringOrStrings", "MagnetList", "QuotaMap"}
SKIP = {"ResponseStatus", "ErrorCode"}
# "Error" clashes with a built-in type in Kotlin, C# and Swift.
RENAME = {"Error": "ApiError"}


# --------------------------------------------------------------------------- IR


@dataclass
class TypeRef:
    kind: str  # string, int, long, double, bool, list, map, ref, flex, any
    item: TypeRef | None = None
    name: str = ""


@dataclass
class Field:
    json_name: str
    type: TypeRef
    required: bool
    description: str


@dataclass
class Model:
    name: str
    description: str
    fields: list[Field]
    alias_of: TypeRef | None = None  # a named map, e.g. ServiceMap


def type_of(schema: dict[str, Any]) -> TypeRef:
    if "$ref" in schema:
        name = schema["$ref"].rsplit("/", 1)[-1]
        return TypeRef("flex" if name in FLEX else "ref", name=RENAME.get(name, name))
    t = schema.get("type")
    if t == "string":
        return TypeRef("string")
    if t == "integer":
        return TypeRef("long" if schema.get("format") == "int64" else "int")
    if t == "number":
        return TypeRef("double")
    if t == "boolean":
        return TypeRef("bool")
    if t == "array":
        return TypeRef("list", item=type_of(schema.get("items", {})))
    if t == "object" and isinstance(schema.get("additionalProperties"), dict):
        return TypeRef("map", item=type_of(schema["additionalProperties"]))
    return TypeRef("any")


def load_models() -> tuple[list[Model], list[str]]:
    spec = yaml.safe_load(SPEC.read_text())
    schemas: dict[str, Any] = spec["components"]["schemas"]
    models: list[Model] = []
    for name, schema in schemas.items():
        if name in FLEX or name in SKIP or name.endswith(("Request", "Response")):
            continue
        description = (schema.get("description") or "").strip()
        name = RENAME.get(name, name)
        if schema.get("type") == "object" and "properties" not in schema:
            models.append(Model(name, description, [], alias_of=type_of(schema)))
            continue
        required = set(schema.get("required", []))
        fields = [
            Field(json_name, type_of(prop), json_name in required, (prop.get("description") or "").strip())
            for json_name, prop in schema.get("properties", {}).items()
        ]
        models.append(Model(name, description, fields))
    error_codes = [str(c) for c in schemas["ErrorCode"]["enum"]]
    return models, error_codes


def camel(name: str) -> str:
    parts = re.split(r"[_\W]+", name)
    return parts[0] + "".join(p[:1].upper() + p[1:] for p in parts[1:])


def pascal(name: str) -> str:
    c = camel(name)
    return c[:1].upper() + c[1:]


def code_constant(code: str) -> str:
    """GENERIC -> GENERIC, 404 -> NOT_FOUND."""
    return "NOT_FOUND" if code == "404" else code


def doc_lines(text: str) -> list[str]:
    return [line.rstrip() for line in text.splitlines()] if text else []


# ----------------------------------------------------------------------- Kotlin

KOTLIN_FLEX = {
    "StringOrInteger": ("FlexString", "FlexStringSerializer"),
    "StringOrStrings": ("List<String>", "StringsSerializer"),
    "MagnetList": ("List<Magnet>", "MagnetsSerializer"),
    "QuotaMap": ("Map<String, Long>", "QuotasSerializer"),
}
KOTLIN_PRIM = {"string": "String", "int": "Int", "long": "Long", "double": "Double", "bool": "Boolean", "any": "JsonElement"}
KOTLIN_KEYWORDS = {"as", "break", "class", "continue", "do", "else", "false", "for", "fun", "if", "in", "interface",
                   "is", "null", "object", "package", "return", "super", "this", "throw", "true", "try", "typealias",
                   "typeof", "val", "var", "when", "while"}


def kotlin_type(t: TypeRef) -> str:
    if t.kind in KOTLIN_PRIM:
        return KOTLIN_PRIM[t.kind]
    if t.kind == "list":
        assert t.item
        return f"List<{kotlin_type(t.item)}>"
    if t.kind == "map":
        assert t.item
        return f"Map<String, {kotlin_type(t.item)}>"
    if t.kind == "flex":
        return KOTLIN_FLEX[t.name][0]
    return t.name


def kotlin_kdoc(text: str, indent: str) -> list[str]:
    lines = doc_lines(text.replace("*/", "*&#47;"))
    if not lines:
        return []
    if len(lines) == 1:
        return [f"{indent}/** {lines[0]} */"]
    return [f"{indent}/**", *[f"{indent} * {line}".rstrip() for line in lines], f"{indent} */"]


def emit_kotlin(models: list[Model], codes: list[str], package: str) -> str:
    out = [
        "// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.",
        "@file:Suppress(\"ktlint\")",
        "",
        f"package {package}",
        "",
        "import kotlinx.serialization.SerialName",
        "import kotlinx.serialization.Serializable",
        "import kotlinx.serialization.json.JsonElement",
        "",
    ]
    for m in models:
        if m.alias_of is not None:
            out += kotlin_kdoc(m.description, "")
            out += [f"public typealias {m.name} = {kotlin_type(m.alias_of)}", ""]
            continue
        out += kotlin_kdoc(m.description, "")
        out.append("@Serializable")
        out.append(f"public data class {m.name}(")
        for f in m.fields:
            prop = camel(f.json_name)
            if prop in KOTLIN_KEYWORDS:
                prop = f"`{prop}`"
            out += kotlin_kdoc(f.description, "    ")
            annotations = []
            if f.json_name != prop.strip("`"):
                annotations.append(f'@SerialName("{f.json_name}")')
            if f.type.kind == "flex":
                annotations.append(f"@Serializable(with = {KOTLIN_FLEX[f.type.name][1]}::class)")
            typ = kotlin_type(f.type)
            default = "" if f.required else " = null"
            nullable = "" if f.required else "?"
            prefix = " ".join(annotations) + " " if annotations else ""
            out.append(f"    {prefix}val {prop}: {typ}{nullable}{default},")
        out += [")", ""]
    out += [
        "/** Error codes the API is known to return. The API may return others. */",
        "public object ErrorCodes {",
        *[f'    public const val {code_constant(c)}: String = "{c}"' for c in codes],
        "",
        "    /** Every code of the spec. */",
        "    public val ALL: List<String> = listOf(",
        *[f"        {code_constant(c)}," for c in codes],
        "    )",
        "}",
        "",
    ]
    return "\n".join(out)


# -------------------------------------------------------------------------- PHP

PHP_HEADER = "<?php\n\n// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.\n\ndeclare(strict_types=1);\n"


def php_resolve(t: TypeRef, aliases: dict[str, TypeRef]) -> TypeRef:
    """Replaces references to named maps (ServiceMap) with the map itself."""
    if t.kind == "ref" and t.name in aliases:
        return aliases[t.name]
    return t


def php_types(t: TypeRef, aliases: dict[str, TypeRef]) -> tuple[str, str]:
    """The native type and the PHPDoc type of a value."""
    t = php_resolve(t, aliases)
    prim = {"string": "string", "int": "int", "long": "int", "double": "float", "bool": "bool", "any": "mixed"}
    if t.kind in prim:
        return prim[t.kind], prim[t.kind]
    if t.kind == "list":
        assert t.item
        return "array", f"list<{php_types(t.item, aliases)[1]}>"
    if t.kind == "map":
        assert t.item
        return "array", f"array<string, {php_types(t.item, aliases)[1]}>"
    if t.kind == "flex":
        return {
            "StringOrInteger": ("string", "string"),
            "StringOrStrings": ("array", "list<string>"),
            "MagnetList": ("array", "list<Magnet>"),
            "QuotaMap": ("array", "array<string, int>"),
        }[t.name]
    return t.name, t.name


def php_reader(t: TypeRef, aliases: dict[str, TypeRef], value: str) -> str:
    """An expression that converts a decoded JSON value to the field type."""
    t = php_resolve(t, aliases)
    simple = {"string": "string", "int": "int", "long": "int", "double": "float", "bool": "bool"}
    if t.kind in simple:
        return f"Hydrate::{simple[t.kind]}({value})"
    if t.kind == "any":
        return value
    if t.kind == "list":
        assert t.item
        return f"Hydrate::list({value}, static fn (mixed $v): {php_types(t.item, aliases)[0]} => {php_reader(t.item, aliases, '$v')})"
    if t.kind == "map":
        assert t.item
        return f"Hydrate::map({value}, static fn (mixed $v): {php_types(t.item, aliases)[0]} => {php_reader(t.item, aliases, '$v')})"
    if t.kind == "flex":
        return {
            "StringOrInteger": f"Hydrate::flexString({value})",
            "StringOrStrings": f"Hydrate::strings({value})",
            "MagnetList": f"Hydrate::list(Hydrate::oneOrMany({value}), static fn (mixed $v): Magnet => Magnet::fromArray(Hydrate::object($v)))",
            "QuotaMap": f"Hydrate::map(Hydrate::emptyListAsMap({value}), static fn (mixed $v): int => Hydrate::int($v))",
        }[t.name]
    return f"{t.name}::fromArray(Hydrate::object({value}))"


def php_doc(lines: list[str], indent: str) -> list[str]:
    lines = [line.replace("*/", "*&#47;") for line in lines]
    if not lines:
        return []
    if len(lines) == 1:
        return [f"{indent}/** {lines[0]} */"]
    return [f"{indent}/**", *[f"{indent} * {line}".rstrip() for line in lines], f"{indent} */"]


def emit_php(models: list[Model], codes: list[str], namespace: str) -> dict[str, str]:
    aliases = {m.name: m.alias_of for m in models if m.alias_of is not None}
    files: dict[str, str] = {}
    for m in models:
        if m.alias_of is not None:
            continue
        fields = sorted(m.fields, key=lambda f: not f.required)  # required parameters first
        out = [PHP_HEADER, f"namespace {namespace}\\Model;", "", f"use {namespace}\\Internal\\Hydrate;", ""]
        out += php_doc(doc_lines(m.description), "")
        out += [f"final class {m.name}", "{"]
        params = []
        for f in fields:
            native, doc = php_types(f.type, aliases)
            nullable = "" if f.required else "?"
            default = "" if f.required else " = null"
            prop = camel(f.json_name)
            desc = doc_lines(f.description)
            if native != doc:
                desc = [*desc, f"@var {doc}{'' if f.required else '|null'}"]
            params += php_doc(desc, "        ")
            params.append(f"        public readonly {nullable}{native} ${prop}{default},")
        out += ["    public function __construct(", *params, "    ) {", "    }", ""]
        out += [
            "    /** @param array<array-key, mixed> $data A decoded JSON object. */",
            "    public static function fromArray(array $data): self",
            "    {",
            "        return new self(",
        ]
        for f in fields:
            prop = camel(f.json_name)
            reader = php_reader(f.type, aliases, f"$data['{f.json_name}']")
            if f.required:
                out.append(f"            {prop}: {php_reader(f.type, aliases, f'Hydrate::required($data, {json_str(f.json_name)})')},")
            else:
                out.append(f"            {prop}: array_key_exists({json_str(f.json_name)}, $data) && $data[{json_str(f.json_name)}] !== null ? {reader} : null,")
        out += ["        );", "    }", "}", ""]
        files[f"Model/{m.name}.php"] = "\n".join(out)

    consts = [f"    public const {code_constant(c)} = {json_str(c)};" for c in codes]
    files["ErrorCode.php"] = "\n".join([
        PHP_HEADER,
        f"namespace {namespace};",
        "",
        "/** Error codes the API is known to return. The API may return others. */",
        "final class ErrorCode",
        "{",
        *consts,
        "",
        "    /** Every code of the spec. */",
        "    public const ALL = [",
        *[f"        self::{code_constant(c)}," for c in codes],
        "    ];",
        "}",
        "",
    ])
    return files


def json_str(s: str) -> str:
    return "'" + s.replace("\\", "\\\\").replace("'", "\\'") + "'"


# ----------------------------------------------------------------------- C#

CSHARP_FLEX = {
    "StringOrInteger": ("FlexString", None),
    "StringOrStrings": ("IReadOnlyList<string>", "StringsConverter"),
    "MagnetList": ("IReadOnlyList<Magnet>", "MagnetListConverter"),
    "QuotaMap": ("IReadOnlyDictionary<string, long>", "QuotasConverter"),
}


def csharp_type(t: TypeRef) -> str:
    prim = {"string": "string", "int": "int", "long": "long", "double": "double", "bool": "bool",
            "any": "System.Text.Json.JsonElement"}
    if t.kind in prim:
        return prim[t.kind]
    if t.kind == "list":
        assert t.item
        return f"IReadOnlyList<{csharp_type(t.item)}>"
    if t.kind == "map":
        assert t.item
        return f"IReadOnlyDictionary<string, {csharp_type(t.item)}>"
    if t.kind == "flex":
        return CSHARP_FLEX[t.name][0]
    return t.name


def csharp_doc(text: str, indent: str) -> list[str]:
    lines = [line.replace("&", "&amp;").replace("<", "&lt;").replace(">", "&gt;") for line in doc_lines(text)]
    if not lines:
        return []
    return [f"{indent}/// <summary>", *[f"{indent}/// {line}".rstrip() for line in lines], f"{indent}/// </summary>"]


def csharp_const(code: str) -> str:
    return "NotFound" if code == "404" else pascal(code.lower())


def emit_csharp(models: list[Model], codes: list[str], namespace: str) -> str:
    aliases = {m.name: m.alias_of for m in models if m.alias_of is not None}
    out = [
        "// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.",
        "#nullable enable",
        "",
        "using System.Collections.Generic;",
        "using System.Text.Json.Serialization;",
        "",
        f"namespace {namespace};",
        "",
    ]
    for m in models:
        if m.alias_of is not None:
            continue
        out += csharp_doc(m.description, "")
        out += [f"public sealed record {m.name}", "{"]
        for i, f in enumerate(m.fields):
            t = aliases.get(f.type.name, f.type) if f.type.kind == "ref" else f.type
            typ = csharp_type(t)
            if i:
                out.append("")
            out += csharp_doc(f.description, "    ")
            out.append(f'    [JsonPropertyName("{f.json_name}")]')
            if t.kind == "flex" and CSHARP_FLEX[t.name][1]:
                out.append(f"    [JsonConverter(typeof({CSHARP_FLEX[t.name][1]}))]")
            if f.required:
                out.append(f"    public required {typ} {pascal(f.json_name)} {{ get; init; }}")
            else:
                out.append(f"    public {typ}? {pascal(f.json_name)} {{ get; init; }}")
        out += ["}", ""]
    out += [
        "/// <summary>Error codes the API is known to return. The API may return others.</summary>",
        "public static class ErrorCodes",
        "{",
        *[f'    public const string {csharp_const(c)} = "{c}";' for c in codes],
        "",
        "    /// <summary>Every code of the spec.</summary>",
        "    public static readonly IReadOnlyList<string> All = new[]",
        "    {",
        *[f"        {csharp_const(c)}," for c in codes],
        "    };",
        "}",
        "",
    ]
    return "\n".join(out)


# ------------------------------------------------------------------------ Swift

SWIFT_FLEX = {
    "StringOrInteger": ("FlexString", None),
    "StringOrStrings": ("[String]", "Strings"),
    "MagnetList": ("[Magnet]", "Magnets"),
    "QuotaMap": ("[String: Int64]", "Quotas"),
}


def swift_type(t: TypeRef) -> str:
    prim = {"string": "String", "int": "Int", "long": "Int64", "double": "Double", "bool": "Bool"}
    if t.kind in prim:
        return prim[t.kind]
    if t.kind == "list":
        assert t.item
        return f"[{swift_type(t.item)}]"
    if t.kind == "map":
        assert t.item
        return f"[String: {swift_type(t.item)}]"
    if t.kind == "flex":
        return SWIFT_FLEX[t.name][0]
    return t.name


def swift_doc(text: str, indent: str) -> list[str]:
    return [f"{indent}/// {line}".rstrip() for line in doc_lines(text)]


def swift_const(code: str) -> str:
    return "notFound" if code == "404" else camel(code.lower())


def emit_swift(models: list[Model], codes: list[str]) -> str:
    out = ["// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.", "// swift-format-ignore-file", "",
           "import Foundation", ""]
    for m in models:
        if m.alias_of is not None:
            out += swift_doc(m.description, "")
            out += [f"public typealias {m.name} = {swift_type(m.alias_of)}", ""]
            continue
        out += swift_doc(m.description, "")
        out.append(f"public struct {m.name}: Decodable, Sendable, Equatable {{")
        fields = [(f, camel(f.json_name)) for f in m.fields]
        for f, prop in fields:
            out += swift_doc(f.description, "    ")
            out.append(f"    public let {prop}: {swift_type(f.type)}{'' if f.required else '?'}")
        # Memberwise initializer, public so apps and tests can build values.
        params = ", ".join(f"{prop}: {swift_type(f.type)}{'' if f.required else '? = nil'}" for f, prop in fields)
        out += ["", f"    public init({params}) {{", *[f"        self.{prop} = {prop}" for _, prop in fields], "    }", ""]
        out.append("    private enum CodingKeys: String, CodingKey {")
        for f, prop in fields:
            out.append(f"        case {prop}" + (f' = "{f.json_name}"' if prop != f.json_name else ""))
        out += ["    }", "", "    public init(from decoder: any Decoder) throws {",
                "        let c = try decoder.container(keyedBy: CodingKeys.self)"]
        for f, prop in fields:
            flex = SWIFT_FLEX[f.type.name][1] if f.type.kind == "flex" else None
            if flex:
                call = f"decode{flex}(forKey: .{prop})" if f.required else f"decode{flex}IfPresent(forKey: .{prop})"
            else:
                call = f"decode{'' if f.required else 'IfPresent'}({swift_type(f.type)}.self, forKey: .{prop})"
            out.append(f"        {prop} = try c.{call}")
        out += ["    }", "}", ""]
    out += [
        "/// Error codes the API is known to return. The API may return others.",
        "public enum ErrorCode {",
        *[f'    public static let {swift_const(c)} = "{c}"' for c in codes],
        "",
        "    /// Every code of the spec.",
        "    public static let all: [String] = [",
        *[f"        {swift_const(c)}," for c in codes],
        "    ]",
        "}",
        "",
    ]
    return "\n".join(out)


# ------------------------------------------------------------------------- Rust

RUST_KEYWORDS = {"type", "match", "ref", "self", "struct", "fn", "move", "mod", "use", "impl", "in", "loop", "where"}
RUST_FLEX = {
    "StringOrInteger": ("FlexString", None),
    "StringOrStrings": ("Vec<String>", "strings"),
    "MagnetList": ("Vec<Magnet>", "magnets"),
    "QuotaMap": ("HashMap<String, i64>", "quotas"),
}


def snake(name: str) -> str:
    s = re.sub(r"([a-z0-9])([A-Z])", r"\1_\2", name)
    return re.sub(r"\W+", "_", s).lower().strip("_")


def rust_type(t: TypeRef) -> str:
    prim = {"string": "String", "int": "i64", "long": "i64", "double": "f64", "bool": "bool", "any": "serde_json::Value"}
    if t.kind in prim:
        return prim[t.kind]
    if t.kind == "list":
        assert t.item
        return f"Vec<{rust_type(t.item)}>"
    if t.kind == "map":
        assert t.item
        return f"HashMap<String, {rust_type(t.item)}>"
    if t.kind == "flex":
        return RUST_FLEX[t.name][0]
    return t.name


def rust_doc(text: str, indent: str) -> list[str]:
    return [f"{indent}/// {line}".rstrip() for line in doc_lines(text)]


def emit_rust(models: list[Model], codes: list[str]) -> str:
    out = [
        "// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.",
        "",
        "#![allow(clippy::doc_markdown, clippy::doc_lazy_continuation)]",
        "",
        "use std::collections::HashMap;",
        "",
        "use serde::{Deserialize, Serialize};",
        "",
        "use crate::flex::{self, FlexString};",
        "",
    ]
    for m in models:
        if m.alias_of is not None:
            out += rust_doc(m.description, "")
            out += [f"pub type {m.name} = {rust_type(m.alias_of)};", ""]
            continue
        out += rust_doc(m.description, "")
        out += ["#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]", f"pub struct {m.name} {{"]
        for f in m.fields:
            name = snake(f.json_name)
            ident = f"r#{name}" if name in RUST_KEYWORDS else name
            out += rust_doc(f.description, "    ")
            attrs = []
            if name != f.json_name:
                attrs.append(f'rename = "{f.json_name}"')
            typ = rust_type(f.type)
            flex_fn = RUST_FLEX[f.type.name][1] if f.type.kind == "flex" else None
            if not f.required:
                attrs.append("default")
                attrs.append('skip_serializing_if = "Option::is_none"')
                if flex_fn:
                    attrs.append(f'deserialize_with = "flex::opt_{flex_fn}"')
                typ = f"Option<{typ}>"
            elif flex_fn:
                attrs.append(f'deserialize_with = "flex::{flex_fn}"')
            if attrs:
                out.append(f"    #[serde({', '.join(attrs)})]")
            out.append(f"    pub {ident}: {typ},")
        out += ["}", ""]
    out += [
        "/// Error codes the API is known to return. The API may return others.",
        "pub mod error_codes {",
        *[f'    pub const {code_constant(c)}: &str = "{c}";' for c in codes],
        "",
        "    /// Every code of the spec.",
        "    pub const ALL: &[&str] = &[",
        *[f"        {code_constant(c)}," for c in codes],
        "    ];",
        "}",
        "",
    ]
    return "\n".join(out)


# ------------------------------------------------------------------------- main


def main() -> int:
    if len(sys.argv) != 3:
        print(__doc__, file=sys.stderr)
        return 2
    lang, output = sys.argv[1], Path(sys.argv[2])
    models, codes = load_models()
    if lang == "kotlin":
        text = emit_kotlin(models, codes, "io.github.pierreemmanuelj.alldebrid")
    elif lang == "rust":
        text = emit_rust(models, codes)
    elif lang == "swift":
        text = emit_swift(models, codes)
    elif lang == "csharp":
        text = emit_csharp(models, codes, "OpenAllDebrid")
    elif lang == "php":
        # The output is the src directory; models go to src/Model/.
        for old in (output / "Model").glob("*.php"):
            old.unlink()
        for name, content in emit_php(models, codes, "OpenAllDebrid").items():
            path = output / name
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_text(content)
        return 0
    else:
        print(f"unknown language {lang}", file=sys.stderr)
        return 2
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text(text)
    return 0


if __name__ == "__main__":
    sys.exit(main())
