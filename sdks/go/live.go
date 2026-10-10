package alldebrid

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

// LiveMagnets follows the magnets of the account with the API's live mode:
// after a first full list, each Poll only transfers the magnets that changed.
// It keeps the merged state. It is safe for concurrent use.
type LiveMagnets struct {
	c       *Client
	session int64
	mu      sync.Mutex
	counter int64
	state   map[int64]json.RawMessage
}

// LiveMagnets starts a live mode session with a random session ID.
func (c *Client) LiveMagnets() *LiveMagnets {
	return &LiveMagnets{c: c, session: rand.Int64N(1 << 31), state: map[int64]json.RawMessage{}}
}

type liveData struct {
	Magnets  json.RawMessage `json:"magnets"`
	Counter  int64           `json:"counter"`
	Fullsync bool            `json:"fullsync"`
}

// Poll fetches the changes since the last call, applies them, and returns
// every magnet, newest first.
func (l *LiveMagnets) Poll(ctx context.Context) ([]Magnet, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	f := form{
		"session": {strconv.FormatInt(l.session, 10)},
		"counter": {strconv.FormatInt(l.counter, 10)},
	}
	d, err := call[liveData](ctx, l.c, http.MethodPost, "/v4.1/magnet/status", f)
	if err != nil {
		return nil, err
	}
	changes, err := splitMagnets(d.Magnets)
	if err != nil {
		return nil, err
	}
	if d.Fullsync {
		l.state = map[int64]json.RawMessage{}
	}
	for _, raw := range changes {
		if err := l.apply(raw); err != nil {
			return nil, err
		}
	}
	l.counter = d.Counter
	return l.snapshot()
}

// splitMagnets returns the magnets of a status response, which the API sends
// as a list or as a single object.
func splitMagnets(raw json.RawMessage) ([]json.RawMessage, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	if raw[0] == '{' {
		return []json.RawMessage{raw}, nil
	}
	var list []json.RawMessage
	err := json.Unmarshal(raw, &list)
	return list, err
}

// apply merges one changed magnet into the state, field by field.
func (l *LiveMagnets) apply(raw json.RawMessage) error {
	var delta map[string]json.RawMessage
	if err := json.Unmarshal(raw, &delta); err != nil {
		return err
	}
	var id int64
	if err := json.Unmarshal(delta["id"], &id); err != nil {
		return err
	}
	var deleted bool
	if v, ok := delta["deleted"]; ok {
		_ = json.Unmarshal(v, &deleted)
	}
	if deleted {
		delete(l.state, id)
		return nil
	}
	merged := map[string]json.RawMessage{}
	if prev, ok := l.state[id]; ok {
		if err := json.Unmarshal(prev, &merged); err != nil {
			return err
		}
	}
	for k, v := range delta {
		merged[k] = v
	}
	b, err := json.Marshal(merged)
	if err != nil {
		return err
	}
	l.state[id] = b
	return nil
}

func (l *LiveMagnets) snapshot() ([]Magnet, error) {
	out := make([]Magnet, 0, len(l.state))
	for _, raw := range l.state {
		var m Magnet
		if err := json.Unmarshal(raw, &m); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].UploadDate != out[j].UploadDate {
			return out[i].UploadDate > out[j].UploadDate
		}
		return out[i].ID > out[j].ID
	})
	return out, nil
}
