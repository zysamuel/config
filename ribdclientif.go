package main

import (
	"database/sql"
	"models"
	"ribdServices"
)

type RIBDClient struct {
	IPCClientBase
	ClientHdl *ribdServices.RIBDServicesClient
}

func (clnt *RIBDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *RIBDClient) ConnectToServer() bool {

	clnt.Transport, clnt.PtrProtocolFactory = CreateIPCHandles(clnt.Address)
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = ribdServices.NewRIBDServicesClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
	}
	return true
}
func (clnt *RIBDClient) IsConnectedToServer() bool {
	return true
}
func (clnt *RIBDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {
	default:
		break
	}

	return objId, true
}
func (clnt *RIBDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {

	switch obj.(type) {
	default:
		break
	}

	return true
}
func (clnt *RIBDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
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
func (clnt *RIBDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []byte, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called RIBD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.IPV4Route:
		// cast original object
		origdata := dbObj.(models.IPV4Route)
		updatedata := obj.(models.IPV4Route)
		// create new thrift objects
		origconf := ribdServices.NewIPV4Route()
		updateconf := ribdServices.NewIPV4Route()

		origconf.DestinationNw = string(origdata.DestinationNw)
		origconf.OutgoingIntfType = string(origdata.OutgoingIntfType)
		origconf.Protocol = string(origdata.Protocol)
		origconf.OutgoingInterface = string(origdata.OutgoingInterface)
		origconf.NetworkMask = string(origdata.NetworkMask)
		origconf.NextHopIp = string(origdata.NextHopIp)

		updateconf.DestinationNw = string(updatedata.DestinationNw)
		updateconf.OutgoingIntfType = string(updatedata.OutgoingIntfType)
		updateconf.Protocol = string(updatedata.Protocol)
		updateconf.OutgoingInterface = string(updatedata.OutgoingInterface)
		updateconf.NetworkMask = string(updatedata.NetworkMask)
		updateconf.NextHopIp = string(updatedata.NextHopIp)

		//convert attrSet to uint8 list
		newattrset := make([]int8, len(attrSet))
		for i, v := range attrSet {
			newattrset[i] = int8(v)
		}
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateIPV4Route(origconf, updateconf, newattrset)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	default:
		break
	}
	return ok

}
