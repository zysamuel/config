package main

import (
	"asicdServices"
	"database/sql"
	"models"
	"utils/ipcutils"
)

type ASICDClient struct {
	ipcutils.IPCClientBase
	ClientHdl *asicdServices.ASICDServicesClient
}

func (clnt *ASICDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *ASICDClient) ConnectToServer() bool {

	clnt.Transport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = asicdServices.NewASICDServicesClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}
func (clnt *ASICDClient) IsConnectedToServer() bool {
	return clnt.IsConnected
}
func (clnt *ASICDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	switch obj.(type) {

	case models.PortIntfConfig:
		data := obj.(models.PortIntfConfig)
		conf := asicdServices.NewPortIntfConfig()
		models.ConvertasicdPortIntfConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.CreatePortIntfConfig(conf)
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
func (clnt *ASICDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {

	switch obj.(type) {

	case models.PortIntfConfig:
		data := obj.(models.PortIntfConfig)
		conf := asicdServices.NewPortIntfConfig()
		models.ConvertasicdPortIntfConfigObjToThrift(&data, conf)
		_, err := clnt.ClientHdl.DeletePortIntfConfig(conf)
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
func (clnt *ASICDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {

	logger.Println("### Get Bulk request called with", currMarker, count)
	switch obj.(type) {

	case models.PortIntfState:

		if clnt.ClientHdl != nil {
			var ret_obj models.PortIntfState
			bulkInfo, _ := clnt.ClientHdl.GetBulkPortIntfState(asicdServices.Int(currMarker), asicdServices.Int(count))
			if bulkInfo.Count != 0 {
				objCount = int64(bulkInfo.Count)
				more = bool(bulkInfo.More)
				nextMarker = int64(bulkInfo.EndIdx)
				for i := 0; i < int(bulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}
					for _, data := range bulkInfo.PortIntfStateList[i].PortStats {
						ret_obj.PortStats = int64(data)
					}

					ret_obj.PortNum = int32(bulkInfo.PortIntfStateList[i].PortNum)
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
func (clnt *ASICDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {

	logger.Println("### Update Object called ASICD", attrSet, objKey)
	ok := false
	switch obj.(type) {

	case models.PortIntfConfig:
		// cast original object
		origdata := dbObj.(models.PortIntfConfig)
		updatedata := obj.(models.PortIntfConfig)
		// create new thrift objects
		origconf := asicdServices.NewPortIntfConfig()
		updateconf := asicdServices.NewPortIntfConfig()
		models.ConvertasicdPortIntfConfigObjToThrift(&origdata, origconf)
		models.ConvertasicdPortIntfConfigObjToThrift(&updatedata, updateconf)
		if clnt.ClientHdl != nil {
			ok, err := clnt.ClientHdl.UpdatePortIntfConfig(origconf, updateconf, attrSet)
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
