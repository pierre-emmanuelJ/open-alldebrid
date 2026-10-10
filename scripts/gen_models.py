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
Languages: kotlin, csharp, swift
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


# ------------------------------------------------------------------------- main


def main() -> int:
    if len(sys.argv) != 3:
        print(__doc__, file=sys.stderr)
        return 2
    lang, output = sys.argv[1], Path(sys.argv[2])
    models, codes = load_models()
    if lang == "kotlin":
        text = emit_kotlin(models, codes, "io.github.pierreemmanuelj.alldebrid")
    else:
        print(f"unknown language {lang}", file=sys.stderr)
        return 2
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text(text)
    return 0


if __name__ == "__main__":
    sys.exit(main())
