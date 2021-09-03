# JSON-RPC 2.0 Server for Golang

## Example
```go
package main

import (
	"log"
	"net/http"

	"github.com/lapitskyss/jsonrpc"
	"github.com/lapitskyss/jsonrpc/middleware"
	"github.com/lapitskyss/jsonrpc/middleware_global"
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

	s.UseGlobal(middleware_global.RealIP())
	s.Use(middleware.Recovery())

	s.Register("sum", sumService.Sum)

	http.Handle("/rpc", s)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
```

### Curl request for example above
```
curl -H "Content-Type: application/json" \
  --request POST \
  --data '{"jsonrpc":"2.0","method":"sum","params":[1, 2, 3, 4],"id":1}' \
  http://localhost:3000/rpc
```
- response for request:
```
{"jsonrpc":"2.0","id":1,"result":10}
```

### RequestCtx
```go
func (ss *SumService) Sum(ctx *jsonrpc.RequestCtx) (jsonrpc.Result, *jsonrpc.Error) {
	request := ctx.R // *http.Request   
	ID := ctx.ID // json rpc id
	Version := ctx.Version // json rpc version

	// set new key/value for request context
	ctx.Set("user", "1")

	// get value
	user, ok := ctx.Get("user")
}
```