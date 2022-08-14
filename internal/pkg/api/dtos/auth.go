package dtos

type RegisterReq struct {
	Phone string
	Password string
	VerifyCode string
}

type RegisterResp struct {
	Token string
}

type LoginReq struct {
	Phone string
	Password string
	Code string
	IsByCode bool
}

type LoginResp struct {
	Token string
}

type SendVerifyCodeForRegisterReq struct {
	Phone string
}

type SendVerifyCodeForRegisterResp struct {
}

type SendVerifyCodeForLoginReq struct {
	Phone string
}

type SendVerifyCodeForLoginResp struct {
}