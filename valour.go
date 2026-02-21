package valourgo

import (
	"net/http"
)

const (
	baseClientAddress    = "https://app.valour.gg"
	apiBase              = "api"
	apiPlanetBase        = apiBase + "/planets"
	apiMessageBase       = apiBase + "/messages"
	apiPlanetInitialData = "initialData"
)

type Client struct {
	*Node
	baseAddress string
	client      *http.Client
}

type Option func(c *Client)

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseAddress = baseURL
	}
}

func NewClient(token string, opts ...Option) (*Client, error) {
	c := &Client{
		baseAddress: baseClientAddress,
	}

	for _, opt := range opts {
		opt(c)
	}

	primaryNode, err := NewNode(c.baseAddress, "", token)

	if err != nil {
		return nil, err
	}

	c.Node = primaryNode

	return c, nil
}
