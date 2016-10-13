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

package clients

import (
	"errors"
	"models/actions"
	"models/objects"
	"strconv"
	"utils/dbutils"
)

type xponderCfgRecipe struct {
	vlanId        int32
	portList      []string
	untagPortList []string
	adminState    string
}

var xponderGlobal objects.XponderGlobal

const (
	XPONDER_GLOBAL_UPDATE_ATTR_MODE = 1

	XPONDER_CLIENT_PORT_MIN         = 1
	XPONDER_CLIENT_PORT_MAX         = 20
	XPONDER_OVERSUB_CLIENT_PORT_MIN = 9
	XPONDER_OVERSUB_CLIENT_PORT_MAX = 12

	XPONDER_MODE_IN_SVC_WIRE    = "InServiceWire"
	XPONDER_MODE_IN_SVC_REGEN   = "InServiceRegen"
	XPONDER_MODE_IN_SVC_OVERSUB = "InServiceOverSub"
	XPONDER_MODE_IN_SVC_PKT_OPT = "InServicePacketOptical"
	XPONDER_MODE_OUT_OF_SVC     = "OutOfService"
)

func xponderGlobalPreUpdateValidate(dbObj, obj objects.XponderGlobal, attrSet []bool, dbHdl *dbutils.DBUtil) error {
	var err error
	gClientMgr.logger.Debug("Pre config validate called for Xponder Global object")
	return err
}

func xponderGlobalPostUpdateProcessing(dbObj, obj objects.XponderGlobal, attrSet []bool, dbHdl *dbutils.DBUtil) error {
	var err error
	gClientMgr.logger.Info("Post config processing called for Xponder Global object", dbObj, obj, attrSet)
	/* When XponderMode is updated, setup Vlan configuration for the various modes
	(a) In service wire configuration: Clnt ports 1-8 are used and mapped 1:1 to AC400 inputs
	(b) In service oversub configuration: Clnt ports 1-12 are used and mapped as shown below
		Clnt ports 1-8 are 1:1 mapped with AC400 inputs and
		Clnt ports 9, 10 are mapped to AC400_0:3, AC400_0:4
		Clnt ports 11, 12 are mapped to AC400_1:3, AC400_1:4
	(c) In service packet optical, Out of service modes - no default mapping is proviced
	*/
	if attrSet[XPONDER_GLOBAL_UPDATE_ATTR_MODE] {
		switch dbObj.XponderMode {
		case XPONDER_MODE_IN_SVC_WIRE:
			err = xponderModeInSvcWireCfgRemove(dbHdl)
		case XPONDER_MODE_IN_SVC_REGEN:
			err = xponderModeInSvcRegenCfgRemove(dbHdl)
		case XPONDER_MODE_IN_SVC_OVERSUB:
			err = xponderModeInSvcOverSubCfgRemove(dbHdl)
		case XPONDER_MODE_IN_SVC_PKT_OPT:
			err = xponderModeInSvcPktOptCfgRemove(dbHdl)
		case XPONDER_MODE_OUT_OF_SVC:
			err = xponderModeOutOfSvcCfgRemove(dbHdl)
		}
		if err != nil {
			return err
		}
		switch dbObj.XponderMode {
		case XPONDER_MODE_IN_SVC_WIRE, XPONDER_MODE_IN_SVC_REGEN,
			XPONDER_MODE_IN_SVC_OVERSUB, XPONDER_MODE_IN_SVC_PKT_OPT:
			if obj.XponderMode == XPONDER_MODE_OUT_OF_SVC {
				//When going out of service clear all FCAPS data and disable faults/alarms
				err = xponderFCAPSEnable(false)
				if err != nil {
					gClientMgr.logger.Err("Failed to disable FCAPS when transition into Out of Service mode")
				}
			}
		case XPONDER_MODE_OUT_OF_SVC:
			if obj.XponderMode != XPONDER_MODE_OUT_OF_SVC { //Check likely not required
				//When leaving OutOfService mode, enable faults/Alarms
				err = xponderFCAPSEnable(true)
				if err != nil {
					gClientMgr.logger.Err("Failed to enable FCAPS when transition out of Out of Service mode")
				}
			}
		default:
		}
		switch obj.XponderMode {
		case XPONDER_MODE_IN_SVC_WIRE:
			err = xponderModeInSvcWireCfgSet(dbHdl)
		case XPONDER_MODE_IN_SVC_REGEN:
			err = xponderModeInSvcRegenCfgSet(dbHdl)
		case XPONDER_MODE_IN_SVC_OVERSUB:
			err = xponderModeInSvcOverSubCfgSet(dbHdl)
		case XPONDER_MODE_IN_SVC_PKT_OPT:
			err = xponderModeInSvcPktOptCfgSet(dbHdl)
		case XPONDER_MODE_OUT_OF_SVC:
			err = xponderModeOutOfSvcCfgSet(dbHdl)
		}
	}
	return err
}

