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
	"config/objects"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	modelActions "models/actions"
	modelObjs "models/objects"
	"net/http"
	"os"
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
	logger    *logging.Writer
	paramsDir string
	dbHdl     *objects.DbHandler
	//	ExecuteAction(obj actions.ActionObj) error
	ObjHdlMap map[string]ActionObjInfo
	clientMgr *clients.ClientMgr
	objectMgr *objects.ObjectMgr
}

var gActionMgr *ActionMgr
var ApplyConfigOrder = []string{
	"SystemLogging",
	"ComponentLogging",
	"PolicyCondition",
	"PolicyStmt",
	"PolicyDefinition",
	"Port",
	"LaPortChannel",
	"LLDPIntf",
	"Vlan",
	"StpBridgeInstance",
	"StpPort",
	"ArpGlobal",
	"ArpConfig",
	"LogicalIntf",
	"IPv4Intf",
	"SubIPv4Intf",
	"IPv4Route",
	"IpTableAcl",
	"BfdGlobal",
	"BfdInterface",
	"BfdSession",
	"BGPGlobal",
	"BGPNeighbor",
	"BGPPeerGroup",
	"BGPPolicyCondition",
	"BGPPolicyStmt",
	"BGPPolicyDefinition",
	"OspfAreaAggregateEntry",
	"OspfAreaEntry",
	"OspfGlobal",
	"OspfHostEntry",
	"OspfIfEntry",
	"OspfIfMetricEntry",
	"OspfNbrEntry",
	"OspfStubAreaEntry",
	"OspfVirtIfEntry",
	"VrrpIntf",
	"DhcpRelayGlobal",
	"DhcpRelayIntf",
	"VxlanInstance",
	"VxlanVtepInstances",
}

//temporary static list to indicate which objects are auto created
var AutoCreateList = "BGPGlobal," + "ArpGlobal," + "BfdGlobal," + "OspfGlobal," + "Port," + "DhcpRelayGlobal," + "SystemLogging," + "ComponentLogging"

// SR error codes
const (
	SRFail              = 0
	SRSuccess           = 1
	SRSystemNotReady    = 2
	SRRespMarshalErr    = 3
	SRNotFound          = 4
	SRIdStoreFail       = 5
	SRIdDeleteFail      = 6
	SRServerError       = 7
	SRObjHdlError       = 8
	SRObjMapError       = 9
	SRBulkGetTooLarge   = 10
	SRNoContent         = 11
	SRAuthFailed        = 12
	SRAlreadyConfigured = 13
	SRUpdateKeyError    = 14
	SRUpdateNoChange    = 15
)

// This structure represents the json layout for action objects
type ActionObjJson struct {
	Owner string `json:"Owner"`
}

// This structure represents the in memory layout of all the action object handlers
type ActionObjInfo struct {
	Owner clients.ClientIf
}

