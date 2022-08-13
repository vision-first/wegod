package enums

import "fmt"

const (
	MsgTemplateVerifyCodeForRegister = "您的注册验证码为%s"
)

func MakeVerifyCodeMsgForRegister(code string) string {
	return fmt.Sprintf(MsgTemplateVerifyCodeForRegister, code)
}
