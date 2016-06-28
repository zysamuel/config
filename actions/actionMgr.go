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
        "Port",
        "LaPortChannel",
        "LLDPIntf",
        "Vlan",
        "StpBridgeInstance",
        "StpPort",
        "ArpConfig",
        "LogicalIntf",
        "IPv4Intf",
        "SubIPv4Intf",
        "IPv4Route",
        "IpTableAcl",
        "BfdGlobal",
        "BfdInterface",
        "BfdSession",
        "PolicyCondition",
        "PolicyStmt",
        "PolicyDefinition",
        "BGPGlobal",
        "BGPNeighbor",
        "BGPPeerGroup",
        "BGPPolicyAction",
        "BGPPolicyCondition",
        "BGPPolicyDefinition",
        "BGPPolicyDefinitionStmtPrecedence",
        "BGPPolicyStmt",
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

const (
	MAX_JSON_LENGTH = 4096
)

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

// SR error strings
var ErrString = map[int]string{
	SRFail:              "Configuration failed.",
	SRSuccess:           "None.",
	SRSystemNotReady:    "System not ready.",
	SRRespMarshalErr:    "Configuration applied successfully. However, failed to marshal response.",
	SRNotFound:          "Failed to find entry.",
	SRIdStoreFail:       "Failed to store Id in DB. However, configuration has been applied.",
	SRIdDeleteFail:      "Failed to delete Id from DB. However, configuration has been removed.",
	SRServerError:       "Backend server failed to apply configuration.",
	SRObjHdlError:       "Failed to get object handle.",
	SRObjMapError:       "Failed to get object map.",
	SRBulkGetTooLarge:   "More than maximum number of objects requested in a bulkget.",
	SRNoContent:         "Insufficient information.",
	SRAuthFailed:        "User authentication failed.",
	SRAlreadyConfigured: "Already configured. Delete and Update operations are allowed.",
	SRUpdateKeyError:    "Cannot update key in an object.",
	SRUpdateNoChange:    "Nothing to be updated.",
}

//Given a code reurn error string
func SRErrString(errCode int) string {
	return ErrString[errCode]
}

type ActionResponse struct {
	Error string `json:"Error"`
}

type ErrorResponse struct {
	Error string `json:"Error"`
}
type ConfigResponse struct {
	UUId  string `json:"ObjectId"`
	Error string `json:"Error"`
}

//
// This structure represents the json layout for action objects
type ActionObjJson struct {
	Owner string `json:"Owner"`
}

// This structure represents the in memory layout of all the action object handlers
type ActionObjInfo struct {
	Owner clients.ClientIf
}