var dmnListForFCAPS []string = []string{"opticd", "asicd"}

func xponderFCAPSEnable(enable bool) error {
	var err error
	fMgrClntHdl, exist := gClientMgr.Clients["fMgrd"]
	if exist && fMgrClntHdl.IsConnectedToServer() {
		for _, val := range dmnListForFCAPS {
			obj := actions.FaultEnable{
				OwnerName: val,
				EventName: "all",
				Enable:    enable,
			}
			err = fMgrClntHdl.ExecuteAction(obj)
			if err != nil {
				gClientMgr.logger.Err("Failed to change FCAPS state for - " + val)
			}
		}
	}
	return err
}

func xponderGlobalCreate(obj objects.XponderGlobal) (error, bool) {
	var en bool
	opticdClient, exist := gClientMgr.Clients["opticd"]
	if exist && opticdClient.IsConnectedToServer() {
		gClientMgr.logger.Debug("Create received for XponderGlobal")
		xponderGlobal.XponderId = obj.XponderId
		xponderGlobal.XponderMode = obj.XponderMode
		xponderGlobal.XponderDescription = obj.XponderDescription

		switch obj.XponderMode {
		case XPONDER_MODE_OUT_OF_SVC:
			en = false
		default:
			en = true
		}
		err := xponderFCAPSEnable(en)
		if err != nil {
			gClientMgr.logger.Err("Failed to change FCAPS state during xponder auto create")
		}
		return nil, true
	}
	return errors.New("Not supported on this platform"), false
}

func xponderGlobalDelete(obj objects.XponderGlobal) (error, bool) {
	opticdClient, exist := gClientMgr.Clients["opticd"]
	if exist && opticdClient.IsConnectedToServer() {
		return errors.New("Delete operation not supported for XponderGlobal"), false
	}
	return errors.New("Not supported on this platform"), false
}

func xponderGlobalUpdate(obj objects.XponderGlobal) (error, bool) {
	opticdClient, exist := gClientMgr.Clients["opticd"]
	if exist && opticdClient.IsConnectedToServer() {
		gClientMgr.logger.Debug("Update received for XponderGlobal")
		xponderGlobal.XponderMode = obj.XponderMode
		xponderGlobal.XponderDescription = obj.XponderDescription
		return nil, true
	}
	return errors.New("Not supported on this platform"), false
}

func xponderGlobalGet() (error, objects.ConfigObj) {
	opticdClient, exist := gClientMgr.Clients["opticd"]
	if exist && opticdClient.IsConnectedToServer() {
		gClientMgr.logger.Debug("Get received for XponderGlobal")
		return nil, xponderGlobal
	}
	return errors.New("Not supported on this platform"), nil
}

func xponderGlobalGetBulk() (int64, int64, bool, []objects.ConfigObj) {
	opticdClient, exist := gClientMgr.Clients["opticd"]
	if exist && opticdClient.IsConnectedToServer() {
		gClientMgr.logger.Debug("GETBULK xponderGbl : ", xponderGlobal)
		return int64(1), int64(0), false, []objects.ConfigObj{xponderGlobal}
	}
	return int64(0), int64(0), false, []objects.ConfigObj{}
}

