package main

import (
	"git.apache.org/thrift.git/lib/go/thrift"
)

func SocketCloseNotification(clnt interface{}) (err error) {
	logger.Println("### Socket closed for client ", clnt)
	return nil
}

//
// This method gets Thrift related IPC handles.
//
func CreateIPCHandles(address string, clnt interface{}) (thrift.TTransport, *thrift.TBinaryProtocolFactory) {
	var transportFactory thrift.TTransportFactory
	var transport thrift.TTransport
	var protocolFactory *thrift.TBinaryProtocolFactory
	var err error

	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory = thrift.NewTTransportFactory()
	transport, err = thrift.NewTSocketTimeout(address, 2*60*1000, SocketCloseNotification, clnt)
	transport = transportFactory.GetTransport(transport)
	if err = transport.Open(); err != nil {
		//logger.Println("Failed to Open Transport", transport, protocolFactory)
		return nil, nil
	}
	return transport, protocolFactory
}
