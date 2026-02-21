package valourgo

import "net/http"

type headerRoundTripper struct {
	headers http.Header
	rt      http.RoundTripper
}

func (h *headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid mutating the original
	req2 := req.Clone(req.Context())

	for k, values := range h.headers {
		for _, v := range values {
			req2.Header.Add(k, v)
		}
	}

	return h.rt.RoundTrip(req2)
}

func Ref[V any](v V) *V {
	return &v
}
