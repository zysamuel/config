package main

import (
	"database/sql"
	"models"
	"utils/crypto/bcrypt"
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

func (clnt *LocalClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (error, bool) {
	var err error
	switch obj.(type) {
	case models.User:
		data := obj.(models.User)
		// Hashing the password with the default cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Println("Failed to encrypt password for user ", data.UserName)
		}
		// Create user in configmgr's users table
		if ok := gMgr.CreateUser(data.UserName); ok {
			// Store the encrypted password in DB
			data.Password = string(hashedPassword)
			_, err = data.StoreObjectInDb(dbHdl)
		}
		break

	case models.IPV4AddressBlock:
		GetIpBlockMgr().CreateObject(obj, dbHdl)

	default:
		break
	}
	return err, true
}

func (clnt *LocalClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) (error, bool) {
	switch obj.(type) {
	case models.User:
		data := obj.(models.User)
		// Delete user from configmgr's users table
		if ok := gMgr.DeleteUser(data.UserName); ok {
			data.DeleteObjectFromDb(objKey, dbHdl)
		}
		break

	case models.IPV4AddressBlock:
		GetIpBlockMgr().DeleteObject(obj, objKey, dbHdl)

	default:
		break
	}
	return nil, true
}

func (clnt *LocalClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {
	logger.Println("### Get Bulk request called with", currMarker, count)
	switch obj.(type) {
	case models.UserState:
		break
	default:
		break
	}
	return nil, objCount, nextMarker, more, objs
}

func (clnt *LocalClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) (error, bool) {
	logger.Println("### Update Object called CONFD", attrSet, objKey)
	ok := false
	switch obj.(type) {
	case models.User:
		//origdata := dbObj.(models.User)
		updatedata := obj.(models.User)
		updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
		break

	case models.IPV4AddressBlock:
		GetIpBlockMgr().UpdateObject(dbObj, obj, attrSet, objKey, dbHdl)

	default:
		break
	}
	return nil, ok
}

func (clnt *LocalClient) GetObject(obj models.ConfigObj) (error, models.ConfigObj) {
	var retObj models.ConfigObj
	switch obj.(type) {
	case models.UserState:
		logger.Println("### Get request called for UserState")
		break
	case models.SystemStatus:
		logger.Println("### Get request called for SystemStatus")
		retObj = gMgr.GetSystemStatus()
		return nil, retObj
	default:
		break
	}
	return nil, nil
}
