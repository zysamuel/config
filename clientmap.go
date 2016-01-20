package main

var ClientInterfaces = map[string]ClientIf{"ribd": &RibClient{},
	"portd": &PortDClient{},
	"asicd": &AsicDClient{},
	"arpd":  &ArpDClient{},
	"bgpd":  &BgpDClient{},
        "ospfd": &OSPFDClient{}}
