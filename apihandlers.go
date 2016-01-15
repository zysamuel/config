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
	"encoding/base64"
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
		http.Error(w, "Show: "+SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	logger.Println("####  ShowConfigObject called")
}

func ConfigObjectsBulkGet(w http.ResponseWriter, r *http.Request) {
	var errCode int
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, "GetBulk: "+SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
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
		if objCount > MAX_OBJECTS_IN_GETBULK {
			http.Error(w, SRErrString(SRBulkGetTooLarge), http.StatusRequestEntityTooLarge)
			logger.Println("Too many objects requested in bulkget ", objCount)
			return
		}
		resp.CurrentMarker = currentIndex
		switch obj.(type) {
		case models.UserConfig:
			err, resp.ObjCount, resp.NextMarker, resp.MoreExist,
				resp.StateObjects = GetBulkObject(obj, currentIndex, objCount)
		default:
			err, resp.ObjCount, resp.NextMarker, resp.MoreExist,
				resp.StateObjects = gMgr.objHdlMap[resource].owner.GetBulkObject(obj,
				currentIndex,
				objCount)
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

fmt.Println("Create: ", *r)
auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
payload, _ := base64.StdEncoding.DecodeString(auth[1])
pair := strings.SplitN(string(payload), ":", 2)
fmt.Println("UserName: %s Password: %s", pair[0], pair[1])
return

	if CheckIfSystemIsReady(w) != true {
		http.Error(w, "Create: "+SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	resource := strings.TrimPrefix(r.URL.String(), gMgr.apiBase)
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		if body, obj, err := GetConfigObj(r, objHdl); err == nil {
			updateKeys, _ := GetUpdateKeys(body)
			if len(updateKeys) == 0 {
				errCode = SRNoContent
				logger.Println("Nothing to configure")
			} else {
				switch obj.(type) {
				case models.UserConfig:
					_, success = CreateObject(obj, gMgr.dbHdl)
				default:
					_, success = gMgr.objHdlMap[resource].owner.CreateObject(obj, gMgr.dbHdl)
				}
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
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, "Delete: "+SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
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
			switch obj.(type) {
			case models.UserConfig:
				success = DeleteObject(dbObj, objKey, gMgr.dbHdl)
			default:
				success = gMgr.objHdlMap[resource].owner.DeleteObject(dbObj, objKey, gMgr.dbHdl)
			}
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

func ConfigObjectUpdate(w http.ResponseWriter, r *http.Request) {
	var resp ConfigResponse
	var errCode int
	var objKey string
	var success bool
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, "Update: "+SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
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
			switch obj.(type) {
			case models.UserConfig:
				success = UpdateObject(dbObj, mergedObj, diff, objKey, gMgr.dbHdl)
			default:
				success = gMgr.objHdlMap[resource].owner.UpdateObject(dbObj, mergedObj, diff, objKey, gMgr.dbHdl)
			}
			if success == true {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				resp.UUId = vars["objId"]
				js, err := json.Marshal(resp)
				if err != nil {
					errCode = SRRespMarshalErr
				} else {
					w.Write(js)
					errCode = SRSuccess
				}
			} else {
				errCode = SRServerError
				logger.Println("UpdateObject failed for resource ", updateKeys, resource)
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
		http.Error(w, SRErrString(errCode), http.StatusNotModified)
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
