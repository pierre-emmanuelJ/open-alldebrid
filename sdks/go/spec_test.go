package alldebrid

import (
	"bufio"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

const specFile = "../../openapi.yaml"

// methods maps every operationId of the spec to the Client method that
// implements it.
var methods = map[string][]string{
	"ping":                  {"Ping"},
	"getPin":                {"GetPin", "WaitForPin"},
	"checkPin":              {"CheckPin"},
	"getHosts":              {"Hosts"},
	"getHostDomains":        {"HostDomains"},
	"getHostPriority":       {"HostPriority"},
	"getUser":               {"User"},
	"getUserHosts":          {"UserHosts"},
	"getVerificationStatus": {"VerificationStatus"},
	"resendVerification":    {"ResendVerification"},
	"clearNotification":     {"ClearNotification"},
	"getLinkInfos":          {"LinkInfos"},
	"extractRedirector":     {"ExtractRedirector"},
	"unlockLink":            {"Unlock"},
	"getStreamingLink":      {"StreamingLink"},
	"getDelayedLink":        {"DelayedLink", "WaitDelayedLink"},
	"uploadMagnets":         {"UploadMagnets"},
	"uploadTorrentFiles":    {"UploadTorrentFiles"},
	"getMagnetStatus":       {"Magnets", "Magnet", "LiveMagnets"},
	"getMagnetFiles":        {"MagnetFiles"},
	"deleteMagnet":          {"DeleteMagnet"},
	"restartMagnets":        {"RestartMagnet", "RestartMagnets"},
	"getSavedLinks":         {"SavedLinks"},
	"saveLinks":             {"SaveLinks"},
	"deleteSavedLinks":      {"DeleteSavedLinks"},
	"getHistory":            {"History"},
	"purgeHistory":          {"PurgeHistory"},
	"getVoucherBalance":     {"VoucherBalance"},
	"getVouchers":           {"Vouchers"},
	"generateVouchers":      {"GenerateVouchers"},
}

func readSpec(t *testing.T) string {
	t.Helper()
	b, err := os.ReadFile(specFile)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestEveryOperationHasAMethod(t *testing.T) {
	ops := regexp.MustCompile(`(?m)^\s+operationId: (\w+)$`).FindAllStringSubmatch(readSpec(t), -1)
	if len(ops) == 0 {
		t.Fatal("no operationId in the spec")
	}
	client := reflect.TypeOf(&Client{})
	seen := map[string]bool{}
	for _, m := range ops {
		op := m[1]
		seen[op] = true
		names, ok := methods[op]
		if !ok {
			t.Errorf("operation %s has no Client method", op)
			continue
		}
		for _, name := range names {
			if _, ok := client.MethodByName(name); !ok {
				t.Errorf("operation %s: Client.%s does not exist", op, name)
			}
		}
	}
	for op := range methods {
		if !seen[op] {
			t.Errorf("method table lists %s, which is not in the spec", op)
		}
	}
}

// TestErrorCodesMatchSpec keeps the Code* constants in step with the
// ErrorCode enum of the spec.
func TestErrorCodesMatchSpec(t *testing.T) {
	var fromSpec []string
	sc := bufio.NewScanner(strings.NewReader(readSpec(t)))
	in := false
	for sc.Scan() {
		line := sc.Text()
		switch {
		case line == "    ErrorCode:":
			in = true
		case in && strings.HasPrefix(line, "    ") && !strings.HasPrefix(line, "     "):
			in = false
		case in && strings.HasPrefix(line, "        - "):
			fromSpec = append(fromSpec, strings.Trim(strings.TrimPrefix(line, "        - "), `"`))
		}
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "errors.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	var inGo []string
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			vs := spec.(*ast.ValueSpec)
			if id, ok := vs.Type.(*ast.Ident); ok && id.Name == "ErrorCode" {
				inGo = append(inGo, strings.Trim(vs.Values[0].(*ast.BasicLit).Value, `"`))
			}
		}
	}
	sort.Strings(fromSpec)
	sort.Strings(inGo)
	if len(fromSpec) == 0 || !reflect.DeepEqual(fromSpec, inGo) {
		t.Errorf("error codes differ\nspec: %v\ngo:   %v", fromSpec, inGo)
	}
}
