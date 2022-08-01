package config

type Node struct {
	Host string
	Port int
	Password string
}

type Redis struct {
	Nodes []*Node `json:"nodes"`
}
