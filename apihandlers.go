package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"models"
	"net/http"
	"strings"
	//"net/url"
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
