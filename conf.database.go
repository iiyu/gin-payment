package main

import (
	"github.com/cihub/seelog"
	"github.com/dlintw/goconf"
)

var (
	Conn   string
	JwtKey string
)

func init() {
	conf, err := goconf.ReadConfigFile(basePath + "/src/gin-payment/.env")
	if err != nil {
		seelog.Critical("err parsing config log file", err)
		return
	}
	user, _ := conf.GetString("mysql", "user")
	password, _ := conf.GetString("mysql", "password")
	host, _ := conf.GetString("mysql", "host")
	port, _ := conf.GetString("mysql", "port")
	db, _ := conf.GetString("mysql", "db")
	jwtkey, _ := conf.GetString("jwt", "jwtkey")
	Conn = user + ":" + password + "@tcp(" + host + ":" + port + ")/" + db
	JwtKey = jwtkey
}
