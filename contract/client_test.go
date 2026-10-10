package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	demoPremium = "staticDemoApikeyPrem"
	demoTrial   = "staticDemoApikeyTria"
	demoFree    = "staticDemoApikeyFree"
)

func baseURL() string {
	if u := os.Getenv("ALLDEBRID_BASE_URL"); u != "" {
		return strings.TrimRight(u, "/")
	}
	return "https://api.alldebrid.com"
}

var httpClient = &http.Client{Timeout: 30 * time.Second}

// throttle keeps the suite well under the 12 requests per second limit.
var throttle = time.NewTicker(150 * time.Millisecond)

// upload is one file of a multipart request.
type upload struct {
	field, name string
	content     []byte
}

// request describes one API call of a test case.
type request struct {
	method string
	path   string
	key    string
	form   [][2]string
	files  []upload
}

// result is a decoded API response.
type result struct {
	httpStatus int
	body       map[string]any
}

func (r result) status() string { s, _ := r.body["status"].(string); return s }

func (r result) errorCode() string {
	code, _ := dig(r.body, "error", "code").(string)
	return code
}

func (r result) data() map[string]any { d, _ := r.body["data"].(map[string]any); return d }

func do(t *testing.T, req request) result {
	t.Helper()
	<-throttle.C

	var body io.Reader
	contentType := ""
	target := baseURL() + req.path
	switch {
	case len(req.files) > 0:
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)
		for _, kv := range req.form {
			_ = w.WriteField(kv[0], kv[1])
		}
		for _, f := range req.files {
			part, err := w.CreateFormFile(f.field, f.name)
			if err != nil {
				t.Fatal(err)
			}
			_, _ = part.Write(f.content)
		}
		_ = w.Close()
		body, contentType = &buf, w.FormDataContentType()
	case req.method == http.MethodGet && len(req.form) > 0:
		q := url.Values{}
		for _, kv := range req.form {
			q.Add(kv[0], kv[1])
		}
		target += "?" + q.Encode()
	case req.method == http.MethodPost:
		q := url.Values{}
		for _, kv := range req.form {
			q.Add(kv[0], kv[1])
		}
		body, contentType = strings.NewReader(q.Encode()), "application/x-www-form-urlencoded"
	}

	hreq, err := http.NewRequest(req.method, target, body)
	if err != nil {
		t.Fatal(err)
	}
	hreq.Header.Set("User-Agent", "open-alldebrid-contract-tests")
	if contentType != "" {
		hreq.Header.Set("Content-Type", contentType)
	}
	if req.key != "" {
		hreq.Header.Set("Authorization", "Bearer "+req.key)
	}
	resp, err := httpClient.Do(hreq)
	if err != nil {
		// The error may contain the URL but never the key, which is in a header.
		t.Fatalf("%s %s: %v", req.method, req.path, err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	dec := json.NewDecoder(bytes.NewReader(raw))
	if err := dec.Decode(&decoded); err != nil {
		t.Fatalf("%s %s: HTTP %d, body is not JSON (%d bytes)", req.method, req.path, resp.StatusCode, len(raw))
	}
	return result{httpStatus: resp.StatusCode, body: decoded}
}

// check validates a response against the spec and fails on undeclared fields.
func check(t *testing.T, req request, res result) {
	t.Helper()
	s := loadSpec(t)
	if res.httpStatus != http.StatusOK {
		t.Errorf("HTTP %d, want 200", res.httpStatus)
	}
	name, err := s.responseSchema(req.method, req.path)
	if err != nil {
		t.Fatal(err)
	}
	schemaErr, unknown, err := s.validate(name, res.body)
	if err != nil {
		t.Fatal(err)
	}
	if schemaErr != nil {
		t.Errorf("response does not match %s:\n%v", name, schemaErr)
	}
	if len(unknown) > 0 {
		t.Errorf("response has fields %s does not declare: %s", name, strings.Join(unknown, ", "))
	}
}

func expect(t *testing.T, res result, wantError string) {
	t.Helper()
	switch {
	case wantError == "" && res.status() != "success":
		t.Errorf("status = %q (error %q), want success", res.status(), res.errorCode())
	case wantError != "" && res.errorCode() != wantError:
		t.Errorf("error = %q (status %q), want %s", res.errorCode(), res.status(), wantError)
	}
}

func describe(req request) string { return fmt.Sprintf("%s %s", req.method, req.path) }
