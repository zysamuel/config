package main

import (
	"github.com/garyburd/redigo/redis"
	"models"
)

type ClientIf interface {
	Initialize(name string, address string)
	ConnectToServer() bool
	IsConnectedToServer() bool
	CreateObject(obj models.ConfigObj, dbHdl redis.Conn) (error, bool)
	DeleteObject(obj models.ConfigObj, objKey string, dbHdl redis.Conn) (error, bool)
	GetBulkObject(obj models.ConfigObj, dbHdl redis.Conn, currMarker int64, count int64) (err error,
		objcount int64,
		nextMarker int64,
		more bool,
		objs []models.ConfigObj)
	UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl redis.Conn) (error, bool)
	GetObject(obj models.ConfigObj, dbHdl redis.Conn) (error, models.ConfigObj)
	ExecuteAction(obj models.ConfigObj) error
	GetServerName() string
}
