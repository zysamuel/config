package main

import (
	"database/sql"
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
	UUId string `json:"ObjectId"`
}

type ReturnObject struct {
	ObjectId string `json:"ObjectId"`
	models.ConfigObj
}

type GetBulkResponse struct {
	MoreExist     bool  `json:"MoreExist"`
	ObjCount      int64 `json:"ObjCount"`
	CurrentMarker int64 `json:"CurrentMarker"`
	NextMarker    int64 `json:"NextMarker"`
	StateObjects  []ReturnObject
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
	//if err := json.NewEncoder(w).Encode(peers); err != nil {
	//	return
	//}
}

func CheckIfSystemIsReady(w http.ResponseWriter) bool {
	return gMgr.IsReady()
}

func GetOneObjectForId(w http.ResponseWriter, r *http.Request) {
	var errCode int
	var obj models.ConfigObj
	var objKey string
	var retObj ReturnObject
	var err error
	var success bool
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBase), "/")[0]
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		if _, obj, err = GetConfigObj(r, objHdl); err != nil {
			http.Error(w, SRErrString(SRNotFound), http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, SRErrString(SRNotFound), http.StatusNotFound)
		return
	}
	vars := mux.Vars(r)
	uuid := vars["objId"]
	//if objId is provided then read objkey from DB
	err = gMgr.dbHdl.QueryRow("select Key from UuidMap where Uuid = ?", uuid).Scan(&objKey)
	if err != nil {
		http.Error(w, SRErrString(SRNotFound), http.StatusNotFound)
		return
	}
	if strings.Contains(gMgr.objHdlMap[resource].access, "r") {
		if retObj.ConfigObj, success = gMgr.objHdlMap[resource].owner.GetObject(obj); success == false {
			errCode = SRServerError
		}
	} else {
		if retObj.ConfigObj, err = obj.GetObjectFromDb(objKey, gMgr.dbHdl); err != nil {
			errCode = SRServerError
		}
	}
	gMgr.dbHdl.QueryRow("select Uuid from UuidMap where Key = ?", objKey).Scan(&retObj.ObjectId)
	js, err := json.Marshal(retObj)
	if err != nil {
		errCode = SRRespMarshalErr
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
		errCode = SRSuccess
	}
	if errCode != SRSuccess {
		http.Error(w, SRErrString(errCode), http.StatusInternalServerError)
	}
	return
}

func GetOneObject(w http.ResponseWriter, r *http.Request) {
	var errCode int
	var obj models.ConfigObj
	var objKey string
	var retObj ReturnObject
	var err error
	var success bool
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBase), "/")[0]
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		if _, obj, err = GetConfigObj(r, objHdl); err != nil {
			http.Error(w, SRErrString(SRNotFound), http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, SRErrString(SRNotFound), http.StatusNotFound)
		return
	}
	//Get key fields provided in the request.
	objKey, err = obj.GetKey()
	if err != nil {
		http.Error(w, SRErrString(SRNotFound), http.StatusNotFound)
		return
	}
	if strings.Contains(gMgr.objHdlMap[resource].access, "r") {
		if retObj.ConfigObj, success = gMgr.objHdlMap[resource].owner.GetObject(obj); success == false {
			errCode = SRServerError
		}
	} else {
		if retObj.ConfigObj, err = obj.GetObjectFromDb(objKey, gMgr.dbHdl); err != nil {
			errCode = SRServerError
		}
	}
	gMgr.dbHdl.QueryRow("select Uuid from UuidMap where Key = ?", objKey).Scan(&retObj.ObjectId)
	js, err := json.Marshal(retObj)
	if err != nil {
		errCode = SRRespMarshalErr
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
		errCode = SRSuccess
	}
	if errCode != SRSuccess {
		http.Error(w, SRErrString(errCode), http.StatusInternalServerError)
	}
	return
}

func BulkGetObjects(w http.ResponseWriter, r *http.Request) {
	resource := strings.TrimPrefix(r.URL.String(), gMgr.apiBase)
	resource = strings.Split(resource, "?")[0]
	resource = resource[:len(resource)-1]
	if strings.Contains(gMgr.objHdlMap[resource].access, "r") {
		StateObjectsBulkGet(resource, w, r)
	} else {
		ConfigObjectsBulkGet(resource, w, r)
	}
}