var xponderInSvcWireCfgRecipe []xponderCfgRecipe = []xponderCfgRecipe{
	xponderCfgRecipe{
		vlanId:        2,
		untagPortList: []string{"fpPort1", "fpPort13"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        3,
		untagPortList: []string{"fpPort2", "fpPort14"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        4,
		untagPortList: []string{"fpPort3", "fpPort15"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        5,
		untagPortList: []string{"fpPort4", "fpPort16"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        6,
		untagPortList: []string{"fpPort5", "fpPort17"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        7,
		untagPortList: []string{"fpPort6", "fpPort18"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        8,
		untagPortList: []string{"fpPort7", "fpPort19"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        9,
		untagPortList: []string{"fpPort8", "fpPort20"},
		adminState:    "UP",
	},
}

func xponderUpdatePortAdminState(ifName, adminState string, dbHdl *dbutils.DBUtil) error {
	asicdClntHdl := gClientMgr.Clients["asicd"]
	obj := new(objects.Port)
	obj.IntfRef = ifName
	objKey := dbHdl.GetKey(obj)
	dbObj, _ := dbHdl.GetObjectFromDb(obj, objKey)
	*obj = dbObj.(objects.Port)
	//Preserve all attrs and modify admin state alone
	obj.AdminState = adminState
	patchOpInfoSlice := make([]objects.PatchOpInfo, 0)
	err, _ := asicdClntHdl.UpdateObject(dbObj, *obj, []bool{false, false, false,
		false, true, false, false, false, false, false, false, false, false, false,
		false, false, false},
		patchOpInfoSlice, objKey, dbHdl)
	return err
}
func xponderModeInSvcWireCfgSet(dbHdl *dbutils.DBUtil) error {
	var err error
	var adminState string
	asicdClntHdl := gClientMgr.Clients["asicd"]
	for _, val := range xponderInSvcWireCfgRecipe {
		obj := new(objects.Vlan)
		obj.VlanId = val.vlanId
		obj.IntfList = val.portList
		obj.UntagIntfList = val.untagPortList
		obj.AdminState = val.adminState
		err, _ = asicdClntHdl.CreateObject(*obj, dbHdl)
		if err != nil {
			gClientMgr.logger.Err("Failed applying cfg recipe in InSvcWireCfgSet")
			break
		}
	}
	for idx := XPONDER_CLIENT_PORT_MIN; idx <= XPONDER_CLIENT_PORT_MAX; idx++ {
		ifName := "fpPort" + strconv.Itoa(idx)
		//Disable oversub client ports,  enable all other client ports
		if idx >= XPONDER_OVERSUB_CLIENT_PORT_MIN && idx <= XPONDER_OVERSUB_CLIENT_PORT_MAX {
			adminState = "DOWN"
		} else {
			adminState = "UP"
		}
		err = xponderUpdatePortAdminState(ifName, adminState, dbHdl)
		if err != nil {
			gClientMgr.logger.Err("Failed applying cfg recipe in InSvcWireCfgSet")
			break
		}
	}
	return err
}

func xponderModeInSvcWireCfgRemove(dbHdl *dbutils.DBUtil) error {
	var err error
	asicdClntHdl := gClientMgr.Clients["asicd"]
	for _, val := range xponderInSvcWireCfgRecipe {
		obj := new(objects.Vlan)
		obj.VlanId = val.vlanId
		obj.IntfList = val.portList
		obj.UntagIntfList = val.untagPortList
		err, _ = asicdClntHdl.DeleteObject(*obj, "", dbHdl)
		if err != nil {
			gClientMgr.logger.Err("Failed applying cfg recipe in InSvcWireCfgSet")
			break
		}
	}
	//Disable all ports
	for idx := XPONDER_CLIENT_PORT_MIN; idx <= XPONDER_CLIENT_PORT_MAX; idx++ {
		ifName := "fpPort" + strconv.Itoa(idx)
		adminState := "DOWN"
		err = xponderUpdatePortAdminState(ifName, adminState, dbHdl)
		if err != nil {
			gClientMgr.logger.Err("Failed applying cfg recipe in InSvcWireCfgSet")
			break
		}
	}
	return err
}

var xponderInSvcOverSubCfgRecipe []xponderCfgRecipe = []xponderCfgRecipe{
	xponderCfgRecipe{
		vlanId:        2,
		portList:      []string{"fpPort13"},
		untagPortList: []string{"fpPort1"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        3,
		portList:      []string{"fpPort14"},
		untagPortList: []string{"fpPort2"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        4,
		portList:      []string{"fpPort15"},
		untagPortList: []string{"fpPort3"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        5,
		portList:      []string{"fpPort16"},
		untagPortList: []string{"fpPort4"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        6,
		portList:      []string{"fpPort17"},
		untagPortList: []string{"fpPort5"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        7,
		portList:      []string{"fpPort18"},
		untagPortList: []string{"fpPort6"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        8,
		portList:      []string{"fpPort19"},
		untagPortList: []string{"fpPort7"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        9,
		portList:      []string{"fpPort20"},
		untagPortList: []string{"fpPort8"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        10,
		portList:      []string{"fpPort15"},
		untagPortList: []string{"fpPort9"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        11,
		portList:      []string{"fpPort16"},
		untagPortList: []string{"fpPort10"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        12,
		portList:      []string{"fpPort17"},
		untagPortList: []string{"fpPort11"},
		adminState:    "UP",
	},
	xponderCfgRecipe{
		vlanId:        13,
		portList:      []string{"fpPort18"},
		untagPortList: []string{"fpPort12"},
		adminState:    "UP",
	},
}

func xponderModeInSvcOverSubCfgSet(dbHdl *dbutils.DBUtil) error {
	var err error
	asicdClntHdl := gClientMgr.Clients["asicd"]
	for _, val := range xponderInSvcOverSubCfgRecipe {
		obj := new(objects.Vlan)
		obj.VlanId = val.vlanId
		obj.IntfList = val.portList
		obj.UntagIntfList = val.untagPortList
		err, _ = asicdClntHdl.CreateObject(*obj, dbHdl)
		if err != nil {
			gClientMgr.logger.Err("Failed applying cfg recipe in xponderModeInSvcOverSubCfgSet")
			break
		}
	}
	for idx := XPONDER_CLIENT_PORT_MIN; idx <= XPONDER_CLIENT_PORT_MAX; idx++ {
		ifName := "fpPort" + strconv.Itoa(idx)
		adminState := "UP"
		err = xponderUpdatePortAdminState(ifName, adminState, dbHdl)
		if err != nil {
			gClientMgr.logger.Err("Failed applying cfg recipe in InSvcWireCfgSet")
			break
		}
	}
	return err
}

func xponderModeInSvcOverSubCfgRemove(dbHdl *dbutils.DBUtil) error {
	var err error
	asicdClntHdl := gClientMgr.Clients["asicd"]
	for _, val := range xponderInSvcOverSubCfgRecipe {
		obj := new(objects.Vlan)
		obj.VlanId = val.vlanId
		obj.IntfList = val.portList
		obj.UntagIntfList = val.untagPortList
		obj.AdminState = val.adminState
		err, _ = asicdClntHdl.DeleteObject(*obj, "", dbHdl)
		if err != nil {
			gClientMgr.logger.Err("Failed applying cfg recipe in xponderInSvcOverSubCfgRemove")
			break
		}
	}
	//Disable all ports
	for idx := XPONDER_CLIENT_PORT_MIN; idx <= XPONDER_CLIENT_PORT_MAX; idx++ {
		ifName := "fpPort" + strconv.Itoa(idx)
		adminState := "DOWN"
		err = xponderUpdatePortAdminState(ifName, adminState, dbHdl)
		if err != nil {
			gClientMgr.logger.Err("Failed applying cfg recipe in InSvcWireCfgSet")
			break
		}
	}
	return err
}

/* In regen mode, disable all clnt ports and enable N/W loopback on DWDMModule */
func xponderModeInSvcRegenCfgApply(adminState string, nwLb bool, dbHdl *dbutils.DBUtil) error {
	var err error
	asicdClntHdl := gClientMgr.Clients["asicd"]
	for idx := XPONDER_CLIENT_PORT_MIN; idx <= XPONDER_CLIENT_PORT_MAX; idx++ {
		obj := new(objects.Port)
		obj.IntfRef = "fpPort" + strconv.Itoa(idx)
		obj.AdminState = adminState
		objKey := dbHdl.GetKey(obj)
		dbObj, _ := dbHdl.GetObjectFromDb(obj, objKey)
		patchOpInfoSlice := make([]objects.PatchOpInfo, 0)
		err, _ = asicdClntHdl.UpdateObject(dbObj, *obj, []bool{false, false, false,
			false, true, false, false, false, false, false, false, false, false, false,
			false, false, false},
			patchOpInfoSlice, objKey, dbHdl)
		if err != nil {
			gClientMgr.logger.Err("Failed applying recipe in xponderModeInSvcRegenCfgApply")
			break
		}
	}
	opticdClntHdl := gClientMgr.Clients["opticd"]
	for modId := 0; modId < 2; modId++ {
		for clntIntfId := 0; clntIntfId < 4; clntIntfId++ {
			obj := new(objects.DWDMModuleClntIntf)
			obj.ModuleId = uint8(modId)
			obj.ClntIntfId = uint8(clntIntfId)
			obj.EnableIntSerdesNWLoopback = nwLb
			objKey := dbHdl.GetKey(obj)
			dbObj, _ := dbHdl.GetObjectFromDb(obj, objKey)
			patchOpInfoSlice := make([]objects.PatchOpInfo, 0)
			err, _ = opticdClntHdl.UpdateObject(dbObj, *obj, []bool{false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, true, false, false},
				patchOpInfoSlice, objKey, dbHdl)
			if err != nil {
				gClientMgr.logger.Err("Failed applying recipe in xponderModeInSvcRegenCfgApply")
				break
			}
		}
	}
	return err
}
func xponderModeInSvcRegenCfgSet(dbHdl *dbutils.DBUtil) error {
	return xponderModeInSvcRegenCfgApply("DOWN", true, dbHdl)
}
func xponderModeInSvcRegenCfgRemove(dbHdl *dbutils.DBUtil) error {
	return xponderModeInSvcRegenCfgApply("UP", false, dbHdl)
}

func xponderModeInSvcPktOptCfgSet(dbHdl *dbutils.DBUtil) error {
	/*No-Op : User manages all config*/
	return nil
}
func xponderModeInSvcPktOptCfgRemove(dbHdl *dbutils.DBUtil) error {
	/*No-Op : User manages all config*/
	return nil
}

func xponderModeOutOfSvcCfgSet(dbHdl *dbutils.DBUtil) error {
	/*No-Op : User manages all config*/
	return nil
}
func xponderModeOutOfSvcCfgRemove(dbHdl *dbutils.DBUtil) error {
	/*No-Op : User manages all config*/
	return nil
}
