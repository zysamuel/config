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

	case models.VlanConfig:
		data := obj.(models.VlanConfig)
		conf := asicdServices.NewVlanConfig()
		models.ConvertasicdVlanConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateVlanConfig(conf)
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

	case models.VlanConfig:
		data := obj.(models.VlanConfig)
		conf := asicdServices.NewVlanConfig()
		models.ConvertasicdVlanConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteVlanConfig(conf)
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

	logger.Println("### Get Bulk request called with", currMarker, count)
	switch obj.(type) {

	case models.VlanState:

		if clnt.ClientHdl != nil {
			var ret_obj models.VlanState
			bulkInfo, err := clnt.ClientHdl.GetBulkVlanState(asicdServices.Int(currMarker), asicdServices.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.IfIndex = int32(bulkInfo.VlanStateList[i].IfIndex)
					ret_obj.VlanName = string(bulkInfo.VlanStateList[i].VlanName)
					ret_obj.OperState = string(bulkInfo.VlanStateList[i].OperState)
					ret_obj.VlanId = int32(bulkInfo.VlanStateList[i].VlanId)
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

					ret_obj.IfIndex = int32(bulkInfo.PortStateList[i].IfIndex)
					for _, data := range bulkInfo.PortStateList[i].PortStats {
						ret_obj.PortStats = int64(data)
					}

					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.OspfHostEntryState:

		if clnt.ClientHdl != nil {
			var ret_obj models.OspfHostEntryState
			bulkInfo, err := clnt.ClientHdl.GetBulkOspfHostEntryState(asicdServices.Int(currMarker), asicdServices.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.HostIpAddressKey = string(bulkInfo.OspfHostEntryStateList[i].HostIpAddressKey)
					ret_obj.HostAreaID = string(bulkInfo.OspfHostEntryStateList[i].HostAreaID)
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

	case models.VlanConfig:
		// cast original object
		origdata := dbObj.(models.VlanConfig)
		updatedata := obj.(models.VlanConfig)
		// create new thrift objects
		origconf := asicdServices.NewVlanConfig()
		updateconf := asicdServices.NewVlanConfig()
		models.ConvertasicdVlanConfigObjToThrift(&origdata, origconf)
		models.ConvertasicdVlanConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateVlanConfig(origconf, updateconf, attrSet)
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
