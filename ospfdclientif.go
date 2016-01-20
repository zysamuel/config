package main

import (
	"database/sql"
	"fmt"
	"models"
	"ospfd"
	"utils/ipcutils"
)

type OSPFDClient struct {
	ipcutils.IPCClientBase
	ClientHdl *ospfd.OSPFDServicesClient
}

func (clnt *OSPFDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *OSPFDClient) ConnectToServer() bool {

	clnt.TTransport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	if clnt.TTransport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = ospfd.NewOSPFDServicesClientFactory(clnt.TTransport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}
func (clnt *OSPFDClient) IsConnectedToServer() bool {
	return clnt.IsConnected
}
func (clnt *OSPFDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {

	case models.OspfAreaEntryConfig:
		data := obj.(models.OspfAreaEntryConfig)
		conf := ospfd.NewOspfAreaEntryConfig()
		models.ConvertospfdOspfAreaEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateOspfAreaEntryConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.OspfStubAreaEntryConfig:
		data := obj.(models.OspfStubAreaEntryConfig)
		conf := ospfd.NewOspfStubAreaEntryConfig()
		models.ConvertospfdOspfStubAreaEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateOspfStubAreaEntryConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.OspfAreaRangeEntryConfig:
		data := obj.(models.OspfAreaRangeEntryConfig)
		conf := ospfd.NewOspfAreaRangeEntryConfig()
		models.ConvertospfdOspfAreaRangeEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateOspfAreaRangeEntryConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.OspfHostEntryConfig:
		data := obj.(models.OspfHostEntryConfig)
		conf := ospfd.NewOspfHostEntryConfig()
		models.ConvertospfdOspfHostEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateOspfHostEntryConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.OspfIfEntryConfig:
		data := obj.(models.OspfIfEntryConfig)
		conf := ospfd.NewOspfIfEntryConfig()
		models.ConvertospfdOspfIfEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateOspfIfEntryConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.OspfIfMetricEntryConfig:
		data := obj.(models.OspfIfMetricEntryConfig)
		conf := ospfd.NewOspfIfMetricEntryConfig()
		models.ConvertospfdOspfIfMetricEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateOspfIfMetricEntryConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.OspfVirtIfEntryConfig:
		data := obj.(models.OspfVirtIfEntryConfig)
		conf := ospfd.NewOspfVirtIfEntryConfig()
		models.ConvertospfdOspfVirtIfEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateOspfVirtIfEntryConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.OspfNbrEntryConfig:
		data := obj.(models.OspfNbrEntryConfig)
		conf := ospfd.NewOspfNbrEntryConfig()
		models.ConvertospfdOspfNbrEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateOspfNbrEntryConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.OspfAreaAggregateEntryConfig:
		data := obj.(models.OspfAreaAggregateEntryConfig)
		conf := ospfd.NewOspfAreaAggregateEntryConfig()
		models.ConvertospfdOspfAreaAggregateEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateOspfAreaAggregateEntryConfig(conf)
		if err != nil {
			return int64(0), false
		}
		objId, _ = data.StoreObjectInDb(dbHdl)
		break

	case models.OspfGlobalConfig:
		data := obj.(models.OspfGlobalConfig)
		conf := ospfd.NewOspfGlobalConfig()
		models.ConvertospfdOspfGlobalConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreateOspfGlobalConfig(conf)
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
func (clnt *OSPFDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {

	switch obj.(type) {

	case models.OspfAreaEntryConfig:
		data := obj.(models.OspfAreaEntryConfig)
		conf := ospfd.NewOspfAreaEntryConfig()
		models.ConvertospfdOspfAreaEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteOspfAreaEntryConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.OspfStubAreaEntryConfig:
		data := obj.(models.OspfStubAreaEntryConfig)
		conf := ospfd.NewOspfStubAreaEntryConfig()
		models.ConvertospfdOspfStubAreaEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteOspfStubAreaEntryConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.OspfAreaRangeEntryConfig:
		data := obj.(models.OspfAreaRangeEntryConfig)
		conf := ospfd.NewOspfAreaRangeEntryConfig()
		models.ConvertospfdOspfAreaRangeEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteOspfAreaRangeEntryConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.OspfHostEntryConfig:
		data := obj.(models.OspfHostEntryConfig)
		conf := ospfd.NewOspfHostEntryConfig()
		models.ConvertospfdOspfHostEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteOspfHostEntryConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.OspfIfEntryConfig:
		data := obj.(models.OspfIfEntryConfig)
		conf := ospfd.NewOspfIfEntryConfig()
		models.ConvertospfdOspfIfEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteOspfIfEntryConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.OspfIfMetricEntryConfig:
		data := obj.(models.OspfIfMetricEntryConfig)
		conf := ospfd.NewOspfIfMetricEntryConfig()
		models.ConvertospfdOspfIfMetricEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteOspfIfMetricEntryConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.OspfVirtIfEntryConfig:
		data := obj.(models.OspfVirtIfEntryConfig)
		conf := ospfd.NewOspfVirtIfEntryConfig()
		models.ConvertospfdOspfVirtIfEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteOspfVirtIfEntryConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.OspfNbrEntryConfig:
		data := obj.(models.OspfNbrEntryConfig)
		conf := ospfd.NewOspfNbrEntryConfig()
		models.ConvertospfdOspfNbrEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteOspfNbrEntryConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.OspfAreaAggregateEntryConfig:
		data := obj.(models.OspfAreaAggregateEntryConfig)
		conf := ospfd.NewOspfAreaAggregateEntryConfig()
		models.ConvertospfdOspfAreaAggregateEntryConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteOspfAreaAggregateEntryConfig(conf)
		if err != nil {
			return false
		}
		data.DeleteObjectFromDb(objKey, dbHdl)
		break

	case models.OspfGlobalConfig:
		data := obj.(models.OspfGlobalConfig)
		conf := ospfd.NewOspfGlobalConfig()
		models.ConvertospfdOspfGlobalConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeleteOspfGlobalConfig(conf)
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
func (clnt *OSPFDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {

	logger.Println("### Get Bulk request called with", currMarker, count)
	switch obj.(type) {

	case models.OspfAreaEntryState:

		if clnt.ClientHdl != nil {
			var ret_obj models.OspfAreaEntryState
			bulkInfo, err := clnt.ClientHdl.GetBulkOspfAreaEntryState(ospfd.Int(currMarker), ospfd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.AreaIdKey = string(bulkInfo.OspfAreaEntryStateList[i].AreaIdKey)
					ret_obj.AsBdrRtrCount = uint32(bulkInfo.OspfAreaEntryStateList[i].AsBdrRtrCount)
					ret_obj.AreaNssaTranslatorEvents = uint32(bulkInfo.OspfAreaEntryStateList[i].AreaNssaTranslatorEvents)
					ret_obj.SpfRuns = uint32(bulkInfo.OspfAreaEntryStateList[i].SpfRuns)
					ret_obj.AreaLsaCksumSum = int32(bulkInfo.OspfAreaEntryStateList[i].AreaLsaCksumSum)
					ret_obj.AreaNssaTranslatorState = int32(bulkInfo.OspfAreaEntryStateList[i].AreaNssaTranslatorState)
					ret_obj.AreaLsaCount = uint32(bulkInfo.OspfAreaEntryStateList[i].AreaLsaCount)
					ret_obj.AreaBdrRtrCount = uint32(bulkInfo.OspfAreaEntryStateList[i].AreaBdrRtrCount)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.OspfLsdbEntryState:

		if clnt.ClientHdl != nil {
			var ret_obj models.OspfLsdbEntryState
			bulkInfo, err := clnt.ClientHdl.GetBulkOspfLsdbEntryState(ospfd.Int(currMarker), ospfd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.LsdbLsidKey = string(bulkInfo.OspfLsdbEntryStateList[i].LsdbLsidKey)
					ret_obj.LsdbAreaIdKey = string(bulkInfo.OspfLsdbEntryStateList[i].LsdbAreaIdKey)
					ret_obj.LsdbChecksum = int32(bulkInfo.OspfLsdbEntryStateList[i].LsdbChecksum)
					ret_obj.LsdbAdvertisement = string(bulkInfo.OspfLsdbEntryStateList[i].LsdbAdvertisement)
					ret_obj.LsdbAge = int32(bulkInfo.OspfLsdbEntryStateList[i].LsdbAge)
					ret_obj.LsdbRouterIdKey = string(bulkInfo.OspfLsdbEntryStateList[i].LsdbRouterIdKey)
					ret_obj.LsdbSequence = int32(bulkInfo.OspfLsdbEntryStateList[i].LsdbSequence)
					ret_obj.LsdbTypeKey = int32(bulkInfo.OspfLsdbEntryStateList[i].LsdbTypeKey)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.OspfIfEntryState:

		if clnt.ClientHdl != nil {
			var ret_obj models.OspfIfEntryState
			bulkInfo, err := clnt.ClientHdl.GetBulkOspfIfEntryState(ospfd.Int(currMarker), ospfd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.IfIpAddressKey = string(bulkInfo.OspfIfEntryStateList[i].IfIpAddressKey)
					ret_obj.AddressLessIfKey = int32(bulkInfo.OspfIfEntryStateList[i].AddressLessIfKey)
					ret_obj.IfLsaCksumSum = uint32(bulkInfo.OspfIfEntryStateList[i].IfLsaCksumSum)
					ret_obj.IfLsaCount = uint32(bulkInfo.OspfIfEntryStateList[i].IfLsaCount)
					ret_obj.IfDesignatedRouter = string(bulkInfo.OspfIfEntryStateList[i].IfDesignatedRouter)
					ret_obj.IfDesignatedRouterId = string(bulkInfo.OspfIfEntryStateList[i].IfDesignatedRouterId)
					ret_obj.IfBackupDesignatedRouterId = string(bulkInfo.OspfIfEntryStateList[i].IfBackupDesignatedRouterId)
					ret_obj.IfBackupDesignatedRouter = string(bulkInfo.OspfIfEntryStateList[i].IfBackupDesignatedRouter)
					ret_obj.IfEvents = uint32(bulkInfo.OspfIfEntryStateList[i].IfEvents)
					ret_obj.IfState = int32(bulkInfo.OspfIfEntryStateList[i].IfState)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.OspfNbrEntryState:

		if clnt.ClientHdl != nil {
			var ret_obj models.OspfNbrEntryState
			bulkInfo, err := clnt.ClientHdl.GetBulkOspfNbrEntryState(ospfd.Int(currMarker), ospfd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.NbrIpAddrKey = string(bulkInfo.OspfNbrEntryStateList[i].NbrIpAddrKey)
					ret_obj.NbrAddressLessIndexKey = int32(bulkInfo.OspfNbrEntryStateList[i].NbrAddressLessIndexKey)
					ret_obj.NbmaNbrPermanence = int32(bulkInfo.OspfNbrEntryStateList[i].NbmaNbrPermanence)
					ret_obj.NbrRestartHelperStatus = int32(bulkInfo.OspfNbrEntryStateList[i].NbrRestartHelperStatus)
					ret_obj.NbrOptions = int32(bulkInfo.OspfNbrEntryStateList[i].NbrOptions)
					ret_obj.NbrRtrId = string(bulkInfo.OspfNbrEntryStateList[i].NbrRtrId)
					ret_obj.NbrLsRetransQLen = uint32(bulkInfo.OspfNbrEntryStateList[i].NbrLsRetransQLen)
					ret_obj.NbrEvents = uint32(bulkInfo.OspfNbrEntryStateList[i].NbrEvents)
					ret_obj.NbrRestartHelperAge = uint32(bulkInfo.OspfNbrEntryStateList[i].NbrRestartHelperAge)
					ret_obj.NbrRestartHelperExitReason = int32(bulkInfo.OspfNbrEntryStateList[i].NbrRestartHelperExitReason)
					ret_obj.NbrState = int32(bulkInfo.OspfNbrEntryStateList[i].NbrState)
					ret_obj.NbrHelloSuppressed = bool(bulkInfo.OspfNbrEntryStateList[i].NbrHelloSuppressed)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.OspfVirtNbrEntryState:

		if clnt.ClientHdl != nil {
			var ret_obj models.OspfVirtNbrEntryState
			bulkInfo, err := clnt.ClientHdl.GetBulkOspfVirtNbrEntryState(ospfd.Int(currMarker), ospfd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.VirtNbrRestartHelperExitReason = int32(bulkInfo.OspfVirtNbrEntryStateList[i].VirtNbrRestartHelperExitReason)
					ret_obj.VirtNbrRtrIdKey = string(bulkInfo.OspfVirtNbrEntryStateList[i].VirtNbrRtrIdKey)
					ret_obj.VirtNbrOptions = int32(bulkInfo.OspfVirtNbrEntryStateList[i].VirtNbrOptions)
					ret_obj.VirtNbrState = int32(bulkInfo.OspfVirtNbrEntryStateList[i].VirtNbrState)
					ret_obj.VirtNbrLsRetransQLen = uint32(bulkInfo.OspfVirtNbrEntryStateList[i].VirtNbrLsRetransQLen)
					ret_obj.VirtNbrRestartHelperStatus = int32(bulkInfo.OspfVirtNbrEntryStateList[i].VirtNbrRestartHelperStatus)
					ret_obj.VirtNbrHelloSuppressed = bool(bulkInfo.OspfVirtNbrEntryStateList[i].VirtNbrHelloSuppressed)
					ret_obj.VirtNbrAreaKey = string(bulkInfo.OspfVirtNbrEntryStateList[i].VirtNbrAreaKey)
					ret_obj.VirtNbrRestartHelperAge = uint32(bulkInfo.OspfVirtNbrEntryStateList[i].VirtNbrRestartHelperAge)
					ret_obj.VirtNbrIpAddr = string(bulkInfo.OspfVirtNbrEntryStateList[i].VirtNbrIpAddr)
					ret_obj.VirtNbrEvents = uint32(bulkInfo.OspfVirtNbrEntryStateList[i].VirtNbrEvents)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.OspfExtLsdbEntryState:

		if clnt.ClientHdl != nil {
			var ret_obj models.OspfExtLsdbEntryState
			bulkInfo, err := clnt.ClientHdl.GetBulkOspfExtLsdbEntryState(ospfd.Int(currMarker), ospfd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.ExtLsdbTypeKey = int32(bulkInfo.OspfExtLsdbEntryStateList[i].ExtLsdbTypeKey)
					ret_obj.ExtLsdbLsidKey = string(bulkInfo.OspfExtLsdbEntryStateList[i].ExtLsdbLsidKey)
					ret_obj.ExtLsdbChecksum = int32(bulkInfo.OspfExtLsdbEntryStateList[i].ExtLsdbChecksum)
					ret_obj.ExtLsdbRouterIdKey = string(bulkInfo.OspfExtLsdbEntryStateList[i].ExtLsdbRouterIdKey)
					ret_obj.ExtLsdbSequence = int32(bulkInfo.OspfExtLsdbEntryStateList[i].ExtLsdbSequence)
					ret_obj.ExtLsdbAdvertisement = string(bulkInfo.OspfExtLsdbEntryStateList[i].ExtLsdbAdvertisement)
					ret_obj.ExtLsdbAge = int32(bulkInfo.OspfExtLsdbEntryStateList[i].ExtLsdbAge)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.OspfLocalLsdbEntryState:

		if clnt.ClientHdl != nil {
			var ret_obj models.OspfLocalLsdbEntryState
			bulkInfo, err := clnt.ClientHdl.GetBulkOspfLocalLsdbEntryState(ospfd.Int(currMarker), ospfd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.LocalLsdbChecksum = int32(bulkInfo.OspfLocalLsdbEntryStateList[i].LocalLsdbChecksum)
					ret_obj.LocalLsdbAdvertisement = string(bulkInfo.OspfLocalLsdbEntryStateList[i].LocalLsdbAdvertisement)
					ret_obj.LocalLsdbLsidKey = string(bulkInfo.OspfLocalLsdbEntryStateList[i].LocalLsdbLsidKey)
					ret_obj.LocalLsdbIpAddressKey = string(bulkInfo.OspfLocalLsdbEntryStateList[i].LocalLsdbIpAddressKey)
					ret_obj.LocalLsdbRouterIdKey = string(bulkInfo.OspfLocalLsdbEntryStateList[i].LocalLsdbRouterIdKey)
					ret_obj.LocalLsdbAddressLessIfKey = int32(bulkInfo.OspfLocalLsdbEntryStateList[i].LocalLsdbAddressLessIfKey)
					ret_obj.LocalLsdbTypeKey = int32(bulkInfo.OspfLocalLsdbEntryStateList[i].LocalLsdbTypeKey)
					ret_obj.LocalLsdbSequence = int32(bulkInfo.OspfLocalLsdbEntryStateList[i].LocalLsdbSequence)
					ret_obj.LocalLsdbAge = int32(bulkInfo.OspfLocalLsdbEntryStateList[i].LocalLsdbAge)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.OspfVirtLocalLsdbEntryState:

		if clnt.ClientHdl != nil {
			var ret_obj models.OspfVirtLocalLsdbEntryState
			bulkInfo, err := clnt.ClientHdl.GetBulkOspfVirtLocalLsdbEntryState(ospfd.Int(currMarker), ospfd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.VirtLocalLsdbAge = int32(bulkInfo.OspfVirtLocalLsdbEntryStateList[i].VirtLocalLsdbAge)
					ret_obj.VirtLocalLsdbSequence = int32(bulkInfo.OspfVirtLocalLsdbEntryStateList[i].VirtLocalLsdbSequence)
					ret_obj.VirtLocalLsdbNeighborKey = string(bulkInfo.OspfVirtLocalLsdbEntryStateList[i].VirtLocalLsdbNeighborKey)
					ret_obj.VirtLocalLsdbTypeKey = int32(bulkInfo.OspfVirtLocalLsdbEntryStateList[i].VirtLocalLsdbTypeKey)
					ret_obj.VirtLocalLsdbAdvertisement = string(bulkInfo.OspfVirtLocalLsdbEntryStateList[i].VirtLocalLsdbAdvertisement)
					ret_obj.VirtLocalLsdbChecksum = int32(bulkInfo.OspfVirtLocalLsdbEntryStateList[i].VirtLocalLsdbChecksum)
					ret_obj.VirtLocalLsdbRouterIdKey = string(bulkInfo.OspfVirtLocalLsdbEntryStateList[i].VirtLocalLsdbRouterIdKey)
					ret_obj.VirtLocalLsdbTransitAreaKey = string(bulkInfo.OspfVirtLocalLsdbEntryStateList[i].VirtLocalLsdbTransitAreaKey)
					ret_obj.VirtLocalLsdbLsidKey = string(bulkInfo.OspfVirtLocalLsdbEntryStateList[i].VirtLocalLsdbLsidKey)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.OspfAsLsdbEntryState:

		if clnt.ClientHdl != nil {
			var ret_obj models.OspfAsLsdbEntryState
			bulkInfo, err := clnt.ClientHdl.GetBulkOspfAsLsdbEntryState(ospfd.Int(currMarker), ospfd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.AsLsdbAdvertisement = string(bulkInfo.OspfAsLsdbEntryStateList[i].AsLsdbAdvertisement)
					ret_obj.AsLsdbChecksum = int32(bulkInfo.OspfAsLsdbEntryStateList[i].AsLsdbChecksum)
					ret_obj.AsLsdbLsidKey = string(bulkInfo.OspfAsLsdbEntryStateList[i].AsLsdbLsidKey)
					ret_obj.AsLsdbSequence = int32(bulkInfo.OspfAsLsdbEntryStateList[i].AsLsdbSequence)
					ret_obj.AsLsdbTypeKey = int32(bulkInfo.OspfAsLsdbEntryStateList[i].AsLsdbTypeKey)
					ret_obj.AsLsdbRouterIdKey = string(bulkInfo.OspfAsLsdbEntryStateList[i].AsLsdbRouterIdKey)
					ret_obj.AsLsdbAge = int32(bulkInfo.OspfAsLsdbEntryStateList[i].AsLsdbAge)
					objs = append(objs, ret_obj)
				}

			} else {
				fmt.Println(err)
			}
		}
		break

	case models.OspfAreaLsaCountEntryState:

		if clnt.ClientHdl != nil {
			var ret_obj models.OspfAreaLsaCountEntryState
			bulkInfo, err := clnt.ClientHdl.GetBulkOspfAreaLsaCountEntryState(ospfd.Int(currMarker), ospfd.Int(count))
			if bulkInfo != nil && bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}

					ret_obj.AreaLsaCountAreaIdKey = string(bulkInfo.OspfAreaLsaCountEntryStateList[i].AreaLsaCountAreaIdKey)
					ret_obj.AreaLsaCountNumber = uint32(bulkInfo.OspfAreaLsaCountEntryStateList[i].AreaLsaCountNumber)
					ret_obj.AreaLsaCountLsaTypeKey = int32(bulkInfo.OspfAreaLsaCountEntryStateList[i].AreaLsaCountLsaTypeKey)
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
func (clnt *OSPFDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called OSPFD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.OspfAreaEntryConfig:
		// cast original object
		origdata := dbObj.(models.OspfAreaEntryConfig)
		updatedata := obj.(models.OspfAreaEntryConfig)
		// create new thrift objects
		origconf := ospfd.NewOspfAreaEntryConfig()
		updateconf := ospfd.NewOspfAreaEntryConfig()
		models.ConvertospfdOspfAreaEntryConfigObjToThrift(&origdata, origconf)
		models.ConvertospfdOspfAreaEntryConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateOspfAreaEntryConfig(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.OspfStubAreaEntryConfig:
		// cast original object
		origdata := dbObj.(models.OspfStubAreaEntryConfig)
		updatedata := obj.(models.OspfStubAreaEntryConfig)
		// create new thrift objects
		origconf := ospfd.NewOspfStubAreaEntryConfig()
		updateconf := ospfd.NewOspfStubAreaEntryConfig()
		models.ConvertospfdOspfStubAreaEntryConfigObjToThrift(&origdata, origconf)
		models.ConvertospfdOspfStubAreaEntryConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateOspfStubAreaEntryConfig(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.OspfAreaRangeEntryConfig:
		// cast original object
		origdata := dbObj.(models.OspfAreaRangeEntryConfig)
		updatedata := obj.(models.OspfAreaRangeEntryConfig)
		// create new thrift objects
		origconf := ospfd.NewOspfAreaRangeEntryConfig()
		updateconf := ospfd.NewOspfAreaRangeEntryConfig()
		models.ConvertospfdOspfAreaRangeEntryConfigObjToThrift(&origdata, origconf)
		models.ConvertospfdOspfAreaRangeEntryConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateOspfAreaRangeEntryConfig(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.OspfHostEntryConfig:
		// cast original object
		origdata := dbObj.(models.OspfHostEntryConfig)
		updatedata := obj.(models.OspfHostEntryConfig)
		// create new thrift objects
		origconf := ospfd.NewOspfHostEntryConfig()
		updateconf := ospfd.NewOspfHostEntryConfig()
		models.ConvertospfdOspfHostEntryConfigObjToThrift(&origdata, origconf)
		models.ConvertospfdOspfHostEntryConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateOspfHostEntryConfig(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.OspfIfEntryConfig:
		// cast original object
		origdata := dbObj.(models.OspfIfEntryConfig)
		updatedata := obj.(models.OspfIfEntryConfig)
		// create new thrift objects
		origconf := ospfd.NewOspfIfEntryConfig()
		updateconf := ospfd.NewOspfIfEntryConfig()
		models.ConvertospfdOspfIfEntryConfigObjToThrift(&origdata, origconf)
		models.ConvertospfdOspfIfEntryConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateOspfIfEntryConfig(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.OspfIfMetricEntryConfig:
		// cast original object
		origdata := dbObj.(models.OspfIfMetricEntryConfig)
		updatedata := obj.(models.OspfIfMetricEntryConfig)
		// create new thrift objects
		origconf := ospfd.NewOspfIfMetricEntryConfig()
		updateconf := ospfd.NewOspfIfMetricEntryConfig()
		models.ConvertospfdOspfIfMetricEntryConfigObjToThrift(&origdata, origconf)
		models.ConvertospfdOspfIfMetricEntryConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateOspfIfMetricEntryConfig(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.OspfVirtIfEntryConfig:
		// cast original object
		origdata := dbObj.(models.OspfVirtIfEntryConfig)
		updatedata := obj.(models.OspfVirtIfEntryConfig)
		// create new thrift objects
		origconf := ospfd.NewOspfVirtIfEntryConfig()
		updateconf := ospfd.NewOspfVirtIfEntryConfig()
		models.ConvertospfdOspfVirtIfEntryConfigObjToThrift(&origdata, origconf)
		models.ConvertospfdOspfVirtIfEntryConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateOspfVirtIfEntryConfig(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.OspfNbrEntryConfig:
		// cast original object
		origdata := dbObj.(models.OspfNbrEntryConfig)
		updatedata := obj.(models.OspfNbrEntryConfig)
		// create new thrift objects
		origconf := ospfd.NewOspfNbrEntryConfig()
		updateconf := ospfd.NewOspfNbrEntryConfig()
		models.ConvertospfdOspfNbrEntryConfigObjToThrift(&origdata, origconf)
		models.ConvertospfdOspfNbrEntryConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateOspfNbrEntryConfig(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.OspfAreaAggregateEntryConfig:
		// cast original object
		origdata := dbObj.(models.OspfAreaAggregateEntryConfig)
		updatedata := obj.(models.OspfAreaAggregateEntryConfig)
		// create new thrift objects
		origconf := ospfd.NewOspfAreaAggregateEntryConfig()
		updateconf := ospfd.NewOspfAreaAggregateEntryConfig()
		models.ConvertospfdOspfAreaAggregateEntryConfigObjToThrift(&origdata, origconf)
		models.ConvertospfdOspfAreaAggregateEntryConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateOspfAreaAggregateEntryConfig(origconf, updateconf, attrSet)
			if ok {
				updatedata.UpdateObjectInDb(dbObj, attrSet, dbHdl)
			} else {
				panic(err)
			}
		}
		break

	case models.OspfGlobalConfig:
		// cast original object
		origdata := dbObj.(models.OspfGlobalConfig)
		updatedata := obj.(models.OspfGlobalConfig)
		// create new thrift objects
		origconf := ospfd.NewOspfGlobalConfig()
		updateconf := ospfd.NewOspfGlobalConfig()
		models.ConvertospfdOspfGlobalConfigObjToThrift(&origdata, origconf)
		models.ConvertospfdOspfGlobalConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateOspfGlobalConfig(origconf, updateconf, attrSet)
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
