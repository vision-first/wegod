package services

import "github.com/995933447/log-go"

type SMS struct {
	logger log.Logger
}

func NewSMS(logger log.Logger) *SMS {
	return &SMS{
		logger: logger,
	}
}

func (s *SMS) SendMsg(msg string) error {
	return nil
}