package xerror

import (
	"bytes"
	"fmt"
)

type NewError struct {
	RawError error
	Code     uint32
	Message  string
	stack    []Error
}

func (e *NewError) Error() string {
	return e.String()
}

func (e *NewError) GetRawError() error {
	return e.RawError
}

func (e *NewError) SetRawError(err error) {
	e.RawError = err
}

func (e *NewError) GetCode() uint32 {
	return e.Code
}

func (e *NewError) SetCode(code uint32) {
	e.Code = code
}

func (e *NewError) GetMsg() string {
	return e.Message
}

func (e *NewError) SetMsg(msg string) {
	e.Message = msg
}

func (e *NewError) GetStack() []Error {
	return e.stack
}

func (e *NewError) pushStack(err Error) {
	if len(e.stack) == 0 {
		e.stack = make([]Error, 0, 10)
	}
	e.stack = append(e.stack, err)
	e.SetRawError(err.GetRawError())
	e.SetCode(err.GetCode())
	e.SetMsg(err.GetMsg())
}

func (e *NewError) Is(err error) bool {
	return e.GetRawError() == err
}

func (e *NewError) Contain(err error) bool {
	for _, e := range e.GetStack() {
		if e.Is(err) {
			return true
		}
	}
	return false
}

func (e *NewError) String() string {
	var errMessage bytes.Buffer
	for _, err := range e.GetStack() {
		if err.GetRawError() == nil {
			errMessage.WriteString(fmt.Sprintf(`{xcode:%v, xmessage:%v, xerror:nil}`, err.GetCode(), err.GetMsg()))
		} else {
			errMessage.WriteString(fmt.Sprintf(`{xcode:%v, xmessage:%v, xerror:%v}`, err.GetCode(), err.GetMsg(), err.GetRawError().Error()))
		}
	}
	return errMessage.String()
}

// Wrap 老的错误信息包裹新的错误信息
//
// @params
//
//	originalError	Error			原始Error
//	newErrors		[]Error			新的Error
//
// @return
//
//	Error
func Wrap(originalError, newError Error) Error {
	if originalError == nil {
		panic("the parameter originalError is nil")
	}
	if originalError != nil {
		if _, ok := originalError.GetRawError().(Error); ok {
			panic("xerror not wrap self")
		}
	}
	if newError != nil {
		if _, ok := newError.GetRawError().(Error); ok {
			panic("xerror not wrap self")
		}
	}

	if newError == nil {
		originalError.pushStack(&NewError{
			RawError: originalError.GetRawError(),
			Code:     originalError.GetCode(),
			Message:  originalError.GetMsg(),
		})
	} else {
		originalError.pushStack(newError)
	}

	//return
	return originalError
}
