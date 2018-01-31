package main

import (
	"github.com/cihub/seelog"
	"github.com/dlintw/goconf"
)

var (
	DomainUrl     string
	SessionSecret string
	SignupEnabled bool
)

func init() {
	conf, err := goconf.ReadConfigFile(basePath + "/src/gin-payment/.env")
	if err != nil {
		seelog.Critical("err parsing config log file", err)
		return
	}
	DomainUrl, _ = conf.GetString("app", "domainurl")
	SessionSecret, _ = conf.GetString("app", "session_secret")
	SignupEnabled, _ = conf.GetBool("app", "signup_enabled")
}
