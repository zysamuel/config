package main

import (
	"database/sql"
	"fmt"
	"models"
	"stpd"
	"utils/ipcutils"
)

type STPDClient struct {
	ipcutils.IPCClientBase
	ClientHdl *stpd.STPDServicesClient
}

func (clnt *STPDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *STPDClient) ConnectToServer() bool {

	clnt.TTransport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	if clnt.TTransport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = stpd.NewSTPDServicesClientFactory(clnt.TTransport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}
func (clnt *STPDClient) IsConnectedToServer() bool {
	return clnt.IsConnected
}
func (clnt *STPDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {

	case models.Dot1dStpPortEntryConfig:
		data := obj.(models.Dot1dStpPortEntryConfig)
		conf := stpd.NewDot1dStpPortEntryConfig()
		models.ConvertstpdDot1dStpPortEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateDot1dStpPortEntryConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.Dot1dStpBridgeConfig:
		data := obj.(models.Dot1dStpBridgeConfig)
		conf := stpd.NewDot1dStpBridgeConfig()
		models.ConvertstpdDot1dStpBridgeConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateDot1dStpBridgeConfig(conf)
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
func (clnt *STPDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {

	switch obj.(type) {

	case models.Dot1dStpPortEntryConfig:
		data := obj.(models.Dot1dStpPortEntryConfig)
		conf := stpd.NewDot1dStpPortEntryConfig()
		models.ConvertstpdDot1dStpPortEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteDot1dStpPortEntryConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.Dot1dStpBridgeConfig:
		data := obj.(models.Dot1dStpBridgeConfig)
		conf := stpd.NewDot1dStpBridgeConfig()
		models.ConvertstpdDot1dStpBridgeConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteDot1dStpBridgeConfig(conf)
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
func (clnt *STPDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {

	logger.Println("### Get Bulk request called with", currMarker, count)
	switch obj.(type) {

	case models.Dot1dStpPortEntryStateCounters:

		if clnt.ClientHdl != nil {
			var ret_obj models.Dot1dStpPortEntryStateCounters
			bulkInfo, err := clnt.ClientHdl.GetBulkDot1dStpPortEntryStateCounters(stpd.Int(currMarker), stpd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.Dot1dStpPortPriority = int32(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortPriority)
					ret_obj.Dot1dStpPortDesignatedBridge = string(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortDesignatedBridge)
					ret_obj.TcInPkts = uint64(bulkInfo.Dot1dStpPortEntryStateCountersList[i].TcInPkts)
					ret_obj.PvstOutPkts = uint64(bulkInfo.Dot1dStpPortEntryStateCountersList[i].PvstOutPkts)
					ret_obj.StpOutPkts = uint64(bulkInfo.Dot1dStpPortEntryStateCountersList[i].StpOutPkts)
					ret_obj.BpduInPkts = uint64(bulkInfo.Dot1dStpPortEntryStateCountersList[i].BpduInPkts)
					ret_obj.Dot1dStpPortProtocolMigration = int32(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortProtocolMigration)
					ret_obj.Dot1dStpPortState = int32(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortState)
					ret_obj.Dot1dStpPortEnable = int32(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortEnable)
					ret_obj.Dot1dStpPortDesignatedRoot = string(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortDesignatedRoot)
					ret_obj.Dot1dStpPortAdminPointToPoint = int32(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortAdminPointToPoint)
					ret_obj.Dot1dStpPortDesignatedCost = int32(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortDesignatedCost)
					ret_obj.Dot1dStpPortAdminPathCost = int32(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortAdminPathCost)
					ret_obj.BpduOutPkts = uint64(bulkInfo.Dot1dStpPortEntryStateCountersList[i].BpduOutPkts)
					ret_obj.Dot1dStpPortPathCost32 = int32(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortPathCost32)
					ret_obj.PvstInPkts = uint64(bulkInfo.Dot1dStpPortEntryStateCountersList[i].PvstInPkts)
					ret_obj.StpInPkts = uint64(bulkInfo.Dot1dStpPortEntryStateCountersList[i].StpInPkts)
					ret_obj.Dot1dStpPortOperPointToPoint = int32(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortOperPointToPoint)
					ret_obj.Dot1dBrgIfIndex = int32(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dBrgIfIndex)
					ret_obj.RstpInPkts = uint64(bulkInfo.Dot1dStpPortEntryStateCountersList[i].RstpInPkts)
					ret_obj.Dot1dStpPortOperEdgePort = int32(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortOperEdgePort)
					ret_obj.TcOutPkts = uint64(bulkInfo.Dot1dStpPortEntryStateCountersList[i].TcOutPkts)
					ret_obj.Dot1dStpPortDesignatedPort = string(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortDesignatedPort)
					ret_obj.Dot1dStpPortAdminEdgePort = int32(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortAdminEdgePort)
					ret_obj.Dot1dStpPortForwardTransitions = uint32(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortForwardTransitions)
					ret_obj.RstpOutPkts = uint64(bulkInfo.Dot1dStpPortEntryStateCountersList[i].RstpOutPkts)
					ret_obj.Dot1dStpPort = int32(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPort)
					ret_obj.Dot1dStpPortPathCost = int32(bulkInfo.Dot1dStpPortEntryStateCountersList[i].Dot1dStpPortPathCost)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.Dot1dStpBridgeState:

		if clnt.ClientHdl != nil {
			var ret_obj models.Dot1dStpBridgeState
			bulkInfo, err := clnt.ClientHdl.GetBulkDot1dStpBridgeState(stpd.Int(currMarker), stpd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.Dot1dBrgIfIndex = int32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dBrgIfIndex)
					ret_obj.Dot1dStpDesignatedRoot = string(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpDesignatedRoot)
					ret_obj.Dot1dStpBridgeForceVersion = int32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpBridgeForceVersion)
					ret_obj.Dot1dBridgeAddress = string(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dBridgeAddress)
					ret_obj.Dot1dStpBridgeHelloTime = int32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpBridgeHelloTime)
					ret_obj.Dot1dStpHelloTime = int32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpHelloTime)
					ret_obj.Dot1dStpPriority = int32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpPriority)
					ret_obj.Dot1dStpProtocolSpecification = int32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpProtocolSpecification)
					ret_obj.Dot1dStpForwardDelay = int32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpForwardDelay)
					ret_obj.Dot1dStpRootPort = int32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpRootPort)
					ret_obj.Dot1dStpRootCost = int32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpRootCost)
					ret_obj.Dot1dStpBridgeTxHoldCount = int32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpBridgeTxHoldCount)
					ret_obj.Dot1dStpTimeSinceTopologyChange = uint32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpTimeSinceTopologyChange)
					ret_obj.Dot1dStpMaxAge = int32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpMaxAge)
					ret_obj.Dot1dStpTopChanges = uint32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpTopChanges)
					ret_obj.Dot1dStpBridgeForwardDelay = int32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpBridgeForwardDelay)
					ret_obj.Dot1dStpBridgeMaxAge = int32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpBridgeMaxAge)
					ret_obj.Dot1dStpVlan = uint16(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpVlan)
					ret_obj.Dot1dStpHoldTime = int32(bulkInfo.Dot1dStpBridgeStateList[i].Dot1dStpHoldTime)
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
func (clnt *STPDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called STPD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.Dot1dStpPortEntryConfig:
		// cast original object
		origdata := dbObj.(models.Dot1dStpPortEntryConfig)
		updatedata := obj.(models.Dot1dStpPortEntryConfig)
		// create new thrift objects
		origconf := stpd.NewDot1dStpPortEntryConfig()
		updateconf := stpd.NewDot1dStpPortEntryConfig()
		models.ConvertstpdDot1dStpPortEntryConfigObjToThrift(&origdata, origconf)
		models.ConvertstpdDot1dStpPortEntryConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateDot1dStpPortEntryConfig(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.Dot1dStpBridgeConfig:
		// cast original object
		origdata := dbObj.(models.Dot1dStpBridgeConfig)
		updatedata := obj.(models.Dot1dStpBridgeConfig)
		// create new thrift objects
		origconf := stpd.NewDot1dStpBridgeConfig()
		updateconf := stpd.NewDot1dStpBridgeConfig()
		models.ConvertstpdDot1dStpBridgeConfigObjToThrift(&origdata, origconf)
		models.ConvertstpdDot1dStpBridgeConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateDot1dStpBridgeConfig(origconf, updateconf, attrSet)
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
