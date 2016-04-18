package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
	"io"
	"io/ioutil"
	"models"
	"net/http"
	"strconv"
	"strings"
	"utils/dbutils"
	//"net/url"
	//"path"
)

const (
	MAX_OBJECTS_IN_GETBULK = 1024
	MAX_JSON_LENGTH        = 4096
)

type ConfigResponse struct {
	UUId  string `json:"ObjectId"`
	Error string `json:"Error"`
}

type ReturnObject struct {
	ObjectId         string `json:"ObjectId"`
	models.ConfigObj `json:"Object"`
}

type GetBulkResponse struct {
	MoreExist     bool  `json:"MoreExist"`
	ObjCount      int64 `json:"ObjCount"`
	CurrentMarker int64 `json:"CurrentMarker"`
	NextMarker    int64 `json:"NextMarker"`
	Objects       []ReturnObject
}

type ActionResponse struct {
	Error string `json:"Error"`
}

type ErrorResponse struct {
	Error string `json:"Error"`
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
	if err != nil {
		err = errors.New("Failed to decode input json data")
	}
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

func CheckIfSystemIsReady() bool {
	return gMgr.IsReady()
}

func RespondErrorForApiCall(w http.ResponseWriter, errCode int, errString string) error {
	var errResp ErrorResponse
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if errCode == SRBulkGetTooLarge {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
	} else if errCode == SRSystemNotReady {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	errResp.Error = SRErrString(errCode) + " " + errString
	js, _ := json.Marshal(errResp)
	w.Write(js)
	return nil
}

func GetOneConfigObjectForId(w http.ResponseWriter, r *http.Request) {
	var obj models.ConfigObj
	var dbObj models.ConfigObj
	var objKey string
	var retObj ReturnObject
	var err error

	gMgr.apiCallStats.NumGetCalls++
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBaseConfig), "/")[0]
	objHdl, ok := models.ConfigObjectMap[resource]
	if !ok {
		RespondErrorForApiCall(w, SRNotFound, "")
	}
	_, obj, err = GetConfigObj(r, objHdl)
	if err != nil {
		RespondErrorForApiCall(w, SRNotFound, err.Error())
		return
	}
	vars := mux.Vars(r)
	uuid := vars["objId"]
	//if objId is provided then read objkey from DB
	err = gMgr.dbHdl.QueryRow("select Key from UuidMap where Uuid = ?", uuid).Scan(&objKey)
	if err != nil {
		RespondErrorForApiCall(w, SRNotFound, err.Error())
		return
	}
	if dbObj, err = obj.GetObjectFromDb(objKey, gMgr.dbHdl); err != nil {
		RespondErrorForApiCall(w, SRNotFound, err.Error())
		return
	} else {
		retObj.ConfigObj = dbObj
	}
	retObj.ObjectId = uuid
	js, err := json.Marshal(retObj)
	if err == nil {
		gMgr.apiCallStats.NumGetCallsSuccess++
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}
	return
}

func GetOneConfigObject(w http.ResponseWriter, r *http.Request) {
	var obj models.ConfigObj
	var objKey string
	var retObj ReturnObject
	var err error
	var uuid string

	gMgr.apiCallStats.NumGetCalls++
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBaseConfig), "/")[0]
	objHdl, ok := models.ConfigObjectMap[resource]
	if !ok {
		RespondErrorForApiCall(w, SRNotFound, "")
	}
	_, obj, err = GetConfigObj(r, objHdl)
	if err != nil {
		RespondErrorForApiCall(w, SRNotFound, err.Error())
		return
	}
	//Get key fields provided in the request.
	objKey, err = obj.GetKey()
	if err != nil {
		RespondErrorForApiCall(w, SRNotFound, err.Error())
		return
	}
	if retObj.ConfigObj, err = obj.GetObjectFromDb(objKey, gMgr.dbHdl); err != nil {
		RespondErrorForApiCall(w, SRNotFound, err.Error())
		return
	}
	gMgr.dbHdl.QueryRow("select Uuid from UuidMap where Key = ?", objKey).Scan(&uuid)
	retObj.ObjectId = uuid
	js, err := json.Marshal(retObj)
	if err == nil {
		gMgr.apiCallStats.NumGetCallsSuccess++
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}
	return
}

