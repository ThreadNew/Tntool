package util

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"
)

/**
通用重试逻辑封装
1、支持设置重试时间
2、支持设置重试次数
3、支持根据特定错误码重试
*/

const (
	DefaultInterval = time.Duration(10) // 10 ms
	DefaultTimes    = int32(3)          // 重试3次
)

// RetryConfig 定义retry的配置类
type RetryConfig struct {
	interval  time.Duration //重试的间隔时长 (注休眠时长 interval*time.Millisecond ms)
	times     int32         // 重试次数
	errHandle ErrHandleFunc // 非特定错误可以排除重试
}

//ErrHandleFunc 错误处理句柄
type ErrHandleFunc func(err error) bool

//EmptyErrHandleFunc 空处理逻辑 默认都需要重试
func EmptyErrHandleFunc(err error) bool {
	return true
}

// 定义retry的optional
type RetryConfigOptional func(c *RetryConfig)

// WithInterval 设置重试的间隔时长
func WithInterval(interval time.Duration) RetryConfigOptional {
	return func(c *RetryConfig) {
		c.interval = interval
	}
}

// WithTimes 设置重试次数
func WithTimes(times int32) RetryConfigOptional {
	return func(c *RetryConfig) {
		c.times = times
	}
}

//WithErrHandle 设置错误处理handle
func WithErrHandle(ehFunc ErrHandleFunc) RetryConfigOptional {
	return func(c *RetryConfig) {
		if ehFunc == nil {
			return
		}
		c.errHandle = ehFunc
	}
}

// 定义通用的func
type RetryFunc func(ctx context.Context) error

// Retry 重试逻辑
func Retry(ctx context.Context, fun RetryFunc, opts ...RetryConfigOptional) error {
	if fun == nil {
		return errors.New(" please set the RetryFunc")
	}
	// panic 捕获
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(string(debug.Stack()))
		}
	}()
	var err error
	// 重试配置初始化
	rConfig := &RetryConfig{
		interval:  DefaultInterval,
		times:     DefaultTimes,
		errHandle: EmptyErrHandleFunc,
	}
	for _, opt := range opts {
		opt(rConfig)
	}
	for i := int32(0); i < rConfig.times; i++ {
		err = fun(ctx)
		if err == nil {
			return err
		}
		retry := rConfig.errHandle(err)
		if !retry {
			return err
		}
		time.Sleep(rConfig.interval * time.Millisecond)
	}
	return err
}
