package jsonrpc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	rpc := NewServer(Options{})

	sumService := SumService{}
	rpc.Register("sum", sumService.sum)

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
			name: "OKBatch",
			in:   `[{"jsonrpc":"2.0","method":"sum","params":[1, 2, 3, 4],"id":1}, {"jsonrpc":"2.0","method":"sum","params":[1, 2],"id":2}]`,
			out:  `[{"jsonrpc":"2.0","id":2,"result":3}, {"jsonrpc":"2.0","id":1,"result":10}]`,
		},
		//{
		//	name: "Notification",
		//	in:   `{"jsonrpc": "2.0", "method": "sum", "params": [1, 2, 3, 4] }`,
		//	out:  ``,
		//},
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
			res, err := http.Post(ts.URL, "application/json", bytes.NewBufferString(c.in))
			if err != nil {
				t.Errorf("Received unexpected error:\n%+v", err)
				t.FailNow()
			}

			resp, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Received unexpected error:\n%+v", err)
				t.FailNow()
			}
			err = res.Body.Close()
			if err != nil {
				t.Errorf("Received unexpected error:\n%+v", err)
				t.FailNow()
			}

			if !IsJSONEqual(c.out, string(resp)) {
				t.Errorf("Unexpected result. Expected %v. Got %v", c.out, string(resp))
				t.FailNow()
			}
		})
	}
}

func BenchmarkServeHTTP(b *testing.B) {
	rpc := NewServer(Options{})

	sumService := SumService{}
	rpc.Register("sum", sumService.sum)

	var tc = []struct {
		name, in string
	}{
		{
			name: "OK",
			in:   `{"jsonrpc": "2.0", "method": "sum", "params": [1, 2, 3, 4], "id": "1" }`,
		},
		{
			name: "OKBatch",
			in:   `[{"jsonrpc":"2.0","method":"sum","params":[1, 2, 3, 4],"id":1}, {"jsonrpc":"2.0","method":"sum","params":[1, 2],"id":2}]`,
		},
	}

	for _, c := range tc {
		b.Run("route:"+c.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				b.StopTimer()
				w := httptest.NewRecorder()
				r, _ := http.NewRequest("POST", "/", bytes.NewBufferString(c.in))
				r.Header.Set("Content-Type", "application/json")
				b.StartTimer()

				rpc.ServeHTTP(w, r)
			}
		})
	}
}

func Test_handleRequest(t *testing.T) {
	rpc := NewServer(Options{})
	sumService := SumService{}
	rpc.Register("sum", sumService.sum)

	r, _ := http.NewRequest("POST", "/", nil)
	j := []byte(`{"jsonrpc": "2.0", "method": "sum", "params": [1, 2, 3, 4], "id": "1" }`)

	res := rpc.handleRequest(r, j)
	expected := `{"jsonrpc":"2.0","result":10,"id":"1"}`

	if !IsJSONEqual(expected, string(res)) {
		t.Errorf("Unexpected result. Expected %v. Got %v", expected, string(res))
		t.FailNow()
	}
}

func Benchmark_handleRequest(b *testing.B) {
	rpc := NewServer(Options{})
	sumService := SumService{}
	rpc.Register("sum", sumService.sum)

	r, _ := http.NewRequest("POST", "/", nil)
	j := []byte(`{"jsonrpc": "2.0", "method": "sum", "params": [1, 2, 3, 4], "id": "1" }`)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rpc.handleRequest(r, j)
	}
}

func IsJSONEqual(expected string, actual string) bool {
	if expected == "" {
		return expected == actual
	}

	var e, a interface{}

	err := json.Unmarshal([]byte(expected), &e)
	if err != nil {
		return false
	}

	err = json.Unmarshal([]byte(actual), &a)
	if err != nil {
		return false
	}

	exp, ok := e.([]byte)
	if !ok {
		return reflect.DeepEqual(e, a)
	}

	act, ok := a.([]byte)
	if !ok {
		return false
	}

	return string(exp) == string(act)
}

type SumService struct {
}

func (ss *SumService) sum(ctx *RequestCtx) (Result, Error) {
	var sumRequest []int
	err := ctx.GetParams(&sumRequest)
	if err != nil {
		return nil, ErrInvalidParamsJSON()
	}

	s := 0
	for _, item := range sumRequest {
		s += item
	}

	return ctx.Result(s)
}
