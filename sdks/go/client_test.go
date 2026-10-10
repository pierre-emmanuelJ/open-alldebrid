package alldebrid

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// fakeAPI records requests and answers with canned JSON bodies.
type fakeAPI struct {
	t       *testing.T
	mu      sync.Mutex
	reqs    []recorded
	answers []answer
}

type recorded struct {
	method, path, auth, contentType string
	query, form                     url.Values
	files                           map[string][]string // field -> file names
}

type answer struct {
	status int
	body   string
}

func newFake(t *testing.T, answers ...answer) (*fakeAPI, *Client) {
	f := &fakeAPI{t: t, answers: answers}
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)
	return f, New("test-key", WithBaseURL(srv.URL), WithUserAgent("unit-test"))
}

func ok(data string) answer { return answer{200, `{"status":"success","data":` + data + `}`} }

func apiErr(code string) answer {
	return answer{200, `{"status":"error","error":{"code":"` + code + `","message":"m"}}`}
}

func (f *fakeAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	rec := recorded{
		method: r.Method, path: r.URL.Path, query: r.URL.Query(),
		auth: r.Header.Get("Authorization"), contentType: r.Header.Get("Content-Type"),
	}
	mt, params, _ := mime.ParseMediaType(rec.contentType)
	switch mt {
	case "application/x-www-form-urlencoded":
		b, _ := io.ReadAll(r.Body)
		rec.form, _ = url.ParseQuery(string(b))
	case "multipart/form-data":
		mr := multipart.NewReader(r.Body, params["boundary"])
		rec.files = map[string][]string{}
		for {
			p, err := mr.NextPart()
			if err != nil {
				break
			}
			rec.files[p.FormName()] = append(rec.files[p.FormName()], p.FileName())
		}
	}
	if r.Header.Get("User-Agent") != "unit-test" {
		f.t.Errorf("User-Agent = %q", r.Header.Get("User-Agent"))
	}
	f.reqs = append(f.reqs, rec)
	a := answer{200, `{"status":"success","data":{}}`}
	if len(f.answers) > 0 {
		a, f.answers = f.answers[0], f.answers[1:]
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(a.status)
	_, _ = io.WriteString(w, a.body)
}

func (f *fakeAPI) last() recorded {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.reqs[len(f.reqs)-1]
}

var ctx = context.Background()

func TestFormListsArePHPStyle(t *testing.T) {
	f, c := newFake(t, ok(`{"magnets":[{"magnet":"a","id":1},{"magnet":"b","error":{"code":"MAGNET_INVALID_URI","message":"m"}}]}`))
	got, err := c.UploadMagnets(ctx, "a", "b")
	if err != nil {
		t.Fatal(err)
	}
	req := f.last()
	if req.method != http.MethodPost || req.path != "/v4/magnet/upload" {
		t.Errorf("request = %s %s", req.method, req.path)
	}
	if req.auth != "Bearer test-key" {
		t.Errorf("Authorization = %q", req.auth)
	}
	if !reflect.DeepEqual(req.form["magnets[]"], []string{"a", "b"}) {
		t.Errorf("form = %v, want magnets[]=a&magnets[]=b", req.form)
	}
	if len(got) != 2 || got[0].ID != 1 || got[1].Error.Code != "MAGNET_INVALID_URI" {
		t.Errorf("result = %+v", got)
	}
}

func TestEveryListParameter(t *testing.T) {
	f, c := newFake(t)
	cases := []struct {
		call  func() error
		path  string
		field string
		want  []string
	}{
		{func() error { _, err := c.LinkInfos(ctx, "", "l1", "l2"); return err }, "/v4/link/infos", "link[]", []string{"l1", "l2"}},
		{func() error { _, err := c.MagnetFiles(ctx, 1, 2); return err }, "/v4/magnet/files", "id[]", []string{"1", "2"}},
		{func() error { _, err := c.RestartMagnets(ctx, 3, 4); return err }, "/v4/magnet/restart", "ids[]", []string{"3", "4"}},
		{func() error { return c.SaveLinks(ctx, "s1", "s2") }, "/v4/user/links/save", "links[]", []string{"s1", "s2"}},
		{func() error { return c.DeleteSavedLinks(ctx, "d1") }, "/v4/user/links/delete", "links[]", []string{"d1"}},
	}
	for _, tc := range cases {
		if err := tc.call(); err != nil {
			t.Fatalf("%s: %v", tc.path, err)
		}
		req := f.last()
		if req.path != tc.path || !reflect.DeepEqual(req.form[tc.field], tc.want) {
			t.Errorf("%s: got %s %v", tc.path, req.path, req.form)
		}
	}
}

func TestOptionalPasswordAndHostsOnly(t *testing.T) {
	f, c := newFake(t)
	if _, err := c.Unlock(ctx, "https://example.com/f", ""); err != nil {
		t.Fatal(err)
	}
	if _, ok := f.last().form["password"]; ok {
		t.Error("empty password was sent")
	}
	if _, err := c.Unlock(ctx, "https://example.com/f", "pw"); err != nil {
		t.Fatal(err)
	}
	if got := f.last().form.Get("password"); got != "pw" {
		t.Errorf("password = %q", got)
	}
	if _, err := c.Hosts(ctx, true); err != nil {
		t.Fatal(err)
	}
	if req := f.last(); req.method != http.MethodGet || req.query.Get("hostsOnly") != "true" {
		t.Errorf("hosts request = %s %v", req.method, req.query)
	}
}

func TestNoKeyNoAuthorization(t *testing.T) {
	f, c := newFake(t, ok(`{"ping":"pong"}`))
	if err := c.WithAPIKey("").Ping(ctx); err != nil {
		t.Fatal(err)
	}
	if f.last().auth != "" {
		t.Errorf("Authorization sent without a key: %q", f.last().auth)
	}
}

func TestAPIErrorsAreComparable(t *testing.T) {
	_, c := newFake(t, apiErr("MAGNET_INVALID_ID"))
	err := c.DeleteMagnet(ctx, 1)
	if !errors.Is(err, CodeMagnetInvalidID) {
		t.Fatalf("errors.Is(%v, CodeMagnetInvalidID) = false", err)
	}
	if errors.Is(err, CodeLinkDown) {
		t.Error("matched the wrong code")
	}
	var apiError *Error
	if !errors.As(err, &apiError) || apiError.Message != "m" {
		t.Errorf("errors.As = %v", err)
	}
}

func TestBlockedLoginCarriesToken(t *testing.T) {
	_, c := newFake(t, answer{200, `{"status":"error","error":{"code":"AUTH_BLOCKED","message":"m","token":"tok"}}`})
	_, err := c.User(ctx)
	var apiError *Error
	if !errors.As(err, &apiError) || apiError.Token != "tok" || !errors.Is(err, CodeAuthBlocked) {
		t.Fatalf("err = %#v", err)
	}
}

func TestHTTPErrors(t *testing.T) {
	_, c := newFake(t,
		answer{429, `{"status":"error","error":{"code":"TOO_MANY","message":"slow down"}}`},
		answer{404, `{"status":"error","error":{"code":"404","message":"Endpoint doesn't exist"}}`},
		answer{502, `<html>bad gateway</html>`},
	)
	err := c.Ping(ctx)
	if !IsRateLimited(err) {
		t.Errorf("429 not reported as rate limited: %v", err)
	}
	err = c.Ping(ctx)
	if !errors.Is(err, CodeNotFound) {
		t.Errorf("404 body code lost: %v", err)
	}
	err = c.Ping(ctx)
	var h *HTTPError
	if !errors.As(err, &h) || h.StatusCode != 502 || h.Err != nil {
		t.Errorf("502 = %#v", err)
	}
}

func TestSuccessWithoutJSON(t *testing.T) {
	_, c := newFake(t, answer{200, "not json"})
	if err := c.Ping(ctx); err == nil {
		t.Fatal("no error on a non-JSON body")
	}
}

func TestMagnetSingleObject(t *testing.T) {
	_, c := newFake(t, ok(`{"magnets":{"id":7,"filename":"f","statusCode":4,"files":[{"n":"dir","e":[{"n":"a","s":1,"l":"L"}]}]}}`))
	m, err := c.Magnet(ctx, 7)
	if err != nil {
		t.Fatal(err)
	}
	if m.ID != 7 || !m.IsReady() || m.Files[0].E[0].L != "L" {
		t.Errorf("magnet = %+v", m)
	}
}

func TestFlexibleTypes(t *testing.T) {
	var u UnlockData
	if err := json.Unmarshal([]byte(`{"delayed":"123","streaming":[{"id":"mp3","quality":"mp3"},{"id":"hd","quality":1080}]}`), &u); err != nil {
		t.Fatal(err)
	}
	if u.Delayed != "123" || u.Streaming[0].Quality != "mp3" || u.Streaming[1].Quality != "1080" {
		t.Errorf("unlock = %+v", u)
	}
	var s StreamingData
	if err := json.Unmarshal([]byte(`{"delayed":2277564}`), &s); err != nil {
		t.Fatal(err)
	}
	if n, err := s.Delayed.Int64(); err != nil || n != 2277564 {
		t.Errorf("delayed = %q", s.Delayed)
	}
	var h Service
	if err := json.Unmarshal([]byte(`{"name":"a","type":"free","domains":[],"regexp":["x","y"]}`), &h); err != nil {
		t.Fatal(err)
	}
	if len(h.Regexp) != 2 {
		t.Errorf("regexp = %v", h.Regexp)
	}
	if err := json.Unmarshal([]byte(`{"name":"a","type":"free","domains":[],"regexp":"x"}`), &h); err != nil || len(h.Regexp) != 1 {
		t.Errorf("regexp = %v, %v", h.Regexp, err)
	}
	var user User
	if err := json.Unmarshal([]byte(`{"limitedHostersQuotas":[]}`), &user); err != nil || user.LimitedHostersQuotas == nil {
		t.Errorf("empty quotas = %v, %v", user.LimitedHostersQuotas, err)
	}
	if err := json.Unmarshal([]byte(`{"limitedHostersQuotas":{"ddl":2000}}`), &user); err != nil || user.LimitedHostersQuotas["ddl"] != 2000 {
		t.Errorf("quotas = %v, %v", user.LimitedHostersQuotas, err)
	}
	if err := json.Unmarshal([]byte(`{"limitedHostersQuotas":[1]}`), &user); err == nil {
		t.Error("a non-empty quota list was accepted")
	}
}

func TestUploadTorrentFiles(t *testing.T) {
	f, c := newFake(t, ok(`{"files":[{"magnet":"m","id":9},{"file":"b.zip","error":{"code":"MAGNET_INVALID_FILE","message":"m"}}]}`))
	got, err := c.UploadTorrentFiles(ctx,
		TorrentFile{Name: "a.torrent", Content: strings.NewReader("d4:infod4:name1:aee")},
		TorrentFile{Name: "b.zip", Content: strings.NewReader("zip")},
	)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"a.torrent", "b.zip"}; !reflect.DeepEqual(f.last().files["files[]"], want) {
		t.Errorf("files = %v, want %v", f.last().files, want)
	}
	if got[0].ID != 9 || got[1].File != "b.zip" {
		t.Errorf("result = %+v", got)
	}
}

