package main

import (
	"asicdServices"
	"database/sql"
	"models"
)

type ASICDClient struct {
	IPCClientBase
	ClientHdl *asicdServices.ASICDServicesClient
}

func (clnt *ASICDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *ASICDClient) ConnectToServer() bool {

	clnt.Transport, clnt.PtrProtocolFactory = CreateIPCHandles(clnt.Address)
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = asicdServices.NewASICDServicesClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
	}
	return true
}
func (clnt *ASICDClient) IsConnectedToServer() bool {
	return true
}
func (clnt *ASICDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {

	case models.PortIntfConfig:
		data := obj.(models.PortIntfConfig)
		conf := asicdServices.NewPortIntfConfig()
		conf.OperState = string(data.OperState)
		conf.MacAddr = string(data.MacAddr)
		conf.PortNum = int32(data.PortNum)
		conf.Name = string(data.Name)
		conf.Duplex = string(data.Duplex)
		conf.Type = string(data.Type)
		conf.MediaType = string(data.MediaType)
		conf.Mtu = int32(data.Mtu)
		conf.AdminState = string(data.AdminState)
		conf.Autoneg = string(data.Autoneg)
		conf.Speed = int32(data.Speed)
		conf.Description = string(data.Description)

		_, err := clnt.ClientHdl.CreatePortIntfConfig(conf)
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
func (clnt *ASICDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {

	switch obj.(type) {

	case models.PortIntfConfig:
		data := obj.(models.PortIntfConfig)
		conf := asicdServices.NewPortIntfConfig()
		conf.OperState = string(data.OperState)
		conf.MacAddr = string(data.MacAddr)
		conf.PortNum = int32(data.PortNum)
		conf.Name = string(data.Name)
		conf.Duplex = string(data.Duplex)
		conf.Type = string(data.Type)
		conf.MediaType = string(data.MediaType)
		conf.Mtu = int32(data.Mtu)
		conf.AdminState = string(data.AdminState)
		conf.Autoneg = string(data.Autoneg)
		conf.Speed = int32(data.Speed)
		conf.Description = string(data.Description)

		_, err := clnt.ClientHdl.DeletePortIntfConfig(conf)
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
func (clnt *ASICDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
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
func (clnt *ASICDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []byte, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called ASICD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.IPV4Route:
		// cast original object
		origdata := dbObj.(models.IPV4Route)
		updatedata := obj.(models.IPV4Route)
		// create new thrift objects
		origconf := asicdServices.NewIPV4Route()
		updateconf := asicdServices.NewIPV4Route()

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

	case models.IPv4Intf:
		// cast original object
		origdata := dbObj.(models.IPv4Intf)
		updatedata := obj.(models.IPv4Intf)
		// create new thrift objects
		origconf := asicdServices.NewIPv4Intf()
		updateconf := asicdServices.NewIPv4Intf()

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
		origconf := asicdServices.NewIPv4Neighbor()
		updateconf := asicdServices.NewIPv4Neighbor()

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

	case models.PortIntfConfig:
		// cast original object
		origdata := dbObj.(models.PortIntfConfig)
		updatedata := obj.(models.PortIntfConfig)
		// create new thrift objects
		origconf := asicdServices.NewPortIntfConfig()
		updateconf := asicdServices.NewPortIntfConfig()

		origconf.OperState = string(origdata.OperState)
		origconf.MacAddr = string(origdata.MacAddr)
		origconf.PortNum = int32(origdata.PortNum)
		origconf.Name = string(origdata.Name)
		origconf.Duplex = string(origdata.Duplex)
		origconf.Type = string(origdata.Type)
		origconf.MediaType = string(origdata.MediaType)
		origconf.Mtu = int32(origdata.Mtu)
		origconf.AdminState = string(origdata.AdminState)
		origconf.Autoneg = string(origdata.Autoneg)
		origconf.Speed = int32(origdata.Speed)
		origconf.Description = string(origdata.Description)

		updateconf.OperState = string(updatedata.OperState)
		updateconf.MacAddr = string(updatedata.MacAddr)
		updateconf.PortNum = int32(updatedata.PortNum)
		updateconf.Name = string(updatedata.Name)
		updateconf.Duplex = string(updatedata.Duplex)
		updateconf.Type = string(updatedata.Type)
		updateconf.MediaType = string(updatedata.MediaType)
		updateconf.Mtu = int32(updatedata.Mtu)
		updateconf.AdminState = string(updatedata.AdminState)
		updateconf.Autoneg = string(updatedata.Autoneg)
		updateconf.Speed = int32(updatedata.Speed)
		updateconf.Description = string(updatedata.Description)

		//convert attrSet to uint8 list
		newattrset := make([]int8, len(attrSet))
		for i, v := range attrSet {
			newattrset[i] = int8(v)
		}
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdatePortIntfConfig(origconf, updateconf, newattrset)
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