func InitializeActionMgr(paramsDir string, infoFiles []string, logger *logging.Writer, dbHdl *objects.DbHandler, objectMgr *objects.ObjectMgr, clientMgr *clients.ClientMgr) *ActionMgr {
	mgr := new(ActionMgr)
	mgr.paramsDir = paramsDir
	if logger == nil {
		gActionMgr.logger.Err("logger nil")
		return nil
	}
	mgr.logger = logger
	if clientMgr == nil {
		gActionMgr.logger.Err("clientMgr nil")
		return nil
	}
	mgr.clientMgr = clientMgr
	if objectMgr == nil {
		gActionMgr.logger.Err("objectMgr nil")
		return nil
	}
	mgr.objectMgr = objectMgr
	if dbHdl == nil {
		gActionMgr.logger.Err("dbHdl nil")
		return nil
	}
	mgr.dbHdl = dbHdl
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
			mgr.logger.Debug(fmt.Sprintln("Error in reading Action configuration file", objFile))
			return false
		}
		err = json.Unmarshal(bytes, &actionMap)
		if err != nil {
			mgr.logger.Debug(fmt.Sprintln("Error in unmarshaling data from ", objFile))
		}

		for k, v := range actionMap {
			mgr.logger.Debug(fmt.Sprintln("For Action [", k, "] Primary owner is [", v.Owner, "] "))
			entry := new(ActionObjInfo)
			if mgr.clientMgr != nil {
				entry.Owner = mgr.clientMgr.Clients[v.Owner]
			}
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

func GetActionObj(r *http.Request, obj modelActions.ActionObj) (body []byte, retobj modelActions.ActionObj, err error) {
	//var ret_obj map[string]modelActions.DummyStruct
	gActionMgr.logger.Debug(fmt.Sprintln("GetActionObj r:", r, " obj:", obj))
	if obj == nil {
		err = errors.New("Action Object is nil")
		return body, retobj, err
	}
	if r != nil {
		body, err = ioutil.ReadAll(io.LimitReader(r.Body, r.ContentLength))
		gActionMgr.logger.Debug(fmt.Sprintln("err:", err, " body:", body))
		if err != nil {
			return body, retobj, err
		}
		if err = r.Body.Close(); err != nil {
			return body, retobj, err
		}
	} else {
		fmt.Println("r nil")
		return body, retobj, err
	}
	retobj, err = obj.UnmarshalAction(body)
	if err != nil {
		fmt.Println("UnmarshalObject returned error", err, " for ojbect info", retobj)
	}
	return body, retobj, err
}
func UpdateConfig(resource string, body json.RawMessage) { //[]byte) {
	var success bool
	var err error
	var obj modelObjs.ConfigObj
	var objKey string

	gActionMgr.logger.Debug(fmt.Sprintln("update config resource:", resource))
	if objHdl, ok := modelObjs.ConfigObjectMap[resource]; ok {
		if obj, err = objHdl.UnmarshalObject(body); err == nil {
			objKey = obj.GetKey()
			updateKeys, _ := objects.GetUpdateKeys(body)
			dbObj, gerr := obj.GetObjectFromDb(objKey, gActionMgr.dbHdl.DBUtil)
			if gerr != nil {
				gActionMgr.logger.Err("entry not found in DB")
				return
			}
			_, err = gActionMgr.dbHdl.GetUUIDFromObjKey(objKey)
			diff, _ := obj.CompareObjectsAndDiff(updateKeys, dbObj)
			anyUpdated := false
			for _, updated := range diff {
				if updated == true {
					anyUpdated = true
					break
				}
			}
			if anyUpdated == false {
				gActionMgr.logger.Err("No updates to be made")
				return
			}

			mergedObj, _ := obj.MergeDbAndConfigObj(dbObj, diff)
			mergedObjKey := mergedObj.GetKey()
			if objKey == mergedObjKey {
				resourceOwner := gActionMgr.objectMgr.ObjHdlMap[resource].Owner
				if resourceOwner.IsConnectedToServer() == false {
					return
				}

				err, success = resourceOwner.UpdateObject(dbObj, mergedObj, diff, nil, objKey, gActionMgr.dbHdl.DBUtil)
				if err == nil && success == true {
					_, dbErr := gActionMgr.dbHdl.StoreUUIDToObjKeyMap(objKey)
					if dbErr == nil {
					} else {
						gActionMgr.logger.Err(fmt.Sprintln("Failed to store UuidToKey map ", obj, dbErr))
					}
				} else {
					gActionMgr.logger.Err(fmt.Sprintln("Failed to update object: ", obj, " due to error: ", err))
				}
			} else {
				gActionMgr.logger.Err(fmt.Sprintln("Failed to get object handle from http request ", objHdl, resource, err))
			}
		} else {
			fmt.Println("Failed to get object map")
			gActionMgr.logger.Err(fmt.Sprintln("Failed to get ObjectMap ", resource))
		}
	}
}
func CreateConfig(resource string, body json.RawMessage) {
	var errCode int
	var success bool
	var err error
	var obj modelObjs.ConfigObj
	var objKey string
	errCode = SRSuccess

	gActionMgr.logger.Debug(fmt.Sprintln("Create config resource:", resource))
	if objHdl, ok := modelObjs.ConfigObjectMap[resource]; ok {
		if obj, err = objHdl.UnmarshalObject(body); err == nil {
			updateKeys, _ := objects.GetUpdateKeys(body)
			if len(updateKeys) == 0 {
				errCode = SRNoContent
				gActionMgr.logger.Err("Nothing to configure")
			} else {
				objKey = obj.GetKey()
				fmt.Println("objKey derived")
				_, err = gActionMgr.dbHdl.GetUUIDFromObjKey(objKey)
				if err == nil {
					gActionMgr.logger.Err("Config object is present")
					UpdateConfig(resource, body)
					return
				}
			}

			if errCode != SRSuccess {
				fmt.Println("errcode not success, return")
				return
			}
			if gActionMgr.objectMgr.ObjHdlMap == nil {
				fmt.Println("objHdlMap nil")
				return
			}
			_, ok = gActionMgr.objectMgr.ObjHdlMap[resource]
			if !ok {
				fmt.Println("objhdlmap for resource:", resource, " nil")
				return
			}
			resourceOwner := gActionMgr.objectMgr.ObjHdlMap[resource].Owner
			if resourceOwner.IsConnectedToServer() == false {
				fmt.Println("Not connected to resourceOwner:", resourceOwner)
				return
			}
			fmt.Println("resource:", resource, " resourceOwner:", resourceOwner, " obj:", obj)

			err, success = resourceOwner.CreateObject(obj, gActionMgr.dbHdl.DBUtil)
			if err == nil && success == true {
				_, dbErr := gActionMgr.dbHdl.StoreUUIDToObjKeyMap(objKey)
				if dbErr == nil {
					errCode = SRSuccess
				} else {
					errCode = SRIdStoreFail
					gActionMgr.logger.Err(fmt.Sprintln("Failed to store UuidToKey map ", obj, dbErr))
				}
			} else {
				errCode = SRServerError
				gActionMgr.logger.Err(fmt.Sprintln("Failed to create object: ", obj, " due to error: ", err))
			}
		} else {
			errCode = SRObjHdlError
			gActionMgr.logger.Err(fmt.Sprintln("Failed to get object handle from http request ", objHdl, resource, err))
		}
	} else {
		fmt.Println("Failed to get object map")
		errCode = SRObjMapError
		gActionMgr.logger.Err(fmt.Sprintln("Failed to get ObjectMap ", resource))
	}

}
func ApplyConfigObject(data modelActions.ApplyConfig, resource string) {
	for key, value := range data.ConfigData {
		gActionMgr.logger.Debug(fmt.Sprintln("key:", key, "value:", value, " resoure:", resource))
		if resource != key {
			continue
		}
		fmt.Println("key:", key, "value:", value, " resoure:", resource)
		for _, v := range value {
			if _, err := json.Marshal(v); err == nil {
				CreateConfig(key, v)
			}
		}
	}
}

func SaveConfigObject(data modelActions.SaveConfigObj, resource string) error {
	gActionMgr.logger.Debug(fmt.Sprintln("SaveConfigObject for resource:", resource))
	objHdl, ok := modelObjs.ConfigObjectMap[resource]
	if !ok {
		gActionMgr.logger.Err("objHdl nil")
		return errors.New("objHdl Nil")
	}
	_, obj, err := objects.GetConfigObj(nil, objHdl)
	if err != nil {
		gActionMgr.logger.Err(fmt.Sprintln("GetConfigObj return err: ", err))
		return errors.New("getConfigObj return err")
	}
	var configObjects []modelObjs.ConfigObj
	err, objCount, _, _, configObjects := obj.GetBulkObjFromDb(0, 100, gActionMgr.dbHdl.DBUtil)
	if err != nil {
		gActionMgr.logger.Err(fmt.Sprintln("GetBulkObjFromDB returned error:", err))
		return errors.New("GetBulkObjFromDb returned error")
	}
	if objCount == 0 {
		gActionMgr.logger.Debug(fmt.Sprintln("No objects of type:", resource, " configured"))
		return nil
	}
	if data.ConfigData[resource] == nil {
		data.ConfigData[resource] = make([]interface{}, 0)
	}
	for _, configObject := range configObjects {
		data.ConfigData[resource] = append(data.ConfigData[resource], configObject)
	}
	return nil

}
func OpenFile(cfgFileName string) (fo *os.File, err error) {
	gActionMgr.logger.Debug(fmt.Sprintln("Full config file : ", cfgFileName))
	_, err = os.Stat(cfgFileName)
	if os.IsNotExist(err) {
		gActionMgr.logger.Debug(fmt.Sprintln(cfgFileName, " not present, create it"))
		fo, err = os.Create(cfgFileName)
		if err != nil {
			gActionMgr.logger.Err(fmt.Sprintln("Error :", err, " when creating file:", cfgFileName))
			return fo, err
		}
	} else if err == nil {
		// open cfg file
		gActionMgr.logger.Debug("cfgFile present, open it for update")
		fo, err = os.OpenFile(cfgFileName, os.O_RDWR, 0666)
		if err != nil {
			gActionMgr.logger.Err(fmt.Sprintln("Error:", err, "when opening cfgFile:", cfgFileName))
			return fo, err
		}
	} else {
		gActionMgr.logger.Err(fmt.Sprintln("Error:", err, " when handling the cfgFile:", cfgFileName))
		return fo, err
	}
	return fo, err
}

func ResetConfigObject(data modelActions.ResetConfig) (err error) {
	gActionMgr.logger.Debug(fmt.Sprintln("Start config reset"))

	configCount := len(ApplyConfigOrder)
	configCount = configCount - 1
	/* 1) Get all config objects */

	gActionMgr.logger.Debug(fmt.Sprintln("Get all object owners : "))
	//for key, objMap := range gActionMgr.objectMgr.ObjHdlMap {
	for index := configCount; index > -1; index-- {
		key := ApplyConfigOrder[index]
		objMap, ok := gActionMgr.objectMgr.ObjHdlMap[key]
		if !ok {
			gActionMgr.logger.Debug(fmt.Sprintln("Key ", key, " doesnt exist in ObjHdlMap"))
			continue
		}
		gActionMgr.logger.Debug(fmt.Sprintln("***************************************"))
		gActionMgr.logger.Debug(fmt.Sprintln("name ", key, "Access ", objMap.Access,
			"autocreate ", objMap.AutoCreate,
			"Owner ", objMap.Owner))
		if objMap.Owner.IsConnectedToServer() == false {
			gActionMgr.logger.Err(fmt.Sprintln("ResetConfig: Not connected to daemon ", key))
			continue
		}

		if objMap.Access == "w" && !objMap.AutoCreate {
			gActionMgr.logger.Debug(fmt.Sprintln("Get db objects for  ", key))
			//resource := objMap.Owner

			//get  object handle
			if objHdl, ok := modelObjs.ConfigObjectMap[key]; ok {
				_, obj, _ := objects.GetConfigObj(nil, objHdl)
				currentIndex := int64(0)
				objCount := int64(1024)
				err, _, _, _, objs := obj.GetBulkObjFromDb(currentIndex, objCount, gActionMgr.dbHdl.DBUtil)
				if err != nil {
					gActionMgr.logger.Debug(fmt.Sprintln("Failed to do getBulk object ", objMap.Owner))
				}
				gActionMgr.logger.Debug(fmt.Sprintln("No of objects collected ", len(objs)))
				for index := range objs {
					objKey := objs[index].GetKey()
					gActionMgr.logger.Debug(fmt.Sprintln("Obj ", objs[index], " key ", objKey))
					err, success := objMap.Owner.DeleteObject(objs[index], objKey, gActionMgr.dbHdl.DBUtil)
					if err == nil && success == true {
						gActionMgr.logger.Debug(fmt.Sprintln("Delete UUID to objectKeyMap"))
						uuid, er := gActionMgr.dbHdl.GetUUIDFromObjKey(objKey)
						if er == nil {
							err = gActionMgr.dbHdl.DeleteUUIDToObjKeyMap(uuid, objKey)
							if err != nil {
								gActionMgr.logger.Err(fmt.Sprintln("Failed to delete uuid map ", uuid))
							}
						}
					}
				}

			}
		}

	}
	return nil
}

func ExecutePerformAction(obj modelActions.ActionObj) (err error) {
	gActionMgr.logger.Debug(fmt.Sprintln("local client Execute action obj: ", obj))
	if gActionMgr == nil {
		gActionMgr.logger.Err("Action mgr not initialized")
		return err
	}
	switch obj.(type) {
	case modelActions.ApplyConfig:
		gActionMgr.logger.Debug("ApplyConfig")
		fmt.Println("ApplyConfig")
		data := obj.(modelActions.ApplyConfig)
		for _, applyResource := range ApplyConfigOrder {
			ApplyConfigObject(data, applyResource)
		}
	case modelActions.SaveConfig:
		gActionMgr.logger.Debug("SaveConfig")
		var fo *os.File
		var err error
		data := obj.(modelActions.SaveConfig)
		fileName := data.FileName
		gActionMgr.logger.Debug(fmt.Sprintln("FileName:", fileName))
		if fileName == "" {
			gActionMgr.logger.Debug("FileName not set, setting it to default startup-config")
			fileName = "startup-config"
		}
		// open config file
		cfgFileName := gActionMgr.paramsDir + "/" + fileName + ".json"
		fo, err = OpenFile(cfgFileName)
		if err != nil {
			gActionMgr.logger.Err(fmt.Sprintln("error with OpenFile, err:", err))
			return err
		}
		// close fo on exit and check for its returned error
		defer func() {
			if err := fo.Close(); err != nil {
				panic(err)
			}
		}()
		var wdata modelActions.SaveConfigObj
		wdata.ConfigData = make(map[string][]interface{})
		for _, applyResource := range ApplyConfigOrder {
			SaveConfigObject(wdata, applyResource)
		}
		js, err := json.Marshal(wdata)
		if err != nil {
			gActionMgr.logger.Err(fmt.Sprintln("json marshal returned error:", err))
			return err
		}
		gActionMgr.logger.Debug(fmt.Sprintln("js:", string(js)))
		_, err = fo.Write(js)
		if err != nil {
			gActionMgr.logger.Err(fmt.Sprintln("Error writing:", err))
			return err
		}

	case modelActions.ResetConfig:
		gActionMgr.logger.Debug("Action resolved as ResetConfig")
		data := obj.(modelActions.ResetConfig)
		ResetConfigObject(data)
	}
	return err
}
