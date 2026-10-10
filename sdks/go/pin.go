package alldebrid

import (
	"context"
	"net/http"
	"time"
)

// GetPin starts a PIN authentication. Show PinGetData.Pin (or UserURL) to the
// user, then wait with WaitForPin. It needs no API key.
func (c *Client) GetPin(ctx context.Context) (*PinGetData, error) {
	d, err := call[PinGetData](ctx, c, http.MethodGet, "/v4.1/pin/get", nil)
	return &d, err
}

// CheckPin returns the state of a PIN authentication. Once Activated, APIKey
// holds the new key.
func (c *Client) CheckPin(ctx context.Context, pin, check string) (*PinCheckData, error) {
	d, err := call[PinCheckData](ctx, c, http.MethodPost, "/v4/pin/check", form{"pin": {pin}, "check": {check}})
	return &d, err
}

// WaitForPin polls a PIN authentication every interval (5 seconds if zero)
// until the user confirms it, and returns the new API key. It stops with a
// CodePINExpired error when the PIN expires, or when ctx is done.
//
//	pin, _ := client.GetPin(ctx)
//	fmt.Println("Open", pin.UserURL)
//	key, err := client.WaitForPin(ctx, pin, 0)
func (c *Client) WaitForPin(ctx context.Context, pin *PinGetData, interval time.Duration) (string, error) {
	if interval <= 0 {
		interval = 5 * time.Second
	}
	for {
		d, err := c.CheckPin(ctx, pin.Pin, pin.Check)
		if err != nil {
			return "", err
		}
		if d.Activated && d.Apikey != "" {
			return d.Apikey, nil
		}
		if err := sleep(ctx, interval); err != nil {
			return "", err
		}
	}
}
