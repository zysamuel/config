package main

import (
	"bgpd"
	"database/sql"
	"fmt"
	"models"
	"utils/ipcutils"
)

type BGPDClient struct {
	ipcutils.IPCClientBase
	ClientHdl *bgpd.BGPDServicesClient
}

func (clnt *BGPDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *BGPDClient) ConnectToServer() bool {

	clnt.TTransport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	if clnt.TTransport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = bgpd.NewBGPDServicesClientFactory(clnt.TTransport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}
func (clnt *BGPDClient) IsConnectedToServer() bool {
	return clnt.IsConnected
}
func (clnt *BGPDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {

	case models.BGPGlobalConfig:
		data := obj.(models.BGPGlobalConfig)
		conf := bgpd.NewBGPGlobalConfig()
		models.ConvertbgpdBGPGlobalConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateBGPGlobalConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.BGPNeighborConfig:
		data := obj.(models.BGPNeighborConfig)
		conf := bgpd.NewBGPNeighborConfig()
		models.ConvertbgpdBGPNeighborConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateBGPNeighborConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.BGPPeerGroup:
		data := obj.(models.BGPPeerGroup)
		conf := bgpd.NewBGPPeerGroup()
		models.ConvertbgpdBGPPeerGroupObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateBGPPeerGroup(conf)
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
func (clnt *BGPDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {

	switch obj.(type) {

	case models.BGPGlobalConfig:
		data := obj.(models.BGPGlobalConfig)
		conf := bgpd.NewBGPGlobalConfig()
		models.ConvertbgpdBGPGlobalConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteBGPGlobalConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.BGPNeighborConfig:
		data := obj.(models.BGPNeighborConfig)
		conf := bgpd.NewBGPNeighborConfig()
		models.ConvertbgpdBGPNeighborConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteBGPNeighborConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.BGPPeerGroup:
		data := obj.(models.BGPPeerGroup)
		conf := bgpd.NewBGPPeerGroup()
		models.ConvertbgpdBGPPeerGroupObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteBGPPeerGroup(conf)
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
func (clnt *BGPDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {

	logger.Println("### Get Bulk request called with", currMarker, count)
	switch obj.(type) {

	case models.BGPGlobalState:

		if clnt.ClientHdl != nil {
			var ret_obj models.BGPGlobalState
			bulkInfo, err := clnt.ClientHdl.GetBulkBGPGlobalState(bgpd.Int(currMarker), bgpd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.RouterId = string(bulkInfo.BGPGlobalStateList[i].RouterId)
					ret_obj.TotalPaths = uint32(bulkInfo.BGPGlobalStateList[i].TotalPaths)
					ret_obj.AS = uint32(bulkInfo.BGPGlobalStateList[i].AS)
					ret_obj.IBGPMaxPaths = uint32(bulkInfo.BGPGlobalStateList[i].IBGPMaxPaths)
					ret_obj.EBGPMaxPaths = uint32(bulkInfo.BGPGlobalStateList[i].EBGPMaxPaths)
					ret_obj.UseMultiplePaths = bool(bulkInfo.BGPGlobalStateList[i].UseMultiplePaths)
					ret_obj.EBGPAllowMultipleAS = bool(bulkInfo.BGPGlobalStateList[i].EBGPAllowMultipleAS)
					ret_obj.TotalPrefixes = uint32(bulkInfo.BGPGlobalStateList[i].TotalPrefixes)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.BGPNeighborState:

		if clnt.ClientHdl != nil {
			var ret_obj models.BGPNeighborState
			bulkInfo, err := clnt.ClientHdl.GetBulkBGPNeighborState(bgpd.Int(currMarker), bgpd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.RouteReflectorClusterId = uint32(bulkInfo.BGPNeighborStateList[i].RouteReflectorClusterId)
					ret_obj.RouteReflectorClient = bool(bulkInfo.BGPNeighborStateList[i].RouteReflectorClient)
					ret_obj.Description = string(bulkInfo.BGPNeighborStateList[i].Description)
					ret_obj.MultiHopTTL = uint8(bulkInfo.BGPNeighborStateList[i].MultiHopTTL)
					ret_obj.PeerAS = uint32(bulkInfo.BGPNeighborStateList[i].PeerAS)
					ret_obj.AddPathsRx = bool(bulkInfo.BGPNeighborStateList[i].AddPathsRx)
					ret_obj.KeepaliveTime = uint32(bulkInfo.BGPNeighborStateList[i].KeepaliveTime)
					ret_obj.AuthPassword = string(bulkInfo.BGPNeighborStateList[i].AuthPassword)
					ret_obj.AddPathsMaxTx = uint8(bulkInfo.BGPNeighborStateList[i].AddPathsMaxTx)
					ret_obj.MultiHopEnable = bool(bulkInfo.BGPNeighborStateList[i].MultiHopEnable)
					ret_obj.SessionState = uint32(bulkInfo.BGPNeighborStateList[i].SessionState)
					ret_obj.NeighborAddress = string(bulkInfo.BGPNeighborStateList[i].NeighborAddress)
					ret_obj.HoldTime = uint32(bulkInfo.BGPNeighborStateList[i].HoldTime)
					ret_obj.LocalAS = uint32(bulkInfo.BGPNeighborStateList[i].LocalAS)
					ret_obj.ConnectRetryTime = uint32(bulkInfo.BGPNeighborStateList[i].ConnectRetryTime)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.BGPRoute:

		if clnt.ClientHdl != nil {
			var ret_obj models.BGPRoute
			bulkInfo, err := clnt.ClientHdl.GetBulkBGPRoute(bgpd.Int(currMarker), bgpd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.Updated = string(bulkInfo.BGPRouteList[i].Updated)
					ret_obj.NextHop = string(bulkInfo.BGPRouteList[i].NextHop)
					ret_obj.Network = string(bulkInfo.BGPRouteList[i].Network)
					for _, data := range bulkInfo.BGPRouteList[i].Path {
						ret_obj.Path = uint32(data)
					}

					ret_obj.Metric = uint32(bulkInfo.BGPRouteList[i].Metric)
					ret_obj.LocalPref = uint32(bulkInfo.BGPRouteList[i].LocalPref)
					ret_obj.Mask = string(bulkInfo.BGPRouteList[i].Mask)
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
func (clnt *BGPDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called BGPD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.BGPGlobalConfig:
		// cast original object
		origdata := dbObj.(models.BGPGlobalConfig)
		updatedata := obj.(models.BGPGlobalConfig)
		// create new thrift objects
		origconf := bgpd.NewBGPGlobalConfig()
		updateconf := bgpd.NewBGPGlobalConfig()
		models.ConvertbgpdBGPGlobalConfigObjToThrift(&origdata, origconf)
		models.ConvertbgpdBGPGlobalConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateBGPGlobalConfig(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.BGPNeighborConfig:
		// cast original object
		origdata := dbObj.(models.BGPNeighborConfig)
		updatedata := obj.(models.BGPNeighborConfig)
		// create new thrift objects
		origconf := bgpd.NewBGPNeighborConfig()
		updateconf := bgpd.NewBGPNeighborConfig()
		models.ConvertbgpdBGPNeighborConfigObjToThrift(&origdata, origconf)
		models.ConvertbgpdBGPNeighborConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateBGPNeighborConfig(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.BGPPeerGroup:
		// cast original object
		origdata := dbObj.(models.BGPPeerGroup)
		updatedata := obj.(models.BGPPeerGroup)
		// create new thrift objects
		origconf := bgpd.NewBGPPeerGroup()
		updateconf := bgpd.NewBGPPeerGroup()
		models.ConvertbgpdBGPPeerGroupObjToThrift(&origdata, origconf)
		models.ConvertbgpdBGPPeerGroupObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateBGPPeerGroup(origconf, updateconf, attrSet)
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
