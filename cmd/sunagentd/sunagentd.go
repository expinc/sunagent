package main

import (
	"expinc/sunagent/common"
	"expinc/sunagent/http"
	"expinc/sunagent/log"
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

func exit(status int, msg interface{}) {
	log.Fatal(msg)
	os.Exit(status)
}

func main() {
	defer func() {
		if p := recover(); nil != p {
			log.Fatal(p)
		}
	}()

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
	log.SetRotateFileOutput("logs/"+common.ProcName+".log", fileLimitMb)
	log.Info(fmt.Sprintf("%s started: pid=%d", common.ProcName, common.Pid))

	// start HTTP server
	ip := config.Section("HTTP").Key("ip").String()
	port := config.Section("HTTP").Key("port").MustUint()
	server := http.New(ip, port)
	if err := server.Run(); err != nil {
		exit(1, err)
	}

	// successful exit
	log.Info(fmt.Sprintf("%s stopped: pid=%d", common.ProcName, common.Pid))
}
