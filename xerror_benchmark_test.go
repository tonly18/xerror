package xerror_test

import (
	"errors"
	"fmt"
	"github.com/tonly18/xerror"
	"testing"
)

func BenchmarkBXError(b *testing.B) {
	err := &xerror.NewError{
		Code:    500000000,
		Err:     errors.New("test handler bag.query"),
		Message: "test bag.query",
	}

	//run b.N times
	for n := 0; n < b.N; n++ {
		xerror.Wrap(err, &xerror.NewError{
			Code:    uint32(n * 100000000),
			Err:     fmt.Errorf(`test handler bag.query:%d`, n),
			Message: "test stack",
		})
	}

	fmt.Println("over::::::", len(err.GetStack()), err.GetCode(), err.GetMsg())
}
