package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"github.com/juju/ratelimit"
	"time"
)

func main() {

	ctx := context.Background()
	errChan := make(chan error)

	var svc Service
	svc = ArithmeticService{}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	//日志中间件
	svc = LoggingMiddleware(logger)(svc)
	//创建Endpoint对象
	endpoint := MakeArithmeticEndpoint(svc)

	//添加限流，每秒刷新一次，容量为3
	ratebucket := ratelimit.NewBucket(time.Second*1, 3)
	endpoint = NewTokenBucketLimitterWithJuju(ratebucket)(endpoint)

	//创建http处理对象handlder
	r := MakeHttpHandler(ctx, endpoint, logger)

	go func() {
		fmt.Println("Http Server start at port:9000")
		handler := r
		errChan <- http.ListenAndServe(":9000", handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	fmt.Println(<-errChan)
}