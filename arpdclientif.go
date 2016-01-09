package main

import (
	"arpdServices"
	"database/sql"
	"models"
	"utils/ipcutils"
)

type ARPDClient struct {
	IPCClientBase
	ClientHdl *arpdServices.ARPDServicesClient
}

func (clnt *ARPDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *ARPDClient) ConnectToServer() bool {

	clnt.Transport, clnt.PtrProtocolFactory = ipcutils.CreateIPCHandles(clnt.Address)
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = arpdServices.NewARPDServicesClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
	}
	return true
}
func (clnt *ARPDClient) IsConnectedToServer() bool {
	return true
}
func (clnt *ARPDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {

	case models.ArpConfig:
		data := obj.(models.ArpConfig)
		conf := arpdServices.NewArpConfig()
		conf.ArpConfigKey = string(data.ArpConfigKey)
		conf.Timeout = int32(data.Timeout)

		_, err := clnt.ClientHdl.CreateArpConfig(conf)
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
func (clnt *ARPDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {

	switch obj.(type) {

	case models.ArpConfig:
		data := obj.(models.ArpConfig)
		conf := arpdServices.NewArpConfig()
		conf.ArpConfigKey = string(data.ArpConfigKey)
		conf.Timeout = int32(data.Timeout)

		_, err := clnt.ClientHdl.DeleteArpConfig(conf)
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
func (clnt *ARPDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {

	logger.Println("### Get Bulk request called with", currMarker, count)
	switch obj.(type) {

	case models.ArpEntry:

		if clnt.ClientHdl != nil {
			var ret_obj models.ArpEntry
			bulkInfo, _ := clnt.ClientHdl.GetBulkArpEntry(arpdServices.Int(currMarker), arpdServices.Int(count))
			if bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}
					ret_obj.Intf = string(bulkInfo.ArpEntryList[i].Intf)
					ret_obj.MacAddr = string(bulkInfo.ArpEntryList[i].MacAddr)
					ret_obj.IpAddr = string(bulkInfo.ArpEntryList[i].IpAddr)
					ret_obj.ExpiryTimeLeft = string(bulkInfo.ArpEntryList[i].ExpiryTimeLeft)
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
func (clnt *ARPDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []byte, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called ARPD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.ArpConfig:
		// cast original object
		origdata := dbObj.(models.ArpConfig)
		updatedata := obj.(models.ArpConfig)
		// create new thrift objects
		origconf := arpdServices.NewArpConfig()
		updateconf := arpdServices.NewArpConfig()

		origconf.ArpConfigKey = string(origdata.ArpConfigKey)
		origconf.Timeout = int32(origdata.Timeout)

		updateconf.ArpConfigKey = string(updatedata.ArpConfigKey)
		updateconf.Timeout = int32(updatedata.Timeout)

		//convert attrSet to uint8 list
		newattrset := make([]int8, len(attrSet))
		for i, v := range attrSet {
			newattrset[i] = int8(v)
		}
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateArpConfig(origconf, updateconf, newattrset)
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
