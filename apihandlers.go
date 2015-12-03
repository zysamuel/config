package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"models"
	"net/http"
	"strings"
	//"net/url"
	"strconv"
)

func GetConfigObj(r *http.Request, obj models.ConfigObj) (models.ConfigObj, error) {
	var retObj models.ConfigObj
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return retObj, err
	}
	if err := r.Body.Close(); err != nil {
		return retObj, err
	}
	return obj.UnmarshalObject(body)
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/jsoni;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(peers); err != nil {
		panic(err)
	}
}

func ShowConfigObject(w http.ResponseWriter, r *http.Request) {
	logger.Println("####  ShowConfigObject called")
}

func ConfigObjectCreate(w http.ResponseWriter, r *http.Request) {
	resource := strings.TrimPrefix(r.URL.String(), "/")
	logger.Println("####  CreateObject  called")
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		obj, _ := GetConfigObj(r, objHdl)
		objectId, success := gMgr.objHdlMap[resource].owner.CreateObject(obj, gMgr.dbHdl)
		if success == true {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(objectId); err != nil {
				logger.Println("### Failed to encode the objectId for object ", resource, objectId)
			}
		}
	}
	return
}

func ConfigObjectDelete(w http.ResponseWriter, r *http.Request) {
	resource := strings.Split(r.URL.String(), "/")[1]
	vars := mux.Vars(r)
	objId, err := strconv.ParseInt(vars["objId"], 10, 64)
	if err != nil {
		logger.Println("### Failure in deleting object with Id ", resource, vars["objId"], err)
		return
	}
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		obj, _ := GetConfigObj(r, objHdl)
		success := gMgr.objHdlMap[resource].owner.DeleteObject(obj, objId, gMgr.dbHdl)
		if success == true {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
		}
	}
	return
}
