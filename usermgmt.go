package main

import (
	"database/sql"
	"models"
	"utils/crypto/bcrypt"
)

func CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {
	case models.UserConfig:
		data := obj.(models.UserConfig)
		// Hashing the password with the default cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Println("Failed to encrypt password for user ", data.UserName)
		}
		// Comparing the password with the hash
		/*err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(data.Password))
		if err != nil {
			fmt.Println("Password didn't match ", err)
		} else {
			fmt.Println("Password matched ")
		}*/
		// Store the encrypted password in DB
		data.Password = string(hashedPassword)
		objId, _ = data.StoreObjectInDb(dbHdl)
		break
	default:
		break
	}
	return objId, true
}

func DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {
	switch obj.(type) {
	case models.UserConfig:
		data := obj.(models.UserConfig)
		data.DeleteObjectFromDb(objKey, dbHdl)
		break
	default:
		break
	}
	return true
}

func GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
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

func UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {
	logger.Println("### Update Object called CONFD", attrSet, objKey)
	ok := false
	switch obj.(type) {
	case models.UserConfig:
		//origdata := dbObj.(models.UserConfig)
		updatedata := obj.(models.UserConfig)
		updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
		break

	default:
		break
	}
	return ok
}
