package errs

import "errors"

var (
	//ErrNoDefined 未定义
	ErrNoDefined = errors.New("error no defined")
	//ErrNotImplement 未实现
	ErrNotImplemented = errors.New("Not implemented")
	//ErrNotExist 记录不存在
	ErrNotExist = errors.New("Record does not exist")
	// zero row affect
	ErrZeroRowsAffected = errors.New("zero rows affect")
)

type Error struct {
	Code int
	Msg  string
}

//NewCustomError ...
func NewCustomError(code int, msg string) error {
	return Error{Code: code, Msg: msg}
}

//NewError ...
func NewError(code int) error {
	return Error{Code: code}
}

func (e Error) Error() string {
	return e.Msg
}

func (e Error) Unwrap() error {
	return e
}
