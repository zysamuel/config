package main

var ClientInterfaces = map[string]ClientIf{"ribd": &RibClient{},
	"asicd":      &ASICDClient{},
	"arpd":       &ArpDClient{},
	"bgpd":       &BgpDClient{},
	"lacpd":      &LACPDClient{},
	"dhcprelayd": &DHCPRELAYDClient{},
	"local":      &LocalClient{},
	"ospfd":      &OSPFDClient{},
	"stpd":       &STPDClient{},
	"bfdd":       &BFDDClient{},
	"vrrpd":      &VRRPDClient{},
	"sysd":       &SYSDClient{},
}
