package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"models"
	"net/http"
	"strings"
	//"net/url"
	"github.com/nu7hatch/gouuid"
	"strconv"
	"utils/dbutils"
)

const (
	MAX_OBJECTS_IN_GETBULK = 30
	MAX_JSON_LENGTH        = 4096
)

type GetBulkResponse struct {
	MoreExist     bool               `json:"MoreExist"`
	ObjCount      int64              `json:"ObjCount"`
	CurrentMarker int64              `json:"CurrentMarker"`
	NextMarker    int64              `json:"NextMarker"`
	StateObjects  []models.ConfigObj `json:"StateObjects"`
}

func GetConfigObj(r *http.Request, obj models.ConfigObj) (models.ConfigObj, error) {
	var retObj models.ConfigObj
	var err error
	var body []byte
	if r != nil {
		body, err = ioutil.ReadAll(io.LimitReader(r.Body, MAX_JSON_LENGTH))
		if err != nil {
			return retObj, err
		}
		if err = r.Body.Close(); err != nil {
			return retObj, err
		}
	}
	return obj.UnmarshalObject(body)
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(peers); err != nil {
		return
	}
}

func ShowConfigObject(w http.ResponseWriter, r *http.Request) {
	logger.Println("####  ShowConfigObject called")
}

func ConfigObjectsBulkGet(w http.ResponseWriter, r *http.Request) {
	resource := strings.TrimPrefix(r.URL.String(), "/")
	resource = strings.Split(resource, "?")[0]
	resource = resource[:len(resource)-1]

	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		var resp GetBulkResponse
		var err error
		obj, _ := GetConfigObj(nil, objHdl)
		currentIndex, objCount := ExtractGetBulkParams(r)
		resp.CurrentMarker = currentIndex
		err, resp.ObjCount, resp.NextMarker, resp.MoreExist,
			resp.StateObjects = gMgr.objHdlMap[resource].owner.GetBulkObject(obj,
			currentIndex,
			objCount)
		js, err := json.Marshal(resp)
		if err != nil {
			logger.Println("### Error in marshalling JSON in getBulk for object ", resource, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}
	return
}

func ExtractGetBulkParams(r *http.Request) (currentIndex int64, objectCount int64) {
	valueMap := r.URL.Query()
	if currentIndexStr, ok := valueMap["CurrentMarker"]; ok {
		currentIndex, _ = strconv.ParseInt(currentIndexStr[0], 10, 64)
	} else {
		currentIndex = 0
	}

	if objectCountStr, ok := valueMap["Count"]; ok {
		objectCount, _ = strconv.ParseInt(objectCountStr[0], 10, 64)
	} else {
		objectCount = MAX_OBJECTS_IN_GETBULK
	}
	return currentIndex, objectCount
}

func StoreUuidToKeyMapInDb(obj models.ConfigObj) (*uuid.UUID, error) {
	UUId, err := uuid.NewV4()
	if err != nil {
		logger.Println("Failed to get UUID ", UUId, err)
	}
	objKey, err := obj.GetKey()
	if err != nil || len(objKey) == 0 {
		logger.Println("Failed to get objKey after executing ", objKey, err)
	}
	dbCmd := fmt.Sprintf(`INSERT INTO UuidMap (Uuid, Key) VALUES ('%v', '%v') ;`, UUId.String(), objKey)
	_, err = dbutils.ExecuteSQLStmt(dbCmd, gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to insert uuid entry in db ", dbCmd, err)
	}
	return UUId, err
}

func ConfigObjectCreate(w http.ResponseWriter, r *http.Request) {
	resource := strings.TrimPrefix(r.URL.String(), "/")
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		obj, _ := GetConfigObj(r, objHdl)
		_, success := gMgr.objHdlMap[resource].owner.CreateObject(obj, gMgr.dbHdl)
		if success == true {
			UUId, err := StoreUuidToKeyMapInDb(obj)
			if err != nil {
				logger.Println("### Failed to store UuidMap ", err)
			}
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
			if err = json.NewEncoder(w).Encode(UUId.String()); err != nil {
				logger.Println("### Failed to encode the UUId for object ", resource, UUId.String())
			}
		} else {
			logger.Println("### Failed to CreateObject")
		}
	}
	return
}

func ConfigObjectDelete(w http.ResponseWriter, r *http.Request) {
	var objKey string
	var objKeySqlStr string
	resource := strings.Split(r.URL.String(), "/")[1]
	vars := mux.Vars(r)
	err := gMgr.dbHdl.QueryRow("select Key from UuidMap where Uuid = ?", vars["objId"]).Scan(&objKey)
	if err != nil {
		logger.Println("### Failure in getting objKey for Uuid ", resource, vars["objId"], err)
		return
	}
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		obj, _ := GetConfigObj(nil, objHdl)
		objKeySqlStr, err = obj.GetSqlKeyStr(objKey)
		success := gMgr.objHdlMap[resource].owner.DeleteObject(obj, objKeySqlStr, gMgr.dbHdl)
		if success == true {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

			dbCmd := "delete from " + "UuidMap" + " where Uuid = " + "\"" + vars["objId"] + "\""
			_, err := dbutils.ExecuteSQLStmt(dbCmd, gMgr.dbHdl)
			if err != nil {
				logger.Println("### Failure in deleting Uuid map entry for ", vars["objId"], err)
			}
		}
	}
	return
}

func ConfigObjectUpdate(w http.ResponseWriter, r *http.Request) {
	var objKey string
	var objKeySqlStr string
	resource := strings.Split(r.URL.String(), "/")[1]
	vars := mux.Vars(r)
	err := gMgr.dbHdl.QueryRow("select Key from UuidMap where Uuid = ?", vars["objId"]).Scan(&objKey)
	if err != nil {
		logger.Println("### Failure in getting objKey for Uuid ", resource, vars["objId"], err)
		return
	}
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		obj, _ := GetConfigObj(nil, objHdl)
		objKeySqlStr, err = obj.GetSqlKeyStr(objKey)
		logger.Println("ConfigObjectUpdate", objKeySqlStr, err)
		success := gMgr.objHdlMap[resource].owner.UpdateObject(obj, objKeySqlStr, gMgr.dbHdl)
		if success == true {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			if err = json.NewEncoder(w).Encode(vars["objId"]); err != nil {
				logger.Println("### Failed to encode the UUId for object ", resource, vars["objId"])
			}
		}
	}
	return
}
