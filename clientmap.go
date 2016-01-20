package main

var ClientInterfaces = map[string]ClientIf{"ribd": &RibClient{},
	"asicd":      &AsicDClient{},
	"arpd":       &ArpDClient{},
	"bgpd":       &BgpDClient{},
	"lacpd":      &LACPDClient{},
	"dhcprelayd": &DHCPRELAYDClient{},
	"local":      &LocalClient{},
	"ospfd":      &OSPFDClient{}}
