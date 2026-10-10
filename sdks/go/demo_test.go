package alldebrid

import (
	"errors"
	"strings"
	"testing"
)

// TestDemoAPI calls the real API with the public demo keys. Skip it with -short.
func TestDemoAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("calls the live API")
	}
	c := New(DemoPremiumKey, WithUserAgent("open-alldebrid-go-tests"))
	public := New("", WithUserAgent("open-alldebrid-go-tests"))

	step := func(name string, f func() error) {
		t.Helper()
		if err := f(); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}

	step("Ping", func() error { return public.Ping(ctx) })
	step("Hosts", func() error {
		h, err := public.Hosts(ctx, false)
		if err == nil && (len(h.Hosts) == 0 || len(h.Redirectors) == 0) {
			err = errors.New("empty hosts")
		}
		return err
	})
	step("HostDomains", func() error { _, err := public.HostDomains(ctx); return err })
	step("HostPriority", func() error { _, err := public.HostPriority(ctx); return err })

	step("User", func() error {
		u, err := c.User(ctx)
		if err == nil && (!u.IsPremium || u.Username != "demoUserPremium") {
			err = errors.New("not the premium demo user")
		}
		return err
	})
	step("UserHosts", func() error { _, err := c.UserHosts(ctx, false); return err })
	step("ClearNotification", func() error { return c.ClearNotification(ctx, "NOTIF_CODE") })
	step("VerificationStatus", func() error {
		v, err := public.VerificationStatus(ctx, "verificationTokenAllowed3f9eb0db")
		if err == nil && (v.Verif != VerificationDataVerifAllowed || v.Apikey == "") {
			err = errors.New("verification not allowed")
		}
		return err
	})
	step("ResendVerification", func() error {
		err := public.ResendVerification(ctx, "verificationTokenAlreadySentb0db")
		if !errors.Is(err, CodeAlreadySent) {
			return errors.New("want ALREADY_SENT, got " + errString(err))
		}
		return nil
	})

	step("GetPin+CheckPin", func() error {
		p, err := public.GetPin(ctx)
		if err != nil {
			return err
		}
		if p.Pin == "" || !strings.HasPrefix(p.UserURL, "https://") {
			return errors.New("bad pin")
		}
		d, err := public.CheckPin(ctx, "123B", "664c3ca2635c99f291d28e11ea18e154750bd21a")
		if err == nil && (!d.Activated || d.Apikey == "") {
			err = errors.New("demo pin not activated")
		}
		return err
	})

	step("LinkInfos", func() error {
		infos, err := c.LinkInfos(ctx, "", "https://example.com/file2MB.txt", "https://example.com/down")
		if err == nil && (len(infos) == 0 || infos[0].Filename == "") {
			err = errors.New("no file info")
		}
		return err
	})
	step("ExtractRedirector", func() error { _, err := c.ExtractRedirector(ctx, "https://example.com/redir/folder"); return err })
	step("Unlock", func() error {
		u, err := c.Unlock(ctx, "https://example.com/streaming", "")
		if err == nil && len(u.Streaming) == 0 {
			err = errors.New("no streams")
		}
		return err
	})
	step("Unlock delayed", func() error {
		u, err := c.Unlock(ctx, "https://example.com/delayed", "")
		if err == nil && u.Delayed == "" {
			err = errors.New("no delayed id")
		}
		return err
	})
	step("Unlock down", func() error {
		if _, err := c.Unlock(ctx, "https://example.com/down", ""); !errors.Is(err, CodeLinkDown) {
			return errors.New("want LINK_DOWN, got " + errString(err))
		}
		return nil
	})
	step("StreamingLink", func() error { _, err := c.StreamingLink(ctx, "demov1bc797", "hd-720"); return err })
	step("DelayedLink", func() error {
		d, err := c.DelayedLink(ctx, "1234567")
		if err == nil && (int(d.Status) != DelayedReady || d.Link == "") {
			err = errors.New("delayed link not ready")
		}
		return err
	})

	step("UploadMagnets", func() error {
		m, err := c.UploadMagnets(ctx, "magnet:?xt=urn:btih:3A842783e3005495d5d1637f5364b59343c7844707")
		if err == nil && (len(m) != 1 || m[0].ID == 0) {
			err = errors.New("no magnet id")
		}
		return err
	})
	step("UploadTorrentFiles", func() error {
		_, err := c.UploadTorrentFiles(ctx, TorrentFile{Name: "demo.torrent", Content: strings.NewReader("d4:infod4:name4:demoee")})
		return err
	})
	step("Magnets", func() error {
		ms, err := c.Magnets(ctx)
		if err == nil && len(ms) == 0 {
			err = errors.New("no magnets")
		}
		return err
	})
	step("Magnet", func() error {
		m, err := c.Magnet(ctx, 56789)
		if err == nil && (!m.IsReady() || len(m.Files) == 0) {
			err = errors.New("magnet not ready or without files")
		}
		return err
	})
	step("MagnetFiles", func() error { _, err := c.MagnetFiles(ctx, 56789); return err })
	step("DeleteMagnet", func() error { return c.DeleteMagnet(ctx, 56789) })
	step("RestartMagnet", func() error { return c.RestartMagnet(ctx, 99999) })
	step("LiveMagnets", func() error {
		live := c.LiveMagnets()
		ms, err := live.Poll(ctx)
		if err == nil && len(ms) == 0 {
			err = errors.New("empty full sync")
		}
		return err
	})

	step("SavedLinks", func() error { _, err := c.SavedLinks(ctx); return err })
	step("SaveLinks", func() error { return c.SaveLinks(ctx, "https://example.com/file2MB.txt") })
	step("DeleteSavedLinks", func() error { return c.DeleteSavedLinks(ctx, "https://example.com/file2MB.txt") })
	step("History", func() error { _, err := c.History(ctx); return err })
	step("PurgeHistory", func() error { return c.PurgeHistory(ctx) })

	step("bad key", func() error {
		if _, err := c.WithAPIKey("not-a-real-key").User(ctx); !errors.Is(err, CodeAuthBadAPIKey) {
			return errors.New("want AUTH_BAD_APIKEY, got " + errString(err))
		}
		return nil
	})
}

func errString(err error) string {
	if err == nil {
		return "nil"
	}
	return err.Error()
}
