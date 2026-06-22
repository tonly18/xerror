package xerror

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestTXError(t *testing.T) {
	err := HandleUserRequest(404)

	fmt.Printf("err1:%+v\n", err)
	fmt.Printf("err2:%+v\n", FirstXError(err))
	fmt.Printf("err3:%+v\n", FirstXError(err).Msg)

	fmt.Println("---res:", errors.Is(err, io.EOF))
	fmt.Println("---res:", errors.Is(err, io.ErrClosedPipe))

	//fmt.Println("捕获到错误：")
	//fmt.Println(FormatStack(err))

	//errs := strings.Split(FormatStack(err), "\n")
	//for _, v := range errs {
	//	fmt.Println(v)
	//}

	//first error
	firstError := FirstXError(err)
	fmt.Println("------firstError:", firstError.Code, firstError.Msg)

	Range(err, func(er error) {
		fmt.Println("------er:", er)
	})
}

func HandleUserRequest(id int) *XError {
	err := GetUserInfo(id)
	if err != nil {
		return Wrap(err, "333")
	}
	return nil
}

func GetUserInfo(id int) *XError {
	err := GetUserByID(id)
	if err != nil {
		return Wrap(err, "222", 777777)
		//return Wrap(io.ErrClosedPipe, "222")
	}
	return nil
}

func GetUserByID(id int) *XError {
	//return NewXError("111")
	return Wrap(io.EOF, "111", 666666)
	//return Wrap(nil, "111")
}
