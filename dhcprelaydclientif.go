package main

import (
	"database/sql"
	"dhcprelayd"
	"fmt"
	"models"
	"utils/ipcutils"
)

type DHCPRELAYDClient struct {
	ipcutils.IPCClientBase
	ClientHdl *dhcprelayd.DHCPRELAYDServicesClient
}

func (clnt *DHCPRELAYDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *DHCPRELAYDClient) ConnectToServer() bool {

	clnt.TTransport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	if clnt.TTransport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = dhcprelayd.NewDHCPRELAYDServicesClientFactory(clnt.TTransport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}
func (clnt *DHCPRELAYDClient) IsConnectedToServer() bool {
	return clnt.IsConnected
}
func (clnt *DHCPRELAYDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {

	case models.DhcpRelayGlobalConfig:
		data := obj.(models.DhcpRelayGlobalConfig)
		conf := dhcprelayd.NewDhcpRelayGlobalConfig()
		models.ConvertdhcprelaydDhcpRelayGlobalConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateDhcpRelayGlobalConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.DhcpRelayIntfConfig:
		data := obj.(models.DhcpRelayIntfConfig)
		conf := dhcprelayd.NewDhcpRelayIntfConfig()
		models.ConvertdhcprelaydDhcpRelayIntfConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateDhcpRelayIntfConfig(conf)
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
func (clnt *DHCPRELAYDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {

	switch obj.(type) {

	case models.DhcpRelayGlobalConfig:
		data := obj.(models.DhcpRelayGlobalConfig)
		conf := dhcprelayd.NewDhcpRelayGlobalConfig()
		models.ConvertdhcprelaydDhcpRelayGlobalConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteDhcpRelayGlobalConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.DhcpRelayIntfConfig:
		data := obj.(models.DhcpRelayIntfConfig)
		conf := dhcprelayd.NewDhcpRelayIntfConfig()
		models.ConvertdhcprelaydDhcpRelayIntfConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteDhcpRelayIntfConfig(conf)
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
func (clnt *DHCPRELAYDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {

	logger.Println("### Get Bulk request called with", currMarker, count)
	switch obj.(type) {

	case models.DhcpRelayHostDhcpState:

		if clnt.ClientHdl != nil {
			var ret_obj models.DhcpRelayHostDhcpState
			bulkInfo, err := clnt.ClientHdl.GetBulkDhcpRelayHostDhcpState(dhcprelayd.Int(currMarker), dhcprelayd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.ClientResponse = string(bulkInfo.DhcpRelayHostDhcpStateList[i].ClientResponse)
					ret_obj.ServerRequest = string(bulkInfo.DhcpRelayHostDhcpStateList[i].ServerRequest)
					ret_obj.OfferedIp = string(bulkInfo.DhcpRelayHostDhcpStateList[i].OfferedIp)
					ret_obj.ServerResponse = string(bulkInfo.DhcpRelayHostDhcpStateList[i].ServerResponse)
					ret_obj.MacAddr = string(bulkInfo.DhcpRelayHostDhcpStateList[i].MacAddr)
					ret_obj.LeaseDuration = string(bulkInfo.DhcpRelayHostDhcpStateList[i].LeaseDuration)
					ret_obj.GatewayIp = string(bulkInfo.DhcpRelayHostDhcpStateList[i].GatewayIp)
					ret_obj.AcceptedIp = string(bulkInfo.DhcpRelayHostDhcpStateList[i].AcceptedIp)
					ret_obj.ServerIp = string(bulkInfo.DhcpRelayHostDhcpStateList[i].ServerIp)
					ret_obj.ClientRequest = string(bulkInfo.DhcpRelayHostDhcpStateList[i].ClientRequest)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.DhcpRelayIntfState:

		if clnt.ClientHdl != nil {
			var ret_obj models.DhcpRelayIntfState
			bulkInfo, err := clnt.ClientHdl.GetBulkDhcpRelayIntfState(dhcprelayd.Int(currMarker), dhcprelayd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.TotalDhcpServerRx = int32(bulkInfo.DhcpRelayIntfStateList[i].TotalDhcpServerRx)
					ret_obj.TotalDhcpServerTx = int32(bulkInfo.DhcpRelayIntfStateList[i].TotalDhcpServerTx)
					ret_obj.IntfId = int32(bulkInfo.DhcpRelayIntfStateList[i].IntfId)
					ret_obj.TotalDrops = int32(bulkInfo.DhcpRelayIntfStateList[i].TotalDrops)
					ret_obj.TotalDhcpClientRx = int32(bulkInfo.DhcpRelayIntfStateList[i].TotalDhcpClientRx)
					ret_obj.TotalDhcpClientTx = int32(bulkInfo.DhcpRelayIntfStateList[i].TotalDhcpClientTx)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.DhcpRelayIntfServerState:

		if clnt.ClientHdl != nil {
			var ret_obj models.DhcpRelayIntfServerState
			bulkInfo, err := clnt.ClientHdl.GetBulkDhcpRelayIntfServerState(dhcprelayd.Int(currMarker), dhcprelayd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.Request = int32(bulkInfo.DhcpRelayIntfServerStateList[i].Request)
					ret_obj.ServerIp = string(bulkInfo.DhcpRelayIntfServerStateList[i].ServerIp)
					ret_obj.IntfId = int32(bulkInfo.DhcpRelayIntfServerStateList[i].IntfId)
					ret_obj.Responses = int32(bulkInfo.DhcpRelayIntfServerStateList[i].Responses)
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
func (clnt *DHCPRELAYDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called DHCPRELAYD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.DhcpRelayGlobalConfig:
		// cast original object
		origdata := dbObj.(models.DhcpRelayGlobalConfig)
		updatedata := obj.(models.DhcpRelayGlobalConfig)
		// create new thrift objects
		origconf := dhcprelayd.NewDhcpRelayGlobalConfig()
		updateconf := dhcprelayd.NewDhcpRelayGlobalConfig()
		models.ConvertdhcprelaydDhcpRelayGlobalConfigObjToThrift(&origdata, origconf)
		models.ConvertdhcprelaydDhcpRelayGlobalConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateDhcpRelayGlobalConfig(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.DhcpRelayIntfConfig:
		// cast original object
		origdata := dbObj.(models.DhcpRelayIntfConfig)
		updatedata := obj.(models.DhcpRelayIntfConfig)
		// create new thrift objects
		origconf := dhcprelayd.NewDhcpRelayIntfConfig()
		updateconf := dhcprelayd.NewDhcpRelayIntfConfig()
		models.ConvertdhcprelaydDhcpRelayIntfConfigObjToThrift(&origdata, origconf)
		models.ConvertdhcprelaydDhcpRelayIntfConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateDhcpRelayIntfConfig(origconf, updateconf, attrSet)
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
