package main

import (
	"bfdd"
	"database/sql"
	"fmt"
	"models"
	"utils/ipcutils"
)

type BFDDClient struct {
	ipcutils.IPCClientBase
	ClientHdl *bfdd.BFDDServicesClient
}

func (clnt *BFDDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *BFDDClient) ConnectToServer() bool {

	clnt.TTransport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	if clnt.TTransport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = bfdd.NewBFDDServicesClientFactory(clnt.TTransport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}
func (clnt *BFDDClient) IsConnectedToServer() bool {
	return clnt.IsConnected
}
func (clnt *BFDDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {

	case models.BfdGlobalConfig:
		data := obj.(models.BfdGlobalConfig)
		conf := bfdd.NewBfdGlobalConfig()
		models.ConvertbfddBfdGlobalConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateBfdGlobalConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.BfdIntfConfig:
		data := obj.(models.BfdIntfConfig)
		conf := bfdd.NewBfdIntfConfig()
		models.ConvertbfddBfdIntfConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateBfdIntfConfig(conf)
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
func (clnt *BFDDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {

	switch obj.(type) {

	case models.BfdGlobalConfig:
		data := obj.(models.BfdGlobalConfig)
		conf := bfdd.NewBfdGlobalConfig()
		models.ConvertbfddBfdGlobalConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteBfdGlobalConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.BfdIntfConfig:
		data := obj.(models.BfdIntfConfig)
		conf := bfdd.NewBfdIntfConfig()
		models.ConvertbfddBfdIntfConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteBfdIntfConfig(conf)
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
func (clnt *BFDDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {

	logger.Println("### Get Bulk request called with", currMarker, count)
	switch obj.(type) {

	case models.BfdGlobalState:

		if clnt.ClientHdl != nil {
			var ret_obj models.BfdGlobalState
			bulkInfo, err := clnt.ClientHdl.GetBulkBfdGlobalState(bfdd.Int(currMarker), bfdd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.NumUpSessions = uint32(bulkInfo.BfdGlobalStateList[i].NumUpSessions)
					ret_obj.Enable = bool(bulkInfo.BfdGlobalStateList[i].Enable)
					ret_obj.NumDownSessions = uint32(bulkInfo.BfdGlobalStateList[i].NumDownSessions)
					ret_obj.NumAdminDownSessions = uint32(bulkInfo.BfdGlobalStateList[i].NumAdminDownSessions)
					ret_obj.NumInterfaces = uint32(bulkInfo.BfdGlobalStateList[i].NumInterfaces)
					ret_obj.NumTotalSessions = uint32(bulkInfo.BfdGlobalStateList[i].NumTotalSessions)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.BfdIntfState:

		if clnt.ClientHdl != nil {
			var ret_obj models.BfdIntfState
			bulkInfo, err := clnt.ClientHdl.GetBulkBfdIntfState(bfdd.Int(currMarker), bfdd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.InterfaceId = int32(bulkInfo.BfdIntfStateList[i].InterfaceId)
					ret_obj.DemandEnabled = bool(bulkInfo.BfdIntfStateList[i].DemandEnabled)
					ret_obj.AuthenticationType = int32(bulkInfo.BfdIntfStateList[i].AuthenticationType)
					ret_obj.RequiredMinRxInterval = int32(bulkInfo.BfdIntfStateList[i].RequiredMinRxInterval)
					ret_obj.Enabled = bool(bulkInfo.BfdIntfStateList[i].Enabled)
					ret_obj.DesiredMinTxInterval = int32(bulkInfo.BfdIntfStateList[i].DesiredMinTxInterval)
					ret_obj.AuthenticationEnabled = bool(bulkInfo.BfdIntfStateList[i].AuthenticationEnabled)
					ret_obj.NumSessions = int32(bulkInfo.BfdIntfStateList[i].NumSessions)
					ret_obj.AuthenticationKeyId = int32(bulkInfo.BfdIntfStateList[i].AuthenticationKeyId)
					ret_obj.RequiredMinEchoRxInterval = int32(bulkInfo.BfdIntfStateList[i].RequiredMinEchoRxInterval)
					ret_obj.AuthenticationData = string(bulkInfo.BfdIntfStateList[i].AuthenticationData)
					ret_obj.LocalMultiplier = int32(bulkInfo.BfdIntfStateList[i].LocalMultiplier)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.BfdSessionState:

		if clnt.ClientHdl != nil {
			var ret_obj models.BfdSessionState
			bulkInfo, err := clnt.ClientHdl.GetBulkBfdSessionState(bfdd.Int(currMarker), bfdd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.InterfaceId = int32(bulkInfo.BfdSessionStateList[i].InterfaceId)
					ret_obj.AuthType = uint32(bulkInfo.BfdSessionStateList[i].AuthType)
					ret_obj.DetectionMultiplier = uint32(bulkInfo.BfdSessionStateList[i].DetectionMultiplier)
					ret_obj.RegisteredProtocols = string(bulkInfo.BfdSessionStateList[i].RegisteredProtocols)
					ret_obj.RequiredMinRxInterval = int32(bulkInfo.BfdSessionStateList[i].RequiredMinRxInterval)
					ret_obj.RemoteMinRxInterval = int32(bulkInfo.BfdSessionStateList[i].RemoteMinRxInterval)
					ret_obj.RemoteDiscriminator = uint32(bulkInfo.BfdSessionStateList[i].RemoteDiscriminator)
					ret_obj.SentAuthSeq = uint32(bulkInfo.BfdSessionStateList[i].SentAuthSeq)
					ret_obj.RemoteSessionState = int32(bulkInfo.BfdSessionStateList[i].RemoteSessionState)
					ret_obj.LocalIpAddr = string(bulkInfo.BfdSessionStateList[i].LocalIpAddr)
					ret_obj.DesiredMinTxInterval = int32(bulkInfo.BfdSessionStateList[i].DesiredMinTxInterval)
					ret_obj.SessionId = int32(bulkInfo.BfdSessionStateList[i].SessionId)
					ret_obj.LocalDiscriminator = uint32(bulkInfo.BfdSessionStateList[i].LocalDiscriminator)
					ret_obj.SessionState = int32(bulkInfo.BfdSessionStateList[i].SessionState)
					ret_obj.AuthSeqKnown = bool(bulkInfo.BfdSessionStateList[i].AuthSeqKnown)
					ret_obj.RemoteDemandMode = bool(bulkInfo.BfdSessionStateList[i].RemoteDemandMode)
					ret_obj.ReceivedAuthSeq = uint32(bulkInfo.BfdSessionStateList[i].ReceivedAuthSeq)
					ret_obj.RemoteIpAddr = string(bulkInfo.BfdSessionStateList[i].RemoteIpAddr)
					ret_obj.DemandMode = bool(bulkInfo.BfdSessionStateList[i].DemandMode)
					ret_obj.LocalDiagType = int32(bulkInfo.BfdSessionStateList[i].LocalDiagType)
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
func (clnt *BFDDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called BFDD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.BfdGlobalConfig:
		// cast original object
		origdata := dbObj.(models.BfdGlobalConfig)
		updatedata := obj.(models.BfdGlobalConfig)
		// create new thrift objects
		origconf := bfdd.NewBfdGlobalConfig()
		updateconf := bfdd.NewBfdGlobalConfig()
		models.ConvertbfddBfdGlobalConfigObjToThrift(&origdata, origconf)
		models.ConvertbfddBfdGlobalConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateBfdGlobalConfig(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.BfdIntfConfig:
		// cast original object
		origdata := dbObj.(models.BfdIntfConfig)
		updatedata := obj.(models.BfdIntfConfig)
		// create new thrift objects
		origconf := bfdd.NewBfdIntfConfig()
		updateconf := bfdd.NewBfdIntfConfig()
		models.ConvertbfddBfdIntfConfigObjToThrift(&origdata, origconf)
		models.ConvertbfddBfdIntfConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateBfdIntfConfig(origconf, updateconf, attrSet)
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
