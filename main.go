package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var logger *log.Logger
var gMgr *ConfigMgr

func main() {
	logger = log.New(os.Stdout, "ConfigMgr:", log.Ldate|log.Ltime|log.Lshortfile)

	paramsDir := flag.String("params", "", "Directory Location for config files")
	flag.Parse()
	gMgr = NewConfigMgr(*paramsDir)
	if gMgr == nil {
		return
	}
	go gMgr.ConnectToAllClients()
	restRtr := gMgr.GetRestRtr()
	http.ListenAndServe(":8080", restRtr)
}
