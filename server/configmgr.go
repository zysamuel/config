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

package server

import (
	"config/actions"
	"config/apis"
	"config/clients"
	"config/objects"
	"encoding/json"
	"fmt"
	"io/ioutil"
	modelObjs "models/objects"
	"os"
	"os/signal"
	"syscall"
	"time"
	"utils/logging"
)

type ConfigMgr struct {
	logger       *logging.Writer
	paramsDir    string
	dbHdl        *objects.DbHandler
	bringUpTime  time.Time
	swVersion    SwVersion
	ApiMgr       *apis.ApiMgr
	clientMgr    *clients.ClientMgr
	objectMgr    *objects.ObjectMgr
	actionMgr    *actions.ActionMgr
	clientNameCh chan string
}

var gConfigMgr *ConfigMgr

const (
	MAX_COUNT_AUTO_DISCOVER_OBJ int64 = 200
)

type ConfdGlobals struct {
	Name  string `json: "Name"`
	Value string `json: "Value"`
}

// Get the http port on which rest api calls will be received
func GetConfigHandlerPort(paramsDir string) (bool, string) {
	var globals []ConfdGlobals
	var port string

	globalsFile := paramsDir + "/globals.json"
	bytes, err := ioutil.ReadFile(globalsFile)
	if err != nil {
		gConfigMgr.logger.Err(fmt.Sprintln("Error in reading globals file", globalsFile))
		return false, port
	}

	err = json.Unmarshal(bytes, &globals)
	if err != nil {
		gConfigMgr.logger.Err("Failed to Unmarshall Json")
		return false, port
	}
	for _, global := range globals {
		if global.Name == "httpport" {
			port = global.Value
			return true, port
		}
	}
	return false, port
}

//
// This function would work as a classical constructor for the
// configMgr object
//
func NewConfigMgr(paramsDir string, logger *logging.Writer) *ConfigMgr {
	mgr := new(ConfigMgr)
	mgr.logger = logger
	mgr.paramsDir = paramsDir

	mgr.dbHdl = objects.InstantiateDbIf(logger)

	paramsFile := paramsDir + "/clients.json"

	mgr.clientMgr = clients.InitializeClientMgr(paramsFile, logger, GetSystemStatus, GetSystemSwVersion, actions.ExecuteConfigurationAction)

	objects.CreateObjectMap()
	objectConfigFiles := [...]string{paramsDir + "/genObjectConfig.json"}
	mgr.objectMgr = objects.InitializeObjectMgr(objectConfigFiles[:], logger, mgr.clientMgr)

	actions.CreateActionMap()
	actionConfigFiles := [...]string{paramsDir + "/genObjectAction.json"}
	mgr.actionMgr = actions.InitializeActionMgr(paramsDir, actionConfigFiles[:], logger, mgr.dbHdl, mgr.objectMgr, mgr.clientMgr)

	mgr.ApiMgr = apis.InitializeApiMgr(paramsDir, logger, mgr.dbHdl, mgr.clientMgr, mgr.objectMgr, mgr.actionMgr)

	mgr.ApiMgr.InitializeRestRoutes()
	mgr.ApiMgr.InitializeActionRestRoutes()
	mgr.ApiMgr.InitializeEventRestRoutes()
	mgr.ApiMgr.InstantiateRestRtr()

	mgr.bringUpTime = time.Now()
	// Initialize channel to receive connected client name.
	// When confd connects to a client, it creates global objects owned by that client and
	// stores default logging level in DB, if it does not exist.
	// Global objects and logging objects can only be updated by user.
	mgr.clientNameCh = make(chan string, 10)
	logger.Info("Initialization Done!")

	mgr.ReadSystemSwVersion()
	go mgr.AutoCreateConfigObjects()
	go mgr.clientMgr.ListenToClientStateChanges()
	go mgr.clientMgr.ConnectToAllClients(mgr.clientNameCh)
	go mgr.SigHandler()
	gConfigMgr = mgr

	return mgr
}

