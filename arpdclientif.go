package main

import (
	"arpdServices"
	"database/sql"
	"models"
)

type ARPDClient struct {
	IPCClientBase
	ClientHdl *arpdServices.ARPDServicesClient
}

func (clnt *ARPDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *ARPDClient) ConnectToServer() bool {

	clnt.Transport, clnt.PtrProtocolFactory = CreateIPCHandles(clnt.Address)
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = arpdServices.NewARPDServicesClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
	}
	return true
}
func (clnt *ARPDClient) IsConnectedToServer() bool {
	return true
}
func (clnt *ARPDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {
	default:
		break
	}

	return objId, true
}
func (clnt *ARPDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {

	switch obj.(type) {
	default:
		break
	}

	return true
}
func (clnt *ARPDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {

	logger.Println("### Get Bulk request called with", currMarker, count)
	switch obj.(type) {

	default:
		break
	}
	return nil, objCount, nextMarker, more, objs

}
func (clnt *ARPDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []byte, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called ARPD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	default:
		break
	}
	return ok

}
