package main

import (
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

/*
   @Time : 2019/6/18 15:53 
   @Author : ff
   @DESC   : Transport层用于接收用户网络请求并将其转为Endpoint可以处理的对象，然后交由Endpoint执行，最后将处理结果转为响应对象向用户响应。
*/
var (
	ErrorBadRequest = errors.New("invalid request parameter")
)
func MakeHttpHandler(ctx context.Context, endpoint endpoint.Endpoint, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
	}

	r.Methods("POST").Path("/calculate/{type}/{a}/{b}").Handler(kithttp.NewServer(
		endpoint,
		decodeArithmeticRequest,
		encodeArithmeticResponse,
		options...,
	))

	return r
}
//解码器
func decodeArithmeticRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	requestType, ok := vars["type"]
	if !ok {
		return nil, ErrorBadRequest
	}

	pa, ok := vars["a"]
	if !ok {
		return nil, ErrorBadRequest
	}

	pb, ok := vars["b"]
	if !ok {
		return nil, ErrorBadRequest
	}

	a, _ := strconv.Atoi(pa)
	b, _ := strconv.Atoi(pb)

	return ArithmeticRequest{
		RequestType: requestType,
		A:           a,
		B:           b,
	}, nil
}

// 编码器
func encodeArithmeticResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}