package main

import (
	"database/sql"
	"fmt"
	"lacpd"
	"models"
	"utils/ipcutils"
)

type LACPDClient struct {
	ipcutils.IPCClientBase
	ClientHdl *lacpd.LACPDServicesClient
}

func (clnt *LACPDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *LACPDClient) ConnectToServer() bool {

	clnt.TTransport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	if clnt.TTransport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = lacpd.NewLACPDServicesClientFactory(clnt.TTransport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}
func (clnt *LACPDClient) IsConnectedToServer() bool {
	return clnt.IsConnected
}
func (clnt *LACPDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {

	case models.EthernetConfig:
		data := obj.(models.EthernetConfig)
		conf := lacpd.NewEthernetConfig()
		models.ConvertlacpdEthernetConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateEthernetConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.AggregationLacpConfig:
		data := obj.(models.AggregationLacpConfig)
		conf := lacpd.NewAggregationLacpConfig()
		models.ConvertlacpdAggregationLacpConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateAggregationLacpConfig(conf)
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
func (clnt *LACPDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {

	switch obj.(type) {

	case models.EthernetConfig:
		data := obj.(models.EthernetConfig)
		conf := lacpd.NewEthernetConfig()
		models.ConvertlacpdEthernetConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteEthernetConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.AggregationLacpConfig:
		data := obj.(models.AggregationLacpConfig)
		conf := lacpd.NewAggregationLacpConfig()
		models.ConvertlacpdAggregationLacpConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteAggregationLacpConfig(conf)
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
func (clnt *LACPDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {

	logger.Println("### Get Bulk request called with", currMarker, count)
	switch obj.(type) {

	case models.AggregationLacpState:

		if clnt.ClientHdl != nil {
			var ret_obj models.AggregationLacpState
			bulkInfo, err := clnt.ClientHdl.GetBulkAggregationLacpState(lacpd.Int(currMarker), lacpd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.Description = string(bulkInfo.AggregationLacpStateList[i].Description)
					ret_obj.MinLinks = uint16(bulkInfo.AggregationLacpStateList[i].MinLinks)
					ret_obj.SystemPriority = uint16(bulkInfo.AggregationLacpStateList[i].SystemPriority)
					ret_obj.NameKey = string(bulkInfo.AggregationLacpStateList[i].NameKey)
					ret_obj.Interval = int32(bulkInfo.AggregationLacpStateList[i].Interval)
					ret_obj.Enabled = bool(bulkInfo.AggregationLacpStateList[i].Enabled)
					ret_obj.Mtu = uint16(bulkInfo.AggregationLacpStateList[i].Mtu)
					ret_obj.SystemIdMac = string(bulkInfo.AggregationLacpStateList[i].SystemIdMac)
					ret_obj.LagType = int32(bulkInfo.AggregationLacpStateList[i].LagType)
					ret_obj.Ifindex = uint32(bulkInfo.AggregationLacpStateList[i].Ifindex)
					ret_obj.LagHash = int32(bulkInfo.AggregationLacpStateList[i].LagHash)
					ret_obj.Type = string(bulkInfo.AggregationLacpStateList[i].Type)
					ret_obj.LacpMode = int32(bulkInfo.AggregationLacpStateList[i].LacpMode)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.AggregationLacpMemberStateCounters:

		if clnt.ClientHdl != nil {
			var ret_obj models.AggregationLacpMemberStateCounters
			bulkInfo, err := clnt.ClientHdl.GetBulkAggregationLacpMemberStateCounters(lacpd.Int(currMarker), lacpd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.PartnerCdsChurnMachine = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].PartnerCdsChurnMachine)
					ret_obj.MinLinks = uint16(bulkInfo.AggregationLacpMemberStateCountersList[i].MinLinks)
					ret_obj.RxMachine = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].RxMachine)
					ret_obj.SystemPriority = uint16(bulkInfo.AggregationLacpMemberStateCountersList[i].SystemPriority)
					ret_obj.MuxMachine = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].MuxMachine)
					ret_obj.LampOutPdu = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LampOutPdu)
					ret_obj.LacpRxErrors = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LacpRxErrors)
					ret_obj.Distributing = bool(bulkInfo.AggregationLacpMemberStateCountersList[i].Distributing)
					ret_obj.ActorCdsChurnMachine = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].ActorCdsChurnMachine)
					ret_obj.SystemId = string(bulkInfo.AggregationLacpMemberStateCountersList[i].SystemId)
					ret_obj.LampInResponsePdu = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LampInResponsePdu)
					ret_obj.ActorCdsChurnCount = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].ActorCdsChurnCount)
					ret_obj.Type = string(bulkInfo.AggregationLacpMemberStateCountersList[i].Type)
					ret_obj.PartnerSyncTransitionCount = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].PartnerSyncTransitionCount)
					ret_obj.LampInPdu = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LampInPdu)
					ret_obj.LacpErrors = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LacpErrors)
					ret_obj.Description = string(bulkInfo.AggregationLacpMemberStateCountersList[i].Description)
					ret_obj.PartnerChangeCount = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].PartnerChangeCount)
					ret_obj.PartnerChurnCount = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].PartnerChurnCount)
					ret_obj.LacpTxErrors = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LacpTxErrors)
					ret_obj.LampOutResponsePdu = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LampOutResponsePdu)
					ret_obj.LacpOutPkts = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LacpOutPkts)
					ret_obj.Timeout = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].Timeout)
					ret_obj.Synchronization = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].Synchronization)
					ret_obj.PartnerId = string(bulkInfo.AggregationLacpMemberStateCountersList[i].PartnerId)
					ret_obj.ActorChurnMachine = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].ActorChurnMachine)
					ret_obj.NameKey = string(bulkInfo.AggregationLacpMemberStateCountersList[i].NameKey)
					ret_obj.Interface = string(bulkInfo.AggregationLacpMemberStateCountersList[i].Interface)
					ret_obj.DebugId = uint32(bulkInfo.AggregationLacpMemberStateCountersList[i].DebugId)
					ret_obj.Collecting = bool(bulkInfo.AggregationLacpMemberStateCountersList[i].Collecting)
					ret_obj.MuxReason = string(bulkInfo.AggregationLacpMemberStateCountersList[i].MuxReason)
					ret_obj.ActorChangeCount = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].ActorChangeCount)
					ret_obj.LacpUnknownErrors = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LacpUnknownErrors)
					ret_obj.Interval = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].Interval)
					ret_obj.Enabled = bool(bulkInfo.AggregationLacpMemberStateCountersList[i].Enabled)
					ret_obj.PartnerCdsChurnCount = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].PartnerCdsChurnCount)
					ret_obj.OperKey = uint16(bulkInfo.AggregationLacpMemberStateCountersList[i].OperKey)
					ret_obj.LagHash = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].LagHash)
					ret_obj.PartnerKey = uint16(bulkInfo.AggregationLacpMemberStateCountersList[i].PartnerKey)
					ret_obj.SystemIdMac = string(bulkInfo.AggregationLacpMemberStateCountersList[i].SystemIdMac)
					ret_obj.Activity = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].Activity)
					ret_obj.ActorChurnCount = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].ActorChurnCount)
					ret_obj.RxTime = uint32(bulkInfo.AggregationLacpMemberStateCountersList[i].RxTime)
					ret_obj.PartnerChurnMachine = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].PartnerChurnMachine)
					ret_obj.LacpInPkts = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LacpInPkts)
					ret_obj.Mtu = uint16(bulkInfo.AggregationLacpMemberStateCountersList[i].Mtu)
					ret_obj.ActorSyncTransitionCount = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].ActorSyncTransitionCount)
					ret_obj.LagType = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].LagType)
					ret_obj.Aggregatable = bool(bulkInfo.AggregationLacpMemberStateCountersList[i].Aggregatable)
					ret_obj.LacpMode = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].LacpMode)
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
func (clnt *LACPDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called LACPD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.EthernetConfig:
		// cast original object
		origdata := dbObj.(models.EthernetConfig)
		updatedata := obj.(models.EthernetConfig)
		// create new thrift objects
		origconf := lacpd.NewEthernetConfig()
		updateconf := lacpd.NewEthernetConfig()
		models.ConvertlacpdEthernetConfigObjToThrift(&origdata, origconf)
		models.ConvertlacpdEthernetConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateEthernetConfig(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.AggregationLacpConfig:
		// cast original object
		origdata := dbObj.(models.AggregationLacpConfig)
		updatedata := obj.(models.AggregationLacpConfig)
		// create new thrift objects
		origconf := lacpd.NewAggregationLacpConfig()
		updateconf := lacpd.NewAggregationLacpConfig()
		models.ConvertlacpdAggregationLacpConfigObjToThrift(&origdata, origconf)
		models.ConvertlacpdAggregationLacpConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateAggregationLacpConfig(origconf, updateconf, attrSet)
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
