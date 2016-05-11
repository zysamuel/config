package main

import (
	"config/server"
	"flag"
	"fmt"
	"net/http"
	"utils/keepalive"
	"utils/logging"
)

func main() {
	fmt.Println("Starting ConfigMgr daemon")
	paramsDir := flag.String("params", "./params", "Directory Location for config files")
	flag.Parse()
	paramsDirName := *paramsDir
	if paramsDirName[len(paramsDirName)-1] != '/' {
		paramsDirName = paramsDirName + "/"
	}

	fmt.Println("ConfigMgr: Start logger")
	logger, err := logging.NewLogger("confd", "ConfigMgr", true)
	if err != nil {
		fmt.Println("Failed to start logger. Nothing will be logged ...")
	}

	configMgr := server.NewConfigMgr(*paramsDir, logger)
	if configMgr == nil {
		logger.Err("Failed to initialize CONF Mgr. Exiting!!!")
		return
	}

	foundConfPort, confPort := server.GetConfigHandlerPort(*paramsDir)
	restRtr := configMgr.ApiMgr.GetRestRtr()

	// Start keepalive routine
	go keepalive.InitKeepAlive("confd", paramsDirName)
	/*
		// TODO: uncomment this section for https server
		certFile := *paramsDir+"/cert.pem"
		keyFile := *paramsDir+"/key.pem"
		err = ConfigMgrCheck(certFile, keyFile)
		if err != nil {
			err = ConfigMgrGenerate(certFile, keyFile)
			if err != nil {
				syslogger.Info("### CONF Mgr Failed to generate certs")
			}
		}
		if foundConfPort {
			http.ListenAndServeTLS(":"+confPort, certFile, keyFile, restRtr)
		} else
			http.ListenAndServeTLS(":8080", certFile, keyFile, restRtr)
		}
	*/
	if foundConfPort {
		http.ListenAndServe(":"+confPort, restRtr)
	} else {
		http.ListenAndServe(":8080", restRtr)
	}
	logger.Info("CONF Mgr. Exiting!!!")
}
