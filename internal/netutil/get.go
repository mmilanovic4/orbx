package netutil

import (
	"io"
	"net/http"
	"time"
)

type HTTPResult struct {
	Status  string
	Headers http.Header
	Body    []byte
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

func Get(url string, opts ...Option) (*HTTPResult, error) {
	url = NormalizeURL(url)
	client := &http.Client{}

	options := &GetOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if options.Timeout > 0 {
		client.Timeout = options.Timeout
	}

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

	return &HTTPResult{Status: resp.Status, Headers: resp.Header, Body: body}, nil
}