func TestWaitForPin(t *testing.T) {
	f, c := newFake(t,
		ok(`{"activated":false,"expires_in":600}`),
		ok(`{"activated":true,"expires_in":590,"apikey":"new-key"}`),
	)
	key, err := c.WaitForPin(ctx, &PinGetData{Pin: "ABCD", Check: "chk"}, time.Millisecond)
	if err != nil || key != "new-key" {
		t.Fatalf("key = %q, %v", key, err)
	}
	if req := f.last(); req.form.Get("pin") != "ABCD" || req.form.Get("check") != "chk" {
		t.Errorf("form = %v", req.form)
	}

	_, c = newFake(t, apiErr("PIN_EXPIRED"))
	if _, err := c.WaitForPin(ctx, &PinGetData{Pin: "A", Check: "c"}, time.Millisecond); !errors.Is(err, CodePINExpired) {
		t.Errorf("expired pin: %v", err)
	}
}

func TestWaitDelayedLink(t *testing.T) {
	_, c := newFake(t, ok(`{"status":1,"time_left":5}`), ok(`{"status":2,"time_left":0,"link":"L"}`))
	link, err := c.WaitDelayedLink(ctx, "1", time.Millisecond)
	if err != nil || link != "L" {
		t.Fatalf("link = %q, %v", link, err)
	}
	_, c = newFake(t, ok(`{"status":3,"time_left":0}`))
	if _, err := c.WaitDelayedLink(ctx, "1", time.Millisecond); err == nil {
		t.Error("failed delayed link did not error")
	}
}