func RespondErrorForApiCall(w http.ResponseWriter, errCode int, errString string) error {
	var errResp ErrorResponse
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if errCode == SRBulkGetTooLarge {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
	} else if errCode == SRSystemNotReady {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	errResp.Error = SRErrString(errCode) + " " + errString
	js, _ := json.Marshal(errResp)
	w.Write(js)
	return nil
}

func InitializeActionMgr(paramsDir string, infoFiles []string, logger *logging.Writer, dbHdl *objects.DbHandler, objectMgr *objects.ObjectMgr, clientMgr *clients.ClientMgr) *ActionMgr {
	mgr := new(ActionMgr)
	mgr.paramsDir = paramsDir
	mgr.logger = logger
	mgr.clientMgr = clientMgr
	mgr.objectMgr = objectMgr
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
			mgr.logger.Info(fmt.Sprintln("Error in reading Action configuration file", objFile))
			return false
		}
		err = json.Unmarshal(bytes, &actionMap)
		if err != nil {
			mgr.logger.Info(fmt.Sprintln("Error in unmarshaling data from ", objFile))
		}

		for k, v := range actionMap {
			mgr.logger.Info(fmt.Sprintln("For Action [", k, "] Primary owner is [", v.Owner, "] "))
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
	fmt.Println("GetActionObj r:", r, " obj:", obj)
	if obj == nil {
		err = errors.New("Action Object is nil")
		return body, retobj, err
	}
	if r != nil {
		body, err = ioutil.ReadAll(io.LimitReader(r.Body, r.ContentLength))
		fmt.Println("err:", err, " body:", body)
		if err != nil {
			return body, retobj, err
		}
		if err = r.Body.Close(); err != nil {
			return body, retobj, err
		}
	} else {
		fmt.Println("r nil, test case")
		objFile := "/home/madhavi/testCfg1.json"
		body, err = ioutil.ReadFile(objFile)
		if err != nil {
			fmt.Println("Error in reading Action configuration file", objFile)
			return body, retobj, err
		}
	}
	retobj, err = obj.UnmarshalAction(body)
	//err = json.Unmarshal(body,&ret_obj)
	if err != nil {
		fmt.Println("UnmarshalObject returnexd error", err, "for ojbect info", retobj)
	}
	//fmt.Println("ret_obj:",ret_obj)
	return body, retobj, err
}
func CreateConfig(resource string, body []byte) {
	//var w http.ResponseWriter
	var resp ConfigResponse
	var errCode int
	var success bool
	var uuid string
	var err error
	var obj modelObjs.ConfigObj
	var objKey string
	errCode = SRSuccess
	//w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	gActionMgr.logger.Info(fmt.Sprintln("logger print ; Create config resource:", resource))
	fmt.Println("Create Config resource:", resource)
	if objHdl, ok := modelObjs.ConfigObjectMap[resource]; ok {
		fmt.Println("objHdl:", objHdl)
		//if body, obj, err := objects.GetConfigObj(r, objHdl); err == nil {
		if obj, err = objHdl.UnmarshalObject(body); err == nil {
			updateKeys, _ := objects.GetUpdateKeys(body)
			if len(updateKeys) == 0 {
				errCode = SRNoContent
				fmt.Println("nothing to configure")
				gActionMgr.logger.Err("Nothing to configure")
			} else {
				objKey = obj.GetKey()
				uuid, err = gActionMgr.dbHdl.GetUUIDFromObjKey(objKey)
				if err == nil {
					errCode = SRAlreadyConfigured
					fmt.Println("config object is present")
					gActionMgr.logger.Err("Config object is present")
				}
			}
			if errCode != SRSuccess {
				//		w.WriteHeader(http.StatusInternalServerError)
				resp.UUId = uuid
				resp.Error = SRErrString(errCode)
				//	js, _ := json.Marshal(resp)
				//	w.Write(js)
				fmt.Println("errcode not success, return")
				return
			}
			resourceOwner := gActionMgr.objectMgr.ObjHdlMap[resource].Owner
			if resourceOwner.IsConnectedToServer() == false {
				//			errString := "Confd not connected to " + resourceOwner.GetServerName()
				//RespondErrorForApiCall(w, SRSystemNotReady, errString)
				return
			}
			fmt.Println("resource:", resource, " resourceOwner:", resourceOwner, " obj:", obj)
			err, success = resourceOwner.CreateObject(obj, gActionMgr.dbHdl.DBUtil)
			if err == nil && success == true {
				uuid, dbErr := gActionMgr.dbHdl.StoreUUIDToObjKeyMap(objKey)
				if dbErr == nil {
					//gActionMgr.ApiCallStats.NumCreateCallsSuccess++
					/*w.WriteHeader(http.StatusCreated)*/
					resp.UUId = uuid
					errCode = SRSuccess
				} else {
					errCode = SRIdStoreFail
					gActionMgr.logger.Err(fmt.Sprintln("Failed to store UuidToKey map ", obj, dbErr))
				}
			} else {
				resp.Error = err.Error()
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
		for _, v := range value {
			if vbyte, err := json.Marshal(v); err == nil {
				CreateConfig(key, vbyte)
			}
		}
	}
}
func SaveConfigObject(data modelActions.ApplyConfig, resource string) error {
	gActionMgr.logger.Info(fmt.Sprintln("SaveConfigObject for resource:",resource))
	objHdl, ok := modelObjs.ConfigObjectMap[resource]
	if !ok {
		gActionMgr.logger.Err("objHdl nil")
		return errors.New("objHdl Nil")
	}
	_, obj, err := objects.GetConfigObj(nil, objHdl)
	if err != nil {
		gActionMgr.logger.Err(fmt.Sprintln("GetConfigObj return err: ",err))
		return errors.New("getConfigObj return err")
	}
	var configObjects []modelObjs.ConfigObj
	err, objCount, _, _,configObjects := obj.GetBulkObjFromDb(0, 100, gActionMgr.dbHdl.DBUtil)
	if err != nil {
		gActionMgr.logger.Err(fmt.Sprintln("GetBulkObjFromDB returned error:",err))
		return errors.New("GetBulkObjFromDb returned error")
	}
	if objCount == 0 {
		gActionMgr.logger.Info(fmt.Sprintln("No objects of type:",resource, " configured"))
		return nil
	}
	if data.ConfigData[resource] == nil {
		data.ConfigData[resource] = make([]interface{},0)
	}
    for _, configObject := range configObjects {
        data.ConfigData[resource] = append(data.ConfigData[resource],configObject)
	}
	gActionMgr.logger.Info(fmt.Sprintln("data at the end of SaveConfig:",data))
	return nil

}
func OpenFile(cfgFileName string) (fo *os.File, err error) {
		gActionMgr.logger.Info(fmt.Sprintln("Full config file : ", cfgFileName))
		_,err = os.Stat(cfgFileName)
		if os.IsNotExist(err) {
			gActionMgr.logger.Info(fmt.Sprintln(cfgFileName, " not present, create it"))
			fo, err = os.Create(cfgFileName)
			if err != nil {
				gActionMgr.logger.Err(fmt.Sprintln("Error :", err, " when creating file:", cfgFileName))
				return fo,err
			}
		} else if err == nil {
			// open cfg file
			gActionMgr.logger.Info("cfgFile present, open it for update")
			fo, err = os.OpenFile(cfgFileName, os.O_RDWR, 0666)
			if err != nil {
				gActionMgr.logger.Err(fmt.Sprintln("Error:", err, "when opening cfgFile:", cfgFileName))
				return fo,err
			}
		} else {
			gActionMgr.logger.Err(fmt.Sprintln("Error:", err, " when handling the cfgFile:", cfgFileName))
			return fo,err
		}
		return fo,err
}
func ResetConfigObject(data modelActions.ApplyConfig) (err error) {
    gActionMgr.logger.Debug(fmt.Sprintln("Start config reset"))
	
    /* 1) Get all config objects */
     for key, objHandle := range modelObjs.ConfigObjectMap {
	gActionMgr.logger.Debug(fmt.Sprintln("ResetConfig: Got object ", key , " : ", objHandle))
    /* 2) Check if the object is Autoconfig if not
	  delete config */
         _, _, err := objects.GetConfigObj(nil, objHandle) 
	if err != nil {
	gActionMgr.logger.Debug(fmt.Sprintln("Config object doesn't exist ", err))
	}
	}
    
return nil	
}

func ExecutePerformAction(obj modelActions.ActionObj) (err error) {
	gActionMgr.logger.Debug(fmt.Sprintln("local client Execute action obj: ", obj))

	switch obj.(type) {
	case modelActions.ApplyConfig:
		gActionMgr.logger.Info("ApplyConfig")
		data := obj.(modelActions.ApplyConfig)
		for _, applyResource := range ApplyConfigOrder {
			ApplyConfigObject(data, applyResource)
		}
	case modelActions.SaveConfig:
		gActionMgr.logger.Info("SaveConfig")
		var fo *os.File
		var err error
		data := obj.(modelActions.SaveConfig)
		fileName := data.FileName
		gActionMgr.logger.Info(fmt.Sprintln("FileName:", fileName))
		if fileName == "" {
			gActionMgr.logger.Info("FileName not set, setting it to default startup-config")
			fileName = "startup-config"
		}
		// open config file
		cfgFileName := gActionMgr.paramsDir + "/" + fileName + ".json"
		fo,err = OpenFile(cfgFileName)
		if err != nil {
		    gActionMgr.logger.Err(fmt.Sprintln("error with OpenFile, err:",err))
			return err	
		}
		// close fo on exit and check for its returned error
		defer func() {
			if err := fo.Close(); err != nil {
				panic(err)
			}
		}()
		var wdata modelActions.ApplyConfig
		wdata.ConfigData = make(map[string] []interface{})
		for _, applyResource := range ApplyConfigOrder {
			SaveConfigObject(wdata, applyResource)
	        gActionMgr.logger.Info(fmt.Sprintln("data after calling SaveConfig for resource:",applyResource, " is:", wdata))
		}
	    js, err := json.Marshal(wdata)
	    if err != nil {
			gActionMgr.logger.Err(fmt.Sprintln("json marshal returned error:",err))
			return err
		}
		gActionMgr.logger.Info(fmt.Sprintln("js:",string(js)))
		_,err = fo.Write(js) 
		if err != nil {
			gActionMgr.logger.Err(fmt.Sprintln("Error writing:",err))
			return err
		}

	case modelActions.ResetConfig:
		gActionMgr.logger.Info("ResetConfig")
		data := obj.(modelActions.ApplyConfig)
		ResetConfigObject(data)
	}
	return err
}
