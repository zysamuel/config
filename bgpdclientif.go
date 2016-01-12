package main

import (
	"bgpd"
	"database/sql"
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

	clnt.Transport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = bgpd.NewBGPDServicesClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
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
			bulkInfo, _ := clnt.ClientHdl.GetBulkBGPGlobalState(bgpd.Int(currMarker), bgpd.Int(count))
			if bulkInfo.Count != 0 {
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
					ret_obj.TotalPrefixes = uint32(bulkInfo.BGPGlobalStateList[i].TotalPrefixes)
					objs = append(objs, ret_obj)
				}
			}
		}
		break

	case models.BGPNeighborState:

		if clnt.ClientHdl != nil {
			var ret_obj models.BGPNeighborState
			bulkInfo, _ := clnt.ClientHdl.GetBulkBGPNeighborState(bgpd.Int(currMarker), bgpd.Int(count))
			if bulkInfo.Count != 0 {
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
					ret_obj.SessionState = uint32(bulkInfo.BGPNeighborStateList[i].SessionState)
					ret_obj.NeighborAddress = string(bulkInfo.BGPNeighborStateList[i].NeighborAddress)
					ret_obj.PeerAS = uint32(bulkInfo.BGPNeighborStateList[i].PeerAS)
					ret_obj.LocalAS = uint32(bulkInfo.BGPNeighborStateList[i].LocalAS)
					ret_obj.AuthPassword = string(bulkInfo.BGPNeighborStateList[i].AuthPassword)
					objs = append(objs, ret_obj)
				}
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

	default:
		break
	}
	return ok

}
