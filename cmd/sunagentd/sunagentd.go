package main

import (
	"expinc/sunagent/common"
	"expinc/sunagent/http"
	"expinc/sunagent/log"
	"expinc/sunagent/ops"
	"fmt"

	"github.com/go-ini/ini"
)

func main() {
	defer func() {
		if p := recover(); nil != p {
			log.Error(fmt.Sprintf("%s panicked: pid=%d", common.ProcName, common.Pid))
			log.Fatal(p)
		}
	}()
	defer func() {
		log.Info(fmt.Sprintf("%s stopped: pid=%d", common.ProcName, common.Pid))
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

	// load grimoire
	err = ops.ReloadGrimoireFromFile()
	if nil != err {
		log.Fatal(err)
	}

	// start HTTP server
	ip := config.Section("HTTP").Key("ip").String()
	port := config.Section("HTTP").Key("port").MustUint()
	server := http.New(ip, port)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
