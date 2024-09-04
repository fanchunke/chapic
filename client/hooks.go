package client

import (
	"context"
	"net/http"
)

type (
	RequestHook  func(context.Context, *http.Request) error
	ResponseHook func(context.Context, *http.Response) error
)
