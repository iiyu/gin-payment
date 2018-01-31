package main

import (
	"log"

	"github.com/dlintw/goconf"
)

var (
	WxAppId     string
	WxMchId     string
	WxAppKey    string
	WxAppSecret string
)

func init() {
	conf, err := goconf.ReadConfigFile(basePath + "/src/gin-payment/.env")
	if err != nil {
		log.Println(err)
		return
	}
	WxAppId, _ = conf.GetString("wxpay", "wxappid")
	WxMchId, _ = conf.GetString("wxpay", "wxmchid")
	WxAppKey, _ = conf.GetString("wxpay", "wxappkey")
	WxAppSecret, _ = conf.GetString("wxpay", "wxappsecret")
}
