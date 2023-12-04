package main

import (
	"context"
	"net/http"

	rochev1 "github.com/fanchunke/chapic/api/roche/v1"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	e := echo.New()

	g := e.Group("/")
	g.GET("/api/v1/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"message": "test"})
	})

	h := e.Group("/")
	h.GET("/api/v2/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"message": "test"})
	})
	rochev1.RegisterHelloServiceHTTPServer(e, HelloService{})
	e.Logger.Fatal(e.Start(":8000"))
}

type HelloService struct {
}

func (h HelloService) Greet(ctx context.Context, req *rochev1.GreetRequest) (*rochev1.GreetResponse, error) {
	tmp, _ := protojson.Marshal(req)
	return &rochev1.GreetResponse{Greet: string(tmp)}, nil
}