func ConfigObjectsBulkGet(resource string, w http.ResponseWriter, r *http.Request) {
	var retObjs []ReturnObject
	var errCode int
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		if _, obj, err := GetConfigObj(r, objHdl); err == nil {
			objList, err := obj.GetAllObjFromDb(gMgr.dbHdl)
			if err == nil {
				retObjs = make([]ReturnObject, len(objList))
				for idx, object := range objList {
					retObjs[idx].ConfigObj = object
					objKey, _ := object.GetKey()
					gMgr.dbHdl.QueryRow("select Uuid from UuidMap where Key = ?", objKey).Scan(&retObjs[idx].ObjectId)
				}
				js, err := json.Marshal(retObjs)
				if err != nil {
					errCode = SRRespMarshalErr
				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusOK)
					w.Write(js)
					errCode = SRSuccess
				}
			}
		}
	}
	if errCode != SRSuccess {
		http.Error(w, SRErrString(errCode), http.StatusInternalServerError)
	}
	return
}

func StateObjectsBulkGet(resource string, w http.ResponseWriter, r *http.Request) {
	var errCode int
	var objKey string
	var stateObjects []models.ConfigObj
	var resp GetBulkResponse
	var err error
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		_, obj, _ := GetConfigObj(nil, objHdl)
		currentIndex, objCount := ExtractGetBulkParams(r)
		if objCount > MAX_OBJECTS_IN_GETBULK {
			http.Error(w, SRErrString(SRBulkGetTooLarge), http.StatusRequestEntityTooLarge)
			logger.Println("Too many objects requested in bulkget ", objCount)
			return
		}
		resp.CurrentMarker = currentIndex
		err, resp.ObjCount, resp.NextMarker, resp.MoreExist,
			stateObjects = gMgr.objHdlMap[resource].owner.GetBulkObject(obj,
			currentIndex,
			objCount)
		if err == nil {
			resp.StateObjects = make([]ReturnObject, resp.ObjCount)
			for idx, stateObject := range stateObjects {
				resp.StateObjects[idx].ConfigObj = stateObject
				objKey, _ = stateObject.GetKey()
				gMgr.dbHdl.QueryRow("select Uuid from UuidMap where Key = ?", objKey).Scan(&resp.StateObjects[idx].ObjectId)
			}
			js, err := json.Marshal(resp)
			if err != nil {
				errCode = SRRespMarshalErr
				logger.Println("### Error in marshalling JSON in getBulk for object ", resource, err)
			} else {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				w.Write(js)
				errCode = SRSuccess
			}
		}
	}
	if errCode != SRSuccess {
		http.Error(w, SRErrString(errCode), http.StatusInternalServerError)
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
	var errCode int
	var success bool
	var uuid string
	logger.Println("Create Object Called")
	resource := strings.TrimPrefix(r.URL.String(), gMgr.apiBase)
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		if body, obj, err := GetConfigObj(r, objHdl); err == nil {
			objKey, _ := obj.GetKey()
			_, err := obj.GetObjectFromDb(objKey, gMgr.dbHdl)
			if err == nil {
				logger.Println("Config object is present")
				gMgr.dbHdl.QueryRow("select Uuid from UuidMap where Key = ?", objKey).Scan(&uuid)
				http.Error(w, SRErrString(SRAlreadyConfigured)+" Existing ObjectId: "+uuid, http.StatusInternalServerError)
				return
			} else if err == sql.ErrNoRows {
				updateKeys, _ := GetUpdateKeys(body)
				if len(updateKeys) == 0 {
					errCode = SRNoContent
					logger.Println("Nothing to configure")
				} else {
					_, success = gMgr.objHdlMap[resource].owner.CreateObject(obj, gMgr.dbHdl)
					if success == true {
						UUId, err := StoreUuidToKeyMapInDb(obj)
						if err == nil {
							w.Header().Set("Content-Type", "application/json; charset=UTF-8")
							w.WriteHeader(http.StatusCreated)
							resp.UUId = UUId.String()
							js, err := json.Marshal(resp)
							if err != nil {
								errCode = SRRespMarshalErr
							} else {
								w.Write(js)
								errCode = SRSuccess
							}
						} else {
							errCode = SRIdStoreFail
							logger.Println("Failed to store UuidToKey map ", obj, err)
						}
					} else {
						errCode = SRServerError
						logger.Println("Failed to create object ", obj)
					}
				}
			}
		} else {
			errCode = SRObjHdlError
			logger.Println("Failed to get object handle from http request ", objHdl, resource, err)
		}
	} else {
		errCode = SRObjMapError
		logger.Println("Failed to get ObjectMap ", resource)
	}
	if errCode != SRSuccess {
		logger.Println("returning error as http request failed")
		http.Error(w, SRErrString(errCode), http.StatusInternalServerError)
	} else {
		logger.Println("Create Object Request Processed successfully")
	}
	return
}

