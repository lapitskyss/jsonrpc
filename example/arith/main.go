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
