# JSON-RPC 2.0 Server for Golang

## Example
```go
package main

import (
	"log"

	"github.com/lapitskyss/jsonrpc"
	"github.com/lapitskyss/jsonrpc/middleware"
	"github.com/valyala/fasthttp"
)

type SumService struct {
}

func (ss *SumService) Sum(ctx *jsonrpc.RequestCtx) (jsonrpc.Result, *jsonrpc.Error) {
	var sumRequest []int
	err := ctx.Params(&sumRequest)
	if err != nil {
		return nil, err
	}

	s := 0
	for _, item := range sumRequest {
		s += item
	}

	return s, nil
}

func main() {
	sumService := SumService{}

	s := jsonrpc.NewServer(jsonrpc.Options{})

	s.Use(middleware.Recovery())
	s.Register("sum", sumService.Sum)

	m := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/rpc":
			s.HandleFastHTTP(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	log.Fatal(fasthttp.ListenAndServe(":3000", m))
}
```

### Curl request for example above
```
curl -H "Content-Type: application/json" \
  --request POST \
  --data '{"jsonrpc":"2.0","method":"sum","params":[1, 2, 3, 4],"id":1}' \
  http://localhost:3000/rpc

response for request:

{"jsonrpc":"2.0","id":1,"result":10}
```

### Batch request
```
curl -H "Content-Type: application/json" \
  --request POST \
  --data '[{"jsonrpc":"2.0","method":"sum","params":[1, 2, 3, 4],"id":1}, {"jsonrpc":"2.0","method":"sum","params":[1, 2],"id":2}]' \
  http://localhost:3000/rpc

response for request:

[{"jsonrpc":"2.0","id":2,"result":3},{"jsonrpc":"2.0","id":1,"result":10}]
```
