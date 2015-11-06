package main

import (
	"log"
	"net/http"
	"os"
)

var logger *log.Logger
var gMgr *ConfigMgr

func main() {
	logger = log.New(os.Stdout, "ConfigMgr:", log.Ldate|log.Ltime|log.Lshortfile)
	configFile := "./params/clients.json"
	gMgr = NewConfigMgr(configFile)
	logger.Println("### GMgr is ", gMgr)
	go gMgr.ConnectToAllClients()
	restRtr := gMgr.GetRestRtr()
	http.ListenAndServe(":8080", restRtr)
}
