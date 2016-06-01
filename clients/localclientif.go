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
	"models"
	"utils/dbutils"
)

type LocalClient struct {
}

func (clnt *LocalClient) Initialize(name string, address string) {
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

func (clnt *LocalClient) GetServerName() string {
	return "local"
}

func (clnt *LocalClient) CreateObject(obj models.ConfigObj, dbHdl *dbutils.DBUtil) (error, bool) {
	var err error
	return err, true
}

func (clnt *LocalClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *dbutils.DBUtil) (error, bool) {
	return nil, true
}

func (clnt *LocalClient) GetBulkObject(obj models.ConfigObj, dbHdl *dbutils.DBUtil, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {
	return nil, objCount, nextMarker, more, objs
}

func (clnt *LocalClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, op []models.PatchOpInfo, objKey string, dbHdl *dbutils.DBUtil) (error, bool) {
	return nil, true
}

func (clnt *LocalClient) GetObject(obj models.ConfigObj, dbHdl *dbutils.DBUtil) (error, models.ConfigObj) {
	var retObj models.ConfigObj
	switch obj.(type) {
	case models.SystemStatusState:
		retObj = gClientMgr.systemStatusCB()
		return nil, retObj
	case models.SystemSwVersionState:
		retObj = gClientMgr.systemSwVersionCB()
		return nil, retObj
	default:
		break
	}
	return nil, nil
}

func (clnt *LocalClient) ExecuteAction(obj models.ConfigObj) error {
	return nil
}
