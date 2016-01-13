package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"models"
	"utils/dbutils"
)

var UsrConfDbName string

//
//  This method creates new rest router interface
//
func (mgr *ConfigMgr) InstantiateDbIf() error {
	var err error
	var DbName string = "/UsrConfDb.db"

	UsrConfDbName = mgr.fullPath + DbName

	mgr.dbHdl, err = sql.Open("sqlite3", UsrConfDbName)
	if err == nil {
		for key, obj := range models.ConfigObjectMap {
			logger.Println("Creating DB for object", key)
			err = obj.CreateDBTable(mgr.dbHdl)
			if err != nil {
				logger.Println("Failed to create DB for object", key)
			}
		}
	} else {
		logger.Println("### Failed to open DB", UsrConfDbName, err)
	}

	/*
	 * Created a table in DB to store UUID to ConfigObject key mapping.
	 */
	logger.Println("Creating table for UUID")
	dbCmd := "CREATE TABLE IF NOT EXISTS UuidMap " +
		"(Uuid varchar(255) PRIMARY KEY ," +
		"Key varchar(255))"

	_, err = dbutils.ExecuteSQLStmt(dbCmd, mgr.dbHdl)
	if err != nil {
		logger.Println("Failed to create DB for object UUID")
	}
	return nil
}