func TestLiveMagnets(t *testing.T) {
	f, c := newFake(t,
		ok(`{"fullsync":true,"counter":1,"magnets":[{"id":1,"filename":"a","uploadDate":10,"statusCode":1,"downloaded":5},{"id":2,"filename":"b","uploadDate":20,"statusCode":4}]}`),
		ok(`{"counter":2,"magnets":{"id":1,"downloaded":9}}`),
		ok(`{"counter":3,"magnets":[{"id":2,"deleted":true},{"id":3,"filename":"c","uploadDate":30,"statusCode":0}]}`),
		ok(`{"counter":4,"magnets":[]}`),
		ok(`{"fullsync":true,"counter":1,"magnets":[{"id":9,"filename":"z","uploadDate":1}]}`),
	)
	live := c.LiveMagnets()
	poll := func() []Magnet {
		t.Helper()
		ms, err := live.Poll(ctx)
		if err != nil {
			t.Fatal(err)
		}
		return ms
	}

	ms := poll()
	if len(ms) != 2 || ms[0].ID != 2 || f.last().form.Get("counter") != "0" {
		t.Fatalf("full sync = %+v, form %v", ms, f.last().form)
	}
	session := f.last().form.Get("session")

	ms = poll()
	if f.last().form.Get("counter") != "1" || f.last().form.Get("session") != session {
		t.Errorf("second poll form = %v", f.last().form)
	}
	if ms[1].ID != 1 || ms[1].Downloaded != 9 || ms[1].Filename != "a" {
		t.Errorf("delta not merged: %+v", ms[1])
	}

	ms = poll()
	if len(ms) != 2 || ms[0].ID != 3 || ms[1].ID != 1 {
		t.Errorf("add and delete not applied: %+v", ms)
	}

	if ms = poll(); len(ms) != 2 {
		t.Errorf("empty delta changed the state: %+v", ms)
	}

	if ms = poll(); len(ms) != 1 || ms[0].ID != 9 {
		t.Errorf("new full sync did not replace the state: %+v", ms)
	}
}

