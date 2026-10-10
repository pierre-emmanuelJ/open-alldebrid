package alldebrid

import (
	"regexp"
	"strings"
	"testing"
)

// specOperations returns "METHOD /path" for every operation of the spec.
func specOperations(t *testing.T) map[string]bool {
	t.Helper()
	ops := map[string]bool{}
	path := ""
	pathRe := regexp.MustCompile(`^  (/\S+):$`)
	methodRe := regexp.MustCompile(`^    (get|post):$`)
	for _, line := range strings.Split(readSpec(t), "\n") {
		if m := pathRe.FindStringSubmatch(line); m != nil {
			path = m[1]
		} else if m := methodRe.FindStringSubmatch(line); m != nil && path != "" {
			ops[strings.ToUpper(m[1])+" "+path] = true
		} else if line == "components:" {
			break
		}
	}
	return ops
}

// TestRouting calls every Client method against a fake server and checks it
// uses an HTTP method and path of the spec, with the expected parameters.
func TestRouting(t *testing.T) {
	ops := specOperations(t)
	if len(ops) < 30 {
		t.Fatalf("only %d operations parsed from the spec", len(ops))
	}
	f, c := newFake(t)
	used := map[string]bool{}

	cases := []struct {
		name string
		call func() error
		want string
		form map[string]string
	}{
		{"Ping", func() error { return c.Ping(ctx) }, "GET /v4/ping", nil},
		{"GetPin", func() error { _, err := c.GetPin(ctx); return err }, "GET /v4.1/pin/get", nil},
		{"CheckPin", func() error { _, err := c.CheckPin(ctx, "P", "C"); return err }, "POST /v4/pin/check", map[string]string{"pin": "P", "check": "C"}},
		{"Hosts", func() error { _, err := c.Hosts(ctx, false); return err }, "GET /v4/hosts", nil},
		{"HostDomains", func() error { _, err := c.HostDomains(ctx); return err }, "GET /v4/hosts/domains", nil},
		{"HostPriority", func() error { _, err := c.HostPriority(ctx); return err }, "GET /v4/hosts/priority", nil},
		{"User", func() error { _, err := c.User(ctx); return err }, "GET /v4/user", nil},
		{"UserHosts", func() error { _, err := c.UserHosts(ctx, false); return err }, "GET /v4.1/user/hosts", nil},
		{"VerificationStatus", func() error { _, err := c.VerificationStatus(ctx, "T"); return err }, "POST /v4/user/verif", map[string]string{"token": "T"}},
		{"ResendVerification", func() error { return c.ResendVerification(ctx, "T") }, "POST /v4/user/verif/resend", map[string]string{"token": "T"}},
		{"ClearNotification", func() error { return c.ClearNotification(ctx, "N") }, "POST /v4/user/notification/clear", map[string]string{"code": "N"}},
		{"LinkInfos", func() error { _, err := c.LinkInfos(ctx, "pw", "L"); return err }, "POST /v4/link/infos", map[string]string{"link[]": "L", "password": "pw"}},
		{"ExtractRedirector", func() error { _, err := c.ExtractRedirector(ctx, "L"); return err }, "POST /v4/link/redirector", map[string]string{"link": "L"}},
		{"Unlock", func() error { _, err := c.Unlock(ctx, "L", ""); return err }, "POST /v4/link/unlock", map[string]string{"link": "L"}},
		{"StreamingLink", func() error { _, err := c.StreamingLink(ctx, "G", "S"); return err }, "POST /v4/link/streaming", map[string]string{"id": "G", "stream": "S"}},
		{"DelayedLink", func() error { _, err := c.DelayedLink(ctx, "D"); return err }, "POST /v4/link/delayed", map[string]string{"id": "D"}},
		{"UploadMagnets", func() error { _, err := c.UploadMagnets(ctx, "M"); return err }, "POST /v4/magnet/upload", map[string]string{"magnets[]": "M"}},
		{"UploadTorrentFiles", func() error {
			_, err := c.UploadTorrentFiles(ctx, TorrentFile{Name: "a.torrent", Content: strings.NewReader("x")})
			return err
		}, "POST /v4/magnet/upload/file", nil},
		{"Magnets", func() error { _, err := c.Magnets(ctx); return err }, "POST /v4.1/magnet/status", nil},
		{"MagnetFiles", func() error { _, err := c.MagnetFiles(ctx, 5); return err }, "POST /v4/magnet/files", map[string]string{"id[]": "5"}},
		{"DeleteMagnet", func() error { return c.DeleteMagnet(ctx, 5) }, "POST /v4/magnet/delete", map[string]string{"id": "5"}},
		{"RestartMagnet", func() error { return c.RestartMagnet(ctx, 5) }, "POST /v4/magnet/restart", map[string]string{"id": "5"}},
		{"RestartMagnets", func() error { _, err := c.RestartMagnets(ctx, 5); return err }, "POST /v4/magnet/restart", map[string]string{"ids[]": "5"}},
		{"SavedLinks", func() error { _, err := c.SavedLinks(ctx); return err }, "GET /v4/user/links", nil},
		{"SaveLinks", func() error { return c.SaveLinks(ctx, "L") }, "POST /v4/user/links/save", map[string]string{"links[]": "L"}},
		{"DeleteSavedLinks", func() error { return c.DeleteSavedLinks(ctx, "L") }, "POST /v4/user/links/delete", map[string]string{"links[]": "L"}},
		{"History", func() error { _, err := c.History(ctx); return err }, "GET /v4/user/history", nil},
		{"PurgeHistory", func() error { return c.PurgeHistory(ctx) }, "POST /v4/user/history/delete", nil},
		{"VoucherBalance", func() error { _, err := c.VoucherBalance(ctx); return err }, "GET /v4/voucher/balance", nil},
		{"Vouchers", func() error { _, err := c.Vouchers(ctx, 30, 2); return err }, "POST /v4/voucher/get", map[string]string{"duration": "30", "nb": "2"}},
		{"GenerateVouchers", func() error { _, err := c.GenerateVouchers(ctx, 90, 1); return err }, "POST /v4/voucher/generate", map[string]string{"duration": "90", "nb": "1"}},
	}
	for _, tc := range cases {
		if err := tc.call(); err != nil {
			t.Errorf("%s: %v", tc.name, err)
			continue
		}
		req := f.last()
		got := req.method + " " + req.path
		if got != tc.want {
			t.Errorf("%s called %s, want %s", tc.name, got, tc.want)
		}
		if !ops[got] {
			t.Errorf("%s called %s, which is not in the spec", tc.name, got)
		}
		used[got] = true
		for k, v := range tc.form {
			if req.form.Get(k) != v {
				t.Errorf("%s: form %s = %q, want %q", tc.name, k, req.form.Get(k), v)
			}
		}
	}
	for op := range ops {
		if !used[op] {
			t.Errorf("no client method calls %s", op)
		}
	}
}

func TestMagnetHelpers(t *testing.T) {
	if !(Magnet{StatusCode: MagnetReady}).IsReady() || (Magnet{StatusCode: MagnetDownloading}).IsReady() {
		t.Error("IsReady")
	}
	if !(Magnet{StatusCode: 7}).IsError() || (Magnet{StatusCode: MagnetReady}).IsError() {
		t.Error("IsError")
	}
}
