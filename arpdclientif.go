package main

import (
	"arpd"
	"database/sql"
	"models"
	"utils/ipcutils"
)

type ARPDClient struct {
	ipcutils.IPCClientBase
	ClientHdl *arpd.ARPDServicesClient
}

func (clnt *ARPDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *ARPDClient) ConnectToServer() bool {

	clnt.TTransport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	if clnt.TTransport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = arpd.NewARPDServicesClientFactory(clnt.TTransport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}
func (clnt *ARPDClient) IsConnectedToServer() bool {
	return clnt.IsConnected
}
func (clnt *ARPDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {

	case models.ArpConfig:
		data := obj.(models.ArpConfig)
		conf := arpd.NewArpConfig()
		models.ConvertarpdArpConfigObjToThrift(&data, conf)
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
		conf := arpd.NewArpConfig()
		models.ConvertarpdArpConfigObjToThrift(&data, conf)
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
			bulkInfo, _ := clnt.ClientHdl.GetBulkArpEntry(arpd.Int(currMarker), arpd.Int(count))
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
func (clnt *ARPDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called ARPD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.ArpConfig:
		// cast original object
		origdata := dbObj.(models.ArpConfig)
		updatedata := obj.(models.ArpConfig)
		// create new thrift objects
		origconf := arpd.NewArpConfig()
		updateconf := arpd.NewArpConfig()
		models.ConvertarpdArpConfigObjToThrift(&origdata, origconf)
		models.ConvertarpdArpConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdateArpConfig(origconf, updateconf, attrSet)
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
