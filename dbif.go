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
func (mgr *ConfigMgr) InstantiateDbIf(params_Dir string) error {
	var err error
        var DbName string = "UsrConfDb.db"
        UsrConfDbName = params_Dir + "/../bin/" + DbName
	mgr.dbHdl, err = sql.Open("sqlite3", UsrConfDbName)
	if err == nil {
		for key, obj := range models.ConfigObjectMap {
			logger.Println("Creating DB for object", key)
			obj.CreateDBTable(mgr.dbHdl)
		}
	} else {
		logger.Println("### Failed to open DB", UsrConfDbName, err)
	}

	/*
	 * Created a table in DB to store UUID to ConfigObject key mapping.
	 */
	dbCmd := "CREATE TABLE IF NOT EXISTS UuidMap " +
		"(Uuid varchar(255) PRIMARY KEY ," +
		"Key varchar(255))"

	_, err = dbutils.ExecuteSQLStmt(dbCmd, mgr.dbHdl)
	if err == nil {
		logger.Println("Created table for UUID")
	}

	return nil
}
