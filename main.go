package main

import ("os"
        "log"
        "net/http")	

var logger *log.Logger 
func main() {
    logger = log.New(os.Stdout, "ConfigMgr:", log.Ldate|log.Ltime|log.Lshortfile)
	 configFile := "./params/clients.json"
    mgr := NewConfigMgr ( configFile)
    go mgr.ConnectToAllClients()
    restRtr := mgr.GetRestRtr()
    http.ListenAndServe(":8080", restRtr)
}
