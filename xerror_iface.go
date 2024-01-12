package xerror

type Error interface {
	Error() string

	GetRawError() error
	SetRawError(error)
	GetCode() uint32
	SetCode(uint32)
	GetMsg() string
	SetMsg(string)

	GetStack() []Error
	pushStack(Error)

	Is(error) bool
	Contain(error) bool
	String() string
}