func TestContextCancel(t *testing.T) {
	_, c := newFake(t, ok(`{"status":1,"time_left":5}`))
	cctx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
	defer cancel()
	if _, err := c.WaitDelayedLink(cctx, "1", time.Hour); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v", err)
	}
}

func TestErrorStrings(t *testing.T) {
	e := &Error{Code: "LINK_DOWN", Message: "gone"}
	if e.Error() != "alldebrid: LINK_DOWN: gone" || CodeLinkDown.Error() != "LINK_DOWN" {
		t.Errorf("Error() = %q", e.Error())
	}
	h := &HTTPError{StatusCode: 429, Err: e}
	if h.Error() != "alldebrid: HTTP 429: LINK_DOWN: gone" || !errors.Is(h, CodeLinkDown) {
		t.Errorf("HTTPError = %q", h.Error())
	}
	if (&HTTPError{StatusCode: 503}).Error() != "alldebrid: HTTP 503 Service Unavailable" {
		t.Error("HTTPError without body")
	}
	if FlexString("12").String() != "12" {
		t.Error("FlexString.String")
	}
}

func TestMagnetsListAndHTTPClient(t *testing.T) {
	var d MagnetStatusData
	if err := json.Unmarshal([]byte(`{"magnets":[{"id":1},{"id":2}]}`), &d); err != nil || len(d.Magnets) != 2 {
		t.Fatalf("list = %v, %v", d.Magnets, err)
	}
	if err := json.Unmarshal([]byte(`{"magnets":null}`), &d); err != nil || d.Magnets != nil {
		t.Fatalf("null = %v, %v", d.Magnets, err)
	}
	h := &http.Client{Timeout: time.Second}
	if New("k", WithHTTPClient(h)).http != h {
		t.Error("WithHTTPClient ignored")
	}
}
