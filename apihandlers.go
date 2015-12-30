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
	//"path"
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

func CheckIfSystemIsReady(w http.ResponseWriter) bool {
	if gMgr.IsReady() == false {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusServiceUnavailable)
		if err := json.NewEncoder(w).Encode("System is not ready"); err != nil {
			logger.Println("### Failed to encode the system not ready message")
		}
		return false
	}
	return true
}

func ShowConfigObject(w http.ResponseWriter, r *http.Request) {
	if CheckIfSystemIsReady(w) != true {
		return
	}
	logger.Println("####  ShowConfigObject called")
}

func ConfigObjectsBulkGet(w http.ResponseWriter, r *http.Request) {
	if CheckIfSystemIsReady(w) != true {
		return
	}
	resource := strings.TrimPrefix(r.URL.String(), gMgr.apiBase)
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
	if CheckIfSystemIsReady(w) != true {
		return
	}
	resource := strings.TrimPrefix(r.URL.String(), gMgr.apiBase)
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
		}
	}
	return
}

func ConfigObjectDelete(w http.ResponseWriter, r *http.Request) {
	var objKey string
	var objKeySqlStr string
	if CheckIfSystemIsReady(w) != true {
		return
	}
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBase), "/")[0]
	vars := mux.Vars(r)
	err := gMgr.dbHdl.QueryRow("select Key from UuidMap where Uuid = ?", vars["objId"]).Scan(&objKey)
	if err != nil {
		logger.Println("### Failure in getting objKey for Uuid ", resource, vars["objId"], err)
		return
	}
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		obj, _ := GetConfigObj(nil, objHdl)
		objKeySqlStr, err = obj.GetSqlKeyStr(objKey)
		dbObj, _ := obj.GetObjectFromDb(objKeySqlStr, gMgr.dbHdl)
		success := gMgr.objHdlMap[resource].owner.DeleteObject(dbObj, objKeySqlStr, gMgr.dbHdl)
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
	if CheckIfSystemIsReady(w) != true {
		logger.Println("Update: System not ready")
		return
	}
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBase), "/")[0]
	vars := mux.Vars(r)
	err := gMgr.dbHdl.QueryRow("select Key from UuidMap where Uuid = ?", vars["objId"]).Scan(&objKey)
	if err != nil {
		logger.Println("### Failure in getting objKey for Uuid ", resource, vars["objId"], err)
		return
	}
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		obj, _ := GetConfigObj(r, objHdl)
		objKeySqlStr, err = obj.GetSqlKeyStr(objKey)
		dbObj, gerr := obj.GetObjectFromDb(objKeySqlStr, gMgr.dbHdl)
		if gerr == nil {
			diff, err := obj.CompareObjectsAndDiff(dbObj)
			mergedObj, _ := obj.MergeDbAndConfigObj(dbObj, diff)
			success := gMgr.objHdlMap[resource].owner.UpdateObject(dbObj, mergedObj, diff, objKeySqlStr, gMgr.dbHdl)
			if success == true {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				if err = json.NewEncoder(w).Encode(vars["objId"]); err != nil {
					logger.Println("### Failed to encode the UUId for object ", resource, vars["objId"])
				}
			} else {
				logger.Println("UpdateObject FAILED for resource ", resource)
			}
		} else {
			logger.Println("Error getting obj via objKeySqlStr ", objKeySqlStr, gerr)
		}

	} else {
		logger.Println("unable to find resource", resource)
	}
}

//func GetAPIDocs(w http.ResponseWriter, r *http.Request) {
//	logger.Println("### GetAPIDocs is called")
//	//fp := path.Join("./", "api-docs.json")

//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
//	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	w.WriteHeader(http.StatusOK)

//	//http.ServeFile(w, r, fp)
//	return
//}

//func GetObjectAPIDocs(w http.ResponseWriter, r *http.Request) {
//	logger.Println("### GetObjectAPIDocs is called")
//	//fp := path.Join("./", "greetings.json")
//	//http.ServeFile(w, r, fp)
//	return
//}
