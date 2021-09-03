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
