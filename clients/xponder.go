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
	"fmt"
	"models/objects"
	"strconv"
	"utils/dbutils"
)

type xponderCfgRecipe struct {
	vlanId   int32
	portList []string
}

var xponderGlobal objects.XponderGlobal

const (
	XPONDER_GLOBAL_UPDATE_ATTR_MODE = 1

	XPONDER_MODE_IN_SVC_WIRE    = "InServiceWire"
	XPONDER_MODE_IN_SVC_REGEN   = "InServiceRegen"
	XPONDER_MODE_IN_SVC_OVERSUB = "InServiceOverSub"
	XPONDER_MODE_IN_SVC_PKT_OPT = "InServicePacketOptical"
	XPONDER_MODE_OUT_OF_SVC     = "OutOfService"
)

func xponderGlobalPreUpdateValidate(dbObj, obj objects.XponderGlobal, attrSet []bool, dbHdl *dbutils.DBUtil) error {
	var err error
	fmt.Println("Pre config validate called for Xponder Global object")
	return err
}

func xponderGlobalPostUpdateProcessing(dbObj, obj objects.XponderGlobal, attrSet []bool, dbHdl *dbutils.DBUtil) error {
	var err error
	fmt.Println("Post config processing called for Xponder Global object", dbObj, obj, attrSet)
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

func xponderGlobalCreate(obj objects.XponderGlobal) (error, bool) {
	fmt.Println("Create received for XponderGlobal")
	xponderGlobal.XponderId = obj.XponderId
	xponderGlobal.XponderMode = obj.XponderMode
	xponderGlobal.XponderDescription = obj.XponderDescription
	return nil, true
}

func xponderGlobalDelete(obj objects.XponderGlobal) (error, bool) {
	return errors.New("Delete operation not supported for XponderGlobal"), false
}

func xponderGlobalUpdate(obj objects.XponderGlobal) (error, bool) {
	fmt.Println("Update received for XponderGlobal")
	xponderGlobal.XponderMode = obj.XponderMode
	xponderGlobal.XponderDescription = obj.XponderDescription
	return nil, true
}

func xponderGlobalGet(obj objects.XponderGlobal) (error, objects.ConfigObj) {
	fmt.Println("Get received for XponderGlobal")
	return nil, xponderGlobal
}

func xponderGlobalGetBulk(obj objects.XponderGlobal) (int64, int64, bool, []objects.ConfigObj) {
	fmt.Println("GETBULK xponderGbl : ", xponderGlobal)
	return int64(1), int64(0), false, []objects.ConfigObj{xponderGlobal}
}

var xponderInSvcWireCfgRecipe []xponderCfgRecipe = []xponderCfgRecipe{
	xponderCfgRecipe{
		vlanId:   2,
		portList: []string{"fpPort1", "fpPort13"},
	},
	xponderCfgRecipe{
		vlanId:   3,
		portList: []string{"fpPort2", "fpPort14"},
	},
	xponderCfgRecipe{
		vlanId:   4,
		portList: []string{"fpPort3", "fpPort15"},
	},
	xponderCfgRecipe{
		vlanId:   5,
		portList: []string{"fpPort4", "fpPort16"},
	},
	xponderCfgRecipe{
		vlanId:   6,
		portList: []string{"fpPort5", "fpPort17"},
	},
	xponderCfgRecipe{
		vlanId:   7,
		portList: []string{"fpPort6", "fpPort18"},
	},
	xponderCfgRecipe{
		vlanId:   8,
		portList: []string{"fpPort7", "fpPort19"},
	},
	xponderCfgRecipe{
		vlanId:   9,
		portList: []string{"fpPort8", "fpPort20"},
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
		false, true, false, false, false, false, false, false, false},
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
		obj.UntagIntfList = val.portList
		err, _ = asicdClntHdl.CreateObject(*obj, dbHdl)
		if err != nil {
			fmt.Println("Failed applying cfg recipe in InSvcWireCfgSet")
			break
		}
	}
	for idx := 1; idx <= 20; idx++ {
		ifName := "fpPort" + strconv.Itoa(idx)
		//Disable client ports fpPort9-fpPort12, enable all other client ports
		if idx > 8 && idx < 13 {
			adminState = "DOWN"
		} else {
			adminState = "UP"
		}
		err = xponderUpdatePortAdminState(ifName, adminState, dbHdl)
		if err != nil {
			fmt.Println("Failed applying cfg recipe in InSvcWireCfgSet")
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
		obj.UntagIntfList = val.portList
		err, _ = asicdClntHdl.DeleteObject(*obj, "", dbHdl)
		if err != nil {
			fmt.Println("Failed applying cfg recipe in InSvcWireCfgSet")
			break
		}
	}
	//Enable all client ports
	for idx := 1; idx <= 20; idx++ {
		ifName := "fpPort" + strconv.Itoa(idx)
		adminState := "UP"
		err = xponderUpdatePortAdminState(ifName, adminState, dbHdl)
		if err != nil {
			fmt.Println("Failed applying cfg recipe in InSvcWireCfgSet")
			break
		}
	}
	return err
}

var xponderInSvcOverSubCfgRecipe []xponderCfgRecipe = append(xponderInSvcWireCfgRecipe,
	[]xponderCfgRecipe{
		xponderCfgRecipe{
			vlanId:   9,
			portList: []string{"fpPort9", "fpPort15"},
		},
		xponderCfgRecipe{
			vlanId:   10,
			portList: []string{"fpPort10", "fpPort16"},
		},
		xponderCfgRecipe{
			vlanId:   11,
			portList: []string{"fpPort11", "fpPort17"},
		},
		xponderCfgRecipe{
			vlanId:   12,
			portList: []string{"fpPort12", "fpPort18"},
		},
	}...)

func xponderModeInSvcOverSubCfgSet(dbHdl *dbutils.DBUtil) error {
	var err error
	asicdClntHdl := gClientMgr.Clients["asicd"]
	for _, val := range xponderInSvcWireCfgRecipe {
		obj := new(objects.Vlan)
		obj.VlanId = val.vlanId
		obj.UntagIntfList = val.portList
		err, _ = asicdClntHdl.CreateObject(*obj, dbHdl)
		if err != nil {
			fmt.Println("Failed applying cfg recipe in xponderModeInSvcOverSubCfgSet")
			break
		}
	}
	return err
}
func xponderModeInSvcOverSubCfgRemove(dbHdl *dbutils.DBUtil) error {
	var err error
	asicdClntHdl := gClientMgr.Clients["asicd"]
	for _, val := range xponderInSvcWireCfgRecipe {
		obj := new(objects.Vlan)
		obj.VlanId = val.vlanId
		obj.UntagIntfList = val.portList
		err, _ = asicdClntHdl.DeleteObject(*obj, "", dbHdl)
		if err != nil {
			fmt.Println("Failed applying cfg recipe in xponderInSvcOverSubCfgRemove")
			break
		}
	}
	return err
}

/* In regen mode, disable all clnt ports and enable N/W loopback on DWDMModule */
func xponderModeInSvcRegenCfgApply(adminState string, nwLb bool, dbHdl *dbutils.DBUtil) error {
	var err error
	asicdClntHdl := gClientMgr.Clients["asicd"]
	for idx := 1; idx <= 20; idx++ {
		obj := new(objects.Port)
		obj.IntfRef = "fpPort" + strconv.Itoa(idx)
		obj.AdminState = adminState
		objKey := dbHdl.GetKey(obj)
		dbObj, _ := dbHdl.GetObjectFromDb(obj, objKey)
		patchOpInfoSlice := make([]objects.PatchOpInfo, 0)
		err, _ = asicdClntHdl.UpdateObject(dbObj, *obj, []bool{false, false, false,
			false, true, false, false, false, false, false, false, false},
			patchOpInfoSlice, objKey, dbHdl)
		if err != nil {
			fmt.Println("Failed applying recipe in xponderModeInSvcRegenCfgApply")
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
				fmt.Println("Failed applying recipe in xponderModeInSvcRegenCfgApply")
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
