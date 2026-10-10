package contract

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"testing"

	"go.yaml.in/yaml/v3"
)

// TestCodeSamples checks docs/code-samples.yaml: one sample per operation
// and language, each calling client methods that exist in that SDK.
func TestCodeSamples(t *testing.T) {
	s := loadSpec(t)
	raw, err := os.ReadFile(filepath.Join("..", "docs", "code-samples.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	var samples map[string]map[string]string
	if err := yaml.Unmarshal(raw, &samples); err != nil {
		t.Fatal(err)
	}

	methods := map[string]map[string]bool{
		"go":         sdkMethods(t, `func \(c \*Client\) (\w+)\(`, "sdks/go/*.go"),
		"typescript": sdkMethods(t, `(?m)^  (?:async )?(\w+)\(`, "sdks/typescript/src/index.ts"),
		"python":     sdkMethods(t, `(?m)^    (?:async )?def (\w+)\(`, "sdks/python/src/open_alldebrid/_async_client.py"),
		"php":        sdkMethods(t, `(?m)^    public function (\w+)\(`, "sdks/php/src/AllDebrid.php"),
		"kotlin":     sdkMethods(t, `(?m)^    public (?:suspend )?fun (\w+)\(`, "sdks/kotlin/src/main/kotlin/io/github/pierreemmanuelj/alldebrid/AllDebrid.kt"),
	}
	call := regexp.MustCompile(`\bclient(?:\.|->)(\w+)\(`)

	var ops []string
	paths, _ := s.doc["paths"].(map[string]any)
	for _, item := range paths {
		for _, op := range item.(map[string]any) {
			if id, ok := dig(op, "operationId").(string); ok {
				ops = append(ops, id)
			}
		}
	}
	sort.Strings(ops)

	for _, op := range ops {
		perLang, ok := samples[op]
		if !ok {
			t.Errorf("%s has no code sample", op)
			continue
		}
		for lang, known := range methods {
			code := perLang[lang]
			calls := call.FindAllStringSubmatch(code, -1)
			if len(calls) == 0 {
				t.Errorf("%s: the %s sample calls no client method", op, lang)
			}
			for _, c := range calls {
				if !known[c[1]] {
					t.Errorf("%s: the %s sample calls client.%s, which the SDK does not have", op, lang, c[1])
				}
			}
		}
	}
	for op := range samples {
		if !contains(ops, op) {
			t.Errorf("code sample for %s, which is not an operation of the spec", op)
		}
	}
}

func sdkMethods(t *testing.T, pattern, glob string) map[string]bool {
	t.Helper()
	files, err := filepath.Glob(filepath.Join("..", glob))
	if err != nil || len(files) == 0 {
		t.Fatalf("no files for %s", glob)
	}
	re := regexp.MustCompile(pattern)
	found := map[string]bool{}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range re.FindAllStringSubmatch(string(b), -1) {
			found[m[1]] = true
		}
	}
	return found
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