func (mgr *ConfigMgr) SigHandler() {
	sigChan := make(chan os.Signal, 1)
	signalList := []os.Signal{syscall.SIGHUP}
	signal.Notify(sigChan, signalList...)

	for {
		select {
		case signal := <-sigChan:
			switch signal {
			case syscall.SIGHUP:
				mgr.logger.Info("Exting!!!")
				os.Exit(0)
			default:
			}
		}
	}
}

func (mgr *ConfigMgr) storeUUID(key string) {
	_, err := mgr.dbHdl.StoreUUIDToObjKeyMap(key)
	if err != nil {
		mgr.logger.Err(fmt.Sprintln(
			"Failed to store uuid map for key ", key, err))
	}
}

func (mgr *ConfigMgr) ConfigureGlobalConfig(clientName string) {
	var obj modelObjs.ConfigObj
	var err error
	if ent, ok := mgr.objectMgr.AutoCreateObjMap[clientName]; ok {
		for _, resource := range ent.ObjList {
			if objHdl, ok := modelObjs.ConfigObjectMap[resource]; ok {
				var body []byte
				obj, _ = objHdl.UnmarshalObject(body)
				gConfigMgr.dbHdl.DbLock.Lock()
				_, err = obj.GetObjectFromDb(obj.GetKey(), mgr.dbHdl)
				gConfigMgr.dbHdl.DbLock.Unlock()
				if err != nil {
					client, exist := mgr.clientMgr.Clients[clientName]
					if exist {
						client.LockApiHandler()
						err, success := client.CreateObject(obj, mgr.dbHdl.DBUtil)
						client.UnlockApiHandler()
						if err == nil && success == true {
							mgr.storeUUID(obj.GetKey())
						} else {
							mgr.logger.Err(fmt.Sprintln("Failed to create "+resource+" ", obj, err))
						}
					}
				}
			}
		}
	}
}

func (mgr *ConfigMgr) AutoCreateConfigObjects() {
	for {
		select {
		case clientName := <-mgr.clientNameCh:
			switch clientName {
			case "Client_Init_Done":
				close(mgr.clientNameCh)
				return
			default:
				mgr.logger.Info("Do Global Init and Discover objects for Client: " + clientName)
				mgr.ConstructSystemParam(clientName)
				mgr.ConfigureGlobalConfig(clientName)
				mgr.AutoDiscoverObjects(clientName)
				mgr.ConfigureComponentLoggingLevel(clientName)
				mgr.logger.Info("Done Global Init and Discover objects for Client: " + clientName)
			}
		}
	}
}

func (mgr *ConfigMgr) AutoDiscoverObjects(clientName string) {
	fmt.Println("AutoDiscover for: ", clientName)
	if ent, ok := mgr.objectMgr.AutoDiscoverObjMap[clientName]; ok {
		for _, resource := range ent.ObjList {
			fmt.Println("AutoDiscover: ", resource)
			if objHdl, ok := modelObjs.ConfigObjectMap[resource]; ok {
				var objs []modelObjs.ConfigObj
				var err error
				_, obj, _ := objects.GetConfigObj(nil, objHdl)
				currentIndex := int64(0)
				objCount := int64(MAX_COUNT_AUTO_DISCOVER_OBJ)
				resourceOwner := mgr.objectMgr.ObjHdlMap[resource].Owner
				resourceOwner.LockApiHandler()
				err, _, _, _, objs = resourceOwner.GetBulkObject(obj, mgr.dbHdl.DBUtil,
					currentIndex, objCount)
				resourceOwner.UnlockApiHandler()
				fmt.Println("AutoDiscover response: ", err, objs)
				if err == nil {
					for _, obj := range objs {
						gConfigMgr.dbHdl.DbLock.Lock()
						_, err := obj.GetObjectFromDb(obj.GetKey(), mgr.dbHdl)
						gConfigMgr.dbHdl.DbLock.Unlock()
						if err != nil {
							gConfigMgr.dbHdl.DbLock.Lock()
							err = obj.StoreObjectInDb(mgr.dbHdl)
							gConfigMgr.dbHdl.DbLock.Unlock()
							if err != nil {
								mgr.logger.Err(fmt.Sprintln("Failed to store"+resource+" config in DB ", obj, err))
							} else {
								mgr.storeUUID(obj.GetKey())
							}
						}
					}
				}
			}
		}
	}
}
