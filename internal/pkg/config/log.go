package config

type File struct {
	Dir string `json:"dir"`
	MaxFileSize int64 `json:"max_file_size"`
}

type Drivers struct {
	*File `json:"file"`
}

type Log struct {
	*Drivers
	Driver string `json:"driver"`
}

const (
	LogDriverFile = "file"
)