package main

import (
	"encoding/json"
	"io/ioutil"
	"models"
)

//
// This structure represents the json layout for config objects
type ConfigObjJson struct {
	Owner     string   `json:"Owner"`
	Access    string   `json: "access"`
	Listeners []string `json:"Listeners"`
}

//
// This structure represents the in memory layout of all the config object handlers
type ConfigObjInfo struct {
	owner     ClientIf
	access    string
	listeners []ClientIf
}

func (mgr *ConfigMgr) CreateObjectMap() {
	//models.ConfigObjectMap
	for objName, obj := range models.GenConfigObjectMap {
		models.ConfigObjectMap[objName] = obj
	}
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
			logger.Printf("For Object [ %s ] Primary owner is [ %s ] access is %s\n", k, v.Owner, v.Access)
			entry := new(ConfigObjInfo)
			entry.owner = mgr.clients[v.Owner]
			entry.access = v.Access
			for _, lsnr := range v.Listeners {
				entry.listeners = append(entry.listeners, mgr.clients[lsnr])
			}
			mgr.objHdlMap[k] = *entry
		}
	}
	return true
}
