package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	BaseURL    string
	httpClient Doer
	logger     *zerolog.Logger
	timeout    time.Duration

	requestHooks  []RequestHook
	responseHooks []ResponseHook
}

func NewClient(opts ...Option) *Client {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	c := &Client{
		httpClient: http.DefaultClient,
		logger:     &logger,
		timeout:    defaultTimeout,
		requestHooks: []RequestHook{
			func(ctx context.Context, request *http.Request) error {
				request.Header.Set("Content-Type", "application/json")
				return nil
			},
		},
		responseHooks: nil,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) WithEndpoint(endpoint string) *Client {
	c.BaseURL = strings.TrimRight(c.BaseURL, "/")
	return c
}

func (c *Client) WithLogger(logger *zerolog.Logger) *Client {
	c.logger = logger
	return c
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.timeout = timeout
	return c
}

func (c *Client) WithHTTPClient(httpClient Doer) *Client {
	c.httpClient = httpClient
	return c
}

func (c *Client) WithRequestHooks(hooks ...RequestHook) *Client {
	c.requestHooks = append(c.requestHooks, hooks...)
	return c
}

func (c *Client) Copy(opts ...Option) *Client {
	if len(opts) == 0 {
		return c
	}

	newClient := &Client{}
	newClient.WithEndpoint(c.BaseURL).
		WithHTTPClient(c.httpClient).
		WithRequestHooks(c.requestHooks...).
		WithLogger(c.logger).
		WithTimeout(c.timeout).
		WithHTTPClient(c.httpClient)
	for _, opt := range opts {
		opt(newClient)
	}
	return newClient
}

func (c *Client) WithResponseHooks(hooks ...ResponseHook) *Client {
	c.responseHooks = append(c.responseHooks, hooks...)
	return c
}

func Request[T, S proto.Message](ctx context.Context, client *Client, method string, url string, req T, result S, opts ...Option) error {
	newClient := client.Copy(opts...)
	data, err := protojson.Marshal(req)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, newClient.timeout)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// run request hooks
	for _, h := range newClient.requestHooks {
		if err := h(ctx, r); err != nil {
			newClient.logger.Error().Ctx(ctx).Err(err).Msg("request hook error")
			return err
		}
	}

	newClient.logger.Debug().Ctx(ctx).RawJSON("request", data).Str("url", r.URL.String()).Interface("header", r.Header).Msg("start http request")
	resp, err := newClient.httpClient.Do(r)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		newClient.logger.Debug().Ctx(ctx).RawJSON("response", body).Int("status", resp.StatusCode).Interface("header", resp.Header).Msg("finish http request")
	} else {
		newClient.logger.Debug().Ctx(ctx).Bytes("response", body).Int("status", resp.StatusCode).Interface("header", resp.Header).Msg("finish http request")
	}

	// run response hook
	for _, h := range newClient.responseHooks {
		if err := h(ctx, resp); err != nil {
			newClient.logger.Error().Ctx(ctx).Err(err).Msg("response hook error")
			return err
		}
	}

	if err := protojson.Unmarshal(body, result); err != nil {
		return err
	}
	return nil
}
