package main

import (
	"context"
	"log"
	"time"

	"github.com/fanchunke/chapic/client"
	"github.com/fanchunke/chapic/example/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	ctx := context.Background()
	c := proto.NewEchoHTTPClient(
		ctx,
		client.WithEndpoint("http://127.0.0.1:8000"),
		client.WithTimeout(5*time.Second),
	)
	resp, err := c.UnaryEcho(ctx, &proto.EchoRequest{
		Id:      1,
		Message: "test",
	})
	if err != nil {
		log.Fatal(err)
	}
	data, err := protojson.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("resp: %s", string(data))
}
