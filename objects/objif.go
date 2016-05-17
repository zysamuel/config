//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package objects

import (
	"config/clients"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"models"
	"net/http"
	"utils/logging"
)

type ObjectMgr struct {
	logger    *logging.Writer
	ObjHdlMap map[string]ConfigObjInfo
	clientMgr *clients.ClientMgr
}

var gObjectMgr *ObjectMgr

//
// This structure represents the json layout for config objects
type ConfigObjJson struct {
	Owner     string   `json:"Owner"`
	Access    string   `json:"Access"`
	Listeners []string `json:"Listeners"`
	PerVRF    bool     `json:"perVRF"`
}

//
// This structure represents the in memory layout of all the config object handlers
type ConfigObjInfo struct {
	Owner     clients.ClientIf
	Access    string
	PerVRF    bool
	Listeners []clients.ClientIf
}

const (
	MAX_JSON_LENGTH = 4096
)

func GetConfigObj(r *http.Request, obj models.ConfigObj) (body []byte, retobj models.ConfigObj, err error) {
	if obj == nil {
		err = errors.New("Config Object is nil")
		return body, retobj, err
	}
	if r != nil {
		body, err = ioutil.ReadAll(io.LimitReader(r.Body, MAX_JSON_LENGTH))
		if err != nil {
			return body, retobj, err
		}
		if err = r.Body.Close(); err != nil {
			return body, retobj, err
		}
	}

	retobj, err = obj.UnmarshalObject(body)
	if err != nil {
		err = errors.New("Failed to decode input json data")
	}
	return body, retobj, err
}

func GetUpdateKeys(body []byte) (map[string]bool, error) {
	var objmap map[string]*json.RawMessage
	var err error
	updateKeys := make(map[string]bool)

	err = json.Unmarshal(body, &objmap)
	if err != nil {
		return updateKeys, err
	}
	for key, _ := range objmap {
		updateKeys[key] = true
	}
	return updateKeys, err
}

func CreateObjectMap() {
	//models.ConfigObjectMap
	for objName, obj := range models.GenConfigObjectMap {
		models.ConfigObjectMap[objName] = obj
	}
}

func InitializeObjectMgr(infoFiles []string, logger *logging.Writer, clientMgr *clients.ClientMgr) *ObjectMgr {
	mgr := new(ObjectMgr)
	mgr.logger = logger
	mgr.clientMgr = clientMgr
	if rc := mgr.InitializeObjectHandles(infoFiles); !rc {
		logger.Err("Error in initializing object handles")
		return nil
	}
	gObjectMgr = mgr
	return mgr
}

//
//  This method reads the config file and connects to all the clients in the list
//
func (mgr *ObjectMgr) InitializeObjectHandles(infoFiles []string) bool {
	var objMap map[string]ConfigObjJson

	mgr.ObjHdlMap = make(map[string]ConfigObjInfo)
	for _, objFile := range infoFiles {
		bytes, err := ioutil.ReadFile(objFile)
		if err != nil {
			mgr.logger.Info(fmt.Sprintln("Error in reading Object configuration file", objFile))
			return false
		}
		err = json.Unmarshal(bytes, &objMap)
		if err != nil {
			mgr.logger.Info(fmt.Sprintln("Error in unmarshaling data from ", objFile))
		}

		for k, v := range objMap {
			mgr.logger.Info(fmt.Sprintf("For Object [", k, "] Primary owner is [", v.Owner, "] access is",
				v.Access, " Global Object ", v.PerVRF))
			entry := new(ConfigObjInfo)
			entry.Owner = mgr.clientMgr.Clients[v.Owner]
			entry.Access = v.Access
			entry.PerVRF = v.PerVRF
			for _, lsnr := range v.Listeners {
				entry.Listeners = append(entry.Listeners, mgr.clientMgr.Clients[lsnr])
			}
			mgr.ObjHdlMap[k] = *entry
		}
	}
	return true
}
func (mgr *ObjectMgr) GetConfigObjHdlMap() map[string]ConfigObjInfo {
	return mgr.ObjHdlMap
}
