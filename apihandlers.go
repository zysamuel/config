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

type ConfigResponse struct {
	UUId    string        `json:"Id"`
}

type GetBulkResponse struct {
	MoreExist     bool               `json:"MoreExist"`
	ObjCount      int64              `json:"ObjCount"`
	CurrentMarker int64              `json:"CurrentMarker"`
	NextMarker    int64              `json:"NextMarker"`
	StateObjects  []models.ConfigObj `json:"StateObjects"`
}

func GetConfigObj(r *http.Request, obj models.ConfigObj) (body []byte, retobj models.ConfigObj, err error) {

	if r != nil {
		body, err = ioutil.ReadAll(io.LimitReader(r.Body, MAX_JSON_LENGTH))
		if err != nil {
			return body, retobj, err
		}
		if err = r.Body.Close(); err != nil {
			return body, retobj, err
		}
	}

	retobj, err = obj.UnmarshalObject(body)
	return body, retobj, err
}

func GetUpdateKeys(body []byte) (map[string]bool, error) {
	var objmap map[string]*json.RawMessage
	var err error
	updateKeys := make(map[string]bool)

	err = json.Unmarshal(body, &objmap)
	if err != nil {
		return updateKeys, err
	}
	for key, _ := range objmap {
		updateKeys[key] = true
	}
	return updateKeys, err
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(peers); err != nil {
		return
	}
}

func CheckIfSystemIsReady(w http.ResponseWriter) bool {
	return gMgr.IsReady()
}

func ShowConfigObject(w http.ResponseWriter, r *http.Request) {
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	logger.Println("####  ShowConfigObject called")
}

func ConfigObjectsBulkGet(w http.ResponseWriter, r *http.Request) {
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	resource := strings.TrimPrefix(r.URL.String(), gMgr.apiBase)
	resource = strings.Split(resource, "?")[0]
	resource = resource[:len(resource)-1]

	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		var resp GetBulkResponse
		var err error
		_, obj, _ := GetConfigObj(nil, objHdl)
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
		return UUId, err
	}
	objKey, err := obj.GetKey()
	if err != nil || len(objKey) == 0 {
		logger.Println("Failed to get objKey after executing ", objKey, err)
		return UUId, err
	}
	dbCmd := fmt.Sprintf(`INSERT INTO UuidMap (Uuid, Key) VALUES ('%v', '%v') ;`, UUId.String(), objKey)
	_, err = dbutils.ExecuteSQLStmt(dbCmd, gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to insert uuid entry in db ", dbCmd, err)
	}
	return UUId, err
}

func ConfigObjectCreate(w http.ResponseWriter, r *http.Request) {
	var resp ConfigResponse
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	resource := strings.TrimPrefix(r.URL.String(), gMgr.apiBase)
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		if _, obj, err := GetConfigObj(r, objHdl); err == nil {
			_, success := gMgr.objHdlMap[resource].owner.CreateObject(obj, gMgr.dbHdl)
			if success == true {
				UUId, err := StoreUuidToKeyMapInDb(obj)
				if err == nil {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusCreated)
					resp.UUId = UUId.String()
					js, err := json.Marshal(resp)
					if err != nil {
						logger.Println("Error in marshalling JSON in config for object ", resource, resp.UUId, err)
						http.Error(w, "Config create successful. Failed to marshal response", http.StatusInternalServerError)
						return
					}
					w.Write(js)
				} else {
					http.Error(w, "Config create failed to store return Id", http.StatusInternalServerError)
					logger.Println("Failed to store UuidToKey map ", obj, err)
				}
			} else {
				http.Error(w, "Config create failed by backend server", http.StatusInternalServerError)
				logger.Println("Failed to create object ", obj)
			}
		} else {
			http.Error(w, "Config create failed to get object handle", http.StatusInternalServerError)
			logger.Println("Failed to get object handle from http request ", objHdl, err)
		}
	} else {
		http.Error(w, "Config create failed to get object map", http.StatusInternalServerError)
		logger.Println("Failed to get ObjectMap ", resource)
	}
	return
}

