package main

import (
	"context"
	"log"
	"time"

	"github.com/fanchunke/chapic/example/proto"
	"github.com/fanchunke/chapic/option"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	ctx := context.Background()
	client := proto.NewEchoHTTPClient(
		ctx,
		option.WithEndpoint("http://127.0.0.1:8000"),
		option.WithTimeout(5*time.Second),
	)
	resp, err := client.UnaryEcho(ctx, &proto.EchoRequest{
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
