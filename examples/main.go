package main

import (
	"fmt"
	"github.com/wenstudiocn/smsgo"
	"github.com/wenstudiocn/smsgo/providers/smsbao"
)

func main() {
	// load config to initialzie sms params
	bao := smsbao.New("username", "password")
	smsgo.Use(bao)

	// there we decide to use smsbao by something
	bao, _ = smsgo.Get("smsbao.com")
	_ = bao.SendTextMessage("13122223333", "test message", nil)
	// async send
	go bao.SendTextMessage("13122223333", "async send", nil)
	// async send and promise the result
	go bao.SendTextMessage("13122223333", "async send again", func(ok bool) {
		if ok {
			fmt.Println("send succeed")
		} else {
			fmt.Println("send failed")
		}
	})
}
