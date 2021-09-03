package jsonrpc_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lapitskyss/jsonrpc"

	"github.com/stretchr/testify/require"
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

func TestServeHTTP(t *testing.T) {
	rpc := jsonrpc.NewServer(jsonrpc.Options{})

	sumService := SumService{}
	rpc.Register("sum", sumService.Sum)

	ts := httptest.NewServer(http.HandlerFunc(rpc.ServeHTTP))
	defer ts.Close()

	var tc = []struct {
		name, in, out string
	}{
		{
			name: "OK",
			in:   `{"jsonrpc":"2.0","method":"sum","params":[1, 2, 3, 4],"id":1}`,
			out:  `{"jsonrpc":"2.0","id":1,"result":10}`,
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
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			res, err := http.Post(ts.URL, "application/json", bytes.NewBufferString(c.in))
			if err != nil {
				require.NoError(t, err)
			}

			resp, err := ioutil.ReadAll(res.Body)
			if err != nil {
				require.NoError(t, err)
			}
			err = res.Body.Close()
			if err != nil {
				require.NoError(t, err)
			}

			if c.out == "" {
				require.Equal(t, c.out, string(resp))
			} else {
				require.JSONEq(t, c.out, string(resp))
			}
		})
	}
}

func BenchmarkServeHTTP(b *testing.B) {
	rpc := jsonrpc.NewServer(jsonrpc.Options{})

	sumService := SumService{}
	rpc.Register("sum", sumService.Sum)

	var tc = []struct {
		name, in string
	}{
		{
			name: "OK",
			in:   `{"jsonrpc": "2.0", "method": "sum", "params": [1, 2, 3, 4], "id": "1" }`,
		},
	}

	for _, c := range tc {
		b.Run("route:"+c.name, func(b *testing.B) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/", bytes.NewBufferString(c.in))
			r.Header.Set("Content-Type", "application/json")

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				rpc.ServeHTTP(w, r)
			}
		})
	}
}
