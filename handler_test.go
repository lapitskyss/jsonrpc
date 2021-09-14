package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func TestServeHTTP(t *testing.T) {
	rpc := NewServer(Options{})

	sumService := SumService{}
	rpc.Register("sum", sumService.sum)

	var tc = []struct {
		name, in, out string
	}{
		{
			name: "OK",
			in:   `{"jsonrpc":"2.0","method":"sum","params":[1, 2, 3, 4],"id":1}`,
			out:  `{"jsonrpc":"2.0","id":1,"result":10}`,
		},
		{
			name: "OKBatch",
			in:   `[{"jsonrpc":"2.0","method":"sum","params":[1, 2, 3, 4],"id":1}, {"jsonrpc":"2.0","method":"sum","params":[1, 2],"id":2}]`,
			out:  `[{"jsonrpc":"2.0","id":2,"result":3}, {"jsonrpc":"2.0","id":1,"result":10}]`,
		},
		{
			name: "Notification",
			in:   `{"jsonrpc": "2.0", "method": "sum", "params": [1, 2, 3, 4] }`,
			out:  ``,
		},
		{
			name: "MethodNotFound",
			in:   `{"jsonrpc": "2.0", "method": "div", "params": [1, 2, 3, 4], "id": 1}`,
			out:  `{"jsonrpc":"2.0","id":1,"error":{"code":-32601,"message":"Method not found"}}`,
		},
		{
			name: "InvalidRequest",
			in:   ``,
			out:  `{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null}`,
		},
		{
			name: "ParseError",
			in:   `{"jsonrpc": "2.0", "method": "foobar, "params": "bar", "baz]`,
			out:  `{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`,
		},
		{
			name: "InvalidParams",
			in:   `{"jsonrpc": "2.0", "method": "sum", "params": "error", "id": 1 }`,
			out:  `{"jsonrpc": "2.0", "error": {"code": -32602, "message": "Invalid params"}, "id": 1}`,
		},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			r, err := http.NewRequest("POST", "http://test/", bytes.NewBufferString(c.in))
			if err != nil {
				require.NoError(t, err)
			}
			r.Header.Set("Content-Type", "application/json")

			res, err := serve(rpc.HandleFastHTTP, r)
			if err != nil {
				require.NoError(t, err)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				require.NoError(t, err)
			}

			result := string(body)

			if result == "" {
				require.Equal(t, c.out, result)
			} else {
				require.JSONEq(t, c.out, result)
			}
		})
	}
}

func BenchmarkServeHTTP(b *testing.B) {
	rpc := NewServer(Options{})

	sumService := SumService{}
	rpc.Register("sum", sumService.sum)

	buf := &bytes.Buffer{}
	buf.WriteString(`[{"jsonrpc":"2.0","method":"sum","params":[1, 2, 3, 4],"id":1}, {"jsonrpc":"2.0","method":"sum","params":[1, 2],"id":2}]`)
	body := buf.Bytes()

	ctx := &fasthttp.RequestCtx{
		Request:  fasthttp.Request{},
		Response: fasthttp.Response{},
	}
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.Header.SetContentType("application/json")
	ctx.Request.SetBody(body)

	b.ReportAllocs()
	b.ResetTimer()

	fasthttp.AcquireRequest()

	for i := 0; i < b.N; i++ {
		rpc.HandleFastHTTP(ctx)
	}
}

func BenchmarkGetRequestId(b *testing.B) {
	message := json.RawMessage(`"123"`)
	id := &message

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		getRequestId(id)
	}
}

type SumService struct {
}

func (ss *SumService) sum(ctx *RequestCtx) (Result, *Error) {
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

func serve(handler fasthttp.RequestHandler, req *http.Request) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	return client.Do(req)
}
