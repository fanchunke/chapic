# chapic

This is a generator for API client and server libraries for APIs specified by protocol buffers. It takes a protocol buffer and uses it to generate a client library and an echo HTTP server.

## Installation

Use `go get` to install the latest version of the generator `protoc-gen-go-chapic`

```shell
go get github.com/fanchunke/chapic/cmd/protoc-gen-go-chapic
```

## Usage

1. define `.proto` file

    ```protobuf
    syntax = "proto3";
    
    package chapic.example.proto;
    
    import "google/api/annotations.proto";
    import "google/protobuf/field_mask.proto";
    import "google/protobuf/struct.proto";
    
    option go_package = "github.com/fanchunke/chapic/example/proto";
    
    // EchoRequest is the request for echo.
    message EchoRequest {
      int32 id = 1;
      string message = 2;
    }
    
    // EchoResponse is the response for echo.
    message EchoResponse {
      int32 id = 1;
      string message = 2;
    }
    
    // Echo is the echo service.
    service Echo {
      // UnaryEcho is unary echo.
      rpc UnaryEcho(EchoRequest) returns (EchoResponse) {
        option (google.api.http) = {
          get: "/v1/example/{id}",
        };
      }
    }
    ```

2. use `buf` generate server and client code

   - define `buf.gen.yaml` and `buf.yaml`
    
    ```text
    # buf.gen.yaml
    
    version: v1
    managed:
      enabled: true
      go_package_prefix:
        default: github.com/fanchunke/chapic/example
        except:
          - buf.build/googleapis/googleapis
    plugins:
      # Use protoc-gen-go at v1.31.0
      - plugin: buf.build/protocolbuffers/go:v1.31.0
        out: .
        opt: paths=source_relative
      - plugin: go-chapic
        out: .
        opt:
          - paths=source_relative
    ```
    
    ```text
    # buf.yaml
    
    version: v1
    breaking:
      use:
        - FILE
    deps:
      - buf.build/googleapis/googleapis
    lint:
      use:
        - DEFAULT
    ```

    - generate code

    ```shell
    buf generate
    ```
   
3. run http server

    ```go
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
    
    func (e *EchoService) UnaryEcho(ctx context.Context, req *proto.EchoRequest) (*proto.EchoResponse, error) {
        resp := &proto.EchoResponse{
            Id:      req.GetId(),
            Message: req.GetMessage(),
        }
        return resp, nil
    }
    
    ```

4. test api

   ```shell
   curl http://127.0.0.1:8000/v1/example/1\?message\=test
   
   # response
   {"id":1,"message":"test"}
   ```

5. use client to call http api

   ```go
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
       log.Printf("%s", string(data))
   }
   
   ```
   
result:

```shell
{"id":1,"message":"test"}
```