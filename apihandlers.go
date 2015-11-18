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
	logger.Println("####  CreateObject  called")
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		obj, _ := GetConfigObj(r, objHdl)
		logger.Println("### Config Object is ", obj)
		gMgr.objHdlMap[resource].owner.CreateObject(obj)
	}
	return
}
