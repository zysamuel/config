package main

import (
	"asicdServices"
	"database/sql"
	"fmt"
	"models"
	"utils/ipcutils"
)

type ASICDClient struct {
	ipcutils.IPCClientBase
	ClientHdl *asicdServices.ASICDServicesClient
}

func (clnt *ASICDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *ASICDClient) ConnectToServer() bool {

	clnt.TTransport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	if clnt.TTransport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = asicdServices.NewASICDServicesClientFactory(clnt.TTransport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}
func (clnt *ASICDClient) IsConnectedToServer() bool {
	return clnt.IsConnected
}
func (clnt *ASICDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {

	case models.Vlan:
		data := obj.(models.Vlan)
		conf := asicdServices.NewVlan()
		models.ConvertasicdVlanObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateVlan(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.IPv4Intf:
		data := obj.(models.IPv4Intf)
		conf := asicdServices.NewIPv4Intf()
		models.ConvertasicdIPv4IntfObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateIPv4Intf(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.PortConfig:
		data := obj.(models.PortConfig)
		conf := asicdServices.NewPortConfig()
		models.ConvertasicdPortConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreatePortConfig(conf)
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

	case models.Vlan:
		data := obj.(models.Vlan)
		conf := asicdServices.NewVlan()
		models.ConvertasicdVlanObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteVlan(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.IPv4Intf:
		data := obj.(models.IPv4Intf)
		conf := asicdServices.NewIPv4Intf()
		models.ConvertasicdIPv4IntfObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteIPv4Intf(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.PortConfig:
		data := obj.(models.PortConfig)
		conf := asicdServices.NewPortConfig()
		models.ConvertasicdPortConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeletePortConfig(conf)
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

	logger.Println("### Get Bulk request called with", currMarker, count, obj)
	switch obj.(type) {

	case models.Vlan:

		if clnt.ClientHdl != nil {
			var ret_obj models.Vlan
			bulkInfo, err := clnt.ClientHdl.GetBulkVlan(asicdServices.Int(currMarker), asicdServices.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.IfIndexList = string(bulkInfo.VlanList[i].IfIndexList)
					ret_obj.VlanName = string(bulkInfo.VlanList[i].VlanName)
					ret_obj.UntagIfIndexList = string(bulkInfo.VlanList[i].UntagIfIndexList)
					ret_obj.IfIndex = int32(bulkInfo.VlanList[i].IfIndex)
					ret_obj.OperState = string(bulkInfo.VlanList[i].OperState)
					ret_obj.VlanId = int32(bulkInfo.VlanList[i].VlanId)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.IPv4Intf:

		if clnt.ClientHdl != nil {
			var ret_obj models.IPv4Intf
			bulkInfo, err := clnt.ClientHdl.GetBulkIPv4Intf(asicdServices.Int(currMarker), asicdServices.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.IfIndex = int32(bulkInfo.IPv4IntfList[i].IfIndex)
					ret_obj.IpAddr = string(bulkInfo.IPv4IntfList[i].IpAddr)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.PortConfig:

		if clnt.ClientHdl != nil {
			var ret_obj models.PortConfig
			bulkInfo, err := clnt.ClientHdl.GetBulkPortConfig(asicdServices.Int(currMarker), asicdServices.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.PhyIntfType = string(bulkInfo.PortConfigList[i].PhyIntfType)
					ret_obj.AdminState = string(bulkInfo.PortConfigList[i].AdminState)
					ret_obj.MacAddr = string(bulkInfo.PortConfigList[i].MacAddr)
					ret_obj.PortNum = int32(bulkInfo.PortConfigList[i].PortNum)
					ret_obj.Description = string(bulkInfo.PortConfigList[i].Description)
					ret_obj.Duplex = string(bulkInfo.PortConfigList[i].Duplex)
					ret_obj.Autoneg = string(bulkInfo.PortConfigList[i].Autoneg)
					ret_obj.Speed = int32(bulkInfo.PortConfigList[i].Speed)
					ret_obj.MediaType = string(bulkInfo.PortConfigList[i].MediaType)
					ret_obj.Mtu = int32(bulkInfo.PortConfigList[i].Mtu)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.PortState:

		if clnt.ClientHdl != nil {
			var ret_obj models.PortState
			bulkInfo, err := clnt.ClientHdl.GetBulkPortState(asicdServices.Int(currMarker), asicdServices.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.IfInDiscards = int64(bulkInfo.PortStateList[i].IfInDiscards)
					ret_obj.OperState = string(bulkInfo.PortStateList[i].OperState)
					ret_obj.IfInErrors = int64(bulkInfo.PortStateList[i].IfInErrors)
					ret_obj.PortNum = int32(bulkInfo.PortStateList[i].PortNum)
					ret_obj.Name = string(bulkInfo.PortStateList[i].Name)
					ret_obj.IfInOctets = int64(bulkInfo.PortStateList[i].IfInOctets)
					ret_obj.IfInUcastPkts = int64(bulkInfo.PortStateList[i].IfInUcastPkts)
					ret_obj.IfOutUcastPkts = int64(bulkInfo.PortStateList[i].IfOutUcastPkts)
					ret_obj.IfOutOctets = int64(bulkInfo.PortStateList[i].IfOutOctets)
					ret_obj.IfOutErrors = int64(bulkInfo.PortStateList[i].IfOutErrors)
					ret_obj.IfInUnknownProtos = int64(bulkInfo.PortStateList[i].IfInUnknownProtos)
					ret_obj.IfIndex = int32(bulkInfo.PortStateList[i].IfIndex)
					ret_obj.IfOutDiscards = int64(bulkInfo.PortStateList[i].IfOutDiscards)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	default:
		break
	}
	return nil, objCount, nextMarker, more, objs

}
func (clnt *ASICDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called ASICD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.Vlan:
		// cast original object
		origdata := dbObj.(models.Vlan)
		updatedata := obj.(models.Vlan)
		// create new thrift objects
		origconf := asicdServices.NewVlan()
		updateconf := asicdServices.NewVlan()
		models.ConvertasicdVlanObjToThrift(&origdata, origconf)
		models.ConvertasicdVlanObjToThrift(&updatedata, updateconf)
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
		origconf := asicdServices.NewIPv4Intf()
		updateconf := asicdServices.NewIPv4Intf()
		models.ConvertasicdIPv4IntfObjToThrift(&origdata, origconf)
		models.ConvertasicdIPv4IntfObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateIPv4Intf(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.PortConfig:
		// cast original object
		origdata := dbObj.(models.PortConfig)
		updatedata := obj.(models.PortConfig)
		// create new thrift objects
		origconf := asicdServices.NewPortConfig()
		updateconf := asicdServices.NewPortConfig()
		models.ConvertasicdPortConfigObjToThrift(&origdata, origconf)
		models.ConvertasicdPortConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdatePortConfig(origconf, updateconf, attrSet)
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
