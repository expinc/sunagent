package main

import (
	"expinc/sunagent/common"
	"expinc/sunagent/http"
	"expinc/sunagent/log"
	"fmt"

	"github.com/go-ini/ini"
)

func main() {
	defer func() {
		if p := recover(); nil != p {
			log.Fatal(p)
		}
	}()

	// load config
	config, err := ini.Load("config.conf")
	if nil != err {
		log.Fatal(err)
	}

	// init log
	log.SetLevel(config.Section("LOG").Key("level").String())
	fileLimitMb, err := config.Section("LOG").Key("filelimitmb").Int()
	if nil != err {
		log.Fatal(err)
	}
	log.SetRotateFileOutput("logs/"+common.ProcName+".log", fileLimitMb)
	log.Info(fmt.Sprintf("%s started: pid=%d", common.ProcName, common.Pid))

	// start HTTP server
	ip := config.Section("HTTP").Key("ip").String()
	port := config.Section("HTTP").Key("port").MustUint()
	server := http.New(ip, port)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

	// successful exit
	log.Info(fmt.Sprintf("%s stopped: pid=%d", common.ProcName, common.Pid))
}
