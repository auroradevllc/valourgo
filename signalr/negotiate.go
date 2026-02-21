package signalr

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type NegotiateResponse struct {
	BaseURL             *url.URL `json:"-"`
	ConnectionToken     string   `json:"connectionToken"`
	AvailableTransports []struct {
		Transport string `json:"transport"`
	} `json:"availableTransports"`
}

func (c *Client) negotiate(ctx context.Context, baseURL string) (*NegotiateResponse, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/negotiate?negotiateVersion=1",
		nil,
	)

	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var n NegotiateResponse

	n.BaseURL, err = url.Parse(baseURL)

	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&n); err != nil {
		return nil, err
	}

	return &n, nil
}

func (n *NegotiateResponse) WebSocketURL() string {
	scheme := "ws"

	// If https, use wss instead
	if n.BaseURL.Scheme == "https" {
		scheme = "wss"
	}

	u := &url.URL{
		Scheme: scheme,
		Host:   n.BaseURL.Host,
		Path:   n.BaseURL.Path,
	}

	q := make(url.Values)
	q.Set("id", n.ConnectionToken)
	u.RawQuery = q.Encode()

	return u.String()
}
