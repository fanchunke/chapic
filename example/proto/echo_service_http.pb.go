// Code generated by protoc-gen-go-chapic. DO NOT EDIT.
// versions:
// - protoc-gen-go-chapic v0.1.0
// - protoc             (unknown)
// source: proto/echo_service.proto

// Echo Service
//
// Echo Service API consists of a single service which returns
// a message.

package proto

import (
	context "context"
	fmt "fmt"
	option "github.com/fanchunke/chapic/option"
	runtime "github.com/fanchunke/chapic/runtime"
	protojson "google.golang.org/protobuf/encoding/protojson"
	http "net/http"
	url "net/url"
)

import (
	echo "github.com/labstack/echo/v4"
	resty "github.com/go-resty/resty/v2"
)

var _ = new(protojson.MarshalOptions)
var _ = new(fmt.State)
var _ = new(url.Values)

// EchoHTTPClient is the client API for Echo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EchoHTTPClient interface {
	// UnaryEcho is unary echo.
	UnaryEcho(ctx context.Context, req *EchoRequest, opts ...option.CallOption) (*EchoResponse, error)
}

type echoHTTPClient struct {
	// The http endpoint to connect to.
	endpoint string

	// The http client.
	cc *resty.Client
	// The AuthFunc.
	authFunc option.AuthFunc
}

func NewEchoHTTPClient(ctx context.Context, opts ...option.ClientOption) EchoHTTPClient {
	o := option.DefaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &echoHTTPClient{cc: o.HTTPClient, endpoint: o.Endpoint, authFunc: o.AuthFunc}
}

func (c *echoHTTPClient) UnaryEcho(ctx context.Context, req *EchoRequest, opts ...option.CallOption) (*EchoResponse, error) {
	baseUrl, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}
	baseUrl.Path += fmt.Sprintf("/v1/example/%v", req.GetId())

	params := url.Values{}
	if req.GetMessage() != "" {
		params.Add("message", fmt.Sprintf("%v", req.GetMessage()))
	}

	baseUrl.RawQuery = params.Encode()

	resp, err := c.authFunc(c.cc.R()).
		Execute("GET", baseUrl.String())
	if err != nil {
		return nil, err
	}

	var result EchoResponse
	um := protojson.UnmarshalOptions{DiscardUnknown: true}
	if err := um.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// EchoHTTPServer is the server API for Echo service.
// All implementations should embed UnimplementedEcho
// for forward compatibility.
type EchoHTTPServer interface {
	// UnaryEcho is unary echo.
	UnaryEcho(ctx context.Context, req *EchoRequest) (*EchoResponse, error)
}

// UnimplementedEchoHTTPServer should be embedded to have forward compatible implementations.
type UnimplementedEchoHTTPServer struct {
}

func (UnimplementedEchoHTTPServer) UnaryEcho(ctx context.Context, req *EchoRequest) (*EchoResponse, error) {
	return nil, fmt.Errorf("method UnaryEcho not implemented")
}

func RegisterEchoHTTPServer(e *echo.Echo, srv EchoHTTPServer, m ...echo.MiddlewareFunc) {
	e.GET("/v1/example/:id", _Echo_UnaryEcho_HTTPHandler(srv), m...)
}

func _Echo_UnaryEcho_HTTPHandler(srv EchoHTTPServer) echo.HandlerFunc {
	return func(c echo.Context) error {
		in := new(EchoRequest)
		pathParams := make(map[string]string, 0)
		for _, key := range c.ParamNames() {
			pathParams[key] = c.Param(key)
		}
		if err := runtime.Bind(in, c.Request(), pathParams); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": err})
		}

		resp, err := srv.UnaryEcho(c.Request().Context(), in)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, resp)
	}
}
