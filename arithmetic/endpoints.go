package main

import (
	"errors"
	"github.com/go-kit/kit/endpoint"
	"strings"
	"context"
)

/*
   @Time : 2019/6/18 15:39 
   @Author : ff
   @DESC   : 请求、响应类型
*/
//endpoint，调用Service中相应的方法处理请求对象（ArithmeticRequest），返回响应对象（ArithmeticResponse）
var (
	ErrInvalidRequestType = errors.New("RequestType has only four type: Add,Subtract,Multiply,Divide")
)

type ArithmeticRequest struct {
	RequestType string `json:"request_type"`
	A           int    `json:"a"`
	B           int    `json:"b"`
}

type ArithmeticResponse struct {
	Result int   `json:"result"`
	Error  error `json:"error"`
}

func MakeArithmeticEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ArithmeticRequest)

		var (
			res, a, b int
			calError  error
		)

		a = req.A
		b = req.B

		if strings.EqualFold(req.RequestType, "Add") {
			res = svc.Add(a, b)
		} else if strings.EqualFold(req.RequestType, "Substract") {
			res = svc.Subtract(a, b)
		} else if strings.EqualFold(req.RequestType, "Multiply") {
			res = svc.Multiply(a, b)
		} else if strings.EqualFold(req.RequestType, "Divide") {
			res, calError = svc.Divide(a, b)
		} else {
			return nil, ErrInvalidRequestType
		}

		return ArithmeticResponse{Result: res, Error: calError}, nil
	}
}