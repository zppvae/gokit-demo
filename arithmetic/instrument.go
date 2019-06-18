package main

import (
	"errors"
	"github.com/juju/ratelimit"
	"github.com/go-kit/kit/endpoint"
	"context"
	"golang.org/x/rate"
)

/*
   @Time : 2019/6/18 17:15 
   @Author : ff
   @DESC   : 
*/

var ErrLimitExceed = errors.New("Rate limit exceed!")

//使用juju/ratelimit创建限流中间件
func NewTokenBucketLimitterWithJuju(bkt *ratelimit.Bucket) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			//获取令牌
			if bkt.TakeAvailable(1) == 0 {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

//使用x/time/rate创建限流中间件
func NewTokenBucketLimitterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}