package option

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	defaultTimeout = time.Second * 5
)

type Options struct {
	Endpoint   string
	Timeout    time.Duration
	HTTPClient *resty.Client
	AuthFunc   AuthFunc
}

func DefaultOptions() *Options {
	cc := defaultHTTPClient()
	return &Options{
		Endpoint:   "",
		Timeout:    defaultTimeout,
		HTTPClient: cc,
		AuthFunc:   func(request *resty.Request) *resty.Request { return request },
	}
}

type Error struct {
	Code    int
	Message string
	Header  http.Header
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error: %d, Message: %s", e.Code, e.Message)
}

func defaultHTTPClient() *resty.Client {
	cc := resty.New()
	cc.SetTimeout(defaultTimeout)
	cc.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		return nil
	})
	cc.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		if response.IsError() {
			return &Error{
				Code:    response.StatusCode(),
				Message: string(response.Body()),
				Header:  response.Header(),
			}
		}
		return nil
	})
	return cc
}

// A ClientOption is an option for an HTTP API client.
type ClientOption func(o *Options)

func WithEndpoint(url string) ClientOption {
	return func(o *Options) {
		o.Endpoint = url
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *Options) {
		o.Timeout = timeout
		o.HTTPClient.SetTimeout(timeout)
	}
}

type AuthFunc func(request *resty.Request) *resty.Request

func WithAuthFunc(f AuthFunc) ClientOption {
	return func(o *Options) {
		o.AuthFunc = f
	}
}

func WithClient(client *resty.Client) ClientOption {
	return func(o *Options) {
		o.HTTPClient = client
	}
}

type CallOptions struct{}

type CallOption func(o *CallOptions)
