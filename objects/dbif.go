package objects

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/nu7hatch/gouuid"
	"strings"
	"utils/dbutils"
	"utils/logging"
)

type DbHandler struct {
	*dbutils.DBUtil
	logger *logging.Writer
}

//  This method initializes the db handler
func InstantiateDbIf(logger *logging.Writer) *DbHandler {
	var err error

	dbHdl := new(DbHandler)
	dbHdl.DBUtil = dbutils.NewDBUtil(nil)
	err = dbHdl.DBUtil.Connect()
	if err != nil {
		logger.Err(fmt.Sprintln("Failed to dial out to Redis server"))
		return nil
	}
	dbHdl.logger = logger
	return dbHdl
}

func (d *DbHandler) StoreUUIDToObjKeyMap(objKey string) (string, error) {
	UUId, err := uuid.NewV4()
	if err != nil {
		d.logger.Err(fmt.Sprintln("Failed to get UUID ", err))
		return "", err
	}
	_, err = d.Do("SET", UUId.String(), objKey)
	if err != nil {
		d.logger.Err(fmt.Sprintln("Failed to insert uuid to objkey entry in db ", err))
		return "", err
	}
	objKeyWithUUIDPrefix := "UUID" + objKey
	_, err = d.Do("SET", objKeyWithUUIDPrefix, UUId.String())
	if err != nil {
		d.logger.Err(fmt.Sprintln("Failed to insert objkey to uuid entry in db ", err))
		return "", err
	}
	return UUId.String(), nil
}

func (d *DbHandler) DeleteUUIDToObjKeyMap(uuid, objKey string) error {
	_, err := d.Do("DEL", uuid)
	if err != nil {
		d.logger.Err(fmt.Sprintln("Failed to delete uuid to objkey entry in db ", err))
		return err
	}
	objKeyWithUUIDPrefix := "UUID" + objKey
	_, err = d.Do("DEL", objKeyWithUUIDPrefix)
	if err != nil {
		d.logger.Err(fmt.Sprintln("Failed to delete objkey to uuid entry in db ", err))
		return err
	}
	return nil
}

func (d *DbHandler) GetUUIDFromObjKey(objKey string) (string, error) {
	objKeyWithUUIDPrefix := "UUID" + objKey
	uuid, err := redis.String(d.Do("GET", objKeyWithUUIDPrefix))
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func (d *DbHandler) GetObjKeyFromUUID(uuid string) (string, error) {
	objKey, err := redis.String(d.Do("GET", uuid))
	if err != nil {
		return "", err
	}
	objKey = strings.TrimRight(objKey, "UUID")
	return objKey, nil
}
