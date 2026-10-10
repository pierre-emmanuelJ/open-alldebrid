package alldebrid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// DefaultBaseURL is the production API.
const DefaultBaseURL = "https://api.alldebrid.com"

// Public demo keys. The API answers them with static sample data.
const (
	DemoPremiumKey = "staticDemoApikeyPrem"
	DemoTrialKey   = "staticDemoApikeyTria"
	DemoFreeKey    = "staticDemoApikeyFree"
)

const defaultUserAgent = "open-alldebrid-go"

// Client calls the AllDebrid API. It is safe for concurrent use.
type Client struct {
	apiKey    string
	baseURL   string
	userAgent string
	http      *http.Client
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets the HTTP client. The default has a 30 second timeout.
func WithHTTPClient(h *http.Client) Option { return func(c *Client) { c.http = h } }

// WithBaseURL sets the API base URL, for tests or a proxy.
func WithBaseURL(u string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(u, "/") }
}

// WithUserAgent identifies your application to AllDebrid.
func WithUserAgent(ua string) Option { return func(c *Client) { c.userAgent = ua } }

// New returns a client for the given API key. The key may be empty to call
// public endpoints, such as the PIN flow.
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:    apiKey,
		baseURL:   DefaultBaseURL,
		userAgent: defaultUserAgent,
		http:      &http.Client{Timeout: 30 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// WithAPIKey returns a copy of the client that uses another API key.
func (c *Client) WithAPIKey(apiKey string) *Client {
	cp := *c
	cp.apiKey = apiKey
	return &cp
}

// envelope is the shape of every API response.
type envelope[T any] struct {
	Status ResponseStatus `json:"status"`
	Data   T              `json:"data"`
	Error  *Error         `json:"error"`
}

// form is a request body. List parameters keep their PHP-style names
// (magnets[]) and are repeated, which is what the API expects.
type form = url.Values

func call[T any](ctx context.Context, c *Client, method, path string, params form) (T, error) {
	var body io.Reader
	contentType := ""
	target := c.baseURL + path
	if len(params) > 0 {
		if method == http.MethodGet {
			target += "?" + params.Encode()
		} else {
			body = strings.NewReader(params.Encode())
			contentType = "application/x-www-form-urlencoded"
		}
	}
	return send[T](ctx, c, method, target, body, contentType)
}

func send[T any](ctx context.Context, c *Client, method, target string, body io.Reader, contentType string) (T, error) {
	var zero T
	req, err := http.NewRequestWithContext(ctx, method, target, body)
	if err != nil {
		return zero, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return zero, err
	}
	defer func() { _ = resp.Body.Close() }()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 64<<20))
	if err != nil {
		return zero, err
	}

	var env envelope[T]
	decodeErr := json.Unmarshal(raw, &env)
	if resp.StatusCode != http.StatusOK {
		herr := &HTTPError{StatusCode: resp.StatusCode}
		if decodeErr == nil && env.Error != nil {
			herr.Err = env.Error
		}
		return zero, herr
	}
	if decodeErr != nil {
		return zero, fmt.Errorf("alldebrid: decode %s %s: %w", method, req.URL.Path, decodeErr)
	}
	if env.Status != ResponseStatusSuccess {
		if env.Error == nil {
			return zero, fmt.Errorf("alldebrid: %s %s: status %q without error", method, req.URL.Path, env.Status)
		}
		return zero, env.Error
	}
	return env.Data, nil
}

func ids(name string, list []int64) form {
	f := form{}
	for _, id := range list {
		f.Add(name, strconv.FormatInt(id, 10))
	}
	return f
}

func list(name string, values []string) form {
	f := form{}
	for _, v := range values {
		f.Add(name, v)
	}
	return f
}

// ---------------------------------------------------------------- general

// Ping checks that the API is up. It needs no API key.
func (c *Client) Ping(ctx context.Context) error {
	_, err := call[PingData](ctx, c, http.MethodGet, "/v4/ping", nil)
	return err
}

// ------------------------------------------------------------------- hosts

