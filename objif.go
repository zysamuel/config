package main

import (
	"encoding/json"
	"io/ioutil"
)

//
// This structure represents the json layout for config objects
type ConfigObjJson struct {
	Owner     string   `json:"Owner"`
	Listeners []string `json:"Listeners"`
}

//
// This structure represents the in memory layout of all the config object handlers
type ConfigObjInfo struct {
	owner     ClientIf
	listeners []ClientIf
}

//
//  This method reads the config file and connects to all the clients in the list
//
func (mgr *ConfigMgr) InitializeObjectHandles(objsFile string) bool {
	var objMap map[string]ConfigObjJson
	bytes, err := ioutil.ReadFile(objsFile)
	if err != nil {
		logger.Println("Error in reading configuration file")
		return false
	}
	err = json.Unmarshal(bytes, &objMap)

	mgr.objHdlMap = make(map[string]ConfigObjInfo)
	for k, v := range objMap {
		entry := new(ConfigObjInfo)
		entry.owner = mgr.clients[v.Owner]
		for _, lsnr := range v.Listeners {
			entry.listeners = append(entry.listeners, mgr.clients[lsnr])
		}
		mgr.objHdlMap[k] = *entry
	}
	mgr.objHdlMap["IPV4Route"].owner.CreateObject()
	return true
}
