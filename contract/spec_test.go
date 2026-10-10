// Package contract checks openapi.yaml against the live AllDebrid API.
//
// Every response is validated against the response schema of its operation,
// and fields the spec does not declare are reported, so the spec cannot drift
// silently from the API.
package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"go.yaml.in/yaml/v3"
)

const schemaURL = "openapi-schemas.json"

// spec is the parsed openapi.yaml, with its component schemas compiled.
type spec struct {
	doc      map[string]any
	defs     map[string]any // components.schemas, refs rewritten to #/$defs/
	compiler *jsonschema.Compiler
	compiled map[string]*jsonschema.Schema
	mu       sync.Mutex
}

var (
	loadOnce   sync.Once
	loadedSpec *spec
	loadErr    error
)

func specPath() string {
	if p := os.Getenv("OPENAPI_SPEC"); p != "" {
		return p
	}
	return filepath.Join("..", "openapi.yaml")
}

func loadSpec(t *testing.T) *spec {
	t.Helper()
	loadOnce.Do(func() { loadedSpec, loadErr = parseSpec(specPath()) })
	if loadErr != nil {
		t.Fatalf("load spec: %v", loadErr)
	}
	return loadedSpec
}

func parseSpec(path string) (*spec, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var y any
	if err := yaml.Unmarshal(raw, &y); err != nil {
		return nil, err
	}
	// Round-trip through JSON so numbers and maps have JSON types.
	j, err := json.Marshal(y)
	if err != nil {
		return nil, err
	}
	var doc map[string]any
	if err := json.Unmarshal(j, &doc); err != nil {
		return nil, err
	}

	components, _ := doc["components"].(map[string]any)
	schemas, ok := components["schemas"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("no components.schemas")
	}
	defsJSON, err := json.Marshal(schemas)
	if err != nil {
		return nil, err
	}
	defsJSON = bytes.ReplaceAll(defsJSON, []byte(`"#/components/schemas/`), []byte(`"#/$defs/`))
	var defs map[string]any
	if err := json.Unmarshal(defsJSON, &defs); err != nil {
		return nil, err
	}

	root, err := jsonschema.UnmarshalJSON(bytes.NewReader(mustJSON(map[string]any{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"$defs":   defs,
	})))
	if err != nil {
		return nil, err
	}
	c := jsonschema.NewCompiler()
	c.DefaultDraft(jsonschema.Draft2020)
	if err := c.AddResource(schemaURL, root); err != nil {
		return nil, err
	}
	return &spec{doc: doc, defs: defs, compiler: c, compiled: map[string]*jsonschema.Schema{}}, nil
}

func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// responseSchema returns the component name of the 200 response of an operation.
func (s *spec) responseSchema(method, path string) (string, error) {
	paths, _ := s.doc["paths"].(map[string]any)
	item, ok := paths[path].(map[string]any)
	if !ok {
		return "", fmt.Errorf("path %s is not in the spec", path)
	}
	op, ok := item[strings.ToLower(method)].(map[string]any)
	if !ok {
		return "", fmt.Errorf("%s %s is not in the spec", method, path)
	}
	ref := dig(op, "responses", "200", "content", "application/json", "schema", "$ref")
	name, ok := ref.(string)
	if !ok {
		return "", fmt.Errorf("%s %s: the 200 response is not a $ref", method, path)
	}
	return strings.TrimPrefix(name, "#/components/schemas/"), nil
}

func dig(v any, keys ...string) any {
	for _, k := range keys {
		m, ok := v.(map[string]any)
		if !ok {
			return nil
		}
		v = m[k]
	}
	return v
}

func (s *spec) schema(name string) (*jsonschema.Schema, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sch, ok := s.compiled[name]; ok {
		return sch, nil
	}
	sch, err := s.compiler.Compile(schemaURL + "#/$defs/" + name)
	if err != nil {
		return nil, err
	}
	s.compiled[name] = sch
	return sch, nil
}

// validate checks a decoded JSON response against the named schema and
// returns the schema errors and the fields the spec does not declare.
func (s *spec) validate(name string, body any) (schemaErr error, unknown []string, err error) {
	sch, err := s.schema(name)
	if err != nil {
		return nil, nil, err
	}
	if verr := sch.Validate(body); verr != nil {
		schemaErr = verr
	}
	unknown = s.unknownFields(map[string]any{"$ref": "#/$defs/" + name}, body, "$", 0)
	sort.Strings(unknown)
	return schemaErr, unknown, nil
}

// resolve follows a local $ref.
func (s *spec) resolve(node map[string]any) map[string]any {
	for i := 0; i < 16; i++ {
		ref, ok := node["$ref"].(string)
		if !ok {
			return node
		}
		next, ok := s.defs[strings.TrimPrefix(ref, "#/$defs/")].(map[string]any)
		if !ok {
			return node
		}
		node = next
	}
	return node
}

