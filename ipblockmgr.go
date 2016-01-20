package main

import (
	"database/sql"
	"fmt"
	"models"
)

var gIpBlockMgr *IpBlockMgr

type IpBlockHostIp struct {
	IpAddress string
	HostName  string
}

type IpBlockEntry struct {
	models.IPV4AddressBlock
	HostMap map[string]IpBlockHostIp
}

type IpBlockMgr struct {
	IpBlockMap map[string]IpBlockEntry
}

func GetIpBlockMgr() *IpBlockMgr {
	if gIpBlockMgr == nil {
		gIpBlockMgr = new(IpBlockMgr)
	}
	return gIpBlockMgr
}

func (ipbMgr *IpBlockMgr) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var ipblk models.IPV4AddressBlock
	ipblk = obj.(models.IPV4AddressBlock)
	fmt.Println(" Create Object called", ipblk)
	return 0, true
}

func (ipbMgr *IpBlockMgr) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {
	return true
}

func (ipbMgr *IpBlockMgr) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {
	return true
}
