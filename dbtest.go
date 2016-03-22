package main

import (
	"models"
	"strconv"
	"time"
)

var readRound int
var writeRound int

func DbTestReadRequest() error {
	readRound = 1
	for {
		go DbTestReadRequestV4Route()
		go DbTestReadRequestV4Intf()
		go DbTestReadRequestV4Neighbor()
		go DbTestReadRequestVlan()
		go DbTestReadRequestUser()
		/*
			go DbTestReadRequestPortIntf()
			go DbTestReadRequestEthernet()
			go DbTestReadRequestDhcpRelayIntf()
			go DbTestReadRequestDhcpRelayGlobal()
			go DbTestReadRequestBgpNeighbor()
			go DbTestReadRequestBgpGlobal()
			go DbTestReadRequestArpEntry()
			go DbTestReadRequestArp()
			go DbTestReadRequestAggregationLacp()
		*/
		go DbTestReadRequestV4Route()
		go DbTestReadRequestV4Intf()
		go DbTestReadRequestV4Neighbor()
		go DbTestReadRequestVlan()
		go DbTestReadRequestUser()
		/*
			go DbTestReadRequestPortIntf()
			go DbTestReadRequestEthernet()
			go DbTestReadRequestDhcpRelayIntf()
			go DbTestReadRequestDhcpRelayGlobal()
			go DbTestReadRequestBgpNeighbor()
			go DbTestReadRequestBgpGlobal()
			go DbTestReadRequestArpEntry()
			go DbTestReadRequestArp()
			go DbTestReadRequestAggregationLacp()
		*/
		time.Sleep(5 * time.Second)
		readRound++
	}
	return nil
}

func DbTestWriteRequest() error {
	writeRound = 1
	for {
		go DbTestWriteRequestV4Route()
		go DbTestWriteRequestV4Intf()
		go DbTestWriteRequestV4Neighbor()
		go DbTestWriteRequestVlan()
		go DbTestWriteRequestUser()
		time.Sleep(20 * time.Second)
		writeRound++
	}
	return nil
}

func DbTestReadRequestV4Route() error {
	var v4Route *models.IPV4Route
	logger.Println("Calling Get IPV4Route - round ", readRound)
	_, err := v4Route.GetAllObjFromDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to get IPV4Route - round ", readRound, err)
	}
	return nil
}

func DbTestReadRequestV4Intf() error {
	var v4Intf *models.IPv4Intf
	logger.Println("Calling Get IPV4Intf - round ", readRound)
	_, err := v4Intf.GetAllObjFromDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to get IPv4Intf - round ", readRound, err)
	}
	return nil
}

func DbTestReadRequestV4Neighbor() error {
	var v4Neighbor *models.IPv4Neighbor
	logger.Println("Calling Get IPV4Neighbor - round ", readRound)
	_, err := v4Neighbor.GetAllObjFromDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to get IPV4Neighbor - round ", readRound, err)
	}
	return nil
}

func DbTestReadRequestVlan() error {
	var vlan *models.Vlan
	logger.Println("Calling Get Vlan - round ", readRound)
	_, err := vlan.GetAllObjFromDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to get Vlan - round ", readRound, err)
	}
	return nil
}

func DbTestReadRequestUser() error {
	var userConfig *models.UserConfig
	logger.Println("Calling Get UserConfig - round ", readRound)
	_, err := userConfig.GetAllObjFromDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to get UserConfig - round ", readRound, err)
	}
	return nil
}

func DbTestReadRequestPortIntf() error {
	var portConfig *models.PortIntfConfig
	logger.Println("Calling Get PortIntfConfig - round ", readRound)
	_, err := portConfig.GetAllObjFromDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to get PortIntfConfig - round ", readRound, err)
	}
	return nil
}

func DbTestReadRequestEthernet() error {
	var ethConfig *models.EthernetConfig
	logger.Println("Calling Get EthernetConfig - round ", readRound)
	_, err := ethConfig.GetAllObjFromDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to get EthernetConfig - round ", readRound, err)
	}
	return nil
}

func DbTestReadRequestDhcpRelayIntf() error {
	var dhcpRelayIntfConfig *models.DhcpRelayIntfConfig
	logger.Println("Calling Get DhcpRelayIntfConfig - round ", readRound)
	_, err := dhcpRelayIntfConfig.GetAllObjFromDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to get DhcpRelayIntfConfig - round ", readRound, err)
	}
	return nil
}

func DbTestReadRequestDhcpRelayGlobal() error {
	var dhcpRelayGlobalConfig *models.DhcpRelayGlobalConfig
	logger.Println("Calling Get DhcpRelayGlobalConfig - round ", readRound)
	_, err := dhcpRelayGlobalConfig.GetAllObjFromDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to get DhcpRelayGlobalConfig - round ", readRound, err)
	}
	return nil
}

func DbTestReadRequestBgpNeighbor() error {
	var bgpNeighborConfig *models.BGPNeighborConfig
	logger.Println("Calling Get BGPNeighborConfig - round ", readRound)
	_, err := bgpNeighborConfig.GetAllObjFromDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to get BGPNeighborConfig - round ", readRound, err)
	}
	return nil
}

