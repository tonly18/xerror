package xerror_test

import (
	"database/sql"
	"fmt"
	"github.com/tonly18/xerror"
	"io"
	"net"
	"os"
	"testing"
)

func TestTXError(t *testing.T) {
	_, err := A(100)

	fmt.Println("err::::::::::::", err.String())

	if err != nil {
		if err.Is(net.ErrClosed) {
			errList := err.GetStack()
			for _, e := range errList {
				fmt.Println("stack::::::", e.GetCode(), e.GetMsg(), e.GetRawError())
			}
		}
	}

}

func A(uid int) (int, xerror.Error) {
	data, err := B(uid)
	if err != nil {
		if err.Is(os.ErrClosed) {
			xerr := xerror.Wrap(err, &xerror.NewError{
				Code:     20005000,
				RawError: net.ErrClosed,
				Message:  "a-message",
			})
			//fmt.Println("a-err:::::::", xerr.GetCode(), xerr.GetRawError(), xerr.GetMsg(), len(xerr.GetStack()))
			return 0, xerr
		}
	}

	return data, nil
}

func B(uid int) (int, xerror.Error) {
	_, err := C(uid)
	if err != nil {
		if err.Is(sql.ErrNoRows) {
			xerr := xerror.Wrap(err, &xerror.NewError{
				Code:     20005010,
				RawError: os.ErrClosed,
				Message:  "b-message",
			})
			//fmt.Println("b-err:::::::", xerr.GetCode(), xerr.GetRawError(), xerr.GetMsg(), len(xerr.GetStack()))
			return 0, xerr
		}
	}

	return 1, nil
}

func C(uid int) (int, xerror.Error) {
	_, err := D(uid)
	if err != nil {
		if err.Is(io.ErrClosedPipe) {
			xerr := xerror.Wrap(err, &xerror.NewError{
				Code:     20005020,
				RawError: sql.ErrNoRows,
				Message:  "c-message",
			})
			//fmt.Println("c-err:::::::", xerr.GetCode(), xerr.GetRawError(), xerr.GetMsg(), len(xerr.GetStack()))
			return 0, xerr
		}
	}

	return 1, nil
}

func D(uid int) (int, xerror.Error) {
	_, err := E(uid)
	if err != nil {
		if err.Is(os.ErrPermission) {
			xerr := xerror.Wrap(err, &xerror.NewError{
				Code:     20005030,
				RawError: io.ErrClosedPipe,
				Message:  "d-message",
			})
			//fmt.Println("d-err:::::::", xerr.GetCode(), xerr.GetRawError(), xerr.GetMsg(), len(xerr.GetStack()))
			return 0, xerr
		}
	}

	return 1, nil
}

func E(uid int) (int, xerror.Error) {
	err := os.ErrPermission
	if err == os.ErrPermission {
		xerr := xerror.Wrap(&xerror.NewError{
			Code:     20005040,
			RawError: err,
			Message:  "e-message",
		}, nil)
		//fmt.Println("e-err:::::::", xerr.GetCode(), xerr.GetRawError(), xerr.GetMsg(), len(xerr.GetStack()))
		return 0, xerr
	}

	return 1, nil
}