// Hosts lists the supported file hosts, streaming sites and redirectors.
// With hostsOnly, only Hosts is filled. It needs no API key; cache the result.
func (c *Client) Hosts(ctx context.Context, hostsOnly bool) (*HostsData, error) {
	var q form
	if hostsOnly {
		q = form{"hostsOnly": {"true"}}
	}
	d, err := call[HostsData](ctx, c, http.MethodGet, "/v4/hosts", q)
	return &d, err
}

// HostDomains lists every supported domain. It needs no API key.
func (c *Client) HostDomains(ctx context.Context) (*HostDomainsData, error) {
	d, err := call[HostDomainsData](ctx, c, http.MethodGet, "/v4/hosts/domains", nil)
	return &d, err
}

// HostPriority ranks host names from the most open (1) to the most restricted.
func (c *Client) HostPriority(ctx context.Context) (map[string]int, error) {
	d, err := call[HostPriorityData](ctx, c, http.MethodGet, "/v4/hosts/priority", nil)
	return d.Hosts, err
}

// -------------------------------------------------------------------- user

// User returns the account information.
func (c *Client) User(ctx context.Context) (*User, error) {
	d, err := call[UserData](ctx, c, http.MethodGet, "/v4/user", nil)
	return &d.User, err
}

// UserHosts lists the hosts this user can use now, with their quotas.
func (c *Client) UserHosts(ctx context.Context, hostsOnly bool) (ServiceMap, error) {
	var q form
	if hostsOnly {
		q = form{"hostsOnly": {"true"}}
	}
	d, err := call[UserHostsData](ctx, c, http.MethodGet, "/v4.1/user/hosts", q)
	return d.Hosts, err
}

// VerificationStatus follows a new-location login verification. The token
// comes with an AUTH_BLOCKED error (Error.Token).
func (c *Client) VerificationStatus(ctx context.Context, token string) (*VerificationData, error) {
	d, err := call[VerificationData](ctx, c, http.MethodPost, "/v4/user/verif", form{"token": {token}})
	return &d, err
}

// ResendVerification sends the verification email again (once per token).
func (c *Client) ResendVerification(ctx context.Context, token string) error {
	_, err := call[ResendVerificationData](ctx, c, http.MethodPost, "/v4/user/verif/resend", form{"token": {token}})
	return err
}

// ClearNotification clears one of the codes in User.Notifications.
func (c *Client) ClearNotification(ctx context.Context, code string) error {
	_, err := call[MessageData](ctx, c, http.MethodPost, "/v4/user/notification/clear", form{"code": {code}})
	return err
}

// ------------------------------------------------------------------- links

// LinkInfos returns the file name and size of each link, in order. A link
// that cannot be read has its Error set; the call itself still succeeds.
func (c *Client) LinkInfos(ctx context.Context, password string, links ...string) ([]LinkInfo, error) {
	f := list("link[]", links)
	if password != "" {
		f.Set("password", password)
	}
	d, err := call[LinkInfosData](ctx, c, http.MethodPost, "/v4/link/infos", f)
	return d.Infos, err
}

// ExtractRedirector extracts the links behind a redirector or link protector.
// The returned links only work with LinkInfos and Unlock.
func (c *Client) ExtractRedirector(ctx context.Context, link string) ([]string, error) {
	d, err := call[RedirectorData](ctx, c, http.MethodPost, "/v4/link/redirector", form{"link": {link}})
	return d.Links, err
}

// Unlock turns a hoster link into a direct download link. When Delayed is
// set, wait for the link with WaitDelayedLink; when Streaming lists several
// qualities, choose one with StreamingLink.
func (c *Client) Unlock(ctx context.Context, link, password string) (*UnlockData, error) {
	f := form{"link": {link}}
	if password != "" {
		f.Set("password", password)
	}
	d, err := call[UnlockData](ctx, c, http.MethodPost, "/v4/link/unlock", f)
	return &d, err
}

