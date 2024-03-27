package util

import (
	"context"
	"errors"
	"fmt"
	bizErr "github.com/ThreadNew/tntool/error"
	"testing"
)

func TestRetry(t *testing.T) {

	Retry(context.Background(), RetrySayHelloFunc)

}

func TestRetry2(t *testing.T) {
	Retry(context.Background(), RetryBadFunc)
	fmt.Println("---------")
	// 重试五次
	Retry(context.Background(), RetryBadFunc, WithTimes(5))
	fmt.Println("---------")
	// 调整间隔时间
	Retry(context.Background(), RetryBadFunc, WithTimes(5), WithInterval(10))
	fmt.Println("---------")
	// 设置错误处理句柄
	bizErr.InitErr(12)
	Retry(context.Background(), SayErrRequest, WithTimes(5), WithErrHandle(SayErrHandleFunc))
}

func TestRetry3(t *testing.T) {
	// panic捕获
	Retry(context.Background(), SayPanicRequest, WithTimes(5))
}

func RetrySayHelloFunc(ctx context.Context) error {
	return SayHello()
}

// 定义一个函数
func SayHello() error {
	fmt.Print("hello world!")
	return nil
}

func RetryBadFunc(ctx context.Context) error {
	return SayBadRequest()
}
func SayBadRequest() error {
	fmt.Println("bad request")
	return errors.New("bad request")
}

func SayErrRequest(ctx context.Context) error {
	fmt.Println("bad request")
	return bizErr.New(100, "badRequest")
}
func SayErrHandleFunc(err error) bool {
	if bizErr.ToError(err).Code == 1200100 {
		return false
	}
	return true
}
func SayPanicRequest(ctx context.Context) error {
	fmt.Println("bad request panic")
	panic("SayPanicRequest")
}
