package main

import (
	"database/sql"
	"models"
)

type ClientIf interface {
	Initialize(name string, address string)
	ConnectToServer() bool
	IsConnectedToServer() bool
	CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (error, bool)
	DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) (error, bool)
	GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
		objcount int64,
		nextMarker int64,
		more bool,
		objs []models.ConfigObj)
	UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) (error, bool)
	GetObject(obj models.ConfigObj) (error, models.ConfigObj)
	ExecuteAction(obj models.ConfigObj) error
	GetServerName() string
}
