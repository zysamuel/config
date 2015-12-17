package main

var ClientInterfaces = map[string]ClientIf{"ribd": &RibClient{},
	"portd": &PortDClient{},
	"asicd": &AsicDClient{},
	"bgpd":  &BgpDClient{},
	"lacpd": &LACPDClient{}}
