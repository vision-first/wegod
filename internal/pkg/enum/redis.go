package enum

import "fmt"

const (
	RedisKeyPhoneVerifyCode = "phone:%s_verifyCode"
)

func MakeRedisKeyPhoneVerifyCode(phone string) string {
	return fmt.Sprintf(RedisKeyPhoneVerifyCode, phone)
}



