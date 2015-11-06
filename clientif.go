package main

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"ribd"
)

//var gRibClient RibClient = RibClient{}

type RibClient struct {
	Address            string
	Transport          thrift.TTransport
	PtrProtocolFactory *thrift.TBinaryProtocolFactory
	ClientHdl          *ribd.RouteServiceClient
	IsConnected        bool
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
	//clnt.ClientHdl.CreateV4Route(0, 0, 0, 0, 0)
	return true
}

func (clnt *RibClient) IsConnectedToServer() bool {
	return true
}

func (clnt *RibClient) CreateObject() bool {
	clnt.ClientHdl.CreateV4Route(0, 0, 0, 0, 0)
	return true
}

type AsicDClient struct {
	Port               int
	PtrTranPort        *thrift.TTransport
	PtrProtocolFactory *thrift.TBinaryProtocolFactory
	IsConnected        bool
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

func (clnt *AsicDClient) CreateObject() bool {
	return true
}
