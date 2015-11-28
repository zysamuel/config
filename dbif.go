package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"models"
)

var UsrConfDbName string = "UsrConfDb.db"

//
//  This method creates new rest router interface
//
func (mgr *ConfigMgr) InstantiateDbIf() error {
	var err error
	mgr.dbHdl, err = sql.Open("sqlite3", UsrConfDbName)
	if err == nil {
		for key, obj := range models.ConfigObjectMap {
			logger.Println("### Creating DB for object", key)
			obj.CreateDBTable(mgr.dbHdl)
		}
	}
	return nil
}
