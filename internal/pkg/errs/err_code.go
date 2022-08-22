package errs

import (
	"github.com/995933447/apperrdef"
)

const (
	ErrCodeInternal apperrdef.ErrCode = -500
	ErrCodeUnauthorized apperrdef.ErrCode = -401
	ErrCodePhoneBeenRegisteredByOtherUser apperrdef.ErrCode = 1001
	ErrCodeBadPhoneVerifyCode apperrdef.ErrCode = 1002
	ErrCodeUserNotFound apperrdef.ErrCode = 1003
	ErrCodeBuddhaNotFound apperrdef.ErrCode = 1004
	ErrCodePasswordNotCorrect apperrdef.ErrCode = 1005
	ErrCodeBuddhaNotPayedRent apperrdef.ErrCode = 1006
	ErrCodeBuddhaRentOrderNotFound apperrdef.ErrCode = 1007
	ErrCodePrayPropNotFound apperrdef.ErrCode = 1008
	ErrCodeBuddhaWorshipPropNotFound apperrdef.ErrCode = 1009
	ErrCodeBuddhaRentExpired apperrdef.ErrCode = 1100
	ErrCodeBuddhaNotRend apperrdef.ErrCode = 1111
	ErrCodeUserPrayPropNotFound apperrdef.ErrCode = 1112
	ErrCodeUserPrayNotFound apperrdef.ErrCode = 1113
	ErrCodePrayPropOrderNotFound apperrdef.ErrCode = 1114
	ErrCodeWorshipOrderNotFound apperrdef.ErrCode = 1115
	ErrCodeUserWorshipPropNotFound apperrdef.ErrCode = 1116
	ErrCodeDonationOrderNotFound apperrdef.ErrCode = 1117
	ErrCodeShopProductNotFound apperrdef.ErrCode = 1118
	ErrCodeCategoryNotFound apperrdef.ErrCode = 1119
	ErrCodePostNotFound apperrdef.ErrCode = 1120
)