func GetOneStateObjectForId(w http.ResponseWriter, r *http.Request) {
	var obj models.ConfigObj
	var objKey string
	var retObj ReturnObject
	var err error

	gMgr.apiCallStats.NumGetCalls++
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBaseState), "/")[0]
	resource = resource + "State"
	objHdl, ok := models.ConfigObjectMap[resource]
	if !ok {
		RespondErrorForApiCall(w, SRNotFound, "")
	}
	_, obj, err = GetConfigObj(r, objHdl)
	if err != nil {
		RespondErrorForApiCall(w, SRNotFound, err.Error())
		return
	}
	vars := mux.Vars(r)
	uuid := vars["objId"]
	//if objId is provided then read objkey from DB
	err = gMgr.dbHdl.QueryRow("select Key from UuidMap where Uuid = ?", uuid).Scan(&objKey)
	if err != nil {
		RespondErrorForApiCall(w, SRNotFound, err.Error())
		return
	}
	resourceOwner := gMgr.objHdlMap[resource].owner
	if resourceOwner.IsConnectedToServer() == false {
		errString := "Confd not connected to " + resourceOwner.GetServerName()
		RespondErrorForApiCall(w, SRSystemNotReady, errString)
		return
	}
	if err, retObj.ConfigObj = resourceOwner.GetObject(obj); err != nil {
		RespondErrorForApiCall(w, SRNotFound, err.Error())
		return
	}
	retObj.ObjectId = uuid
	js, err := json.Marshal(retObj)
	if err == nil {
		gMgr.apiCallStats.NumGetCallsSuccess++
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}
	return
}

func GetOneStateObject(w http.ResponseWriter, r *http.Request) {
	var obj models.ConfigObj
	var objKey string
	var retObj ReturnObject
	var err error
	var uuid string

	gMgr.apiCallStats.NumGetCalls++
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBaseState), "/")[0]
	resource = resource + "State"
	objHdl, ok := models.ConfigObjectMap[resource]
	if !ok {
		RespondErrorForApiCall(w, SRNotFound, "")
	}
	_, obj, err = GetConfigObj(r, objHdl)
	if err != nil {
		RespondErrorForApiCall(w, SRNotFound, err.Error())
		return
	}
	//Get key fields provided in the request.
	objKey, err = obj.GetKey()
	if err != nil {
		RespondErrorForApiCall(w, SRNotFound, err.Error())
		return
	}
	resourceOwner := gMgr.objHdlMap[resource].owner
	if resourceOwner.IsConnectedToServer() == false {
		errString := "Confd not connected to " + resourceOwner.GetServerName()
		RespondErrorForApiCall(w, SRSystemNotReady, errString)
		return
	}
	if err, retObj.ConfigObj = resourceOwner.GetObject(obj); err != nil {
		RespondErrorForApiCall(w, SRNotFound, err.Error())
		return
	}
	gMgr.dbHdl.QueryRow("select Uuid from UuidMap where Key = ?", objKey).Scan(&uuid)
	retObj.ObjectId = uuid
	js, err := json.Marshal(retObj)
	if err == nil {
		gMgr.apiCallStats.NumGetCallsSuccess++
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}
	return
}

func BulkGetConfigObjects(w http.ResponseWriter, r *http.Request) {
	var errCode int
	var objKey string
	var configObjects []models.ConfigObj
	var resp GetBulkResponse
	var err error
	gMgr.apiCallStats.NumGetCalls++
	resource := strings.TrimPrefix(r.URL.String(), gMgr.apiBaseConfig)
	resource = strings.Split(resource, "?")[0]
	resource = resource[:len(resource)-1]
	objHdl, ok := models.ConfigObjectMap[resource]
	if !ok {
		RespondErrorForApiCall(w, SRNotFound, "")
	}
	_, obj, err := GetConfigObj(nil, objHdl)
	if err != nil {
		RespondErrorForApiCall(w, SRNotFound, err.Error())
		return
	}
	currentIndex, objCount := ExtractGetBulkParams(r)
	if objCount > MAX_OBJECTS_IN_GETBULK {
		RespondErrorForApiCall(w, SRBulkGetTooLarge, err.Error())
		logger.Println("Too many objects requested in bulkget ", objCount)
		return
	}
	resp.CurrentMarker = currentIndex
	err, resp.ObjCount, resp.NextMarker, resp.MoreExist,
		configObjects = obj.GetBulkObjFromDb(currentIndex, objCount, gMgr.dbHdl)
	if err == nil {
		resp.Objects = make([]ReturnObject, resp.ObjCount)
		for idx, configObject := range configObjects {
			resp.Objects[idx].ConfigObj = configObject
			objKey, _ = configObject.GetKey()
			gMgr.dbHdl.QueryRow("select Uuid from UuidMap where Key = ?", objKey).Scan(&resp.Objects[idx].ObjectId)
		}
		js, err := json.Marshal(resp)
		if err != nil {
			errCode = SRRespMarshalErr
			logger.Println("### Error in marshalling JSON in getBulk for object ", resource, err)
		} else {
			gMgr.apiCallStats.NumGetCallsSuccess++
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			w.Write(js)
			errCode = SRSuccess
		}
	}
	if errCode != SRSuccess {
		RespondErrorForApiCall(w, errCode, err.Error())
	}
	return
}