// StreamingLink returns the link of one quality of a streaming link. id is
// UnlockData.ID and stream is the ID of one of UnlockData.Streaming.
func (c *Client) StreamingLink(ctx context.Context, id, stream string) (*StreamingData, error) {
	d, err := call[StreamingData](ctx, c, http.MethodPost, "/v4/link/streaming", form{"id": {id}, "stream": {stream}})
	return &d, err
}

// DelayedLink returns the status of a delayed link.
func (c *Client) DelayedLink(ctx context.Context, id string) (*DelayedData, error) {
	d, err := call[DelayedData](ctx, c, http.MethodPost, "/v4/link/delayed", form{"id": {id}})
	return &d, err
}

// Delayed link statuses.
const (
	DelayedProcessing = 1
	DelayedReady      = 2
	DelayedFailed     = 3
)

// WaitDelayedLink polls a delayed link every interval (5 seconds if zero)
// until it is ready, and returns the download link.
func (c *Client) WaitDelayedLink(ctx context.Context, id string, interval time.Duration) (string, error) {
	if interval <= 0 {
		interval = 5 * time.Second
	}
	for {
		d, err := c.DelayedLink(ctx, id)
		if err != nil {
			return "", err
		}
		switch int(d.Status) {
		case DelayedReady:
			return d.Link, nil
		case DelayedFailed:
			return "", fmt.Errorf("alldebrid: delayed link %s failed", id)
		}
		if err := sleep(ctx, interval); err != nil {
			return "", err
		}
	}
}

// ----------------------------------------------------------------- magnets

// UploadMagnets adds magnet URIs or info hashes. Results are in order; an
// invalid magnet has its Error set. A magnet already in the account comes
// back with its existing ID.
func (c *Client) UploadMagnets(ctx context.Context, magnets ...string) ([]UploadedMagnet, error) {
	d, err := call[UploadMagnetsData](ctx, c, http.MethodPost, "/v4/magnet/upload", list("magnets[]", magnets))
	return d.Magnets, err
}

// TorrentFile is a .torrent file to upload.
type TorrentFile struct {
	Name    string
	Content io.Reader
}

