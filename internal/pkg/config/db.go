package config

type MysqlLog struct {
	SlowThreshold int `json:"slow_threshold"`
	Level string `json:"level"`
}

type Mysql struct {
	Database string `json:"database"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Charset string `json:"charset"`
	Host string `json:"host"`
	Port uint32 `json:"port"`
	MaxConns int `json:"max_conns"`
	MaxIdleConns int `json:"max_idle_conns"`
	Log *MysqlLog
}

type DB struct {
	*Mysql `json:"mysql"`
}
