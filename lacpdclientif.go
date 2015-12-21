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
		break
	default:
		break
	}

	return int64(0), true
}
func (clnt *LACPDClient) DeleteObject(obj models.ConfigObj, objId string, dbHdl *sql.DB) bool {
	return true
}
func (clnt *LACPDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrs []byte, objId string, dbHdl *sql.DB) bool {
	return true
}
