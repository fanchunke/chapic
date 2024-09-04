package main

import (
	"context"

	"github.com/fanchunke/chapic/example/proto"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	proto.RegisterEchoHTTPServer(e, &EchoService{})
	e.Logger.Fatal(e.Start(":8000"))
}

type EchoService struct {
}

func (e *EchoService) UnaryEchoPost(ctx context.Context, req *proto.EchoPostRequest) (*proto.EchoResponse, error) {
	resp := &proto.EchoResponse{
		Id:      req.GetId(),
		Message: req.GetMessage(),
	}
	return resp, nil
}

func (e *EchoService) UnaryEcho(ctx context.Context, req *proto.EchoRequest) (*proto.EchoResponse, error) {
	resp := &proto.EchoResponse{
		Id:      req.GetId(),
		Message: req.GetMessage(),
	}
	return resp, nil
}
