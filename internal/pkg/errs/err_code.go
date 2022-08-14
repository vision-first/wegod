package errs

import "github.com/995933447/apperrdef"

const (
	ErrCodeInternal apperrdef.ErrCode = -500
	ErrCodeUnauthorized apperrdef.ErrCode = -401
	ErrPhoneBeenRegisteredByOtherUser apperrdef.ErrCode = 1001
	ErrBadPhoneVerifyCode = 1002
	ErrUserNotFound = 1003
	ErrBuddhaNotFound = 1004
	ErrPasswordNotCorrect = 1005
)