func BulkGetStateObjects(w http.ResponseWriter, r *http.Request) {
	var errCode int
	var objKey string
	var stateObjects []models.ConfigObj
	var resp GetBulkResponse
	var err error
	gMgr.apiCallStats.NumGetCalls++
	resource := strings.TrimPrefix(r.URL.String(), gMgr.apiBaseState)
	resource = strings.Split(resource, "?")[0]
	resource = resource[:len(resource)-1]
	resource = resource + "State"
	objHdl, ok := models.ConfigObjectMap[resource]
	if !ok {
		RespondErrorForApiCall(w, SRNotFound, "")
	}
	_, obj, err := GetConfigObj(nil, objHdl)
	if err != nil {
		RespondErrorForApiCall(w, SRNotFound, err.Error())
		return
	}
	currentIndex, objCount := ExtractGetBulkParams(r)
	if objCount > MAX_OBJECTS_IN_GETBULK {
		RespondErrorForApiCall(w, SRBulkGetTooLarge, err.Error())
		logger.Println("Too many objects requested in bulkget ", objCount)
		return
	}
	resourceOwner := gMgr.objHdlMap[resource].owner
	if resourceOwner.IsConnectedToServer() == false {
		errString := "Confd not connected to " + resourceOwner.GetServerName()
		RespondErrorForApiCall(w, SRSystemNotReady, errString)
		return
	}
	resp.CurrentMarker = currentIndex
	err, resp.ObjCount, resp.NextMarker, resp.MoreExist,
		stateObjects = resourceOwner.GetBulkObject(obj, currentIndex, objCount)
	if err == nil {
		resp.Objects = make([]ReturnObject, resp.ObjCount)
		for idx, stateObject := range stateObjects {
			resp.Objects[idx].ConfigObj = stateObject
			objKey, _ = stateObject.GetKey()
			gMgr.dbHdl.QueryRow("select Uuid from UuidMap where Key = ?", objKey).Scan(&resp.Objects[idx].ObjectId)
		}
		js, err := json.Marshal(resp)
		if err != nil {
			errCode = SRRespMarshalErr
			logger.Println("### Error in marshalling JSON in getBulk for object ", resource, err)
		} else {
			gMgr.apiCallStats.NumGetCallsSuccess++
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			w.Write(js)
			errCode = SRSuccess
		}
	}
	if errCode != SRSuccess {
		RespondErrorForApiCall(w, errCode, err.Error())
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

func ExecuteActionObject(w http.ResponseWriter, r *http.Request) {
	var resp ActionResponse
	var errCode int
	var err error
	var obj models.ConfigObj

	gMgr.apiCallStats.NumActionCalls++
	errCode = SRSuccess
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resource := strings.TrimPrefix(r.URL.String(), gMgr.apiBaseAction)
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		if _, obj, err = GetConfigObj(r, objHdl); err == nil {
			resourceOwner := gMgr.objHdlMap[resource].owner
			if resourceOwner.IsConnectedToServer() == false {
				errString := "Confd not connected to " + resourceOwner.GetServerName()
				RespondErrorForApiCall(w, SRSystemNotReady, errString)
				return
			}
			err = resourceOwner.ExecuteAction(obj)
			if err == nil {
				gMgr.apiCallStats.NumActionCallsSuccess++
				w.WriteHeader(http.StatusOK)
				errCode = SRSuccess
			} else {
				resp.Error = err.Error()
				errCode = SRServerError
				logger.Println("Failed to execute action: ", obj, " due to error: ", err)
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
		w.WriteHeader(http.StatusInternalServerError)
	}
	if errCode != SRServerError {
		resp.Error = SRErrString(errCode)
	}
	js, err := json.Marshal(resp)
	if err != nil {
		logger.Println("ExecuteAction failed to Marshal config response")
	}
	w.Write(js)

	return
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
	var err error
	var obj models.ConfigObj
	var objKey string
	var body []byte

	gMgr.apiCallStats.NumCreateCalls++
	errCode = SRSuccess
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resource := strings.TrimPrefix(r.URL.String(), gMgr.apiBaseConfig)
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		if body, obj, err = GetConfigObj(r, objHdl); err == nil {
			updateKeys, _ := GetUpdateKeys(body)
			if len(updateKeys) == 0 {
				errCode = SRNoContent
				logger.Println("Nothing to configure")
			} else {
				objKey, _ = obj.GetKey()
				dbObj, _ := obj.GetObjectFromDb(objKey, gMgr.dbHdl)
				dbObjKey, _ := dbObj.GetKey()
				if dbObjKey == objKey {
					errCode = SRAlreadyConfigured
					logger.Println("Config object is present")
					gMgr.dbHdl.QueryRow("select Uuid from UuidMap where Key = ?", objKey).Scan(&uuid)
				}
			}
			if errCode != SRSuccess {
				w.WriteHeader(http.StatusInternalServerError)
				resp.UUId = uuid
				resp.Error = SRErrString(errCode)
				js, _ := json.Marshal(resp)
				w.Write(js)
				return
			}
			resourceOwner := gMgr.objHdlMap[resource].owner
			if resourceOwner.IsConnectedToServer() == false {
				errString := "Confd not connected to " + resourceOwner.GetServerName()
				RespondErrorForApiCall(w, SRSystemNotReady, errString)
				return
			}
			err, success = resourceOwner.CreateObject(obj, gMgr.dbHdl)
			if err == nil && success == true {
				UUId, dbErr := StoreUuidToKeyMapInDb(obj)
				if dbErr == nil {
					gMgr.apiCallStats.NumCreateCallsSuccess++
					w.WriteHeader(http.StatusCreated)
					resp.UUId = UUId.String()
					errCode = SRSuccess
				} else {
					errCode = SRIdStoreFail
					logger.Println("Failed to store UuidToKey map ", obj, dbErr)
				}
			} else {
				resp.Error = err.Error()
				errCode = SRServerError
				logger.Println("Failed to create object: ", obj, " due to error: ", err)
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
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err != nil && errCode != SRServerError {
		resp.Error = SRErrString(errCode) + " " + err.Error()
	}
	js, err := json.Marshal(resp)
	if err != nil {
		logger.Println("CreateObject failed to Marshal config response")
	}
	w.Write(js)

	return
}

func ConfigObjectDeleteForId(w http.ResponseWriter, r *http.Request) {
	var resp ConfigResponse
	var errCode int
	var objKey string
	var success bool
	var err error

	gMgr.apiCallStats.NumDeleteCalls++
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBaseConfig), "/")[0]
	vars := mux.Vars(r)
	resp.UUId = vars["objId"]
	err = gMgr.dbHdl.QueryRow("select Key from UuidMap where Uuid = ?", vars["objId"]).Scan(&objKey)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		resp.Error = SRErrString(SRNotFound)
		js, _ := json.Marshal(resp)
		w.Write(js)
		return
	}
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		if _, obj, err := GetConfigObj(nil, objHdl); err == nil {
			dbObj, _ := obj.GetObjectFromDb(objKey, gMgr.dbHdl)
			resourceOwner := gMgr.objHdlMap[resource].owner
			if resourceOwner.IsConnectedToServer() == false {
				errString := "Confd not connected to " + resourceOwner.GetServerName()
				RespondErrorForApiCall(w, SRSystemNotReady, errString)
				return
			}
			err, success = resourceOwner.DeleteObject(dbObj, objKey, gMgr.dbHdl)
			if err == nil && success == true {
				dbCmd := "delete from " + "UuidMap" + " where Uuid = " + "\"" + vars["objId"] + "\""
				_, err = dbutils.ExecuteSQLStmt(dbCmd, gMgr.dbHdl)
				if err != nil {
					errCode = SRIdDeleteFail
					logger.Println("Failure in deleting Uuid map entry for ", vars["objId"], err)
				} else {
					gMgr.apiCallStats.NumDeleteCallsSuccess++
					w.WriteHeader(http.StatusGone)
					errCode = SRSuccess
				}
			} else {
				resp.Error = err.Error()
				errCode = SRServerError
				logger.Println("DeleteObject returned failure ", obj, err)
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
		w.WriteHeader(http.StatusInternalServerError)
	}
	if errCode != SRServerError {
		resp.Error = SRErrString(errCode)
	}
	js, err := json.Marshal(resp)
	if err != nil {
		logger.Println("CreateObject failed to Marshal config response")
	}
	w.Write(js)

	return
}

func ConfigObjectDelete(w http.ResponseWriter, r *http.Request) {
	var resp ConfigResponse
	var errCode int
	var objKey string
	var success bool
	var uuid string
	var err error

	gMgr.apiCallStats.NumDeleteCalls++
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBaseConfig), "/")[0]
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		if _, obj, err := GetConfigObj(r, objHdl); err == nil {
			objKey, _ = obj.GetKey()
			dbObj, err := obj.GetObjectFromDb(objKey, gMgr.dbHdl)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				resp.Error = SRErrString(SRNotFound)
				js, _ := json.Marshal(resp)
				w.Write(js)
				return
			}
			gMgr.dbHdl.QueryRow("select Uuid from UuidMap where Key = ?", objKey).Scan(&uuid)
			resp.UUId = uuid
			resourceOwner := gMgr.objHdlMap[resource].owner
			if resourceOwner.IsConnectedToServer() == false {
				errString := "Confd not connected to " + resourceOwner.GetServerName()
				RespondErrorForApiCall(w, SRSystemNotReady, errString)
				return
			}
			err, success = resourceOwner.DeleteObject(dbObj, objKey, gMgr.dbHdl)
			if err == nil && success == true {
				dbCmd := "delete from " + "UuidMap" + " where Uuid = " + "\"" + uuid + "\""
				_, err = dbutils.ExecuteSQLStmt(dbCmd, gMgr.dbHdl)
				if err != nil {
					errCode = SRIdDeleteFail
					logger.Println("Failure in deleting Uuid map entry for ", uuid, err)
				} else {
					gMgr.apiCallStats.NumDeleteCallsSuccess++
					w.WriteHeader(http.StatusGone)
					errCode = SRSuccess
				}
			} else {
				resp.Error = err.Error()
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
		w.WriteHeader(http.StatusInternalServerError)
	}
	if errCode != SRServerError {
		resp.Error = SRErrString(errCode)
	}
	js, err := json.Marshal(resp)
	if err != nil {
		logger.Println("CreateObject failed to Marshal config response")
	}
	w.Write(js)

	return
}

func ConfigObjectUpdateForId(w http.ResponseWriter, r *http.Request) {
	var resp ConfigResponse
	var errCode int
	var objKey string
	var success bool
	var err error

	gMgr.apiCallStats.NumUpdateCalls++
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBaseConfig), "/")[0]
	vars := mux.Vars(r)
	resp.UUId = vars["objId"]
	err = gMgr.dbHdl.QueryRow("select Key from UuidMap where Uuid = ?", vars["objId"]).Scan(&objKey)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		resp.Error = SRErrString(SRNotFound)
		js, _ := json.Marshal(resp)
		w.Write(js)
		return
	}
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		body, obj, _ := GetConfigObj(r, objHdl)
		updateKeys, _ := GetUpdateKeys(body)
		dbObj, gerr := obj.GetObjectFromDb(objKey, gMgr.dbHdl)
		if gerr == nil {
			diff, _ := obj.CompareObjectsAndDiff(updateKeys, dbObj)
			anyUpdated := false
			for _, updated := range diff {
				if updated == true {
					anyUpdated = true
					break
				}
			}
			if anyUpdated == false {
				w.WriteHeader(http.StatusInternalServerError)
				resp.Error = SRErrString(SRUpdateNoChange)
				js, _ := json.Marshal(resp)
				w.Write(js)
				return
			}
			mergedObj, _ := obj.MergeDbAndConfigObj(dbObj, diff)
			mergedObjKey, _ := mergedObj.GetKey()
			if objKey == mergedObjKey {
				resourceOwner := gMgr.objHdlMap[resource].owner
				if resourceOwner.IsConnectedToServer() == false {
					errString := "Confd not connected to " + resourceOwner.GetServerName()
					RespondErrorForApiCall(w, SRSystemNotReady, errString)
					return
				}
				err, success = resourceOwner.UpdateObject(dbObj, mergedObj, diff, objKey, gMgr.dbHdl)
				if err == nil && success == true {
					gMgr.apiCallStats.NumUpdateCallsSuccess++
					w.WriteHeader(http.StatusOK)
					errCode = SRSuccess
				} else {
					resp.Error = err.Error()
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
		w.WriteHeader(http.StatusInternalServerError)
	}
	if errCode != SRServerError {
		resp.Error = SRErrString(errCode)
	}
	js, err := json.Marshal(resp)
	if err != nil {
		logger.Println("CreateObject failed to Marshal config response")
	}
	w.Write(js)

	return
}

func ConfigObjectUpdate(w http.ResponseWriter, r *http.Request) {
	var resp ConfigResponse
	var errCode int
	var objKey string
	var success bool
	var uuid string
	var err error

	gMgr.apiCallStats.NumUpdateCalls++
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resource := strings.Split(strings.TrimPrefix(r.URL.String(), gMgr.apiBaseConfig), "/")[0]
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		body, obj, _ := GetConfigObj(r, objHdl)
		objKey, _ = obj.GetKey()
		updateKeys, _ := GetUpdateKeys(body)
		dbObj, gerr := obj.GetObjectFromDb(objKey, gMgr.dbHdl)
		if gerr != nil {
			w.WriteHeader(http.StatusNotFound)
			resp.Error = SRErrString(SRNotFound)
			js, _ := json.Marshal(resp)
			w.Write(js)
			return
		}
		gMgr.dbHdl.QueryRow("select Uuid from UuidMap where Key = ?", objKey).Scan(&uuid)
		resp.UUId = uuid
		diff, _ := obj.CompareObjectsAndDiff(updateKeys, dbObj)
		anyUpdated := false
		for _, updated := range diff {
			if updated == true {
				anyUpdated = true
				break
			}
		}
		if anyUpdated == false {
			w.WriteHeader(http.StatusInternalServerError)
			resp.Error = SRErrString(SRUpdateNoChange)
			js, _ := json.Marshal(resp)
			w.Write(js)
			return
		}
		mergedObj, _ := obj.MergeDbAndConfigObj(dbObj, diff)
		mergedObjKey, _ := mergedObj.GetKey()
		if objKey == mergedObjKey {
			resourceOwner := gMgr.objHdlMap[resource].owner
			if resourceOwner.IsConnectedToServer() == false {
				errString := "Confd not connected to " + resourceOwner.GetServerName()
				RespondErrorForApiCall(w, SRSystemNotReady, errString)
				return
			}
			err, success = resourceOwner.UpdateObject(dbObj, mergedObj, diff, objKey, gMgr.dbHdl)
			if err == nil && success == true {
				gMgr.apiCallStats.NumUpdateCallsSuccess++
				w.WriteHeader(http.StatusOK)
				errCode = SRSuccess
			} else {
				resp.Error = err.Error()
				errCode = SRServerError
				logger.Println("UpdateObject failed for resource ", updateKeys, resource)
			}
		} else {
			errCode = SRUpdateKeyError
			logger.Println("Cannot update key ", updateKeys, resource)
		}
	} else {
		errCode = SRObjMapError
		logger.Println("Config update failed cannot get ObjectMap ", resource)
	}

	if errCode != SRSuccess {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if errCode != SRServerError {
		resp.Error = SRErrString(errCode)
	}
	js, err := json.Marshal(resp)
	if err != nil {
		logger.Println("CreateObject failed to Marshal config response")
	}
	w.Write(js)

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
