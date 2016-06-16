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

package actions

import (
	"config/clients"
	"encoding/json"
	"fmt"
	"io/ioutil"
	modelActions "models/actions"
	"utils/logging"
)

//
// Actions are methods exposed by various daemons. These may have an object as parameter.
// The only methods supported on these actions would be POST methods
//

//
// ActionManager provides the following methods for rest of the config manager subsystem
//  -- Initialize
//  -- DeInitialize
//  -- RegisterActions
//  -- PerformAction
//

type ActionMgr struct {
	logger *logging.Writer
	//	ExecuteAction(obj actions.ActionObj) error
	ObjHdlMap map[string]ActionObjInfo
	clientMgr *clients.ClientMgr
}

var gActionMgr *ActionMgr

//
// This structure represents the json layout for action objects
type ActionObjJson struct {
	Owner string `json:"Owner"`
}

// This structure represents the in memory layout of all the action object handlers
type ActionObjInfo struct {
	Owner clients.ClientIf
}

func InitializeActionMgr(infoFiles []string, logger *logging.Writer, clientMgr *clients.ClientMgr) *ActionMgr {
	mgr := new(ActionMgr)
	mgr.logger = logger
	mgr.clientMgr = clientMgr
	if rc := mgr.InitializeActionObjectHandles(infoFiles); !rc {
		logger.Err("Error in initializing action object handles")
		return nil
	}
	gActionMgr = mgr
	return mgr
}

func (mgr *ActionMgr) InitializeActionObjectHandles(infoFiles []string) bool {
	var actionMap map[string]ActionObjJson

	mgr.ObjHdlMap = make(map[string]ActionObjInfo)
	for _, objFile := range infoFiles {
		bytes, err := ioutil.ReadFile(objFile)
		if err != nil {
			mgr.logger.Info(fmt.Sprintln("Error in reading Action configuration file", objFile))
			return false
		}
		err = json.Unmarshal(bytes, &actionMap)
		if err != nil {
			mgr.logger.Info(fmt.Sprintln("Error in unmarshaling data from ", objFile))
		}

		for k, v := range actionMap {
			mgr.logger.Debug(fmt.Sprintln("For Action [", k, "] Primary owner is [", v.Owner, "] "))
			entry := new(ActionObjInfo)
			entry.Owner = mgr.clientMgr.Clients[v.Owner]
			mgr.ObjHdlMap[k] = *entry
		}
	}
	return true
}

func (mgr *ActionMgr) GetAllActions() []string {
	retList := make([]string, 0)
	for key, _ := range modelActions.ActionMap {
		retList = append(retList, key)
	}
	return retList
}
