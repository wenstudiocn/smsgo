package smsgo

import (
	"errors"
	"fmt"
)

var (
	ErrSmsNoMoney                     = errors.New("sms no money")
	ErrSmsProviderReturnBadValue      = errors.New("sms provider return bad value")
	ErrSmsProviderNotFound            = errors.New("sms provider not found")
	ErrSmsProviderReturnBadStatusCode = func(status_code int) error {
		return errors.New(fmt.Sprintf("status code=%d\n", status_code))
	}

	providers map[string]ISms = make(map[string]ISms)
)

type ISms interface {
	SendTextMessage(phone, text string, cb func(bool)) error
	GetBalance() (float64, error)
	GetAvailable() (int64, error)
	Name() string
}

func Use(sms ...ISms) {
	for _, p := range sms {
		providers[p.Name()] = p
	}
}

func All() map[string]ISms {
	return providers
}

func Get(name string) (ISms, error) {
	sms, ok := providers[name]
	if ok {
		return sms, nil
	}
	return nil, ErrSmsProviderNotFound
}

func Clear() {
	providers = map[string]ISms{}
}
