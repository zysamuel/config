package main

import (
	"database/sql"
	"models"
	"portdServices"
	"utils/ipcutils"
)

type PORTDClient struct {
	ipcutils.IPCClientBase
	ClientHdl *portdServices.PORTDServicesClient
}

func (clnt *PORTDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *PORTDClient) ConnectToServer() bool {

	clnt.Transport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = portdServices.NewPORTDServicesClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}
func (clnt *PORTDClient) IsConnectedToServer() bool {
	return clnt.IsConnected
}
func (clnt *PORTDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {

	case models.Vlan:
		data := obj.(models.Vlan)
		conf := portdServices.NewVlan()
		models.ConvertportdVlanObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateVlan(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.IPv4Intf:
		data := obj.(models.IPv4Intf)
		conf := portdServices.NewIPv4Intf()
		models.ConvertportdIPv4IntfObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateIPv4Intf(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.IPv4Neighbor:
		data := obj.(models.IPv4Neighbor)
		conf := portdServices.NewIPv4Neighbor()
		models.ConvertportdIPv4NeighborObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateIPv4Neighbor(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break
	default:
		break
	}

	return objId, true
}
func (clnt *PORTDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {

	switch obj.(type) {

	case models.Vlan:
		data := obj.(models.Vlan)
		conf := portdServices.NewVlan()
		models.ConvertportdVlanObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteVlan(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.IPv4Intf:
		data := obj.(models.IPv4Intf)
		conf := portdServices.NewIPv4Intf()
		models.ConvertportdIPv4IntfObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteIPv4Intf(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.IPv4Neighbor:
		data := obj.(models.IPv4Neighbor)
		conf := portdServices.NewIPv4Neighbor()
		models.ConvertportdIPv4NeighborObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteIPv4Neighbor(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break
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
func (clnt *PORTDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called PORTD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.Vlan:
		// cast original object
		origdata := dbObj.(models.Vlan)
		updatedata := obj.(models.Vlan)
		// create new thrift objects
		origconf := portdServices.NewVlan()
		updateconf := portdServices.NewVlan()
		models.ConvertportdVlanObjToThrift(&origdata, origconf)
		models.ConvertportdVlanObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateVlan(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.IPv4Intf:
		// cast original object
		origdata := dbObj.(models.IPv4Intf)
		updatedata := obj.(models.IPv4Intf)
		// create new thrift objects
		origconf := portdServices.NewIPv4Intf()
		updateconf := portdServices.NewIPv4Intf()
		models.ConvertportdIPv4IntfObjToThrift(&origdata, origconf)
		models.ConvertportdIPv4IntfObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateIPv4Intf(origconf, updateconf, attrSet)
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
		models.ConvertportdIPv4NeighborObjToThrift(&origdata, origconf)
		models.ConvertportdIPv4NeighborObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateIPv4Neighbor(origconf, updateconf, attrSet)
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
