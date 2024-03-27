package error

import (
	"errors"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	InitErr(12)
	fmt.Print(ToError(Get()).ToString())
	fmt.Print(ToError(Set()).ToString())
	err := ToError(Get())
	if err.GetCode() != 1200102 {
		t.Error("code is match")
	}
}
func Get() error {
	return New(102, "rpc failed").NewDetailMsg("调用什么接口失败")
}
func Set() error {
	return errors.New("系统错误")
}
