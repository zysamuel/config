Copyright [2016] [SnapRoute Inc]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

	 Unless required by applicable law or agreed to in writing, software
	 distributed under the License is distributed on an "AS IS" BASIS,
	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	 See the License for the specific language governing permissions and
	 limitations under the License.
package clients

import (
	"models"
	//"utils/crypto/bcrypt"
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

func (clnt *LocalClient) IsConnectedToServer() bool {
	return true
}

func (clnt *LocalClient) GetServerName() string {
	return "local"
}

func (clnt *LocalClient) CreateObject(obj models.ConfigObj, dbHdl *dbutils.DBUtil) (error, bool) {
	var err error
	switch obj.(type) {
	case models.User:
		//data := obj.(models.User)
		// Hashing the password with the default cost of 10
		//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		// Create user in configmgr's users table
		//if ok := gMgr.CreateUser(data.UserName); ok {
		// Store the encrypted password in DB
		//data.Password = string(hashedPassword)
		//err = data.StoreObjectInDb(dbHdl)
		//}
		break
	default:
		break
	}
	return err, true
}

func (clnt *LocalClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *dbutils.DBUtil) (error, bool) {
	switch obj.(type) {
	case models.User:
		//data := obj.(models.User)
		// Delete user from configmgr's users table
		//if ok := gMgr.DeleteUser(data.UserName); ok {
		//	data.DeleteObjectFromDb(dbHdl)
		//}
		break
	default:
		break
	}
	return nil, true
}

func (clnt *LocalClient) GetBulkObject(obj models.ConfigObj, dbHdl *dbutils.DBUtil, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {
	switch obj.(type) {
	case models.UserState:
		break
	default:
		break
	}
	return nil, objCount, nextMarker, more, objs
}

func (clnt *LocalClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *dbutils.DBUtil) (error, bool) {
	ok := false
	switch obj.(type) {
	case models.User:
		//origdata := dbObj.(models.User)
		//updatedata := obj.(models.User)
		//updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
		break
	default:
		break
	}
	return nil, ok
}

func (clnt *LocalClient) GetObject(obj models.ConfigObj, dbHdl *dbutils.DBUtil) (error, models.ConfigObj) {
	var retObj models.ConfigObj
	switch obj.(type) {
	case models.UserState:
		break
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
