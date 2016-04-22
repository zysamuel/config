package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/nu7hatch/gouuid"
	"strings"
)

type dbHandler struct {
	redis.Conn
}

func (d dbHandler) StoreUUIDToObjKeyMap(objKey string) (string, error) {
	UUId, err := uuid.NewV4()
	if err != nil {
		logger.Println("Failed to get UUID ", err)
		return "", err
	}
	_, err = d.Do("SET", UUId.String(), objKey)
	if err != nil {
		logger.Println("Failed to insert uuid to objkey entry in db ", err)
		return "", err
	}
	objKeyWithUUIDPrefix := "UUID" + objKey
	_, err = d.Do("SET", objKeyWithUUIDPrefix, UUId.String())
	if err != nil {
		logger.Println("Failed to insert objkey to uuid entry in db ", err)
		return "", err
	}
	return UUId.String(), nil

}

func (d dbHandler) DeleteUUIDToObjKeyMap(uuid, objKey string) error {
	_, err := d.Do("DEL", uuid)
	if err != nil {
		logger.Println("Failed to delete uuid to objkey entry in db ", err)
		return err
	}
	objKeyWithUUIDPrefix := "UUID" + objKey
	_, err = d.Do("DEL", objKeyWithUUIDPrefix)
	if err != nil {
		logger.Println("Failed to delete objkey to uuid entry in db ", err)
		return err
	}
	return nil
}

func (d dbHandler) GetUUIDFromObjKey(objKey string) (string, error) {
	objKeyWithUUIDPrefix := "UUID" + objKey
	uuid, err := redis.String(d.Do("GET", objKeyWithUUIDPrefix))
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func (d dbHandler) GetObjKeyFromUUID(uuid string) (string, error) {
	objKey, err := redis.String(d.Do("GET", uuid))
	if err != nil {
		return "", err
	}
	objKey = strings.Trim(objKey, "UID")
	return objKey, nil
}

//  This method initializes the db handler
func (mgr *ConfigMgr) InstantiateDbIf() error {
	var err error

	mgr.dbHdl.Conn, err = redis.Dial("tcp", ":6379")
	if err != nil {
		logger.Println("Failed to dial out to Redis server")
	}
	return nil
}