// UploadTorrentFiles adds .torrent files. Results are in order; a file that
// is not a valid torrent has its Error set.
func (c *Client) UploadTorrentFiles(ctx context.Context, files ...TorrentFile) ([]UploadedMagnet, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for _, f := range files {
		part, err := w.CreateFormFile("files[]", f.Name)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(part, f.Content); err != nil {
			return nil, err
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	d, err := send[UploadTorrentFilesData](ctx, c, http.MethodPost, c.baseURL+"/v4/magnet/upload/file", &buf, w.FormDataContentType())
	return d.Files, err
}

// Magnets returns every magnet of the account.
func (c *Client) Magnets(ctx context.Context) ([]Magnet, error) {
	d, err := call[MagnetStatusData](ctx, c, http.MethodPost, "/v4.1/magnet/status", nil)
	return d.Magnets, err
}

// Magnet returns one magnet, with its Files.
func (c *Client) Magnet(ctx context.Context, id int64) (*Magnet, error) {
	d, err := call[MagnetStatusData](ctx, c, http.MethodPost, "/v4.1/magnet/status", form{"id": {strconv.FormatInt(id, 10)}})
	if err != nil {
		return nil, err
	}
	if len(d.Magnets) != 1 {
		return nil, fmt.Errorf("alldebrid: magnet %d: got %d magnets", id, len(d.Magnets))
	}
	return &d.Magnets[0], nil
}

// MagnetFiles returns the file tree of each magnet, in order. An unknown ID
// has its Error set.
func (c *Client) MagnetFiles(ctx context.Context, magnetIDs ...int64) ([]MagnetFiles, error) {
	d, err := call[MagnetFilesData](ctx, c, http.MethodPost, "/v4/magnet/files", ids("id[]", magnetIDs))
	return d.Magnets, err
}

// DeleteMagnet deletes a magnet.
func (c *Client) DeleteMagnet(ctx context.Context, id int64) error {
	_, err := call[MessageData](ctx, c, http.MethodPost, "/v4/magnet/delete", form{"id": {strconv.FormatInt(id, 10)}})
	return err
}

// RestartMagnet restarts a failed magnet.
func (c *Client) RestartMagnet(ctx context.Context, id int64) error {
	_, err := call[RestartMagnetsData](ctx, c, http.MethodPost, "/v4/magnet/restart", form{"id": {strconv.FormatInt(id, 10)}})
	return err
}

// RestartMagnets restarts several failed magnets. Each result has its own
// Error when that magnet could not be restarted.
func (c *Client) RestartMagnets(ctx context.Context, magnetIDs ...int64) ([]RestartedMagnet, error) {
	d, err := call[RestartMagnetsData](ctx, c, http.MethodPost, "/v4/magnet/restart", ids("ids[]", magnetIDs))
	return d.Magnets, err
}

// Magnet status codes.
const (
	MagnetQueued      = 0
	MagnetDownloading = 1
	MagnetCompressing = 2
	MagnetUploading   = 3
	MagnetReady       = 4
)

// IsReady reports whether the magnet's files are available.
func (m Magnet) IsReady() bool { return m.StatusCode == MagnetReady }

// IsError reports whether the magnet failed (status codes 5 and above).
func (m Magnet) IsError() bool { return m.StatusCode > MagnetReady }

// ------------------------------------------------------------- saved links

// SavedLinks lists the links the user saved.
func (c *Client) SavedLinks(ctx context.Context) ([]SavedLink, error) {
	d, err := call[SavedLinksData](ctx, c, http.MethodGet, "/v4/user/links", nil)
	return d.Links, err
}

// SaveLinks saves links for later.
func (c *Client) SaveLinks(ctx context.Context, links ...string) error {
	_, err := call[MessageData](ctx, c, http.MethodPost, "/v4/user/links/save", list("links[]", links))
	return err
}

// DeleteSavedLinks deletes saved links.
func (c *Client) DeleteSavedLinks(ctx context.Context, links ...string) error {
	_, err := call[MessageData](ctx, c, http.MethodPost, "/v4/user/links/delete", list("links[]", links))
	return err
}

// History lists the links unlocked in the last 3 days, when the user has
// enabled history logging (User.History).
func (c *Client) History(ctx context.Context) ([]SavedLink, error) {
	d, err := call[SavedLinksData](ctx, c, http.MethodGet, "/v4/user/history", nil)
	return d.Links, err
}

// PurgeHistory deletes the whole recent links history.
func (c *Client) PurgeHistory(ctx context.Context) error {
	_, err := call[MessageData](ctx, c, http.MethodPost, "/v4/user/history/delete", nil)
	return err
}

// ---------------------------------------------------------------- reseller

// VoucherBalance returns the reseller balance. Reseller accounts only.
func (c *Client) VoucherBalance(ctx context.Context) (float64, error) {
	d, err := call[VoucherBalanceData](ctx, c, http.MethodGet, "/v4/voucher/balance", nil)
	return d.Balance, err
}

func voucherForm(durationDays, count int) form {
	return form{"duration": {strconv.Itoa(durationDays)}, "nb": {strconv.Itoa(count)}}
}

// Vouchers takes up to count vouchers of durationDays (15, 30, 90, 180 or
// 365) from the reseller stock. Reseller accounts only.
func (c *Client) Vouchers(ctx context.Context, durationDays, count int) (*VouchersData, error) {
	d, err := call[VouchersData](ctx, c, http.MethodPost, "/v4/voucher/get", voucherForm(durationDays, count))
	return &d, err
}

// GenerateVouchers generates count new vouchers of durationDays from the
// reseller balance. Reseller accounts only.
func (c *Client) GenerateVouchers(ctx context.Context, durationDays, count int) (*GeneratedVouchersData, error) {
	d, err := call[GeneratedVouchersData](ctx, c, http.MethodPost, "/v4/voucher/generate", voucherForm(durationDays, count))
	return &d, err
}

func sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
