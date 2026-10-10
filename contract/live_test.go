package contract

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"testing"
)

// TestLive runs against a real account when ALLDEBRID_API_KEY is set. It only
// reads, unless ALLDEBRID_LIVE_WRITE=1, which also uploads a test magnet and a
// test torrent and deletes them again. It never purges the history and never
// clears notifications.
//
// Failures name fields and types, never the key; response values may appear
// in schema errors, so do not run it where the output is public.
func TestLive(t *testing.T) {
	key := os.Getenv("ALLDEBRID_API_KEY")
	if key == "" {
		t.Skip("ALLDEBRID_API_KEY is not set")
	}
	loadSpec(t)

	run := func(name string, req request, wantErr string) result {
		t.Helper()
		var res result
		t.Run(name, func(t *testing.T) {
			res = do(t, req)
			expect(t, res, wantErr)
			check(t, req, res)
		})
		return res
	}

	run("ping", get("/v4/ping", ""), "")
	run("pin/get", get("/v4.1/pin/get", ""), "")
	run("hosts", get("/v4/hosts", key), "")
	run("hosts/domains", get("/v4/hosts/domains", ""), "")
	run("hosts/priority", get("/v4/hosts/priority", ""), "")
	run("user", get("/v4/user", key), "")
	run("user/hosts", get("/v4.1/user/hosts", key), "")
	run("user/links", get("/v4/user/links", key), "")
	run("user/history", get("/v4/user/history", key), "")
	run("voucher/balance", get("/v4/voucher/balance", key), "GENERIC")

	all := run("magnet/status", post("/v4.1/magnet/status", key), "")
	var ids []string
	if list, ok := all.data()["magnets"].([]any); ok {
		for _, m := range list {
			if id, ok := dig(m, "id").(float64); ok {
				ids = append(ids, strconv.FormatInt(int64(id), 10))
			}
		}
	}
	if len(ids) > 0 {
		run("magnet/status one", post("/v4.1/magnet/status", key, kv("id", ids[0])), "")
		var form [][2]string
		for _, id := range ids[:min(len(ids), 20)] {
			form = append(form, kv("id[]", id))
		}
		run("magnet/files", post("/v4/magnet/files", key, form...), "")
	}
	session := strconv.Itoa(rand.IntN(1 << 30))
	full := run("magnet/status live fullsync", post("/v4.1/magnet/status", key, kv("session", session), kv("counter", "0")), "")
	if c, ok := full.data()["counter"].(float64); ok {
		run("magnet/status live delta", post("/v4.1/magnet/status", key,
			kv("session", session), kv("counter", strconv.Itoa(int(c)))), "")
	}

	run("magnet/files unknown id", post("/v4/magnet/files", key, kv("id[]", "1")), "")
	run("magnet/delete unknown id", post("/v4/magnet/delete", key, kv("id", "1")), "MAGNET_INVALID_ID")
	run("magnet/restart unknown ids", post("/v4/magnet/restart", key, kv("ids[]", "1"), kv("ids[]", "2")), "")
	run("magnet/upload nothing", post("/v4/magnet/upload", key), "MAGNET_NO_URI")
	run("link/infos unsupported", post("/v4/link/infos", key, kv("link[]", "https://not-supported.invalid/file")), "")
	run("link/infos nothing", post("/v4/link/infos", key), "LINK_IS_MISSING")
	run("link/unlock unsupported", post("/v4/link/unlock", key, kv("link", "https://not-supported.invalid/file")), "LINK_HOST_NOT_SUPPORTED")
	run("link/redirector unsupported", post("/v4/link/redirector", key, kv("link", "https://not-supported.invalid/x")), "REDIRECTOR_NOT_SUPPORTED")
	run("link/delayed unknown", post("/v4/link/delayed", key, kv("id", "1")), "DELAYED_INVALID_ID")
	run("link/streaming unknown", post("/v4/link/streaming", key, kv("id", "x"), kv("stream", "y")), "STREAM_INVALID_GEN_ID")
	run("pin/check invalid", post("/v4/pin/check", "", kv("pin", "ZZZZ"), kv("check", "x")), "PIN_INVALID")
	run("user/verif invalid", post("/v4/user/verif", "", kv("token", "x")), "PIN_EXPIRED")
	run("user/links/delete nothing", post("/v4/user/links/delete", key), "USER_LINK_MISSING")

	if os.Getenv("ALLDEBRID_LIVE_WRITE") != "1" {
		return
	}
	liveWrites(t, key, run)
}

// bigBuckBunny is the info hash of the Creative Commons "Big Buck Bunny"
// torrent used by the WebTorrent project as its sample.
const bigBuckBunny = "dd8255ecdc7ca55fb0bbf81323d87062db1f6d1c"

func liveWrites(t *testing.T, key string, run func(string, request, string) result) {
	var created []string
	defer func() {
		for _, id := range created {
			run("cleanup magnet/delete "+id, post("/v4/magnet/delete", key, kv("id", id)), "")
		}
	}()
	remember := func(items any) {
		list, _ := items.([]any)
		for _, it := range list {
			if id, ok := dig(it, "id").(float64); ok {
				s := strconv.FormatInt(int64(id), 10)
				for _, c := range created {
					if c == s {
						return
					}
				}
				created = append(created, s)
			}
		}
	}

	up := run("magnet/upload", post("/v4/magnet/upload", key,
		kv("magnets[]", bigBuckBunny), kv("magnets[]", "not-a-magnet")), "")
	remember(up.data()["magnets"])

	torrent, hash := testTorrent()
	files := run("magnet/upload/file", request{method: http.MethodPost, path: "/v4/magnet/upload/file", key: key,
		files: []upload{
			{"files[]", "open-alldebrid-contract.torrent", torrent},
			{"files[]", "not-a-torrent.zip", []byte("PK not a torrent")},
		}}, "")
	remember(files.data()["files"])
	t.Logf("test torrent info hash %s", hash)

	for _, id := range created {
		run("magnet/status "+id, post("/v4.1/magnet/status", key, kv("id", id)), "")
		run("magnet/files "+id, post("/v4/magnet/files", key, kv("id[]", id)), "")
	}

	const link = "https://example.com/open-alldebrid-contract-test"
	run("user/links/save", post("/v4/user/links/save", key, kv("links[]", link)), "")
	run("user/links/delete", post("/v4/user/links/delete", key, kv("links[]", link)), "")
}

// testTorrent builds a minimal single-file torrent nobody seeds, and returns
// it with its info hash.
func testTorrent() ([]byte, string) {
	pieces := sha1.Sum([]byte{0})
	info := fmt.Sprintf("d6:lengthi1e4:name32:open-alldebrid-contract-test.bin12:piece lengthi16384e6:pieces20:%se",
		string(pieces[:]))
	sum := sha1.Sum([]byte(info))
	torrent := "d8:announce30:udp://tracker.invalid:1337/ann4:info" + info + "e"
	return []byte(torrent), hex.EncodeToString(sum[:])
}
