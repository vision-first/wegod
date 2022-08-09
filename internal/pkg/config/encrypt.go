package config

type Jwt struct {
	SigningKey string
}

type Encrypt struct {
	*Jwt
}