func ConfigObjectDeleteForId(w http.ResponseWriter, r *http.Request) {
	var resp ConfigResponse
	var errCode int
	var objKey string
	var success bool
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBase), "/")[0]
	vars := mux.Vars(r)
	err := gMgr.dbHdl.QueryRow("select Key from UuidMap where Uuid = ?", vars["objId"]).Scan(&objKey)
	if err != nil {
		http.Error(w, SRErrString(SRNotFound), http.StatusNotFound)
		return
	}
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		if _, obj, err := GetConfigObj(nil, objHdl); err == nil {
			dbObj, _ := obj.GetObjectFromDb(objKey, gMgr.dbHdl)
			success = gMgr.objHdlMap[resource].owner.DeleteObject(dbObj, objKey, gMgr.dbHdl)
			if success == true {
				dbCmd := "delete from " + "UuidMap" + " where Uuid = " + "\"" + vars["objId"] + "\""
				_, err := dbutils.ExecuteSQLStmt(dbCmd, gMgr.dbHdl)
				if err != nil {
					errCode = SRIdDeleteFail
					logger.Println("Failure in deleting Uuid map entry for ", vars["objId"], err)
				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusGone)
					resp.UUId = vars["objId"]
					js, err := json.Marshal(resp)
					if err != nil {
						errCode = SRRespMarshalErr
					} else {
						w.Write(js)
						errCode = SRSuccess
					}
				}
			} else {
				errCode = SRServerError
				logger.Println("DeleteObject returned failure ", obj)
			}
		} else {
			errCode = SRObjHdlError
			logger.Println("Failed to get object handle from http request ", objHdl, err)
		}
	} else {
		errCode = SRObjMapError
		logger.Println("Failed to get ObjectMap ", resource)
	}
	if errCode != SRSuccess {
		http.Error(w, SRErrString(errCode), http.StatusInternalServerError)
	}
	return
}

func ConfigObjectDelete(w http.ResponseWriter, r *http.Request) {
	var resp ConfigResponse
	var errCode int
	var objKey string
	var success bool
	var uuid string
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBase), "/")[0]
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		if _, obj, err := GetConfigObj(r, objHdl); err == nil {
			objKey, _ = obj.GetKey()
			dbObj, _ := obj.GetObjectFromDb(objKey, gMgr.dbHdl)
			success = gMgr.objHdlMap[resource].owner.DeleteObject(dbObj, objKey, gMgr.dbHdl)
			if success == true {
				gMgr.dbHdl.QueryRow("select Uuid from UuidMap where Key = ?", objKey).Scan(&uuid)
				dbCmd := "delete from " + "UuidMap" + " where Uuid = " + "\"" + uuid + "\""
				_, err := dbutils.ExecuteSQLStmt(dbCmd, gMgr.dbHdl)
				if err != nil {
					errCode = SRIdDeleteFail
					logger.Println("Failure in deleting Uuid map entry for ", uuid, err)
				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusGone)
					resp.UUId = uuid
					js, err := json.Marshal(resp)
					if err != nil {
						errCode = SRRespMarshalErr
					} else {
						w.Write(js)
						errCode = SRSuccess
					}
				}
			} else {
				errCode = SRServerError
				logger.Println("DeleteObject returned failure ", obj)
			}
		} else {
			errCode = SRObjHdlError
			logger.Println("Failed to get object handle from http request ", objHdl, err)
		}
	} else {
		errCode = SRObjMapError
		logger.Println("Failed to get ObjectMap ", resource)
	}
	if errCode != SRSuccess {
		http.Error(w, SRErrString(errCode), http.StatusInternalServerError)
	}
	return
}

