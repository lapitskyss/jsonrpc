# JSON-RPC 2.0 Server for Golang

## Example

```go
package main

import (
	"log"
	"net/http"

	"github.com/lapitskyss/jsonrpc"
	"github.com/lapitskyss/jsonrpc/middleware"
)

type SumService struct {
}

func (ss *SumService) Sum(ctx *jsonrpc.RequestCtx) (jsonrpc.Result, jsonrpc.Error) {
	var sumRequest []int
	err := ctx.GetParams(&sumRequest)
	if err != nil {
		return nil, jsonrpc.ErrInvalidParamsJSON()
	}

	s := 0
	for _, item := range sumRequest {
		s += item
	}

	return ctx.Result(s)
}

func main() {
	sumService := SumService{}

	s := jsonrpc.NewServer(jsonrpc.Options{})
	s.Use(middleware.Recovery())

	s.Register("sum", sumService.Sum)

	http.Handle("/rpc", s)

	log.Fatal(http.ListenAndServe(":3000", nil))
}


```

### Curl example

Request

```bash
curl -H "Content-Type: application/json" \
  --request POST \
  --data '{"jsonrpc":"2.0","method":"sum","params":[1, 2, 3, 4],"id":1}' \
  http://localhost:3000/rpc
```

Response

```bash
{"jsonrpc":"2.0","id":1,"result":10}
```

### Curl batch example

Request

```bash
curl -H "Content-Type: application/json" \
  --request POST \
  --data '[{"jsonrpc":"2.0","method":"sum","params":[1, 2, 3, 4],"id":1}, {"jsonrpc":"2.0","method":"sum","params":[1, 2],"id":2}]' \
  http://localhost:3000/rpc
```

Response

```bash
[{"jsonrpc":"2.0","id":2,"result":3},{"jsonrpc":"2.0","id":1,"result":10}]
```
