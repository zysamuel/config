package main

import (
	"database/sql"
	"models"
	"portdServices"
	"utils/ipcutils"
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

	clnt.Transport, clnt.PtrProtocolFactory = ipcutils.CreateIPCHandles(clnt.Address)
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

	case models.Vlan:
		data := obj.(models.Vlan)
		conf := portdServices.NewVlan()
		conf.PortTagType = string(data.PortTagType)
		conf.Ports = string(data.Ports)
		conf.VlanId = int32(data.VlanId)

		_, err := clnt.ClientHdl.CreateVlan(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.IPv4Intf:
		data := obj.(models.IPv4Intf)
		conf := portdServices.NewIPv4Intf()
		conf.RouterIf = int32(data.RouterIf)
		conf.IfType = int32(data.IfType)
		conf.IpAddr = string(data.IpAddr)

		_, err := clnt.ClientHdl.CreateIPv4Intf(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.IPv4Neighbor:
		data := obj.(models.IPv4Neighbor)
		conf := portdServices.NewIPv4Neighbor()
		conf.RouterIf = int32(data.RouterIf)
		conf.MacAddr = string(data.MacAddr)
		conf.IpAddr = string(data.IpAddr)
		conf.VlanId = int32(data.VlanId)

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
		conf.PortTagType = string(data.PortTagType)
		conf.Ports = string(data.Ports)
		conf.VlanId = int32(data.VlanId)

		_, err := clnt.ClientHdl.DeleteVlan(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.IPv4Intf:
		data := obj.(models.IPv4Intf)
		conf := portdServices.NewIPv4Intf()
		conf.RouterIf = int32(data.RouterIf)
		conf.IfType = int32(data.IfType)
		conf.IpAddr = string(data.IpAddr)

		_, err := clnt.ClientHdl.DeleteIPv4Intf(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.IPv4Neighbor:
		data := obj.(models.IPv4Neighbor)
		conf := portdServices.NewIPv4Neighbor()
		conf.RouterIf = int32(data.RouterIf)
		conf.MacAddr = string(data.MacAddr)
		conf.IpAddr = string(data.IpAddr)
		conf.VlanId = int32(data.VlanId)

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

		origconf.PortTagType = string(origdata.PortTagType)
		origconf.Ports = string(origdata.Ports)
		origconf.VlanId = int32(origdata.VlanId)

		updateconf.PortTagType = string(updatedata.PortTagType)
		updateconf.Ports = string(updatedata.Ports)
		updateconf.VlanId = int32(updatedata.VlanId)

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

		origconf.RouterIf = int32(origdata.RouterIf)
		origconf.IfType = int32(origdata.IfType)
		origconf.IpAddr = string(origdata.IpAddr)

		updateconf.RouterIf = int32(updatedata.RouterIf)
		updateconf.IfType = int32(updatedata.IfType)
		updateconf.IpAddr = string(updatedata.IpAddr)

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

		origconf.RouterIf = int32(origdata.RouterIf)
		origconf.MacAddr = string(origdata.MacAddr)
		origconf.IpAddr = string(origdata.IpAddr)
		origconf.VlanId = int32(origdata.VlanId)

		updateconf.RouterIf = int32(updatedata.RouterIf)
		updateconf.MacAddr = string(updatedata.MacAddr)
		updateconf.IpAddr = string(updatedata.IpAddr)
		updateconf.VlanId = int32(updatedata.VlanId)

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
