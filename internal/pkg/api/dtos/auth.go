package dtos

type RegisterReq struct {
	Phone string
	Password string
	NickName string
	Avatar string
	Desc string
	Gender uint8
}

type RegisterResp struct {
}

type LoginReq struct {
}

type LoginResp struct {
}

type SendVerifyCodeForRegisterReq struct {
}

type SendVerifyCodeForRegisterResp struct {
}