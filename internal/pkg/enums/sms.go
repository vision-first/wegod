package enums

import "fmt"

const (
	MsgTemplateVerifyCodeForRegister = "您的注册验证码为%s"
	MsgTemplateVerifyCodeForLogin = "您的登录验证码为%s"
)

func MakeVerifyCodeMsgForRegister(code string) string {
	return fmt.Sprintf(MsgTemplateVerifyCodeForRegister, code)
}

func MakeVerifyCodeMsgForLogin(code string) string {
	return fmt.Sprintf(MsgTemplateVerifyCodeForLogin, code)
}