package contract

import (
	"net/http"
	"testing"
)

// The demo cases below are the "Demo" buttons of https://docs.alldebrid.com:
// fixed inputs that make the API answer static sample data, with the public
// demo keys. They need no account, so they run in public CI.

const (
	pinCheck   = "664c3ca2635c99f291d28e11ea18e154750bd21a"
	demoMagnet = "magnet:?xt=urn:btih:3A842783e3005495d5d1637f5364b59343c7844707&dn=ubuntu-18.04.2-live-server-amd64.iso"
)

type demoCase struct {
	name    string
	req     request
	wantErr string
}

func post(path, key string, form ...[2]string) request {
	return request{method: http.MethodPost, path: path, key: key, form: form}
}

func get(path, key string, query ...[2]string) request {
	return request{method: http.MethodGet, path: path, key: key, form: query}
}

func kv(k, v string) [2]string { return [2]string{k, v} }

var demoCases = []demoCase{
	{"ping", get("/v4/ping", ""), ""},

	{"pin/get", get("/v4.1/pin/get", "", kv("demo", "")), ""},
	{"pin/check waiting", post("/v4/pin/check", "", kv("pin", "123A"), kv("check", pinCheck)), ""},
	{"pin/check activated", post("/v4/pin/check", "", kv("pin", "123B"), kv("check", pinCheck)), ""},
	{"pin/check expired", post("/v4/pin/check", "", kv("pin", "123C"), kv("check", pinCheck)), "PIN_EXPIRED"},
	{"pin/check invalid", post("/v4/pin/check", "", kv("pin", "123D"), kv("check", pinCheck)), "PIN_INVALID"},

	{"hosts", get("/v4/hosts", ""), ""},
	{"hosts hostsOnly", get("/v4/hosts", "", kv("hostsOnly", "true")), ""},
	{"hosts/domains", get("/v4/hosts/domains", ""), ""},
	{"hosts/priority", get("/v4/hosts/priority", ""), ""},

	{"user premium", get("/v4/user", demoPremium), ""},
	{"user trial", get("/v4/user", demoTrial), ""},
	{"user free", get("/v4/user", demoFree), ""},
	{"user/hosts", get("/v4.1/user/hosts", demoPremium), ""},
	{"user/verif waiting", post("/v4/user/verif", "", kv("token", "verificationTokenWaiting3f9eb0db")), ""},
	{"user/verif allowed", post("/v4/user/verif", "", kv("token", "verificationTokenAllowed3f9eb0db")), ""},
	{"user/verif denied", post("/v4/user/verif", "", kv("token", "verificationTokenDenied43f9eb0db")), ""},
	{"user/verif expired", post("/v4/user/verif", "", kv("token", "verificationTokenExpired3f9eb0db")), "PIN_EXPIRED"},
	{"user/verif/resend", post("/v4/user/verif/resend", "", kv("token", "verificationTokenResent43f9eb0db")), ""},
	{"user/verif/resend already sent", post("/v4/user/verif/resend", "", kv("token", "verificationTokenAlreadySentb0db")), "ALREADY_SENT"},
	{"user/verif/resend expired", post("/v4/user/verif/resend", "", kv("token", "verificationTokenExpired3f9eb0db")), "PIN_EXPIRED"},
	{"user/notification/clear", post("/v4/user/notification/clear", demoPremium, kv("code", "NOTIF_CODE")), ""},

	{"link/infos", post("/v4/link/infos", demoPremium, kv("link[]", "https://example.com/file2MB.txt")), ""},
	{"link/infos per-link errors", post("/v4/link/infos", demoPremium,
		kv("link[]", "https://example.com/down"), kv("link[]", "https://not-supported.com/somefile.txt")), ""},
	{"link/redirector", post("/v4/link/redirector", demoPremium, kv("link", "https://example.com/redir/folder")), ""},
	{"link/redirector not supported", post("/v4/link/redirector", demoPremium, kv("link", "https://example.com/redir/notSupported")), "REDIRECTOR_NOT_SUPPORTED"},
	{"link/redirector rate limited", post("/v4/link/redirector", demoPremium, kv("link", "https://example.com/redir/rateLimited")), "REDIRECTOR_RATE_LIMITED"},
	{"link/redirector error", post("/v4/link/redirector", demoPremium, kv("link", "https://example.com/redir/error")), "REDIRECTOR_ERROR"},
	{"link/unlock", post("/v4/link/unlock", demoPremium, kv("link", "https://example.com/file2MB.txt")), ""},
	{"link/unlock delayed", post("/v4/link/unlock", demoPremium, kv("link", "https://example.com/delayed")), ""},
	{"link/unlock streaming", post("/v4/link/unlock", demoPremium, kv("link", "https://example.com/streaming")), ""},
	{"link/unlock not supported", post("/v4/link/unlock", demoPremium, kv("link", "https://notsupported.com/file.txt")), "LINK_HOST_NOT_SUPPORTED"},
	{"link/unlock down", post("/v4/link/unlock", demoPremium, kv("link", "https://example.com/down")), "LINK_DOWN"},
	{"link/streaming", post("/v4/link/streaming", demoPremium, kv("id", "demov1bc797"), kv("stream", "hd-720")), ""},
	{"link/streaming invalid gen", post("/v4/link/streaming", demoPremium, kv("id", "demoInvalid"), kv("stream", "hd-720")), "STREAM_INVALID_GEN_ID"},
	{"link/streaming invalid stream", post("/v4/link/streaming", demoPremium, kv("id", "demov1bc797"), kv("stream", "bad_stream")), "STREAM_INVALID_STREAM_ID"},
	{"link/delayed processing", post("/v4/link/delayed", demoPremium, kv("id", "123456")), ""},
	{"link/delayed ready", post("/v4/link/delayed", demoPremium, kv("id", "1234567")), ""},
	{"link/delayed invalid", post("/v4/link/delayed", demoPremium, kv("id", "0")), "DELAYED_INVALID_ID"},

	{"magnet/upload", post("/v4/magnet/upload", demoPremium, kv("magnets[]", demoMagnet)), ""},
	{"magnet/upload nothing", post("/v4/magnet/upload", demoPremium), "MAGNET_NO_URI"},
	{"magnet/upload/file", request{method: http.MethodPost, path: "/v4/magnet/upload/file", key: demoPremium,
		files: []upload{{"files[]", "ubuntu-18.04.2-live-server-amd64.iso.torrent", []byte("d4:infod4:name4:demoee")}}}, ""},
	{"magnet/upload/file nothing", post("/v4/magnet/upload/file", demoPremium), "MAGNET_FILE_UPLOAD_FAILED"},
	{"magnet/status all", post("/v4.1/magnet/status", demoPremium), ""},
	{"magnet/status one", post("/v4.1/magnet/status", demoPremium, kv("id", "56789")), ""},
	{"magnet/status active", post("/v4.1/magnet/status", demoPremium, kv("status", "active")), ""},
	{"magnet/status live fullsync", post("/v4.1/magnet/status", demoPremium, kv("session", "1234"), kv("counter", "0")), ""},
	{"magnet/status live delta", post("/v4.1/magnet/status", demoPremium, kv("session", "1234"), kv("counter", "1")), ""},
	{"magnet/files", post("/v4/magnet/files", demoPremium, kv("id[]", "56789")), ""},
	{"magnet/files nothing", post("/v4/magnet/files", demoPremium), "MAGNET_INVALID_ID"},
	{"magnet/delete", post("/v4/magnet/delete", demoPremium, kv("id", "56789")), ""},
	{"magnet/delete nothing", post("/v4/magnet/delete", demoPremium), "MAGNET_INVALID_ID"},
	{"magnet/restart", post("/v4/magnet/restart", demoPremium, kv("id", "99999")), ""},
	{"magnet/restart nothing", post("/v4/magnet/restart", demoPremium), "MAGNET_INVALID_ID"},

	{"user/links", get("/v4/user/links", demoPremium), ""},
	{"user/links/save", post("/v4/user/links/save", demoPremium,
		kv("links[]", "https://example.com/file2MB.txt"), kv("links[]", "https://example.com/file10MB.txt")), ""},
	{"user/links/delete", post("/v4/user/links/delete", demoPremium, kv("links[]", "https://example.com/file2MB.txt")), ""},
	{"user/history", get("/v4/user/history", demoPremium), ""},
	{"user/history/delete", post("/v4/user/history/delete", demoPremium), ""},

	{"no key", get("/v4/user", ""), "AUTH_MISSING_APIKEY"},
	{"bad key", get("/v4/user", "not-a-real-key"), "AUTH_BAD_APIKEY"},
}

func TestDemo(t *testing.T) {
	if testing.Short() {
		t.Skip("calls the live API")
	}
	loadSpec(t)
	for _, tc := range demoCases {
		t.Run(tc.name, func(t *testing.T) {
			res := do(t, tc.req)
			expect(t, res, tc.wantErr)
			check(t, tc.req, res)
		})
	}
}

// TestDemoCoversEveryOperation keeps the demo suite in step with the spec:
// adding an operation without a demo case fails here.
func TestDemoCoversEveryOperation(t *testing.T) {
	s := loadSpec(t)
	covered := map[string]bool{}
	for _, tc := range demoCases {
		covered[describe(tc.req)] = true
	}
	// Reseller endpoints have no demo data.
	skip := map[string]bool{
		"GET /v4/voucher/balance":   true,
		"POST /v4/voucher/get":      true,
		"POST /v4/voucher/generate": true,
	}
	paths, _ := s.doc["paths"].(map[string]any)
	for path, item := range paths {
		for method := range item.(map[string]any) {
			op := describe(request{method: map[string]string{"get": "GET", "post": "POST"}[method], path: path})
			if !covered[op] && !skip[op] {
				t.Errorf("no demo case for %s", op)
			}
		}
	}
}
