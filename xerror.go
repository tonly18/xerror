package xerror

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type Error interface {
	Error() string
}

// XError 是一个可携带上下文信息的错误类型。
type XError struct {
	Code  int
	Msg   string
	Cause error
	File  string
	Line  int
	Func  string
}

// NewXError 创建新错误（带文件、行号、函数信息）
func NewXError(msg string, code ...int) *XError {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	if len(code) > 0 {
		return &XError{
			Code: code[0],
			Msg:  msg,
			File: filepath.Base(file),
			Line: line,
			Func: shortFunc(fn.Name()),
		}
	}

	return &XError{
		Msg:  msg,
		File: filepath.Base(file),
		Line: line,
		Func: shortFunc(fn.Name()),
	}
}

// XError 实现 error 接口
func (e *XError) Error() string {
	if e.Cause != nil {
		if e.Msg == "" {
			return fmt.Sprintf("%s:%d [%s] %v", e.File, e.Line, e.Func, e.Cause)
		}
		return fmt.Sprintf("%s:%d [%s] %d %s | %v", e.File, e.Line, e.Func, e.Code, e.Msg, e.Cause)
	}

	return fmt.Sprintf("%s:%d [%s] %d %s", e.File, e.Line, e.Func, e.Code, e.Msg)
}

// Unwrap 支持 errors.Unwrap / errors.Is / errors.As
func (e *XError) Unwrap() error {
	return e.Cause
}

// Wrap 在原有错误上包裹一层上下文信息
func Wrap(err error, msg string, code ...int) *XError {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	errorCode := 0
	if len(code) > 0 {
		errorCode = code[0]
	}

	//err is nil
	if err == nil {
		return &XError{
			Code: errorCode,
			Msg:  msg,
			File: filepath.Base(file),
			Line: line,
			Func: shortFunc(fn.Name()),
		}
	}

	//err not XError
	var xe *XError
	if !errors.As(err, &xe) {
		err = &XError{
			Cause: err,
			File:  filepath.Base(file),
			Line:  line,
			Func:  shortFunc(fn.Name()),
		}
	}

	return &XError{
		Code:  errorCode,
		Msg:   msg,
		Cause: err,
		File:  filepath.Base(file),
		Line:  line,
		Func:  shortFunc(fn.Name()),
	}
}

// shortFunc 去除冗余包路径
func shortFunc(name string) string {
	if i := strings.LastIndex(name, "/"); i != -1 {
		return name[i+1:]
	}

	return name
}

// FormatStack 递归输出错误堆栈（可用于日志打印）
func FormatStack(err error) string {
	if err == nil {
		return ""
	}

	var xerr *XError
	var b strings.Builder
	for e := err; e != nil; e = errors.Unwrap(e) {
		switch {
		case errors.As(e, &xerr):
			b.WriteString(e.Error())
			b.WriteString("\n")
		}
	}

	return b.String()
}

// FirstXError 获取最外层错误
func FirstXError(err error) *XError {
	var xe *XError
	if errors.As(err, &xe) {
		return xe
	}
	return nil
}

// Range 循环处理
func Range(err error, handle func(er error)) {
	for e := err; e != nil; e = errors.Unwrap(e) {
		handle(e)
	}
}
