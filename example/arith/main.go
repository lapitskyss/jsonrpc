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
