package main

import (
	"os"

	"github.com/cihub/seelog"
)

var (
	Logger   seelog.LoggerInterface
	err      error
	basePath = os.Getenv("GOPATH")
)

func init() {
	Logger = seelog.Disabled
	Logger, err = seelog.LoggerFromConfigAsFile(basePath + "/src/gin-payment/seelog.xml")

	if err != nil {
		seelog.Critical("err parsing config log file", err)
	}
}
