package client

import (
	"time"

	"github.com/rs/zerolog"
)

const (
	defaultTimeout = time.Second * 5
)

// Option is an option for an HTTP API client.
type Option func(o *Client)

func WithEndpoint(url string) Option {
	return func(o *Client) {
		o.BaseURL = url
	}
}

func WithHTTPClient(client Doer) Option {
	return func(o *Client) {
		o.httpClient = client
	}
}

func WithLogger(logger *zerolog.Logger) Option {
	return func(o *Client) {
		o.logger = logger
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Client) {
		o.timeout = timeout
	}
}

func WithRequestHooks(hooks ...RequestHook) Option {
	return func(o *Client) {
		o.requestHooks = append(o.requestHooks, hooks...)
	}
}

func WithResponseHooks(hooks ...ResponseHook) Option {
	return func(o *Client) {
		o.responseHooks = append(o.responseHooks, hooks...)
	}
}
