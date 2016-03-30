package main

import (
	"flag"
	"log"
	"log/syslog"
	"net/http"
	"os"
)

var logger *log.Logger
var gMgr *ConfigMgr

func main() {
	logger = log.New(os.Stdout, "ConfigMgr:", log.Ldate|log.Ltime|log.Lshortfile)
	syslogger, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_INFO|syslog.LOG_DAEMON, "ConfigMgr")
	if err == nil {
		syslogger.Info("### CONF Mgr started")
		logger.SetOutput(syslogger)
	}

	paramsDir := flag.String("params", "", "Directory Location for config files")
	flag.Parse()
	gMgr = NewConfigMgr(*paramsDir)
	if gMgr == nil {
		return
	}
	clientsUp := make(chan bool, 1)
	go gMgr.CreateDefaultUser()
	go gMgr.ReadConfiguredUsersFromDb()
	go gMgr.ConnectToAllClients(clientsUp)
	go gMgr.DiscoverSystemObjects(clientsUp)
	go gMgr.MonitorSystemStatus()
	go gMgr.StartUserSessionHandler()

	foundConfPort, confPort := gMgr.GetConfigHandlerPort(*paramsDir)
	restRtr := gMgr.GetRestRtr()
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
}
