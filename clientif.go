package main

import (
	"models"
	"utils/dbutils"
)

type ClientIf interface {
	Initialize(name string, address string)
	ConnectToServer() bool
	IsConnectedToServer() bool
	CreateObject(obj models.ConfigObj, dbHdl *dbutils.DBUtil) (error, bool)
	DeleteObject(obj models.ConfigObj, objKey string, dbHdl *dbutils.DBUtil) (error, bool)
	GetBulkObject(obj models.ConfigObj, dbHdl *dbutils.DBUtil, currMarker int64, count int64) (err error,
		objcount int64,
		nextMarker int64,
		more bool,
		objs []models.ConfigObj)
	UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *dbutils.DBUtil) (error, bool)
	GetObject(obj models.ConfigObj, dbHdl *dbutils.DBUtil) (error, models.ConfigObj)
	ExecuteAction(obj models.ConfigObj) error
	GetServerName() string
}
