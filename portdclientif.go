package main

import (
	"database/sql"
	"models"
	"portdServices"
)

type PORTDClient struct {
	IPCClientBase
	ClientHdl *portdServices.PORTDServicesClient
}

func (clnt *PORTDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *PORTDClient) ConnectToServer() bool {

	clnt.Transport, clnt.PtrProtocolFactory = CreateIPCHandles(clnt.Address)
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = portdServices.NewPORTDServicesClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
	}
	return true
}
func (clnt *PORTDClient) IsConnectedToServer() bool {
	return true
}
func (clnt *PORTDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {
	default:
		break
	}

	return objId, true
}
func (clnt *PORTDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {

	switch obj.(type) {
	default:
		break
	}

	return true
}
func (clnt *PORTDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
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
func (clnt *PORTDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []byte, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called PORTD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.IPv4Intf:
		// cast original object
		origdata := dbObj.(models.IPv4Intf)
		updatedata := obj.(models.IPv4Intf)
		// create new thrift objects
		origconf := portdServices.NewIPv4Intf()
		updateconf := portdServices.NewIPv4Intf()

		origconf.RouterIf = int32(origdata.RouterIf)
		origconf.IfType = int32(origdata.IfType)
		origconf.IpAddr = string(origdata.IpAddr)

		updateconf.RouterIf = int32(updatedata.RouterIf)
		updateconf.IfType = int32(updatedata.IfType)
		updateconf.IpAddr = string(updatedata.IpAddr)

		//convert attrSet to uint8 list
		newattrset := make([]int8, len(attrSet))
		for i, v := range attrSet {
			newattrset[i] = int8(v)
		}
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateIPv4Intf(origconf, updateconf, newattrset)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.IPv4Neighbor:
		// cast original object
		origdata := dbObj.(models.IPv4Neighbor)
		updatedata := obj.(models.IPv4Neighbor)
		// create new thrift objects
		origconf := portdServices.NewIPv4Neighbor()
		updateconf := portdServices.NewIPv4Neighbor()

		origconf.RouterIf = int32(origdata.RouterIf)
		origconf.MacAddr = string(origdata.MacAddr)
		origconf.IpAddr = string(origdata.IpAddr)
		origconf.VlanId = int32(origdata.VlanId)

		updateconf.RouterIf = int32(updatedata.RouterIf)
		updateconf.MacAddr = string(updatedata.MacAddr)
		updateconf.IpAddr = string(updatedata.IpAddr)
		updateconf.VlanId = int32(updatedata.VlanId)

		//convert attrSet to uint8 list
		newattrset := make([]int8, len(attrSet))
		for i, v := range attrSet {
			newattrset[i] = int8(v)
		}
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateIPv4Neighbor(origconf, updateconf, newattrset)
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