func ConfigObjectDelete(w http.ResponseWriter, r *http.Request) {
	var resp ConfigResponse
	var objKey string
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBase), "/")[0]
	vars := mux.Vars(r)
	err := gMgr.dbHdl.QueryRow("select Key from UuidMap where Uuid = ?", vars["objId"]).Scan(&objKey)
	if err != nil {
		http.Error(w, "Config delete failed to find entry", http.StatusNotFound)
		logger.Println("Failure in getting objKey for Uuid ", resource, vars["objId"], err)
		return
	}
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		if _, obj, err := GetConfigObj(nil, objHdl); err == nil {
			dbObj, _ := obj.GetObjectFromDb(objKey, gMgr.dbHdl)
			success := gMgr.objHdlMap[resource].owner.DeleteObject(dbObj, objKey, gMgr.dbHdl)
			if success == true {
				dbCmd := "delete from " + "UuidMap" + " where Uuid = " + "\"" + vars["objId"] + "\""
				_, err := dbutils.ExecuteSQLStmt(dbCmd, gMgr.dbHdl)
				if err != nil {
					logger.Println("Failure in deleting Uuid map entry for ", vars["objId"], err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				resp.UUId = vars["objId"]
				js, err := json.Marshal(resp)
				if err != nil {
					logger.Println("Error in marshalling JSON in update for object ", resource, resp.UUId, err)
					http.Error(w, "Config delete successful. Failed to marshal response", http.StatusInternalServerError)
					return
				}
				w.Write(js)
			} else {
				http.Error(w, "Config delete failed by backend server ", http.StatusInternalServerError)
				logger.Println("DeleteObject returned failure ", obj)
			}
		} else {
			http.Error(w, "Config delete failed to get object handle", http.StatusInternalServerError)
			logger.Println("Failed to get object handle from http request ", objHdl, err)
		}
	} else {
		http.Error(w, "Config delete failed to get object map", http.StatusInternalServerError)
		logger.Println("Failed to get ObjectMap ", resource)
	}
	return
}

func ConfigObjectUpdate(w http.ResponseWriter, r *http.Request) {
	var resp ConfigResponse
	var objKey string
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBase), "/")[0]
	vars := mux.Vars(r)
	err := gMgr.dbHdl.QueryRow("select Key from UuidMap where Uuid = ?", vars["objId"]).Scan(&objKey)
	if err != nil {
		http.Error(w, "Config update failed to find entry", http.StatusNotFound)
		logger.Println("Failure in getting objKey for Uuid ", resource, vars["objId"], err)
		return
	}
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		body, obj, _ := GetConfigObj(r, objHdl)
		updateKeys, _ := GetUpdateKeys(body)
		dbObj, gerr := obj.GetObjectFromDb(objKey, gMgr.dbHdl)
		if gerr == nil {
			diff, _ := obj.CompareObjectsAndDiff(updateKeys, dbObj)
			mergedObj, _ := obj.MergeDbAndConfigObj(dbObj, diff)
			success := gMgr.objHdlMap[resource].owner.UpdateObject(dbObj, mergedObj, diff, objKey, gMgr.dbHdl)
			if success == true {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				resp.UUId = vars["objId"]
				js, err := json.Marshal(resp)
				if err != nil {
					logger.Println("Error in marshalling JSON in update for object ", resource, resp.UUId, err)
					http.Error(w, "Config update successful. Failed to marshal response", http.StatusInternalServerError)
					return
				}
				w.Write(js)
			} else {
				http.Error(w, "Config update failed by backend server", http.StatusInternalServerError)
				logger.Println("UpdateObject failed for resource ", updateKeys, resource)
			}
		} else {
			http.Error(w, "Config update failed in finding object from internal key", http.StatusInternalServerError)
			logger.Println("Config update failed in getting obj via objKeySqlStr ", objKey, gerr)
		}

	} else {
		http.Error(w, "Config update failed to get object map", http.StatusInternalServerError)
		logger.Println("Config update failed t get ObjectMap ", resource)
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
