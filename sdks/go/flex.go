package alldebrid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexString is a value the API sends either as a JSON string or as a JSON
// number, such as a delayed link ID. It holds the value as a string.
type FlexString string

// UnmarshalJSON accepts a string, a number or null.
func (f *FlexString) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if bytes.Equal(b, []byte("null")) {
		*f = ""
		return nil
	}
	if len(b) > 0 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*f = FlexString(s)
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	*f = FlexString(n.String())
	return nil
}

// String returns the value as sent by the API.
func (f FlexString) String() string { return string(f) }

// Int64 returns the value as an integer, when it is one.
func (f FlexString) Int64() (int64, error) { return strconv.ParseInt(string(f), 10, 64) }

// Strings is a value the API sends either as one string or as a list of
// strings, such as the regexp of a host.
type Strings []string

// UnmarshalJSON accepts a string, a list of strings or null.
func (s *Strings) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	switch {
	case bytes.Equal(b, []byte("null")):
		*s = nil
		return nil
	case len(b) > 0 && b[0] == '"':
		var one string
		if err := json.Unmarshal(b, &one); err != nil {
			return err
		}
		*s = Strings{one}
		return nil
	}
	var many []string
	if err := json.Unmarshal(b, &many); err != nil {
		return err
	}
	*s = many
	return nil
}

// Magnets is the magnets field of a magnet status. The API sends a list, or
// a single object when one magnet was requested; both decode to a list.
type Magnets []Magnet

// UnmarshalJSON accepts a list of magnets, one magnet or null.
func (m *Magnets) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	switch {
	case bytes.Equal(b, []byte("null")):
		*m = nil
		return nil
	case len(b) > 0 && b[0] == '{':
		var one Magnet
		if err := json.Unmarshal(b, &one); err != nil {
			return err
		}
		*m = Magnets{one}
		return nil
	}
	var many []Magnet
	if err := json.Unmarshal(b, &many); err != nil {
		return err
	}
	*m = many
	return nil
}

// Quotas maps limited host names to their remaining quota in MB. The API
// sends an empty map as an empty JSON list; it decodes to an empty map.
type Quotas map[string]int64

// UnmarshalJSON accepts an object, an empty list or null.
func (q *Quotas) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if bytes.Equal(b, []byte("null")) || (len(b) > 0 && b[0] == '[') {
		var empty []json.RawMessage
		if b[0] == '[' {
			if err := json.Unmarshal(b, &empty); err != nil {
				return err
			}
			if len(empty) > 0 {
				return fmt.Errorf("alldebrid: quotas sent as a non-empty list")
			}
		}
		*q = Quotas{}
		return nil
	}
	m := map[string]int64{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	*q = m
	return nil
}
