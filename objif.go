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
func (mgr *ConfigMgr) InitializeObjectHandles(infoFiles []string) bool {
	var objMap map[string]ConfigObjJson

	mgr.objHdlMap = make(map[string]ConfigObjInfo)
	for _, objFile := range infoFiles {
		bytes, err := ioutil.ReadFile(objFile)
		if err != nil {
			logger.Println("Error in reading Object configuration file", objFile)
			return false
		}
		err = json.Unmarshal(bytes, &objMap)
		if err != nil {
			logger.Printf("Error in unmarshaling data from ", objFile)
		}

		for k, v := range objMap {
			logger.Printf("For Object [ %s ] Primary owner is [ %s ]\n", k, v.Owner)
			entry := new(ConfigObjInfo)
			entry.owner = mgr.clients[v.Owner]
			for _, lsnr := range v.Listeners {
				entry.listeners = append(entry.listeners, mgr.clients[lsnr])
			}
			mgr.objHdlMap[k] = *entry
		}
	}
	return true
}
