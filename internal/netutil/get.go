package netutil

import (
	"io"
	"net/http"
	"time"
)

type HTTPResult struct {
	Status     string
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// OK reports whether the server answered with a 2xx status.
func (r *HTTPResult) OK() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

type GetOptions struct {
	Headers map[string]string
	Timeout time.Duration
}

type Option func(*GetOptions)

func WithTimeout(d time.Duration) Option {
	return func(o *GetOptions) {
		o.Timeout = d
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(o *GetOptions) {
		o.Headers = headers
	}
}

// DefaultTimeout bounds a request when no other timeout is given, so a
// server that stops answering cannot hang the command forever.
const DefaultTimeout = 30 * time.Second

func Get(url string, opts ...Option) (*HTTPResult, error) {
	url = NormalizeURL(url)

	options := &GetOptions{Timeout: DefaultTimeout}
	for _, opt := range opts {
		opt(options)
	}

	client := &http.Client{Timeout: options.Timeout}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range options.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HTTPResult{Status: resp.Status, StatusCode: resp.StatusCode, Headers: resp.Header, Body: body}, nil
}

// Open sends a GET request and returns the response with its body unread,
// for downloads that should be streamed rather than held in memory. The
// caller closes the body. Connecting and waiting for the response headers
// are time limited, reading the body is not: how long that takes depends
// on the size of the file.
func Open(url string) (*http.Response, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ResponseHeaderTimeout = DefaultTimeout

	client := &http.Client{Transport: transport}
	return client.Get(NormalizeURL(url))
}
