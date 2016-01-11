package main

import (
	"bgpdServices"
	"database/sql"
	"models"
	"utils/ipcutils"
)

type BGPDClient struct {
	IPCClientBase
	ClientHdl *bgpdServices.BGPDServicesClient
}

func (clnt *BGPDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *BGPDClient) ConnectToServer() bool {

	clnt.Transport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = bgpdServices.NewBGPDServicesClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
	}
	return true
}
func (clnt *BGPDClient) IsConnectedToServer() bool {
	return true
}
func (clnt *BGPDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {

	case models.IPV4Route:
		data := obj.(models.IPV4Route)
		conf := bgpdServices.NewIPV4Route()
		conf.DestinationNw = string(data.DestinationNw)
		conf.OutgoingIntfType = string(data.OutgoingIntfType)
		conf.Protocol = string(data.Protocol)
		conf.OutgoingInterface = string(data.OutgoingInterface)
		conf.NetworkMask = string(data.NetworkMask)
		conf.NextHopIp = string(data.NextHopIp)

		_, err := clnt.ClientHdl.CreateIPV4Route(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.BGPGlobalConfig:
		data := obj.(models.BGPGlobalConfig)
		conf := bgpdServices.NewBGPGlobalConfig()
		conf.RouterId = string(data.RouterId)
		conf.ASNum = int32(data.ASNum)

		_, err := clnt.ClientHdl.CreateBGPGlobalConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.BGPNeighborConfig:
		data := obj.(models.BGPNeighborConfig)
		conf := bgpdServices.NewBGPNeighborConfig()
		conf.RouteReflectorClusterId = int32(data.RouteReflectorClusterId)
		conf.RouteReflectorClient = bool(data.RouteReflectorClient)
		conf.Description = string(data.Description)
		conf.NeighborAddress = string(data.NeighborAddress)
		conf.PeerAS = int32(data.PeerAS)
		conf.LocalAS = int32(data.LocalAS)
		conf.AuthPassword = string(data.AuthPassword)

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

	case models.IPV4Route:
		data := obj.(models.IPV4Route)
		conf := bgpdServices.NewIPV4Route()
		conf.DestinationNw = string(data.DestinationNw)
		conf.OutgoingIntfType = string(data.OutgoingIntfType)
		conf.Protocol = string(data.Protocol)
		conf.OutgoingInterface = string(data.OutgoingInterface)
		conf.NetworkMask = string(data.NetworkMask)
		conf.NextHopIp = string(data.NextHopIp)

		_, err := clnt.ClientHdl.DeleteIPV4Route(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.BGPGlobalConfig:
		data := obj.(models.BGPGlobalConfig)
		conf := bgpdServices.NewBGPGlobalConfig()
		conf.RouterId = string(data.RouterId)
		conf.ASNum = int32(data.ASNum)

		_, err := clnt.ClientHdl.DeleteBGPGlobalConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.BGPNeighborConfig:
		data := obj.(models.BGPNeighborConfig)
		conf := bgpdServices.NewBGPNeighborConfig()
		conf.RouteReflectorClusterId = int32(data.RouteReflectorClusterId)
		conf.RouteReflectorClient = bool(data.RouteReflectorClient)
		conf.Description = string(data.Description)
		conf.NeighborAddress = string(data.NeighborAddress)
		conf.PeerAS = int32(data.PeerAS)
		conf.LocalAS = int32(data.LocalAS)
		conf.AuthPassword = string(data.AuthPassword)

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

	case models.IPV4Route:

		if clnt.ClientHdl != nil {
			var ret_obj models.IPV4Route
			bulkInfo, _ := clnt.ClientHdl.GetBulkIPV4Route(bgpdServices.Int(currMarker), bgpdServices.Int(count))
			if bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}
					ret_obj.DestinationNw = string(bulkInfo.IPV4RouteList[i].DestinationNw)
					ret_obj.OutgoingIntfType = string(bulkInfo.IPV4RouteList[i].OutgoingIntfType)
					ret_obj.Protocol = string(bulkInfo.IPV4RouteList[i].Protocol)
					ret_obj.OutgoingInterface = string(bulkInfo.IPV4RouteList[i].OutgoingInterface)
					ret_obj.NetworkMask = string(bulkInfo.IPV4RouteList[i].NetworkMask)
					ret_obj.NextHopIp = string(bulkInfo.IPV4RouteList[i].NextHopIp)
					objs = append(objs, ret_obj)
				}
			}
		}
		break

	case models.BGPGlobalState:

		if clnt.ClientHdl != nil {
			var ret_obj models.BGPGlobalState
			bulkInfo, _ := clnt.ClientHdl.GetBulkBGPGlobalState(bgpdServices.Int(currMarker), bgpdServices.Int(count))
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
			bulkInfo, _ := clnt.ClientHdl.GetBulkBGPNeighborState(bgpdServices.Int(currMarker), bgpdServices.Int(count))
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

	case models.IPV4Route:
		// cast original object
		origdata := dbObj.(models.IPV4Route)
		updatedata := obj.(models.IPV4Route)
		// create new thrift objects
		origconf := bgpdServices.NewIPV4Route()
		updateconf := bgpdServices.NewIPV4Route()

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

		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateIPV4Route(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.BGPGlobalConfig:
		// cast original object
		origdata := dbObj.(models.BGPGlobalConfig)
		updatedata := obj.(models.BGPGlobalConfig)
		// create new thrift objects
		origconf := bgpdServices.NewBGPGlobalConfig()
		updateconf := bgpdServices.NewBGPGlobalConfig()

		origconf.RouterId = string(origdata.RouterId)
		origconf.ASNum = int32(origdata.ASNum)

		updateconf.RouterId = string(updatedata.RouterId)
		updateconf.ASNum = int32(updatedata.ASNum)

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
		origconf := bgpdServices.NewBGPNeighborConfig()
		updateconf := bgpdServices.NewBGPNeighborConfig()

		origconf.RouteReflectorClusterId = int32(origdata.RouteReflectorClusterId)
		origconf.RouteReflectorClient = bool(origdata.RouteReflectorClient)
		origconf.Description = string(origdata.Description)
		origconf.NeighborAddress = string(origdata.NeighborAddress)
		origconf.PeerAS = int32(origdata.PeerAS)
		origconf.LocalAS = int32(origdata.LocalAS)
		origconf.AuthPassword = string(origdata.AuthPassword)

		updateconf.RouteReflectorClusterId = int32(updatedata.RouteReflectorClusterId)
		updateconf.RouteReflectorClient = bool(updatedata.RouteReflectorClient)
		updateconf.Description = string(updatedata.Description)
		updateconf.NeighborAddress = string(updatedata.NeighborAddress)
		updateconf.PeerAS = int32(updatedata.PeerAS)
		updateconf.LocalAS = int32(updatedata.LocalAS)
		updateconf.AuthPassword = string(updatedata.AuthPassword)

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
