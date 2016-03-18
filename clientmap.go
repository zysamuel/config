package main

var ClientInterfaces = map[string]ClientIf{"ribd": &RIBDClient{},
	"asicd":      &ASICDClient{},
	"arpd":       &ARPDClient{},
	"bgpd":       &BGPDClient{},
	"lacpd":      &LACPDClient{},
	"dhcprelayd": &DHCPRELAYDClient{},
	"local":      &LocalClient{},
	"ospfd":      &OSPFDClient{},
	"stpd":       &STPDClient{},
	"bfdd":       &BFDDClient{},
	"vrrpd":      &VRRPDClient{},
	"sysd":       &SYSDClient{},
}
