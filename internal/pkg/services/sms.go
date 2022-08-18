package services

import "github.com/995933447/log-go"

type SMS struct {
	logger *log.Logger
}

func NewSMS(logger *log.Logger) *SMS {
	return &SMS{
		logger: logger,
	}
}

func (s *SMS) TransErr(err error) error {
	return err
}

func (s *SMS) SendMsg(phone, msg string) error {
	return nil
}