func ConfigObjectUpdateForId(w http.ResponseWriter, r *http.Request) {
	var resp ConfigResponse
	var errCode int
	var objKey string
	var success bool
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBase), "/")[0]
	vars := mux.Vars(r)
	err := gMgr.dbHdl.QueryRow("select Key from UuidMap where Uuid = ?", vars["objId"]).Scan(&objKey)
	if err != nil {
		http.Error(w, SRErrString(SRNotFound), http.StatusNotFound)
		return
	}
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		body, obj, _ := GetConfigObj(r, objHdl)
		updateKeys, _ := GetUpdateKeys(body)
		dbObj, gerr := obj.GetObjectFromDb(objKey, gMgr.dbHdl)
		if gerr == nil {
			diff, _ := obj.CompareObjectsAndDiff(updateKeys, dbObj)
			mergedObj, _ := obj.MergeDbAndConfigObj(dbObj, diff)
			mergedObjKey, _ := mergedObj.GetKey()
			if objKey == mergedObjKey {
				success = gMgr.objHdlMap[resource].owner.UpdateObject(dbObj, mergedObj, diff, objKey, gMgr.dbHdl)
				if success == true {
					resp.UUId = vars["objId"]
					js, err := json.Marshal(resp)
					if err != nil {
						errCode = SRRespMarshalErr
					} else {
						w.Header().Set("Content-Type", "application/json; charset=UTF-8")
						w.WriteHeader(http.StatusOK)
						w.Write(js)
						errCode = SRSuccess
					}
				} else {
					errCode = SRServerError
					logger.Println("UpdateObject failed for resource ", updateKeys, resource)
				}
			} else {
				errCode = SRUpdateKeyError
				logger.Println("Cannot update key ", updateKeys, resource)
			}
		} else {
			errCode = SRObjHdlError
			logger.Println("Config update failed in getting obj via objKey ", objKey, gerr)
		}
	} else {
		errCode = SRObjMapError
		logger.Println("Config update failed t get ObjectMap ", resource)
	}

	if errCode != SRSuccess {
		errStr := SRErrString(errCode) + " ObjectId: " + vars["objId"]
		http.Error(w, errStr, http.StatusInternalServerError)
	}
	return
}

func ConfigObjectUpdate(w http.ResponseWriter, r *http.Request) {
	var resp ConfigResponse
	var errCode int
	var objKey string
	var success bool
	var uuid string
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBase), "/")[0]
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		body, obj, _ := GetConfigObj(r, objHdl)
		objKey, _ = obj.GetKey()
		gMgr.dbHdl.QueryRow("select Uuid from UuidMap where Key = ?", objKey).Scan(&uuid)
		updateKeys, _ := GetUpdateKeys(body)
		dbObj, gerr := obj.GetObjectFromDb(objKey, gMgr.dbHdl)
		if gerr == nil {
			diff, _ := obj.CompareObjectsAndDiff(updateKeys, dbObj)
			mergedObj, _ := obj.MergeDbAndConfigObj(dbObj, diff)
			mergedObjKey, _ := mergedObj.GetKey()
			if objKey == mergedObjKey {
				success = gMgr.objHdlMap[resource].owner.UpdateObject(dbObj, mergedObj, diff, objKey, gMgr.dbHdl)
				if success == true {
					resp.UUId = uuid
					js, err := json.Marshal(resp)
					if err != nil {
						errCode = SRRespMarshalErr
					} else {
						w.Header().Set("Content-Type", "application/json; charset=UTF-8")
						w.WriteHeader(http.StatusOK)
						w.Write(js)
						errCode = SRSuccess
					}
				} else {
					errCode = SRServerError
					logger.Println("UpdateObject failed for resource ", updateKeys, resource)
				}
			} else {
				errCode = SRUpdateKeyError
				logger.Println("Cannot update key ", updateKeys, resource)
			}
		} else {
			errCode = SRObjHdlError
			logger.Println("Config update failed in getting obj via objKey ", objKey, gerr)
		}
	} else {
		errCode = SRObjMapError
		logger.Println("Config update failed cannot get ObjectMap ", resource)
	}

	if errCode != SRSuccess {
		errStr := SRErrString(errCode) + " ObjectId: " + uuid
		http.Error(w, errStr, http.StatusInternalServerError)
	}
	return
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
