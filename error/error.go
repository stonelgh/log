package error

const (
	ErrOK        = 0
	ErrUndefined = 10000 + iota
	ErrFormat
	ErrTimeout
	ErrUnknown = 9999
)

type Error struct {
	code   int
	msg    string
	detail string
}

func NewError(code int, msg string) *Error {
	return &Error{code, msg, ""}
}

func (e *Error) Is(code int) bool {
	return e.code == code
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Detail() string {
	return e.detail
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) SetDetail(detail string) {
	e.detail = detail
}

func (e *Error) IsOK() bool {
	return e.code == ErrOK
}

func (e *Error) IsTimeout() bool {
	return e.code == ErrTimeout
}
