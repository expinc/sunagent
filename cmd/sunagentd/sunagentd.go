package main

import (
	"context"
	"expinc/sunagent/common"
	"expinc/sunagent/config"
	"expinc/sunagent/http"
	"expinc/sunagent/log"
	"expinc/sunagent/ops"
	"flag"
	"fmt"
	"runtime"

	"github.com/alecthomas/units"
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

	// parse flags
	configFile := flag.String("config", "config.yml", "Configuration file")
	grimoireFolder := flag.String("grimoire", "grimoires", "Grimoire folder path")
	certFilePath := flag.String("certFile", "", "Certificate file to enable HTTPS")
	keyFilePath := flag.String("keyFile", "", "Key file to enable HTTPS")
	flag.Parse()

	// load config
	cfg, err := config.LoadConfig(*configFile)
	if nil != err {
		log.Fatal(err)
	}

	// set process behaviors
	if 0 == cfg.Runtime.MaxCpus {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(int(cfg.Runtime.MaxCpus))
	}

	// init log
	log.SetLevel(cfg.Log.Level)
	fileSizeLimit, err := units.ParseBase2Bytes(cfg.Log.FileSizeLimit)
	if nil != err {
		log.Fatal(err)
	}
	log.SetRotateFileOutput("logs/"+common.ProcName+".log", int(fileSizeLimit/units.MiB))
	log.Info(fmt.Sprintf("%s started: pid=%d", common.ProcName, common.Pid))

	// load grimoire
	ops.GrimoireFolder = *grimoireFolder
	err = ops.ReloadGrimoire(context.Background())
	if nil != err {
		log.Fatal(err)
	}

	// setup core behaviors
	ops.SetJobCleanThreshold(context.Background(), int(cfg.Core.JobCleanThreshold))
	fileTransferSizeLimit, err := units.ParseBase2Bytes(cfg.Core.FileTransferSizeLimit)
	if nil != err {
		log.Fatal(err)
	}
	common.FileTransferSizeLimit = fileTransferSizeLimit

	// start HTTP server
	ip := cfg.Http.Ip
	port := cfg.Http.Port
	actualAuthMethod := ""
	var authCred interface{}
	// FIXME: add more authentication methods
	if "basic" == cfg.Http.Auth.Method {
		actualAuthMethod = cfg.Http.Auth.Method
		authCred = http.BasicAuthCred{
			User:     cfg.Http.Auth.Credential["user"].(string),
			Password: cfg.Http.Auth.Credential["password"].(string),
		}
	}
	server := http.New(ip, port, actualAuthMethod, authCred, *certFilePath, *keyFilePath)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
