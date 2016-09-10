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
	"fmt"
	"models/actions"
	"models/objects"
	"sync"
	"utils/dbutils"
	"utils/ipcutils"
)

type LocalClient struct {
	ipcutils.IPCClientBase
}

func (clnt *LocalClient) Initialize(name string, address string) {
	clnt.Name = name
	clnt.Address = address
	clnt.ApiHandlerMutex = sync.RWMutex{}
	return
}

func (clnt *LocalClient) ConnectToServer() bool {
	return true
}

func (clnt *LocalClient) DisconnectFromServer() bool {
	return true
}

func (clnt *LocalClient) IsConnectedToServer() bool {
	return true
}

func (clnt *LocalClient) DisableServer() bool {
	return true
}

func (clnt *LocalClient) IsServerEnabled() bool {
	return true
}

func (clnt *LocalClient) GetServerName() string {
	return "local"
}

func (clnt *LocalClient) LockApiHandler() {
	clnt.ApiHandlerMutex.Lock()
}

func (clnt *LocalClient) UnlockApiHandler() {
	clnt.ApiHandlerMutex.Unlock()
}

func (clnt *LocalClient) PreUpdateValidation(dbObj, obj objects.ConfigObj, attrSet []bool, dbHdl *dbutils.DBUtil) error {
	var err error
	switch obj.(type) {
	case objects.XponderGlobal:
		err = xponderGlobalPreUpdateValidate(dbObj.(objects.XponderGlobal), obj.(objects.XponderGlobal), attrSet, dbHdl)
	default:
		break
	}
	return err
}

func (clnt *LocalClient) PostUpdateProcessing(dbObj, obj objects.ConfigObj, attrSet []bool, dbHdl *dbutils.DBUtil) error {
	var err error
	switch obj.(type) {
	case objects.XponderGlobal:
		err = xponderGlobalPostUpdateProcessing(dbObj.(objects.XponderGlobal), obj.(objects.XponderGlobal), attrSet, dbHdl)
	default:
		break
	}
	return err
}

func (clnt *LocalClient) CreateObject(obj objects.ConfigObj, dbHdl *dbutils.DBUtil) (error, bool) {
	var err error
	var ok bool = true
	defer clnt.UnlockApiHandler()
	clnt.LockApiHandler()
	switch obj.(type) {
	case objects.XponderGlobal:
		data := obj.(objects.XponderGlobal)
		err, ok = xponderGlobalCreate(data)
		if ok {
			err = dbHdl.StoreObjectInDb(data)
		}
	default:
		break
	}
	return err, ok
}

func (clnt *LocalClient) DeleteObject(obj objects.ConfigObj, objKey string, dbHdl *dbutils.DBUtil) (error, bool) {
	var err error
	var ok bool = true
	defer clnt.UnlockApiHandler()
	clnt.LockApiHandler()
	switch obj.(type) {
	case objects.XponderGlobal:
		err, ok = xponderGlobalDelete(obj.(objects.XponderGlobal))
	default:
		break
	}
	return err, ok
}

func (clnt *LocalClient) GetBulkObject(obj objects.ConfigObj, dbHdl *dbutils.DBUtil, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []objects.ConfigObj) {
	defer clnt.UnlockApiHandler()
	clnt.LockApiHandler()
	switch obj.(type) {
	case objects.ConfigLogState:
		objCount, nextMarker, more, objs = getApiHistory(dbHdl)
	case objects.XponderGlobal, objects.XponderGlobalState:
		objCount, nextMarker, more, objs = xponderGlobalGetBulk()
	default:
		break
	}
	return nil, objCount, nextMarker, more, objs
}

func (clnt *LocalClient) UpdateObject(dbObj objects.ConfigObj, obj objects.ConfigObj, attrSet []bool, op []objects.PatchOpInfo, objKey string, dbHdl *dbutils.DBUtil) (error, bool) {
	var err error
	var ok bool
	defer clnt.UnlockApiHandler()
	clnt.LockApiHandler()
	switch obj.(type) {
	case objects.XponderGlobal:
		updatedata := obj.(objects.XponderGlobal)
		err, ok = xponderGlobalUpdate(obj.(objects.XponderGlobal))
		if ok == true {
			err = dbHdl.UpdateObjectInDb(updatedata, dbObj, attrSet)
			if err != nil {
				fmt.Println("Update object in DB failed:", err)
				return err, false
			}
		} else {
			return err, false
		}
	default:
		break
	}
	return err, ok
}

func (clnt *LocalClient) GetObject(obj objects.ConfigObj, dbHdl *dbutils.DBUtil) (error, objects.ConfigObj) {
	var retObj objects.ConfigObj
	defer clnt.UnlockApiHandler()
	clnt.LockApiHandler()
	switch obj.(type) {
	case objects.SystemStatusState:
		retObj = gClientMgr.systemStatusCB()
	case objects.SystemSwVersionState:
		retObj = gClientMgr.systemSwVersionCB()
	case objects.XponderGlobal, objects.XponderGlobalState:
		_, retObj = xponderGlobalGet()
	default:
		break
	}
	return nil, retObj
}

func (clnt *LocalClient) ExecuteAction(obj actions.ActionObj) error {
	defer clnt.UnlockApiHandler()
	clnt.LockApiHandler()
	switch obj.(type) {
	case actions.SaveConfig, actions.ApplyConfig, actions.ForceApplyConfig, actions.ResetConfig:
		err := gClientMgr.executeConfigurationActionCB(obj)
		return err
	default:
		break
	}
	return nil
}

/*********************************************************************************************/
// Below methods are for localclient's own use only

func getApiHistory(dbHdl *dbutils.DBUtil) (int64, int64, bool, []objects.ConfigObj) {
	var currMarker int64
	var count int64
	ApiCallObj := objects.ConfigLogState{}
	err, count, next, more, apiCalls := dbHdl.GetBulkObjFromDb(ApiCallObj, currMarker, count)
	if err != nil {
		fmt.Println("Failed to get ApiCalls")
	}
	return count, next, more, apiCalls
}
