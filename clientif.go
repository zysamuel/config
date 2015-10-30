package main
import ( "ribd"
		  "git.apache.org/thrift.git/lib/go/thrift"
		 )
var gRibClient RibClient = RibClient{}
type RibClient struct {
	 Address             string
    Transport           thrift.TTransport 
	 PtrProtocolFactory  *thrift.TBinaryProtocolFactory
	 ClientHdl           *ribd.RouteServiceClient
	 IsConnected         bool
}

func ( clnt RibClient ) Initialize( name string, address string) {
	 gRibClient.Address = address 
	 return 
}

func ( clnt RibClient ) ConnectToServer () bool {

	 gRibClient.Transport, gRibClient.PtrProtocolFactory = CreateIPCHandles(gRibClient.Address)
	 gRibClient.ClientHdl =  ribd.NewRouteServiceClientFactory(gRibClient.Transport, gRibClient.PtrProtocolFactory)
	 //gRibClient.ClientHdl.CreateV4Route(0,0,0,0,0)
	 return true
}

func (clnt RibClient) IsConnectedToServer() bool {
	 return true
}

type AsicDClient struct {
	 Port			         int
    PtrTranPort         *thrift.TTransport 
	 PtrProtocolFactory  *thrift.TBinaryProtocolFactory
	 IsConnected         bool
}

func ( clnt AsicDClient) Initialize( name string, address string) {
	 return 
}

func (clnt AsicDClient) ConnectToServer () bool {
	 return true
}

func (clnt AsicDClient) IsConnectedToServer() bool {
	 return true
}
