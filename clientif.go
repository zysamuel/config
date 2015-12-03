package main

import (
	"asicdServices"
	"bgpd"
	"portdServices"
	//"encoding/binary"
	"git.apache.org/thrift.git/lib/go/thrift"
	"models"
	//"net"
	"database/sql"
	"ribd"
	_ "strconv"
)

type IPCClientBase struct {
	Address            string
	Transport          thrift.TTransport
	PtrProtocolFactory *thrift.TBinaryProtocolFactory
	IsConnected        bool
}

type PortDClient struct {
	IPCClientBase
	ClientHdl *portdServices.PortServiceClient
}

func (clnt *PortDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}

func (clnt *PortDClient) ConnectToServer() bool {

	clnt.Transport, clnt.PtrProtocolFactory = CreateIPCHandles(clnt.Address)
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = portdServices.NewPortServiceClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
	}
	return true
}

func (clnt *PortDClient) IsConnectedToServer() bool {
	return true
}

func (clnt *PortDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {

	switch obj.(type) {

	case models.IPv4Intf: //IPv4Intf
		v4Intf := obj.(models.IPv4Intf)
		_, err := clnt.ClientHdl.CreateV4Intf(v4Intf.IpAddr, v4Intf.RouterIf, v4Intf.VlanEnabled)
		if err != nil {
			return int64(0), false
		}
	case models.IPv4Neighbor: //IPv4Neighbor
		v4Nbr := obj.(models.IPv4Neighbor)
		_, err := clnt.ClientHdl.CreateV4Neighbor(v4Nbr.IpAddr, v4Nbr.MacAddr, v4Nbr.VlanId, v4Nbr.RouterIf)
		if err != nil {
			return int64(0), false
		}
		break
	case models.Vlan: //Vlan
		vlanObj := obj.(models.Vlan)
		_, err := clnt.ClientHdl.CreateVlan(vlanObj.VlanId, vlanObj.Ports, vlanObj.PortTagType)
		if err != nil {
			return int64(0), false
		}
	default:
		break
	}

	return int64(0), true
}

func (clnt *PortDClient) DeleteObject(obj models.ConfigObj, objId int64, dbHdl *sql.DB) bool {
	return true
}

type RibClient struct {
	IPCClientBase
	ClientHdl *ribd.RouteServiceClient
}

func (clnt *RibClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}

func (clnt *RibClient) ConnectToServer() bool {

	clnt.Transport, clnt.PtrProtocolFactory = CreateIPCHandles(clnt.Address)
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = ribd.NewRouteServiceClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
	}
	return true
}

func (clnt *RibClient) IsConnectedToServer() bool {
	return true
}

func (clnt *RibClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {

	switch obj.(type) {

	case models.IPV4Route:
		v4Route := obj.(models.IPV4Route)
		outIntf, _ := strconv.Atoi(v4Route.OutgoingInterface)
		proto, _ := strconv.Atoi(v4Route.Protocol)
		if clnt.ClientHdl != nil {
			clnt.ClientHdl.CreateV4Route(
				v4Route.DestinationNw, //ribd.Int(binary.BigEndian.Uint32(net.ParseIP(v4Route.DestinationNw).To4())),
				v4Route.NetworkMask,   //ribd.Int(prefixLen),
				ribd.Int(v4Route.Cost),
				v4Route.NextHopIp, //ribd.Int(binary.BigEndian.Uint32(net.ParseIP(v4Route.NextHopIp).To4())),
				ribd.Int(outIntf),
				ribd.Int(proto))
		}
		objId, _ := v4Route.StoreObjectInDb(dbHdl)
		return objId, true

	default:
		break
	}

	return int64(0), true
}

func (clnt *RibClient) DeleteObject(obj models.ConfigObj, objId int64, dbHdl *sql.DB) bool {
	logger.Println("### Delete Object is called in RIBClient. ObjectId: ", objId, obj)
	switch obj.(type) {
	case models.IPV4Route:
		v4Route := obj.(models.IPV4Route)
		logger.Println("### Delete Object is called in RIBClient. ObjectId: ", objId)
		v4Route.DeleteObjectFromDb(objId, dbHdl)
		//default:
		//	logger.Println("OBJECT Type is ", obj.(type))
	}

	return true
}

type AsicDClient struct {
	IPCClientBase
	ClientHdl *asicdServices.AsicdServiceClient
}

func (clnt *AsicDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}

func (clnt *AsicDClient) ConnectToServer() bool {
	clnt.Transport, clnt.PtrProtocolFactory = CreateIPCHandles(clnt.Address)
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = asicdServices.NewAsicdServiceClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
	}
	return true
}

func (clnt *AsicDClient) IsConnectedToServer() bool {
	return true
}

func (clnt *AsicDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	switch obj.(type) {
	case models.Vlan: //Vlan
		vlanObj := obj.(models.Vlan)
		_, err := clnt.ClientHdl.CreateVlan(vlanObj.VlanId, vlanObj.Ports, vlanObj.PortTagType)
		if err != nil {
			return int64(0), false
		}
	}
	return int64(0), true
}

func (clnt *AsicDClient) DeleteObject(obj models.ConfigObj, objId int64, dbHdl *sql.DB) bool {
	return true
}

type BgpDClient struct {
	IPCClientBase
	ClientHdl *bgpd.BGPServerClient
}

func (clnt *BgpDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}

func (clnt *BgpDClient) ConnectToServer() bool {
	clnt.Transport, clnt.PtrProtocolFactory = CreateIPCHandles(clnt.Address)
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = bgpd.NewBGPServerClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
	}
	return true
}

func (clnt *BgpDClient) IsConnectedToServer() bool {
	return true
}

func (clnt *BgpDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	switch obj.(type) {
	case models.BGPGlobalConfig:
		bgpGlobalConf := obj.(models.BGPGlobalConfig)
		gConf := bgpd.NewBGPGlobal()
		gConf.AS = int32(bgpGlobalConf.AS)
		gConf.RouterId = bgpGlobalConf.RouterId
		_, err := clnt.ClientHdl.CreateBGPGlobal(gConf)
		if err != nil {
			return int64(0), false
		}

	case models.BGPNeighborConfig:
		bgpNeighborConf := obj.(models.BGPNeighborConfig)
		nConf := bgpd.NewBGPNeighbor()
		nConf.PeerAS = int32(bgpNeighborConf.PeerAS)
		nConf.LocalAS = int32(bgpNeighborConf.LocalAS)
		nConf.NeighborAddress = bgpNeighborConf.NeighborAddress
		nConf.Description = bgpNeighborConf.Description
		_, err := clnt.ClientHdl.CreateBGPNeighbor(nConf)
		if err != nil {
			return int64(0), false
		}
	}
	return int64(0), true
}

func (clnt *BgpDClient) DeleteObject(obj models.ConfigObj, objId int64, dbHdl *sql.DB) bool {
	return true
}
