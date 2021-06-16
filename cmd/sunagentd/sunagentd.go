package main

import (
	"expinc/sunagent/log"
	"os"

	"github.com/go-ini/ini"
)

const procName = "sunagentd"

func exit(status int, msg interface{}) {
	log.Fatal(msg)
	os.Exit(status)
}

func main() {
	// load config
	config, err := ini.Load("config.conf")
	if nil != err {
		exit(1, err)
	}

	// init log
	log.SetLevel(config.Section("LOG").Key("level").String())
	fileLimitMb, err := config.Section("LOG").Key("filelimitmb").Int()
	if nil != err {
		exit(1, err)
	}
	log.SetRotateFileOutput("logs/"+procName+".log", fileLimitMb)
	log.Info(procName)
}
