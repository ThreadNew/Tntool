package error

import (
	"fmt"
	"math"
	"sync"
)

// Error 错误结构定义
type Error struct {
	Code      int32  // 错误码
	Msg       string // 错误提示信息 （对外展示）
	DetailMsg string // 详细的错误提示信息 （服务内部使用）
}

const (
	defaultModuleNumber = 0 // 服务模版序号起始值 系统默认
	systemErrCode       = 999
)

var (
	moduleNumber int32 = defaultModuleNumber // 服务模块序号
	errCodeDigit int32 = 5                   // 错误码位数

	once sync.Once
)

// InitErr 统一注册服务模版号和错误码位数
func InitErr(moduleSeq int32) {
	once.Do(func() {
		moduleNumber = moduleSeq
	})
}

// New 初始化
func New(code int32, msg string) *Error {
	if moduleNumber == defaultModuleNumber {
		panic("Uninitialized serial number")
	}
	mn := moduleNumber * int32(math.Pow(float64(10), float64(errCodeDigit)))
	if code > 0 {
		code = mn + code
	}
	if code < 0 {
		code = code - mn
	}
	err := &Error{
		Code: code,
		Msg:  msg,
	}
	return err
}

// GetCode 获取错误码
func (e *Error) GetCode() int32 {
	return e.Code
}

// GetMsg 错误提示信息 （对外展示）
func (e *Error) GetMsg() string {
	return e.Msg
}

// GetDetailMsg 详细的错误提示信息 （服务内部使用）
func (e *Error) GetDetailMsg() string {
	return e.DetailMsg
}

// NewMsg 设置msg
func (e *Error) NewMsg(msg string) *Error {
	e.Msg = msg
	return e
}

// NewDetailMsg 设置详细的msg
func (e *Error) NewDetailMsg(detailMsg string) *Error {
	e.DetailMsg = detailMsg
	return e
}

// ToString 错误信息打印
func (e *Error) ToString() string {
	// 按照特定的格式输出
	formatStr := `{"code":%v,"msg":"%v","detail":"%v"}`
	return fmt.Sprintf(formatStr, e.Code, e.Msg, e.DetailMsg)
}

// Error() 实现error接口
func (e *Error) Error() string {
	return e.ToString()
}

// ToError error接口和Error转换
func ToError(err error) *Error {
	e, ok := err.(*Error)
	if ok {
		return e
	} else {
		return New(systemErrCode, err.Error())
	}
}
