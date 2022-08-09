package encrypt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

//JWT 定义一个jwt对象
type JWT struct {
	// 声明签名信息
	SigningKey []byte
}

//NewJWT 初始化jwt对象
func NewJWT(signingKey string) *JWT {
	return &JWT{
		[]byte(signingKey),
	}
}

type Claims struct {
	UserId 	uint64 `json:"user_id"`
	jwt.StandardClaims
}

func (j *JWT) GenerateToken(userId uint64) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(2 * time.Minute)
	claims := Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(token string) (*Claims, error,bool) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if validationError,ok := err.(*jwt.ValidationError);ok {
			expiredOk := validationError.Errors&jwt.ValidationErrorExpired==jwt.ValidationErrorExpired
			if expiredOk {
				return nil, nil, expiredOk
			}
		}
		return nil, err, false
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil, false
		}
	}

	return nil, err,false
}
