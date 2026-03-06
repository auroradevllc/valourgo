package valour

import (
	"net/http"
)

const (
	baseClientAddress    = "https://api.valour.gg"
	apiBase              = "api"
	apiPlanetBase        = apiBase + "/planets"
	apiMessageBase       = apiBase + "/messages"
	apiPlanetInitialData = "initialData"
)

type BaseClient struct {
	*Node
	baseAddress string
	client      *http.Client
}

type Option func(c *BaseClient)

func WithBaseURL(baseURL string) Option {
	return func(c *BaseClient) {
		c.baseAddress = baseURL
	}
}

func NewClient(token string, opts ...Option) (Client, error) {
	c := &BaseClient{
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
