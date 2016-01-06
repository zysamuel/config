package main

import (
	"database/sql"
	"models"
	"ribdServices"
)

type RIBDClient struct {
	IPCClientBase
	ClientHdl *ribdServices.RIBDServicesClient
}

func (clnt *RIBDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *RIBDClient) ConnectToServer() bool {

	clnt.Transport, clnt.PtrProtocolFactory = CreateIPCHandles(clnt.Address)
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = ribdServices.NewRIBDServicesClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
	}
	return true
}
func (clnt *RIBDClient) IsConnectedToServer() bool {
	return true
}
func (clnt *RIBDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {

	case models.IPV4Route:
		data := obj.(models.IPV4Route)
		conf := ribdServices.NewIPV4Route()
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
	default:
		break
	}

	return objId, true
}
func (clnt *RIBDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {

	switch obj.(type) {

	case models.IPV4Route:
		data := obj.(models.IPV4Route)
		conf := ribdServices.NewIPV4Route()
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
	default:
		break
	}

	return true
}
func (clnt *RIBDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {

	logger.Println("### Get Bulk request called with", currMarker, count)
	switch obj.(type) {

	case models.IPV4Route:

		if clnt.ClientHdl != nil {
			var ret_obj models.IPV4Route
			bulkInfo, _ := clnt.ClientHdl.GetBulkIPV4Route(ribdServices.Int(currMarker), ribdServices.Int(count))
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

	default:
		break
	}
	return nil, objCount, nextMarker, more, objs

}
func (clnt *RIBDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called RIBD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.IPV4Route:
		// cast original object
		origdata := dbObj.(models.IPV4Route)
		updatedata := obj.(models.IPV4Route)
		// create new thrift objects
		origconf := ribdServices.NewIPV4Route()
		updateconf := ribdServices.NewIPV4Route()

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

	default:
		break
	}
	return ok

}
