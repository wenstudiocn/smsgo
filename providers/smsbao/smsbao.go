package smsbao

import (
	"fmt"
	"github.com/wenstudiocn/smsgo"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	url_SENDSMS   = "http://api.smsbao.com/sms"
	query_SENDSMS = "u=%s&p=%s&m=%s&c=%s"
	url_BALANCE   = "http://www.smsbao.com/query"
	query_BALANCE = "u=%s&p=%s"
)

func New(username, password string) smsgo.ISms {
	return &smsBao{
		username: username,
		password: password,
	}
}

type smsBao struct {
	username, password string
}

func (self *smsBao) SendTextMessage(phone, text string, cb func(bool)) error {
	uri := fmt.Sprintf(url_SENDSMS+"?"+query_SENDSMS, url_SENDSMS,
		self.username, self.password, phone, url.QueryEscape(text))
	var err error = nil
	if cb != nil {
		defer func() {
			cb(err == nil)
		}()
	}
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if string(body) != "0" {
			err = smsgo.ErrSmsNoMoney
		}
	} else {
		err = smsgo.ErrSmsProviderReturnBadStatusCode(resp.StatusCode)
	}
	return err
}

func (self *smsBao) GetBalance() (float64, error) {
	return 0, nil
}

func (self *smsBao) GetAvailable() (int64, error) {
	uri := fmt.Sprintf(url_BALANCE+"?"+query_BALANCE, self.username, self.password)
	resp, err := http.Get(uri)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		sbody := fmt.Sprintf("%s", body)

		lines := strings.Split(sbody, "\n")
		if len(lines) < 2 {
			return 0, smsgo.ErrSmsProviderReturnBadValue
		}
		if lines[0] != "0" {
			return 0, smsgo.ErrSmsProviderReturnBadValue
		}
		parts := strings.Split(lines[1], ",")
		if len(parts) < 2 {
			return 0, smsgo.ErrSmsProviderReturnBadValue
		}
		remains, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, err
		}
		return remains, nil
	} else {
		return 0, smsgo.ErrSmsProviderReturnBadStatusCode(resp.StatusCode)
	}
}

func (self *smsBao) Name() string {
	return "smsbao.com"
}
