package main

import (
	"database/sql"
	"lacpdServices"
	"models"
)

type LACPDClient struct {
	IPCClientBase
	ClientHdl *lacpdServices.LACPDServicesClient
}

func (clnt *LACPDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *LACPDClient) ConnectToServer() bool {

	clnt.Transport, clnt.PtrProtocolFactory = CreateIPCHandles(clnt.Address)
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = lacpdServices.NewLACPDServicesClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
	}
	return true
}
func (clnt *LACPDClient) IsConnectedToServer() bool {
	return true
}
func (clnt *LACPDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {

	case models.EthernetConfig:
		data := obj.(models.EthernetConfig)
		conf := lacpdServices.NewEthernetConfig()
		conf.MacAddress = string(data.MacAddress)
		conf.Description = string(data.Description)
		conf.AggregateId = string(data.AggregateId)
		conf.NameKey = string(data.NameKey)
		conf.Enabled = bool(data.Enabled)
		conf.Speed = string(data.Speed)
		conf.Mtu = int16(data.Mtu)
		conf.DuplexMode = int32(data.DuplexMode)
		conf.EnableFlowControl = bool(data.EnableFlowControl)
		conf.Auto = bool(data.Auto)
		conf.Type = string(data.Type)

		_, err := clnt.ClientHdl.CreateEthernetConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.AggregationLacpConfig:
		data := obj.(models.AggregationLacpConfig)
		conf := lacpdServices.NewAggregationLacpConfig()
		conf.Description = string(data.Description)
		conf.MinLinks = int16(data.MinLinks)
		conf.SystemPriority = int16(data.SystemPriority)
		conf.NameKey = string(data.NameKey)
		conf.Interval = int32(data.Interval)
		conf.Enabled = bool(data.Enabled)
		conf.Mtu = int16(data.Mtu)
		conf.SystemIdMac = string(data.SystemIdMac)
		conf.LagType = int32(data.LagType)
		conf.Type = string(data.Type)
		conf.LacpMode = int32(data.LacpMode)

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
		conf := lacpdServices.NewEthernetConfig()
		conf.MacAddress = string(data.MacAddress)
		conf.Description = string(data.Description)
		conf.AggregateId = string(data.AggregateId)
		conf.NameKey = string(data.NameKey)
		conf.Enabled = bool(data.Enabled)
		conf.Speed = string(data.Speed)
		conf.Mtu = int16(data.Mtu)
		conf.DuplexMode = int32(data.DuplexMode)
		conf.EnableFlowControl = bool(data.EnableFlowControl)
		conf.Auto = bool(data.Auto)
		conf.Type = string(data.Type)

		_, err := clnt.ClientHdl.DeleteEthernetConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.AggregationLacpConfig:
		data := obj.(models.AggregationLacpConfig)
		conf := lacpdServices.NewAggregationLacpConfig()
		conf.Description = string(data.Description)
		conf.MinLinks = int16(data.MinLinks)
		conf.SystemPriority = int16(data.SystemPriority)
		conf.NameKey = string(data.NameKey)
		conf.Interval = int32(data.Interval)
		conf.Enabled = bool(data.Enabled)
		conf.Mtu = int16(data.Mtu)
		conf.SystemIdMac = string(data.SystemIdMac)
		conf.LagType = int32(data.LagType)
		conf.Type = string(data.Type)
		conf.LacpMode = int32(data.LacpMode)

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
			bulkInfo, _ := clnt.ClientHdl.GetBulkAggregationLacpState(lacpdServices.Int(currMarker), lacpdServices.Int(count))
			if bulkInfo.Count != 0 {
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
					ret_obj.Type = string(bulkInfo.AggregationLacpStateList[i].Type)
					ret_obj.LacpMode = int32(bulkInfo.AggregationLacpStateList[i].LacpMode)
					objs = append(objs, ret_obj)
				}
			}
		}
		break

	case models.AggregationLacpMemberStateCounters:

		if clnt.ClientHdl != nil {
			var ret_obj models.AggregationLacpMemberStateCounters
			bulkInfo, _ := clnt.ClientHdl.GetBulkAggregationLacpMemberStateCounters(lacpdServices.Int(currMarker), lacpdServices.Int(count))
			if bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}
					ret_obj.SystemIdMac = string(bulkInfo.AggregationLacpMemberStateCountersList[i].SystemIdMac)
					ret_obj.MinLinks = uint16(bulkInfo.AggregationLacpMemberStateCountersList[i].MinLinks)
					ret_obj.SystemPriority = uint16(bulkInfo.AggregationLacpMemberStateCountersList[i].SystemPriority)
					ret_obj.LacpUnknownErrors = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LacpUnknownErrors)
					ret_obj.Interval = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].Interval)
					ret_obj.Enabled = bool(bulkInfo.AggregationLacpMemberStateCountersList[i].Enabled)
					ret_obj.Aggregatable = bool(bulkInfo.AggregationLacpMemberStateCountersList[i].Aggregatable)
					ret_obj.OperKey = uint16(bulkInfo.AggregationLacpMemberStateCountersList[i].OperKey)
					ret_obj.Mtu = uint16(bulkInfo.AggregationLacpMemberStateCountersList[i].Mtu)
					ret_obj.Distributing = bool(bulkInfo.AggregationLacpMemberStateCountersList[i].Distributing)
					ret_obj.PartnerKey = uint16(bulkInfo.AggregationLacpMemberStateCountersList[i].PartnerKey)
					ret_obj.LacpErrors = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LacpErrors)
					ret_obj.SystemId = string(bulkInfo.AggregationLacpMemberStateCountersList[i].SystemId)
					ret_obj.Timeout = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].Timeout)
					ret_obj.Activity = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].Activity)
					ret_obj.LacpRxErrors = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LacpRxErrors)
					ret_obj.Type = string(bulkInfo.AggregationLacpMemberStateCountersList[i].Type)
					ret_obj.Collecting = bool(bulkInfo.AggregationLacpMemberStateCountersList[i].Collecting)
					ret_obj.LagType = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].LagType)
					ret_obj.Description = string(bulkInfo.AggregationLacpMemberStateCountersList[i].Description)
					ret_obj.LacpTxErrors = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LacpTxErrors)
					ret_obj.LacpOutPkts = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LacpOutPkts)
					ret_obj.LacpInPkts = uint64(bulkInfo.AggregationLacpMemberStateCountersList[i].LacpInPkts)
					ret_obj.Synchronization = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].Synchronization)
					ret_obj.PartnerId = string(bulkInfo.AggregationLacpMemberStateCountersList[i].PartnerId)
					ret_obj.NameKey = string(bulkInfo.AggregationLacpMemberStateCountersList[i].NameKey)
					ret_obj.Interface = string(bulkInfo.AggregationLacpMemberStateCountersList[i].Interface)
					ret_obj.LacpMode = int32(bulkInfo.AggregationLacpMemberStateCountersList[i].LacpMode)
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
func (clnt *LACPDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called LACPD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.EthernetConfig:
		// cast original object
		origdata := dbObj.(models.EthernetConfig)
		updatedata := obj.(models.EthernetConfig)
		// create new thrift objects
		origconf := lacpdServices.NewEthernetConfig()
		updateconf := lacpdServices.NewEthernetConfig()

		origconf.MacAddress = string(origdata.MacAddress)
		origconf.Description = string(origdata.Description)
		origconf.AggregateId = string(origdata.AggregateId)
		origconf.NameKey = string(origdata.NameKey)
		origconf.Enabled = bool(origdata.Enabled)
		origconf.Speed = string(origdata.Speed)
		origconf.Mtu = int16(origdata.Mtu)
		origconf.DuplexMode = int32(origdata.DuplexMode)
		origconf.EnableFlowControl = bool(origdata.EnableFlowControl)
		origconf.Auto = bool(origdata.Auto)
		origconf.Type = string(origdata.Type)

		updateconf.MacAddress = string(updatedata.MacAddress)
		updateconf.Description = string(updatedata.Description)
		updateconf.AggregateId = string(updatedata.AggregateId)
		updateconf.NameKey = string(updatedata.NameKey)
		updateconf.Enabled = bool(updatedata.Enabled)
		updateconf.Speed = string(updatedata.Speed)
		updateconf.Mtu = int16(updatedata.Mtu)
		updateconf.DuplexMode = int32(updatedata.DuplexMode)
		updateconf.EnableFlowControl = bool(updatedata.EnableFlowControl)
		updateconf.Auto = bool(updatedata.Auto)
		updateconf.Type = string(updatedata.Type)

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
		origconf := lacpdServices.NewAggregationLacpConfig()
		updateconf := lacpdServices.NewAggregationLacpConfig()

		origconf.Description = string(origdata.Description)
		origconf.MinLinks = int16(origdata.MinLinks)
		origconf.SystemPriority = int16(origdata.SystemPriority)
		origconf.NameKey = string(origdata.NameKey)
		origconf.Interval = int32(origdata.Interval)
		origconf.Enabled = bool(origdata.Enabled)
		origconf.Mtu = int16(origdata.Mtu)
		origconf.SystemIdMac = string(origdata.SystemIdMac)
		origconf.LagType = int32(origdata.LagType)
		origconf.Type = string(origdata.Type)
		origconf.LacpMode = int32(origdata.LacpMode)

		updateconf.Description = string(updatedata.Description)
		updateconf.MinLinks = int16(updatedata.MinLinks)
		updateconf.SystemPriority = int16(updatedata.SystemPriority)
		updateconf.NameKey = string(updatedata.NameKey)
		updateconf.Interval = int32(updatedata.Interval)
		updateconf.Enabled = bool(updatedata.Enabled)
		updateconf.Mtu = int16(updatedata.Mtu)
		updateconf.SystemIdMac = string(updatedata.SystemIdMac)
		updateconf.LagType = int32(updatedata.LagType)
		updateconf.Type = string(updatedata.Type)
		updateconf.LacpMode = int32(updatedata.LacpMode)

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