func DbTestReadRequestBgpGlobal() error {
	var bgpGlobalConfig *models.BGPGlobalConfig
	logger.Println("Calling Get BGPGlobalConfig - round ", readRound)
	_, err := bgpGlobalConfig.GetAllObjFromDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to get BGPGlobalConfig - round ", readRound, err)
	}
	return nil
}

func DbTestReadRequestArpEntry() error {
	var arpEntry *models.ArpEntry
	logger.Println("Calling Get ArpEntry - round ", readRound)
	_, err := arpEntry.GetAllObjFromDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to get ArpEntry - round ", readRound, err)
	}
	return nil
}

func DbTestReadRequestArp() error {
	var arpConfig *models.ArpConfig
	logger.Println("Calling Get ArpConfig - round ", readRound)
	_, err := arpConfig.GetAllObjFromDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to get ArpConfig - round ", readRound, err)
	}
	return nil
}

func DbTestReadRequestAggregationLacp() error {
	var aggregationLacpConfig *models.AggregationLacpConfig
	logger.Println("Calling Get AggregationLacpConfig - round ", readRound)
	_, err := aggregationLacpConfig.GetAllObjFromDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to get AggregationLacpConfig - round ", readRound, err)
	}
	return nil
}

func DbTestWriteRequestV4Route() error {
	logger.Println("Calling Store IPV4Route - round ", writeRound)
	obj := models.IPV4Route{DestinationNw: "10.1.1." + strconv.Itoa(writeRound),
		NetworkMask:       "255.255.255.0",
		Cost:              10,
		NextHopIp:         "11.1.1.1",
		OutgoingIntfType:  "Eth",
		OutgoingInterface: "if1",
		Protocol:          "Static"}
	_, err := obj.StoreObjectInDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to store IPV4Route - round ", writeRound, err)
	}
	/*
		time.Sleep(1 * time.Second)
		logger.Println("Calling Delete IPV4Route - round ", writeRound)
		objKey, _ := obj.GetKey()
		err = obj.DeleteObjectFromDb(objKey, gMgr.dbHdl)
		if err != nil {
			logger.Println("Failed to delete IPV4Route - round ", writeRound, err)
		}
	*/
	return nil
}

func DbTestWriteRequestVlan() error {
	logger.Println("Calling Store Vlan - round ", writeRound)
	obj := models.Vlan{VlanId: int32(writeRound),
		Ports:       "fp1, fp2",
		PortTagType: "tag1"}
	_, err := obj.StoreObjectInDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to store Vlan - round ", writeRound, err)
	}
	/*
		time.Sleep(1 * time.Second)
		logger.Println("Calling Delete Vlan - round ", writeRound)
		objKey, _ := obj.GetKey()
		err = obj.DeleteObjectFromDb(objKey, gMgr.dbHdl)
		if err != nil {
			logger.Println("Failed to delete Vlan - round ", writeRound, err)
		}
	*/
	return nil
}

func DbTestWriteRequestUser() error {
	logger.Println("Calling Store User - round ", writeRound)
	obj := models.UserConfig{UserName: "abcd" + strconv.Itoa(writeRound),
		Password:    "password",
		Description: "admin",
		Privilege:  "prev"}
	_, err := obj.StoreObjectInDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to store User - round ", writeRound, err)
	}
	/*
		time.Sleep(1 * time.Second)
		logger.Println("Calling Delete User - round ", writeRound)
		objKey, _ := obj.GetKey()
		err = obj.DeleteObjectFromDb(objKey, gMgr.dbHdl)
		if err != nil {
			logger.Println("Failed to delete User - round ", writeRound, err)
		}
	*/
	return nil
}

func DbTestWriteRequestV4Intf() error {
	logger.Println("Calling Store V4Intf - round ", writeRound)
	obj := models.IPv4Intf{IpAddr: "100.1.1." + strconv.Itoa(writeRound),
		RouterIf: 1,
		IfType:   1}
	_, err := obj.StoreObjectInDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to store V4Intf - round ", writeRound, err)
	}
	/*
		time.Sleep(1 * time.Second)
		logger.Println("Calling Delete V4Intf - round ", writeRound)
		objKey, _ := obj.GetKey()
		err = obj.DeleteObjectFromDb(objKey, gMgr.dbHdl)
		if err != nil {
			logger.Println("Failed to delete V4Intf - round ", writeRound, err)
		}
	*/
	return nil
}

func DbTestWriteRequestV4Neighbor() error {
	logger.Println("Calling Store V4Neighbor - round ", writeRound)
	obj := models.IPv4Neighbor{IpAddr: "100.1.1." + strconv.Itoa(writeRound),
		MacAddr:  "11:22:33:44:55:66",
		VlanId:   100,
		RouterIf: 1}
	_, err := obj.StoreObjectInDb(gMgr.dbHdl)
	if err != nil {
		logger.Println("Failed to store V4Neighbor - round ", writeRound, err)
	}
	/*
		time.Sleep(1 * time.Second)
		logger.Println("Calling Delete V4Neighbor - round ", writeRound)
		objKey, _ := obj.GetKey()
		err = obj.DeleteObjectFromDb(objKey, gMgr.dbHdl)
		if err != nil {
			logger.Println("Failed to delete V4Neighbor - round ", writeRound, err)
		}
	*/
	return nil
}