func jsonKind(v any) string {
	switch v.(type) {
	case map[string]any:
		return "object"
	case []any:
		return "array"
	case string:
		return "string"
	case float64, json.Number:
		return "number"
	case bool:
		return "boolean"
	case nil:
		return "null"
	}
	return "unknown"
}

func typeAccepts(node map[string]any, kind string) bool {
	switch t := node["type"].(type) {
	case string:
		return t == kind || (t == "integer" && kind == "number") || (t == "number" && kind == "number")
	case []any:
		for _, x := range t {
			if x == kind || (x == "integer" && kind == "number") {
				return true
			}
		}
		return false
	}
	// No type: accept objects when it declares properties.
	_, hasProps := node["properties"]
	return !hasProps || kind == "object"
}

// unknownFields walks a value alongside its schema and lists object keys the
// schema does not declare.
func (s *spec) unknownFields(node map[string]any, v any, at string, depth int) []string {
	if depth > 64 {
		return nil
	}
	node = s.resolve(node)
	kind := jsonKind(v)

	for _, comb := range []string{"oneOf", "anyOf"} {
		if branches, ok := node[comb].([]any); ok {
			for _, b := range branches {
				bn, _ := b.(map[string]any)
				bn = s.resolve(bn)
				if typeAccepts(bn, kind) {
					return s.unknownFields(bn, v, at, depth+1)
				}
			}
			return nil
		}
	}

	var out []string
	switch val := v.(type) {
	case map[string]any:
		props, _ := node["properties"].(map[string]any)
		addl := node["additionalProperties"]
		for k, child := range val {
			if p, ok := props[k].(map[string]any); ok {
				out = append(out, s.unknownFields(p, child, at+"."+k, depth+1)...)
				continue
			}
			if a, ok := addl.(map[string]any); ok {
				out = append(out, s.unknownFields(a, child, at+"."+k, depth+1)...)
				continue
			}
			if props != nil {
				out = append(out, at+"."+k)
			}
		}
	case []any:
		if items, ok := node["items"].(map[string]any); ok {
			seen := map[string]bool{}
			for _, child := range val {
				for _, u := range s.unknownFields(items, child, at+"[]", depth+1) {
					if !seen[u] {
						seen[u] = true
						out = append(out, u)
					}
				}
			}
		}
	}
	return out
}

func TestSpecLoads(t *testing.T) {
	s := loadSpec(t)
	paths, _ := s.doc["paths"].(map[string]any)
	if len(paths) == 0 {
		t.Fatal("spec has no paths")
	}
	for path, item := range paths {
		for method := range item.(map[string]any) {
			name, err := s.responseSchema(method, path)
			if err != nil {
				t.Errorf("%v", err)
				continue
			}
			if _, err := s.schema(name); err != nil {
				t.Errorf("%s %s: compile %s: %v", method, path, name, err)
			}
		}
	}
}

func TestUnknownFieldsWalker(t *testing.T) {
	s := loadSpec(t)
	body := map[string]any{
		"status": "success",
		"data": map[string]any{
			"magnets": map[string]any{"id": 1.0, "surprise": true, "files": []any{map[string]any{"n": "a", "x": 1.0}}},
			"extra":   1.0,
		},
	}
	_, unknown, err := s.validate("MagnetStatusResponse", body)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"$.data.extra", "$.data.magnets.files[].x", "$.data.magnets.surprise"}
	if strings.Join(unknown, ",") != strings.Join(want, ",") {
		t.Fatalf("unknown fields = %v, want %v", unknown, want)
	}
}

// TestExamplesMatchSchemas validates every response example of the spec
// against its schema, so the examples cannot drift from the shapes.
func TestExamplesMatchSchemas(t *testing.T) {
	s := loadSpec(t)
	paths, _ := s.doc["paths"].(map[string]any)
	count := 0
	for path, item := range paths {
		for method, op := range item.(map[string]any) {
			media, _ := dig(op, "responses", "200", "content", "application/json").(map[string]any)
			name, err := s.responseSchema(method, path)
			if err != nil {
				t.Fatal(err)
			}
			values := map[string]any{}
			if v, ok := media["example"]; ok {
				values["example"] = v
			}
			if exs, ok := media["examples"].(map[string]any); ok {
				for k, ex := range exs {
					values[k] = dig(ex, "value")
				}
			}
			if len(values) == 0 {
				t.Errorf("%s %s has no response example", method, path)
			}
			for k, v := range values {
				count++
				schemaErr, unknown, err := s.validate(name, v)
				if err != nil {
					t.Fatal(err)
				}
				if schemaErr != nil {
					t.Errorf("%s %s example %s does not match %s:\n%v", method, path, k, name, schemaErr)
				}
				if len(unknown) > 0 {
					t.Errorf("%s %s example %s has undeclared fields: %v", method, path, k, unknown)
				}
			}
		}
	}
	t.Logf("%d examples checked", count)
}
