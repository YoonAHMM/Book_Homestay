package errx

import "fmt"

const OK uint32 = 200

const SERVER_ERROR uint32=100001
const REDIS_ERROR=100002
const DB_ERROR=100003
const WXMINI_ERROR=100004
const JWT_ERROR=100005
const LOGIN_ERROR=100006

type CodeError struct {
	errCode uint32
	errMsg  string
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.errCode, e.errMsg)
}

func NewErrCode(errCode uint32,s string) *CodeError {
	return &CodeError{errCode: errCode, errMsg: s}
}