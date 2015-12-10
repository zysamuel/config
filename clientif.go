package main

import (
	"asicdServices"
	"bgpd"
	"portdServices"
	//"encoding/binary"
	"infra/portd/portdCommonDefs"
	"git.apache.org/thrift.git/lib/go/thrift"
	"models"
	//"net"
	"database/sql"
	"ribd"
	"strconv"
)

type IPCClientBase struct {
	Address            string
	Transport          thrift.TTransport
	PtrProtocolFactory *thrift.TBinaryProtocolFactory
	IsConnected        bool
}

func (clnt *IPCClientBase) IsConnectedToServer() bool {
	return clnt.IsConnected
}

func (clnt *IPCClientBase) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more       bool,
	objs []models.ConfigObj) {
	//logger.Println("### Get Bulk request called with", currMarker, count)
	return nil, 0, 0, false,make([]models.ConfigObj, 0)
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

	if clnt.Transport == nil && clnt.PtrProtocolFactory == nil {
		clnt.Transport, clnt.PtrProtocolFactory = CreateIPCHandles(clnt.Address)
	}
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = portdServices.NewPortServiceClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
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

	if clnt.Transport == nil && clnt.PtrProtocolFactory == nil {
		clnt.Transport, clnt.PtrProtocolFactory = CreateIPCHandles(clnt.Address)
	}
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = ribd.NewRouteServiceClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}

func (clnt *RibClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more       bool,
	objs []models.ConfigObj) {
	logger.Println("### Get Bulk request called with", currMarker, count)
	var ret_obj models.IPV4Route
	switch obj.(type) {
		case models.IPV4Route:
		     if(clnt.ClientHdl != nil) {
				routesInfo, _ := clnt.ClientHdl.GetBulkRoutes(ribd.Int(currMarker), ribd.Int(count))
                if(routesInfo.Count != 0){
					objCount = int64(routesInfo.Count)
					more = bool(routesInfo.More)
					nextMarker = int64(routesInfo.EndIdx)
					for i:=0;i<int(routesInfo.Count);i++ {
						if(len(objs) == 0) {
							objs = make([]models.ConfigObj, 0)
						}
						ret_obj.DestinationNw = routesInfo.RouteList[i].Ipaddr
						ret_obj.NetworkMask = routesInfo.RouteList[i].Mask
						ret_obj.NextHopIp = routesInfo.RouteList[i].NextHopIp
						ret_obj.Cost = int(routesInfo.RouteList[i].Metric)
						ret_obj.Protocol = ""
                        if(routesInfo.RouteList[i].NextHopIfType == portdCommonDefs.VLAN) {
							ret_obj.OutgoingIntfType = "VLAN"
						} else {
							ret_obj.OutgoingIntfType = "PHY"
						}
						ret_obj.OutgoingInterface = strconv.Itoa(int(routesInfo.RouteList[i].IfIndex))
						objs = append(objs, ret_obj)
					}
				}
			}
	}
	return nil, objCount, nextMarker, more, objs
}

func (clnt *RibClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {

	switch obj.(type) {

	case models.IPV4Route:
		v4Route := obj.(models.IPV4Route)
		outIntf, _ := strconv.Atoi(v4Route.OutgoingInterface)
		var outIntfType ribd.Int
		if(v4Route.OutgoingIntfType == "VLAN") {
			outIntfType = portdCommonDefs.VLAN
		} else {
			outIntfType = portdCommonDefs.PHY
		}
		proto, _ := strconv.Atoi(v4Route.Protocol)
		if clnt.ClientHdl != nil {
			clnt.ClientHdl.CreateV4Route(
				v4Route.DestinationNw, //ribd.Int(binary.BigEndian.Uint32(net.ParseIP(v4Route.DestinationNw).To4())),
				v4Route.NetworkMask,   //ribd.Int(prefixLen),
				ribd.Int(v4Route.Cost),
				v4Route.NextHopIp, //ribd.Int(binary.BigEndian.Uint32(net.ParseIP(v4Route.NextHopIp).To4())),
               outIntfType,
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
	if clnt.Transport == nil && clnt.PtrProtocolFactory == nil {
		clnt.Transport, clnt.PtrProtocolFactory = CreateIPCHandles(clnt.Address)
	}
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = asicdServices.NewAsicdServiceClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
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
	if clnt.Transport == nil && clnt.PtrProtocolFactory == nil {
		clnt.Transport, clnt.PtrProtocolFactory = CreateIPCHandles(clnt.Address)
	}
	if clnt.Transport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = bgpd.NewBGPServerClientFactory(clnt.Transport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
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

func (clnt *BgpDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error, objCount int64,
	nextMarker int64, more bool, objs []models.ConfigObj) {

	logger.Println("BgpDClient: GetBulkObject called - start")
	switch obj.(type) {
	case models.BGPNeighborState:
		var bgpNeighborStateBulk *bgpd.BGPNeighborStateBulk
		bgpNeighborStateBulk, err = clnt.ClientHdl.BulkGetBGPNeighbors(currMarker, count)
		if err != nil {
			break
		}

		for _, item := range bgpNeighborStateBulk.StateList {
			bgpNeighborState := models.BGPNeighborState{
				PeerAS: uint32(item.PeerAS),
				LocalAS: uint32(item.LocalAS),
				PeerType: models.PeerType(item.PeerType),
				AuthPassword: item.AuthPassword,
				Description: item.Description,
				NeighborAddress: item.NeighborAddress,
				SessionState: uint32(item.SessionState),
				Messages: models.BGPMessages{
					Sent: models.BgpCounters{
						Update: uint64(item.Messages.Sent.Update),
						Notification: uint64(item.Messages.Sent.Notification),
					},
					Received: models.BgpCounters{
						Update: uint64(item.Messages.Received.Update),
						Notification: uint64(item.Messages.Received.Notification),
					},
				},
				Queues: models.BGPQueues{
					Input: uint32(item.Queues.Input),
					Output: uint32(item.Queues.Output),
				},
			}
			objs = append(objs, bgpNeighborState)
		}
		nextMarker = bgpNeighborStateBulk.NextIndex
		objCount = bgpNeighborStateBulk.Count
		more = bgpNeighborStateBulk.More
	}
	return err, objCount, nextMarker, more, objs
}

func (clnt *BgpDClient) DeleteObject(obj models.ConfigObj, objId int64, dbHdl *sql.DB) bool {
	return true
}
