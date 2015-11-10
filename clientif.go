package main

import (
	"asicdServices"
	"bgpd"
	"encoding/binary"
	"git.apache.org/thrift.git/lib/go/thrift"
	"models"
	"net"
	"ribd"
)

type IPCClientBase struct {
	Address            string
	Transport          thrift.TTransport
	PtrProtocolFactory *thrift.TBinaryProtocolFactory
	IsConnected        bool
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

func (clnt *RibClient) CreateObject(obj models.ConfigObj) bool {

	switch obj.(type) {

	case models.IPV4Route:
		v4Route := obj.(models.IPV4Route)
		clnt.ClientHdl.CreateV4Route(
			ribd.Int(binary.BigEndian.Uint32(net.ParseIP(v4Route.DestinationNw).To4())),
			ribd.Int(binary.BigEndian.Uint32(net.ParseIP(v4Route.NetworkMask).To4())),
			ribd.Int(binary.BigEndian.Uint32(net.ParseIP(v4Route.NextHopIp).To4())),
			0,
			//v4Route.OutgoingInterface,
			ribd.Int(v4Route.Cost))
		break
	default:
		break
	}

	return true
}

type AsicDClient struct {
	IPCClientBase
	ClientHdl *asicdServices.AsicdServiceClient
}

func (clnt *AsicDClient) Initialize(name string, address string) {
	return
}

func (clnt *AsicDClient) ConnectToServer() bool {
	return true
}

func (clnt *AsicDClient) IsConnectedToServer() bool {
	return true
}

func (clnt *AsicDClient) CreateObject(obj models.ConfigObj) bool {
	return true
}

type BgpDClient struct {
	IPCClientBase
	ClientHdl *bgpd.BgpServerClient
}

func (clnt *BgpDClient) Initialize(name string, address string) {
	return
}

func (clnt *BgpDClient) ConnectToServer() bool {
	return true
}

func (clnt *BgpDClient) IsConnectedToServer() bool {
	return true
}

func (clnt *BgpDClient) CreateObject(obj models.ConfigObj) bool {
	return true
}
