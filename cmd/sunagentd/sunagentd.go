package main

import (
	"expinc/sunagent/common"
	"expinc/sunagent/http"
	"expinc/sunagent/log"
	"expinc/sunagent/ops"
	"flag"
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
	configFile := flag.String("config", "config.conf", "Configuration file")
	flag.Parse()
	config, err := ini.Load(*configFile)
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

	// setup core behaviors
	ops.SetJobCleanThreshold(config.Section("CORE").Key("jobCleanThreshold").MustInt())

	// start HTTP server
	ip := config.Section("HTTP").Key("ip").String()
	port := config.Section("HTTP").Key("port").MustUint()
	authMethod := config.Section("HTTP").Key("auth").String()
	var authCred interface{}
	if "basic" == authMethod {
		authCred = http.BasicAuthCred{
			User:     config.Section("HTTP").Key("user").String(),
			Password: config.Section("HTTP").Key("password").String(),
		}
	}
	server := http.New(ip, port, authMethod, authCred)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
