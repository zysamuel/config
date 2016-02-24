package main

import (
	"arpd"
	//"asicdServices"
	"bgpd"
	"database/sql"
	"models"
	"ribd"
	"strconv"
	"utils/commonDefs"
	"utils/ipcutils"
)

type ClientIf interface {
	Initialize(name string, address string)
	ConnectToServer() bool
	IsConnectedToServer() bool
	CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool)
	DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool
	GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
		objcount int64,
		nextMarker int64,
		more bool,
		objs []models.ConfigObj)
	UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool
	GetObject(obj models.ConfigObj) (models.ConfigObj, bool)
}

type RibClient struct {
	ipcutils.IPCClientBase
	ClientHdl *ribd.RouteServiceClient
}

func (clnt *RibClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}

func (clnt *RibClient) ConnectToServer() bool {

	if clnt.TTransport == nil && clnt.PtrProtocolFactory == nil {
		clnt.TTransport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	}
	if clnt.TTransport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = ribd.NewRouteServiceClientFactory(clnt.TTransport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}

func (clnt *RibClient) GetObject(obj models.ConfigObj) (models.ConfigObj, bool) {

	switch obj.(type) {

	case models.IPV4Route:
		var retObj models.IPV4Route
		data := obj.(models.IPV4Route)
		routeInfo, err := clnt.ClientHdl.GetRoute(data.DestinationNw, data.NetworkMask)
		if err == nil {
			retObj.DestinationNw = routeInfo.Ipaddr
			retObj.NetworkMask = routeInfo.Mask
			retObj.NextHopIp = routeInfo.NextHopIp
			retObj.Cost = uint32(routeInfo.Metric)
			retObj.Protocol = strconv.Itoa(int(routeInfo.Prototype))
			if routeInfo.NextHopIfType == commonDefs.L2RefTypeVlan {
				retObj.OutgoingIntfType = "VLAN"
			} else {
				retObj.OutgoingIntfType = "PHY"
			}
			retObj.OutgoingInterface = strconv.Itoa(int(routeInfo.IfIndex))
			return retObj, true
		}
		break
	default:
		break
	}
	return nil, false
}

func (clnt *RibClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {
	logger.Println("### Get Bulk request called with", currMarker, count)
	switch obj.(type) {
	case models.IPV4Route:
		if clnt.ClientHdl != nil {
			var ret_obj models.IPV4Route
			routesInfo, _ := clnt.ClientHdl.GetBulkRoutes(ribd.Int(currMarker), ribd.Int(count))
			if routesInfo.Count != 0 {
				objCount = int64(routesInfo.Count)
				more = bool(routesInfo.More)
				nextMarker = int64(routesInfo.EndIdx)
				for i := 0; i < int(routesInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}
					ret_obj.DestinationNw = routesInfo.RouteList[i].Ipaddr
					ret_obj.NetworkMask = routesInfo.RouteList[i].Mask
					ret_obj.NextHopIp = routesInfo.RouteList[i].NextHopIp
					ret_obj.Cost = uint32(routesInfo.RouteList[i].Metric)
					ret_obj.Protocol = routesInfo.RouteList[i].RoutePrototypeString //strconv.Itoa(int(routesInfo.RouteList[i].Prototype))
					if routesInfo.RouteList[i].NextHopIfType == commonDefs.L2RefTypeVlan {
						ret_obj.OutgoingIntfType = "VLAN"
					} else {
						ret_obj.OutgoingIntfType = "PHY"
					}
					ret_obj.OutgoingInterface = strconv.Itoa(int(routesInfo.RouteList[i].IfIndex))
					objs = append(objs, ret_obj)
				}
			}
		}
		break
	case models.IPV4RouteState:
		if clnt.ClientHdl != nil {
			var ret_obj models.IPV4RouteState
			routesInfo, _ := clnt.ClientHdl.GetBulkRoutes(ribd.Int(currMarker), ribd.Int(count))
			if routesInfo.Count != 0 {
				objCount = int64(routesInfo.Count)
				more = bool(routesInfo.More)
				nextMarker = int64(routesInfo.EndIdx)
				for i := 0; i < int(routesInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}
					ret_obj.DestinationNw = routesInfo.RouteList[i].Ipaddr
					ret_obj.NetworkMask = routesInfo.RouteList[i].Mask
					ret_obj.RouteCreatedTime = routesInfo.RouteList[i].RouteCreated
					ret_obj.RouteUpdatedTime = routesInfo.RouteList[i].RouteUpdated
					/*ret_obj.PolicyList = make([]string,0)
					        routePolicyListInfo := ""
					        if routesInfo.RouteList[i].PolicyList != nil {
					          for k,v := range routesInfo.RouteList[i].PolicyList {
						        routePolicyListInfo = k+":"
					            for vv:=0;vv<len(v);vv++ {
					              routePolicyListInfo = routePolicyListInfo + v[vv]+","
					            }
					            ret_obj.PolicyList = append(ret_obj.PolicyList,routePolicyListInfo)
					          }
					        }*/
					ret_obj.PolicyList = make([]string, 0)
					for j := 0; j < len(routesInfo.RouteList[i].PolicyList); j++ {
						ret_obj.PolicyList = append(ret_obj.PolicyList, routesInfo.RouteList[i].PolicyList[j])
					}
					objs = append(objs, ret_obj)
				}
			}
		}
		break
	case models.IPV4EventState:
		if clnt.ClientHdl != nil {
			var ret_obj models.IPV4EventState
			getBulkInfo, _ := clnt.ClientHdl.GetBulkIPV4EventState(ribd.Int(currMarker), ribd.Int(count))
			if getBulkInfo.Count != 0 {
				objCount = int64(getBulkInfo.Count)
				more = bool(getBulkInfo.More)
				nextMarker = int64(getBulkInfo.EndIdx)
				for i := 0; i < int(getBulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}
					ret_obj.TimeStamp = getBulkInfo.IPV4EventStateList[i].TimeStamp
					ret_obj.EventInfo = getBulkInfo.IPV4EventStateList[i].EventInfo
					objs = append(objs, ret_obj)
				}
			}
		}
		break
	case models.RouteDistanceState:
		if clnt.ClientHdl != nil {
			var ret_obj models.RouteDistanceState
			getBulkInfo, _ := clnt.ClientHdl.GetBulkRouteDistanceState(ribd.Int(currMarker), ribd.Int(count))
			if getBulkInfo.Count != 0 {
				objCount = int64(getBulkInfo.Count)
				more = bool(getBulkInfo.More)
				nextMarker = int64(getBulkInfo.EndIdx)
				for i := 0; i < int(getBulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}
					ret_obj.Distance = int(getBulkInfo.RouteDistanceStateList[i].Distance)
					ret_obj.Protocol = getBulkInfo.RouteDistanceStateList[i].Protocol
					objs = append(objs, ret_obj)
				}
			}
		}
		break

		/*    case models.PolicyDefinitionStmtMatchProtocolCondition:
		    logger.Println("PolicyDefinitionStmtMatchProtocolCondition")
			if clnt.ClientHdl != nil {
				var ret_obj models.PolicyDefinitionStmtMatchProtocolCondition
				getBulkInfo, _ := clnt.ClientHdl.GetBulkPolicyDefinitionStmtMatchProtocolConditions(ribd.Int(currMarker), ribd.Int(count))
				if getBulkInfo.Count != 0 {
					objCount = int64(getBulkInfo.Count)
					more = bool(getBulkInfo.More)
					nextMarker = int64(getBulkInfo.EndIdx)
					for i := 0; i < int(getBulkInfo.Count); i++ {
						if len(objs) == 0 {
							objs = make([]models.ConfigObj, 0)
						}
						ret_obj.Name = getBulkInfo.PolicyDefinitionStmtMatchProtocolConditionList[i].Name
						ret_obj.InstallProtocolEq = getBulkInfo.PolicyDefinitionStmtMatchProtocolConditionList[i].InstallProtocolEq
						objs = append(objs, ret_obj)
					}
				}
			}
		    break*/
	case models.PolicyConditionState:
		logger.Println("PolicyConditionState")
		if clnt.ClientHdl != nil {
			var ret_obj models.PolicyConditionState
			getBulkInfo, _ := clnt.ClientHdl.GetBulkPolicyConditionState(ribd.Int(currMarker), ribd.Int(count))
			if getBulkInfo.Count != 0 {
				objCount = int64(getBulkInfo.Count)
				more = bool(getBulkInfo.More)
				nextMarker = int64(getBulkInfo.EndIdx)
				for i := 0; i < int(getBulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}
					ret_obj.Name = getBulkInfo.PolicyConditionStateList[i].Name
					ret_obj.ConditionInfo = getBulkInfo.PolicyConditionStateList[i].ConditionInfo
					ret_obj.PolicyStmtList = make([]string, 0)
					for j := 0; j < len(getBulkInfo.PolicyConditionStateList[i].PolicyStmtList); j++ {
						ret_obj.PolicyStmtList = append(ret_obj.PolicyStmtList, getBulkInfo.PolicyConditionStateList[i].PolicyStmtList[j])
					}
					objs = append(objs, ret_obj)
				}
			}
		}
		break
		/*	case models.PolicyDefinitionStmtRedistributionAction:
			if clnt.ClientHdl != nil {
				var ret_obj models.PolicyDefinitionStmtRedistributionAction
				getBulkInfo, _ := clnt.ClientHdl.GetBulkPolicyDefinitionStmtRedistributionActions(ribd.Int(currMarker), ribd.Int(count))
				if getBulkInfo.Count != 0 {
					objCount = int64(getBulkInfo.Count)
					more = bool(getBulkInfo.More)
					nextMarker = int64(getBulkInfo.EndIdx)
					for i := 0; i < int(getBulkInfo.Count); i++ {
						if len(objs) == 0 {
							objs = make([]models.ConfigObj, 0)
						}
						ret_obj.Name = getBulkInfo.PolicyDefinitionStmtRedistributionActionList[i].Name
						ret_obj.RedistributeTargetProtocol = getBulkInfo.PolicyDefinitionStmtRedistributionActionList[i].RedistributeTargetProtocol
						objs = append(objs, ret_obj)
					}
				}
			}
		    break*/
	case models.PolicyActionState:
		logger.Println("PolicyActionState")
		if clnt.ClientHdl != nil {
			var ret_obj models.PolicyActionState
			getBulkInfo, _ := clnt.ClientHdl.GetBulkPolicyActionState(ribd.Int(currMarker), ribd.Int(count))
			if getBulkInfo.Count != 0 {
				objCount = int64(getBulkInfo.Count)
				more = bool(getBulkInfo.More)
				nextMarker = int64(getBulkInfo.EndIdx)
				for i := 0; i < int(getBulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}
					ret_obj.Name = getBulkInfo.PolicyActionStateList[i].Name
					ret_obj.ActionInfo = getBulkInfo.PolicyActionStateList[i].ActionInfo
					ret_obj.PolicyStmtList = make([]string, 0)
					for j := 0; j < len(getBulkInfo.PolicyActionStateList[i].PolicyStmtList); j++ {
						ret_obj.PolicyStmtList = append(ret_obj.PolicyStmtList, getBulkInfo.PolicyActionStateList[i].PolicyStmtList[j])
					}
					objs = append(objs, ret_obj)
				}
			}
		}
		break
	case models.PolicyStmtState:
		if clnt.ClientHdl != nil {
			var ret_obj models.PolicyStmtState
			getBulkInfo, _ := clnt.ClientHdl.GetBulkPolicyStmtState(ribd.Int(currMarker), ribd.Int(count))
			if getBulkInfo.Count != 0 {
				objCount = int64(getBulkInfo.Count)
				more = bool(getBulkInfo.More)
				nextMarker = int64(getBulkInfo.EndIdx)
				var j int
				for i := 0; i < int(getBulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}
					ret_obj.Name = getBulkInfo.PolicyStmtStateList[i].Name
					ret_obj.MatchConditions = getBulkInfo.PolicyStmtStateList[i].MatchConditions
					ret_obj.Conditions = make([]string, 0)
					for j = 0; j < len(getBulkInfo.PolicyStmtStateList[i].Conditions); j++ {
						ret_obj.Conditions = append(ret_obj.Conditions, getBulkInfo.PolicyStmtStateList[i].Conditions[j])
					}
					ret_obj.Actions = make([]string, 0)
					for j = 0; j < len(getBulkInfo.PolicyStmtStateList[i].Actions); j++ {
						ret_obj.Actions = append(ret_obj.Actions, getBulkInfo.PolicyStmtStateList[i].Actions[j])
					}
					objs = append(objs, ret_obj)
				}
			}
		}
		break
	case models.PolicyDefinitionState:
		if clnt.ClientHdl != nil {
			var ret_obj models.PolicyDefinitionState
			getBulkInfo, _ := clnt.ClientHdl.GetBulkPolicyDefinitionState(ribd.Int(currMarker), ribd.Int(count))
			if getBulkInfo.Count != 0 {
				objCount = int64(getBulkInfo.Count)
				more = bool(getBulkInfo.More)
				nextMarker = int64(getBulkInfo.EndIdx)
				var j int
				for i := 0; i < int(getBulkInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}
					ret_obj.Name = getBulkInfo.PolicyDefinitionStateList[i].Name
					ret_obj.HitCounter = int(getBulkInfo.PolicyDefinitionStateList[i].HitCounter)
					ret_obj.IpPrefixList = make([]string, 0)
					for j = 0; j < len(getBulkInfo.PolicyDefinitionStateList[i].IpPrefixList); j++ {
						ret_obj.IpPrefixList = append(ret_obj.IpPrefixList, getBulkInfo.PolicyDefinitionStateList[i].IpPrefixList[j])
					}
					objs = append(objs, ret_obj)
				}
			}
		}
		break
	}
	return nil, objCount, nextMarker, more, objs
}

func (clnt *RibClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	switch obj.(type) {
	case models.IPV4Route:
		v4Route := obj.(models.IPV4Route)
		outIntf, _ := strconv.Atoi(v4Route.OutgoingInterface)
		var outIntfType ribd.Int
		/*fix me - temporary hack for testing intf dis/ena*/
		/*	if v4Route.OutgoingIntfType == "DIS" {
				if clnt.ClientHdl != nil {
					clnt.ClientHdl.IntfDown("10.1.1.2/24")
				}
				if clnt.ClientHdl != nil {
					clnt.ClientHdl.IntfDown("30.1.1.2/24")
				}
			} else 	if v4Route.OutgoingIntfType == "ENA" {
				if clnt.ClientHdl != nil {
					clnt.ClientHdl.IntfUp("10.1.1.2/24")
				}
				if clnt.ClientHdl != nil {
					clnt.ClientHdl.IntfUp("30.1.1.2/24")
				}
			} else if v4Route.OutgoingIntfType == "VLAN" {
			/* End of hack*/
		if v4Route.OutgoingIntfType == "VLAN" {
			outIntfType = commonDefs.L2RefTypeVlan
		} else {
			outIntfType = commonDefs.L2RefTypePort
		}
		//proto, _ := strconv.Atoi(v4Route.Protocol)
		if clnt.ClientHdl != nil {
			clnt.ClientHdl.CreateV4Route(
				v4Route.DestinationNw, //ribd.Int(binary.BigEndian.Uint32(net.ParseIP(v4Route.DestinationNw).To4())),
				v4Route.NetworkMask,   //ribd.Int(prefixLen),
				ribd.Int(v4Route.Cost),
				v4Route.NextHopIp, //ribd.Int(binary.BigEndian.Uint32(net.ParseIP(v4Route.NextHopIp).To4())),
				outIntfType,
				ribd.Int(outIntf),
				v4Route.Protocol)
			//ribd.Int(proto))
		}
		objId, _ := v4Route.StoreObjectInDb(dbHdl)
		return objId, true
		/*	case models.PolicyDefinitionStmtDstIpMatchPrefixSetCondition:
			logger.Println("PolicyDefinitionStmtDstIpMatchPrefixSetCondition")
			inCfg := obj.(models.PolicyDefinitionStmtDstIpMatchPrefixSetCondition)
			var cfg ribd.PolicyDefinitionStmtDstIpMatchPrefixSetCondition
			if len(inCfg.PrefixSet) > 0 && len(inCfg.Prefix.IpPrefix) > 0 {
				logger.Println("cannot set both prefix set name and a prefix")
				return int64(0), true
			}
			cfg.Name = inCfg.Name
			cfg.PrefixSet = inCfg.PrefixSet
			var cfgIpPrefix ribd.PolicyDefinitionSetsPrefix
			cfgIpPrefix.IpPrefix = inCfg.Prefix.IpPrefix
			cfgIpPrefix.MasklengthRange = inCfg.Prefix.MaskLengthRange
			cfg.Prefix = &cfgIpPrefix
			if clnt.ClientHdl != nil {
				clnt.ClientHdl.CreatePolicyDefinitionStmtDstIpMatchPrefixSetCondition(&cfg)
			}
			objId, _ := inCfg.StoreObjectInDb(dbHdl)
			return objId, true*/
	case models.PolicyPrefixSet:
		logger.Println("PolicyPrefixSet")
		inCfg := obj.(models.PolicyPrefixSet)
		var cfg ribd.PolicyPrefixSet
		//cfg.PrefixSetName = inCfg.PrefixSetName
		//ipPrefixList := strings.Split(inCfg.IpPrefix, ",")
		logger.Println("ipPrefixList len = ", len(inCfg.IpPrefixList))
		cfgIpPrefixList := make([]*ribd.PolicyPrefix, 0)
		cfgIpPrefix := make([]ribd.PolicyPrefix, len(inCfg.IpPrefixList))
		for i := 0; i < len(inCfg.IpPrefixList); i++ {
			cfgIpPrefix[i].IpPrefix = inCfg.IpPrefixList[i].IpPrefix
			cfgIpPrefix[i].MasklengthRange = inCfg.IpPrefixList[i].MaskLengthRange
			cfgIpPrefixList = append(cfgIpPrefixList, &cfgIpPrefix[i])
		}
		cfg.IpPrefixList = cfgIpPrefixList
		if clnt.ClientHdl != nil {
			clnt.ClientHdl.CreatePolicyPrefixSet(&cfg)
		}
		objId, _ := inCfg.StoreObjectInDb(dbHdl)
		return objId, true
		/*	case models.PolicyDefinitionStmtMatchProtocolCondition:
			logger.Println("PolicyDefinitionStmtMatchProtocolCondition")
			inCfg := obj.(models.PolicyDefinitionStmtMatchProtocolCondition)
			var cfg ribd.PolicyDefinitionStmtMatchProtocolCondition
			cfg.Name = inCfg.Name
			cfg.InstallProtocolEq = inCfg.InstallProtocolEq
			if clnt.ClientHdl != nil {
				clnt.ClientHdl.CreatePolicyDefinitionStmtMatchProtocolCondition(&cfg)
			}
			objId, _ := inCfg.StoreObjectInDb(dbHdl)
			return objId, true*/
	case models.PolicyConditionConfig:
		logger.Println("PolicyConditionConfig")
		inCfg := obj.(models.PolicyConditionConfig)
		var cfg ribd.PolicyConditionConfig
		cfg.Name = inCfg.Name
		cfg.ConditionType = inCfg.ConditionType
		switch inCfg.ConditionType {
		case "MatchProtocol":
			logger.Println("MatchProtocol ", inCfg.MatchProtocolConditionInfo)
			matchProto := inCfg.MatchProtocolConditionInfo
			cfg.MatchProtocolConditionInfo = &matchProto
			//dstIpMatchPrefixconditionCfg.Prefix = &cfgIpPrefix
			//cfg.MatchDstIpPrefixConditionInfo = &dstIpMatchPrefixconditionCfg
			break
		case "MatchDstIpPrefix":
			logger.Println("MatchDstIpPrefix")
			inConditionCfg := inCfg.MatchDstIpPrefixConditionInfo
			var cfgIpPrefix ribd.PolicyPrefix
			var dstIpMatchPrefixconditionCfg ribd.PolicyDstIpMatchPrefixSetCondition
			if len(inConditionCfg.PrefixSet) > 0 && len(inConditionCfg.Prefix.IpPrefix) > 0 {
				logger.Println("cannot set both prefix set name and a prefix")
				return int64(0), true
			}
			dstIpMatchPrefixconditionCfg.PrefixSet = inConditionCfg.PrefixSet
			cfgIpPrefix.IpPrefix = inConditionCfg.Prefix.IpPrefix
			cfgIpPrefix.MasklengthRange = inConditionCfg.Prefix.MaskLengthRange
			dstIpMatchPrefixconditionCfg.Prefix = &cfgIpPrefix
			cfg.MatchDstIpPrefixConditionInfo = &dstIpMatchPrefixconditionCfg
			break
		default:
			logger.Println("Invalid condition type")
			return int64(0), true
		}
		if clnt.ClientHdl != nil {
			clnt.ClientHdl.CreatePolicyCondition(&cfg)
		}
		objId, _ := inCfg.StoreObjectInDb(dbHdl)
		return objId, true
		/*	case models.PolicyDefinitionStmtRedistributionAction:
				logger.Println("PolicyDefinitionStmtRedistributionAction")
				inCfg := obj.(models.PolicyDefinitionStmtRedistributionAction)
				var cfg ribd.PolicyDefinitionStmtRedistributionAction
				cfg.Name = inCfg.Name
				cfg.RedistributeTargetProtocol = inCfg.RedistributeTargetProtocol
				cfg.Redistribute = inCfg.Redistribute
				if clnt.ClientHdl != nil {
					clnt.ClientHdl.CreatePolicyDefinitionStmtRedistributionAction(&cfg)
				}
				objId, _ := inCfg.StoreObjectInDb(dbHdl)
				return objId, true
			case models.PolicyDefinitionStmtRouteDispositionAction:
				logger.Println("PolicyDefinitionStmtRouteDispositionAction")
				inCfg := obj.(models.PolicyDefinitionStmtRouteDispositionAction)
				var cfg ribd.PolicyDefinitionStmtRouteDispositionAction
				cfg.Name = inCfg.Name
				cfg.RouteDisposition = inCfg.RouteDisposition
				if clnt.ClientHdl != nil {
					clnt.ClientHdl.CreatePolicyDefinitionStmtRouteDispositionAction(&cfg)
				}
				objId, _ := inCfg.StoreObjectInDb(dbHdl)
				return objId, true
			case models.PolicyDefinitionStmtAdminDistanceAction:
				logger.Println("PolicyDefinitionStmtAdminDistanceAction")
				inCfg := obj.(models.PolicyDefinitionStmtAdminDistanceAction)
				var cfg ribd.PolicyDefinitionStmtAdminDistanceAction
				cfg.Name = inCfg.Name
				cfg.Value = ribd.Int(inCfg.Value)
				if clnt.ClientHdl != nil {
					clnt.ClientHdl.CreatePolicyDefinitionStmtAdminDistanceAction(&cfg)
				}
				objId, _ := inCfg.StoreObjectInDb(dbHdl)
				return objId, true*/
	case models.PolicyActionConfig:
		logger.Println("PolicyActionConfig")
		inCfg := obj.(models.PolicyActionConfig)
		var cfg ribd.PolicyActionConfig
		cfg.Name = inCfg.Name
		cfg.ActionType = inCfg.ActionType
		switch inCfg.ActionType {
		case "RouteDisposition":
			logger.Println("RouteDisposition")
			cfg.Accept = inCfg.Accept
			cfg.Reject = inCfg.Reject
			if inCfg.Accept && inCfg.Reject {
				logger.Println("Cannot set both accept and reject actions to true")
				return int64(0), true
			}
			break
		case "Redistribution":
			logger.Println("Redistribution")
			inActionCfg := inCfg.RedistributeActionInfo
			var actionCfg ribd.PolicyRedistributionAction
			actionCfg.RedistributeTargetProtocol = inActionCfg.RedistributeTargetProtocol
			actionCfg.Redistribute = inActionCfg.Redistribute
			cfg.RedistributeActionInfo = &actionCfg
			break
		case "SetAdminDistance":
			logger.Println("SetSdminDistance to inCfg.SetAdminDistanceValue")
			cfg.SetAdminDistanceValue = ribd.Int(inCfg.SetAdminDistanceValue)
			break
		}
		if clnt.ClientHdl != nil {
			clnt.ClientHdl.CreatePolicyAction(&cfg)
		}
		objId, _ := inCfg.StoreObjectInDb(dbHdl)
		return objId, true
	case models.PolicyStmtConfig:
		logger.Println("PolicyStmtConfig")
		var i int
		inCfg := obj.(models.PolicyStmtConfig)
		var cfg ribd.PolicyStmtConfig
		cfg.Name = inCfg.Name
		logger.Println("Number of conditons = ", len(inCfg.Conditions))
		conditions := make([]string, 0)
		for i = 0; i < len(inCfg.Conditions); i++ {
			conditions = append(conditions, inCfg.Conditions[i])
		}
		cfg.Conditions = conditions
		logger.Println("Number of actions = ", len(inCfg.Actions))
		actions := make([]string, 0)
		for i = 0; i < len(inCfg.Actions); i++ {
			actions = append(actions, inCfg.Actions[i])
		}
		cfg.Actions = actions
		cfg.MatchConditions = inCfg.MatchConditions
		if clnt.ClientHdl != nil {
			clnt.ClientHdl.CreatePolicyStatement(&cfg)
		}
		objId, _ := inCfg.StoreObjectInDb(dbHdl)
		return objId, true
	case models.PolicyDefinitionConfig:
		logger.Println("PolicyDefinitionConfig")
		inCfg := obj.(models.PolicyDefinitionConfig)
		var cfg ribd.PolicyDefinitionConfig
		cfg.Name = inCfg.Name
		cfg.Precedence = ribd.Int(inCfg.Precedence)
		cfg.MatchType = inCfg.MatchType
		cfg.Export = inCfg.Export
		cfg.Import = inCfg.Import
		cfg.Global = inCfg.Global
		if inCfg.Import == false && inCfg.Export == false && inCfg.Global == false {
			logger.Println("Need to set import,export or global to true")
			break
		}
		logger.Println("Number of statements = ", len(inCfg.StatementList))
		policyDefinitionStatements := make([]ribd.PolicyDefinitionStmtPrecedence, len(inCfg.StatementList))
		cfg.PolicyDefinitionStatements = make([]*ribd.PolicyDefinitionStmtPrecedence, 0)
		var i int
		for k, v := range inCfg.StatementList {
			logger.Println("k= ", k, " v= ", v)
			if v == nil {
				logger.Println("Interface nil at key ", k)
				continue
			}
			inCfgStatementIf := v.(map[string]interface{}) //models.PolicyDefinitionStmtPrecedence)
			policyDefinitionStatements[i] = ribd.PolicyDefinitionStmtPrecedence{Precedence: ribd.Int(inCfgStatementIf["Precedence"].(float64)), Statement: inCfgStatementIf["Statement"].(string)}
			cfg.PolicyDefinitionStatements = append(cfg.PolicyDefinitionStatements, &policyDefinitionStatements[i])
			i++
		}
		if clnt.ClientHdl != nil {
			clnt.ClientHdl.CreatePolicyDefinition(&cfg)
		}
		objId, _ := inCfg.StoreObjectInDb(dbHdl)
		return objId, true
		break
	default:
		break
	}
	return int64(0), true
}

func (clnt *RibClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {
	logger.Println("### Delete Object is called in RIBClient. ObjectKey: ", objKey, obj)
	switch obj.(type) {
	case models.IPV4Route:
		v4Route := obj.(models.IPV4Route)
		logger.Println("### DeleteV4Route is called in RIBClient. ", v4Route.DestinationNw, v4Route.NetworkMask, v4Route.OutgoingInterface)
		if clnt.ClientHdl != nil {
			clnt.ClientHdl.DeleteV4Route(
				v4Route.DestinationNw, //ribd.Int(binary.BigEndian.Uint32(net.ParseIP(v4Route.DestinationNw).To4())),
				v4Route.NetworkMask,   //ribd.Int(prefixLen),
				v4Route.Protocol,
				v4Route.NextHopIp)
		}
		v4Route.DeleteObjectFromDb(objKey, dbHdl)
		break
	case models.PolicyStmtConfig:
		logger.Println("PolicyStmtConfig")
		inCfg := obj.(models.PolicyStmtConfig)
		var cfg ribd.PolicyStmtConfig
		cfg.Name = inCfg.Name
		if clnt.ClientHdl != nil {
			clnt.ClientHdl.DeletePolicyStatement(&cfg)
		}
		inCfg.DeleteObjectFromDb(objKey, dbHdl)
		break
	case models.PolicyDefinitionConfig:
		logger.Println("PolicyDefinition")
		inCfg := obj.(models.PolicyDefinitionConfig)
		var cfg ribd.PolicyDefinitionConfig
		cfg.Name = inCfg.Name
		if clnt.ClientHdl != nil {
			clnt.ClientHdl.DeletePolicyDefinition(&cfg)
		}
		inCfg.DeleteObjectFromDb(objKey, dbHdl)
		break

		//default:
		//	logger.Println("OBJECT Type is ", obj.(type))
	}

	return true
}

func (clnt *RibClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {
	logger.Println("### Update Object is called in RIBClient. ", objKey, dbObj, obj, attrSet)
	switch obj.(type) {
	case models.IPV4Route:
		v4Route := obj.(models.IPV4Route)
		outIntf, _ := strconv.Atoi(v4Route.OutgoingInterface)
		logger.Println("### UpdateV4Route is called in RIBClient. ", v4Route.DestinationNw, v4Route.NetworkMask, outIntf)
		/*
			if clnt.ClientHdl != nil {
				clnt.ClientHdl.UpdateV4Route(
					dbObj,
					obj,
					attrSet)
			}
		*/
		v4Route.UpdateObjectInDb(dbObj, attrSet, dbHdl)
		//default:
		//	logger.Println("OBJECT Type is ", obj.(type))
	}
	return true
}

/*
type AsicDClient struct {
	ipcutils.IPCClientBase
	ClientHdl *asicdServices.ASICDServicesClient
}

func (clnt *AsicDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}

func (clnt *AsicDClient) ConnectToServer() bool {
	if clnt.TTransport == nil && clnt.PtrProtocolFactory == nil {
		clnt.TTransport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	}
	if clnt.TTransport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = asicdServices.NewASICDServicesClientFactory(clnt.TTransport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}

func (clnt *AsicDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	var objId int64
	if clnt.ClientHdl != nil {
		switch obj.(type) {
		case models.VlanConfig: //Vlan
			vlanObj := obj.(models.VlanConfig)
			_, err := clnt.ClientHdl.CreateVlan(vlanObj.VlanId, vlanObj.IfIndexList, vlanObj.UntagIfIndexList)
			if err != nil {
				return int64(0), false
			}
			objId, _ = vlanObj.StoreObjectInDb(dbHdl)
			return objId, true
		case models.IPv4Intf: //IPv4Intf
			v4Intf := obj.(models.IPv4Intf)
			_, err := clnt.ClientHdl.CreateIPv4Intf(v4Intf.IpAddr, v4Intf.IfIndex)
			if err != nil {
				return int64(0), false
			}
			objId, _ = v4Intf.StoreObjectInDb(dbHdl)
			return objId, true
		}
	}
	return int64(0), true
}

func (clnt *AsicDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {
	return true
}

func (clnt *AsicDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {

		if clnt.ClientHdl != nil {
			switch obj.(type) {
			case models.PortintfConfig:
				portIntfObj := obj.(models.PortIntfConfig)
				clnt.ClientHdl.UpatePortIntfConfig(dbObj, obj, attrSet)
			}
		}

	return true
}

func (clnt *AsicDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error, objCount int64,
	nextMarker int64, more bool, objs []models.ConfigObj) {
	switch obj.(type) {
	case models.PortConfig:
		portConfigBulk, err := clnt.ClientHdl.GetBulkPortConfig(currMarker, count)
		if err != nil {
			break
		}
		for _, elem := range portConfigBulk.PortConfigList {
			portConfig := models.PortConfig{
				IfIndex:     elem.IfIndex,
				Name:        elem.Name,
				Description: elem.Description,
				Type:        elem.Type,
				AdminState:  elem.AdminState,
				OperState:   elem.OperState,
				MacAddr:     elem.MacAddr,
				Speed:       elem.Speed,
				Duplex:      elem.Duplex,
				Autoneg:     elem.Autoneg,
				MediaType:   elem.MediaType,
				Mtu:         elem.Mtu,
			}
			objs = append(objs, portConfig)
		}
		objCount = portConfigBulk.ObjCount
		nextMarker = portConfigBulk.NextMarker
		more = portConfigBulk.More

	case models.PortState:
		portStateBulk, err := clnt.ClientHdl.GetBulkPortState(currMarker, count)
		if err != nil {
			break
		}
		for _, elem := range portStateBulk.PortStateList {
			portState := models.PortState{
				IfIndex:   elem.IfIndex,
				PortStats: elem.Stats,
			}
			objs = append(objs, portState)
		}
		objCount = portStateBulk.ObjCount
		nextMarker = portStateBulk.NextMarker
		more = portStateBulk.More

	case models.VlanState:
		vlanBulk, err := clnt.ClientHdl.GetBulkVlan(currMarker, count)
		if err != nil {
			break
		}
		for _, elem := range vlanBulk.VlanObjList {
			vlanState := models.VlanState{
				VlanId:    elem.VlanId,
				IfIndex:   elem.IfIndex,
				VlanName:  elem.VlanName,
				OperState: elem.OperState,
			}
			objs = append(objs, vlanState)
		}
		objCount = vlanBulk.ObjCount
		nextMarker = vlanBulk.NextMarker
		more = vlanBulk.More
	}
	return err, objCount, nextMarker, more, objs
}
*/

type BgpDClient struct {
	ipcutils.IPCClientBase
	ClientHdl *bgpd.BGPServerClient
}

func (clnt *BgpDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}

func (clnt *BgpDClient) ConnectToServer() bool {
	if clnt.TTransport == nil && clnt.PtrProtocolFactory == nil {
		clnt.TTransport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	}
	if clnt.TTransport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = bgpd.NewBGPServerClientFactory(clnt.TTransport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}

func convertBGPGlobalConfToThriftObj(bgpGlobalConf models.BGPGlobalConfig) *bgpd.BGPGlobalConfig {
	gConf := bgpd.NewBGPGlobalConfig()
	gConf.ASNum = int32(bgpGlobalConf.ASNum)
	gConf.RouterId = bgpGlobalConf.RouterId
	gConf.UseMultiplePaths = bgpGlobalConf.UseMultiplePaths
	gConf.EBGPMaxPaths = int32(bgpGlobalConf.EBGPMaxPaths)
	gConf.EBGPAllowMultipleAS = bgpGlobalConf.EBGPAllowMultipleAS
	gConf.IBGPMaxPaths = int32(bgpGlobalConf.IBGPMaxPaths)
	return gConf
}

func convertBGPNeighborConfToThriftObj(bgpNeighborConf models.BGPNeighborConfig) *bgpd.BGPNeighborConfig {
	nConf := bgpd.NewBGPNeighborConfig()
	nConf.PeerAS = int32(bgpNeighborConf.PeerAS)
	nConf.LocalAS = int32(bgpNeighborConf.LocalAS)
	nConf.NeighborAddress = bgpNeighborConf.NeighborAddress
	nConf.Description = bgpNeighborConf.Description
	nConf.RouteReflectorClusterId = int32(bgpNeighborConf.RouteReflectorClusterId)
	nConf.RouteReflectorClient = bgpNeighborConf.RouteReflectorClient
	nConf.MultiHopEnable = bgpNeighborConf.MultiHopEnable
	nConf.MultiHopTTL = int8(bgpNeighborConf.MultiHopTTL)
	nConf.ConnectRetryTime = int32(bgpNeighborConf.ConnectRetryTime)
	nConf.HoldTime = int32(bgpNeighborConf.HoldTime)
	nConf.KeepaliveTime = int32(bgpNeighborConf.KeepaliveTime)
	nConf.PeerGroup = bgpNeighborConf.PeerGroup
	return nConf
}

func convertBGPPeerGroupToThriftObj(bgpPeerGroup models.BGPPeerGroup) *bgpd.BGPPeerGroup {
	peerGroup := bgpd.NewBGPPeerGroup()
	peerGroup.PeerAS = int32(bgpPeerGroup.PeerAS)
	peerGroup.LocalAS = int32(bgpPeerGroup.LocalAS)
	peerGroup.Name = bgpPeerGroup.Name
	peerGroup.Description = bgpPeerGroup.Description
	peerGroup.RouteReflectorClusterId = int32(bgpPeerGroup.RouteReflectorClusterId)
	peerGroup.RouteReflectorClient = bgpPeerGroup.RouteReflectorClient
	peerGroup.MultiHopEnable = bgpPeerGroup.MultiHopEnable
	peerGroup.MultiHopTTL = int8(bgpPeerGroup.MultiHopTTL)
	peerGroup.ConnectRetryTime = int32(bgpPeerGroup.ConnectRetryTime)
	peerGroup.HoldTime = int32(bgpPeerGroup.HoldTime)
	peerGroup.KeepaliveTime = int32(bgpPeerGroup.KeepaliveTime)
	return peerGroup
}

func (clnt *BgpDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	retVal := false
	objId := int64(0)

	if clnt.ClientHdl != nil {
		switch obj.(type) {
		case models.BGPGlobalConfig:
			bgpGlobalConf := obj.(models.BGPGlobalConfig)
			gConf := convertBGPGlobalConfToThriftObj(bgpGlobalConf)
			_, err := clnt.ClientHdl.CreateBGPGlobal(gConf)
			if err != nil {
				return int64(0), false
			}
			objId, _ = bgpGlobalConf.StoreObjectInDb(dbHdl)
			retVal = true

		case models.BGPNeighborConfig:
			bgpNeighborConf := obj.(models.BGPNeighborConfig)
			nConf := convertBGPNeighborConfToThriftObj(bgpNeighborConf)
			_, err := clnt.ClientHdl.CreateBGPNeighbor(nConf)
			if err != nil {
				return int64(0), false
			}
			objId, _ = bgpNeighborConf.StoreObjectInDb(dbHdl)
			retVal = true

		case models.BGPPeerGroup:
			bgpPeerGroup := obj.(models.BGPPeerGroup)
			peerGroup := convertBGPPeerGroupToThriftObj(bgpPeerGroup)
			_, err := clnt.ClientHdl.CreateBGPPeerGroup(peerGroup)
			if err != nil {
				return int64(0), false
			}
			objId, _ = bgpPeerGroup.StoreObjectInDb(dbHdl)
			retVal = true
		}
	}

	return objId, retVal
}

func (clnt *BgpDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error, objCount int64,
	nextMarker int64, more bool, objs []models.ConfigObj) {

	logger.Println("BgpDClient: GetBulkObject called - start")
	switch obj.(type) {
	case models.BGPNeighborState:
		var bgpNeighborStateBulk *bgpd.BGPNeighborStateBulk
		bgpNeighborStateBulk, err = clnt.ClientHdl.BulkGetBGPNeighbors(currMarker, count)
		if err != nil {
			break
		}

		for _, item := range bgpNeighborStateBulk.StateList {
			bgpNeighborState := models.BGPNeighborState{
				PeerAS:                  uint32(item.PeerAS),
				LocalAS:                 uint32(item.LocalAS),
				PeerType:                models.PeerType(item.PeerType),
				AuthPassword:            item.AuthPassword,
				Description:             item.Description,
				NeighborAddress:         item.NeighborAddress,
				SessionState:            uint32(item.SessionState),
				RouteReflectorClusterId: uint32(item.RouteReflectorClusterId),
				RouteReflectorClient:    item.RouteReflectorClient,
				MultiHopEnable:          item.MultiHopEnable,
				MultiHopTTL:             uint8(item.MultiHopTTL),
				ConnectRetryTime:        uint32(item.ConnectRetryTime),
				HoldTime:                uint32(item.HoldTime),
				KeepaliveTime:           uint32(item.KeepaliveTime),
				Messages: models.BGPMessages{
					Sent: models.BgpCounters{
						Update:       uint64(item.Messages.Sent.Update),
						Notification: uint64(item.Messages.Sent.Notification),
					},
					Received: models.BgpCounters{
						Update:       uint64(item.Messages.Received.Update),
						Notification: uint64(item.Messages.Received.Notification),
					},
				},
				Queues: models.BGPQueues{
					Input:  uint32(item.Queues.Input),
					Output: uint32(item.Queues.Output),
				},
			}
			objs = append(objs, bgpNeighborState)
		}
		nextMarker = bgpNeighborStateBulk.NextIndex
		objCount = bgpNeighborStateBulk.Count
		more = bgpNeighborStateBulk.More

	case models.BGPRoute:
		var bgpRouteBulk *bgpd.BGPRouteBulk
		bgpRouteBulk, err = clnt.ClientHdl.BulkGetBGPRoutes(currMarker, count)
		if err != nil {
			break
		}

		for _, item := range bgpRouteBulk.RouteList {
			path := make([]uint32, len(item.Path))
			for idx, elem := range item.Path {
				path[idx] = uint32(elem)
			}

			bgpRoute := models.BGPRoute{
				Network:   item.Network,
				Mask:      item.Mask,
				NextHop:   item.NextHop,
				Metric:    uint32(item.Metric),
				LocalPref: uint32(item.LocalPref),
				Path:      path,
				Updated:   item.Updated,
			}
			objs = append(objs, bgpRoute)
		}
		nextMarker = bgpRouteBulk.NextIndex
		objCount = bgpRouteBulk.Count
		more = bgpRouteBulk.More
	}
	return err, objCount, nextMarker, more, objs
}

func (clnt *BgpDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {
	if clnt.ClientHdl != nil {
		switch obj.(type) {
		case models.BGPGlobalConfig:
			return false

		case models.BGPNeighborConfig:
			logger.Println("BgpDClient: BGPNeighborConfig delete")
			bgpNeighborConf := obj.(models.BGPNeighborConfig)
			logger.Println("BgpDClient: BGPNeighborConfig delete - %s", bgpNeighborConf)
			_, err := clnt.ClientHdl.DeleteBGPNeighbor(bgpNeighborConf.NeighborAddress)
			if err != nil {
				return false
			}
			bgpNeighborConf.DeleteObjectFromDb(objKey, dbHdl)

		case models.BGPPeerGroup:
			logger.Println("BgpDClient: BGPPeerGroup delete")
			bgpPeerGroup := obj.(models.BGPPeerGroup)
			logger.Println("BgpDClient: BGPPeerGroup delete - %s", bgpPeerGroup)
			_, err := clnt.ClientHdl.DeleteBGPPeerGroup(bgpPeerGroup.Name)
			if err != nil {
				return false
			}
			bgpPeerGroup.DeleteObjectFromDb(objKey, dbHdl)

		default:
			return false
		}
	}
	return true
}

func (clnt *BgpDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {
	if clnt.ClientHdl != nil {
		switch obj.(type) {
		case models.BGPGlobalConfig:
			logger.Println("BgpDClient: BGPGlobalConfig update")
			origBgpGlobalConf := dbObj.(models.BGPGlobalConfig)
			origGConf := convertBGPGlobalConfToThriftObj(origBgpGlobalConf)
			updatedBgpGlobalConf := obj.(models.BGPGlobalConfig)
			updatedGConf := convertBGPGlobalConfToThriftObj(updatedBgpGlobalConf)
			_, err := clnt.ClientHdl.UpdateBGPGlobal(origGConf, updatedGConf, attrSet)
			if err != nil {
				return false
			}
			origBgpGlobalConf.UpdateObjectInDb(obj, attrSet, dbHdl)

		case models.BGPNeighborConfig:
			logger.Println("BgpDClient: BGPNeighborConfig update")
			origBgpNeighborConf := obj.(models.BGPNeighborConfig)
			origNConf := convertBGPNeighborConfToThriftObj(origBgpNeighborConf)
			updatedBgpNeighborConf := obj.(models.BGPNeighborConfig)
			updatedNConf := convertBGPNeighborConfToThriftObj(updatedBgpNeighborConf)
			_, err := clnt.ClientHdl.UpdateBGPNeighbor(origNConf, updatedNConf, attrSet)
			if err != nil {
				return false
			}
			origBgpNeighborConf.UpdateObjectInDb(obj, attrSet, dbHdl)

		case models.BGPPeerGroup:
			logger.Println("BgpDClient: BGPPeerGroup update")
			origBgpPeerGroup := obj.(models.BGPPeerGroup)
			origGroup := convertBGPPeerGroupToThriftObj(origBgpPeerGroup)
			updatedBgpPeerGroup := obj.(models.BGPPeerGroup)
			updatedGroup := convertBGPPeerGroupToThriftObj(updatedBgpPeerGroup)
			_, err := clnt.ClientHdl.UpdateBGPPeerGroup(origGroup, updatedGroup, attrSet)
			if err != nil {
				return false
			}
			origBgpPeerGroup.UpdateObjectInDb(obj, attrSet, dbHdl)

		default:
			return false
		}
	}
	return true
}

type ArpDClient struct {
	ipcutils.IPCClientBase
	ClientHdl *arpd.ARPDServicesClient
}

func (clnt *ArpDClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}

func (clnt *ArpDClient) ConnectToServer() bool {
	if clnt.TTransport == nil && clnt.PtrProtocolFactory == nil {
		clnt.TTransport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	}
	if clnt.TTransport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = arpd.NewARPDServicesClientFactory(clnt.TTransport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}

func (clnt *ArpDClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
	logger.Println("ArpDClient: CreateObject called - start")
	if clnt.ClientHdl != nil {
		switch obj.(type) {
		case models.ArpConfig: //Arp Timeout
			arpConfigObj := obj.(models.ArpConfig)
			_, err := clnt.ClientHdl.SetArpConfig(arpd.Int(arpConfigObj.Timeout))
			if err != nil {
				return int64(0), false
			}
		}
	}
	return int64(0), true
}

func (clnt *ArpDClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {
	return true
}

func (clnt *ArpDClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {
	return true
}

func (clnt *ArpDClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error, objCount int64,
	nextMarker int64, more bool, objs []models.ConfigObj) {

	logger.Println("ArpDClient: GetBulkObject called - start")
	var ret_obj models.ArpEntry
	switch obj.(type) {
	case models.ArpEntry:
		if clnt.ClientHdl != nil {
			arpEntryBulk, err := clnt.ClientHdl.GetBulkArpEntry(arpd.Int(currMarker), arpd.Int(count))
			if err != nil {
				logger.Println("GetBulkObject call to Arpd failed:", err)
				return nil, objCount, nextMarker, more, objs
			}
			if arpEntryBulk.Count != 0 {
				objCount = int64(arpEntryBulk.Count)
				more = arpEntryBulk.More
				nextMarker = int64(arpEntryBulk.EndIdx)
				cnt := int(arpEntryBulk.Count)
				for i := 0; i < cnt; i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}
					ret_obj.IpAddr = arpEntryBulk.ArpList[i].IpAddr
					ret_obj.MacAddr = arpEntryBulk.ArpList[i].MacAddr
					ret_obj.Vlan = uint32(arpEntryBulk.ArpList[i].Vlan)
					ret_obj.Intf = arpEntryBulk.ArpList[i].Intf
					ret_obj.ExpiryTimeLeft = arpEntryBulk.ArpList[i].ExpiryTimeLeft
					objs = append(objs, ret_obj)
				}
			}
		}
	}
	return nil, objCount, nextMarker, more, objs
}

func (clnt *ArpDClient) GetObject(obj models.ConfigObj) (models.ConfigObj, bool) {
	return nil, false
}

func (clnt *ASICDClient) GetObject(obj models.ConfigObj) (models.ConfigObj, bool) {
	return nil, false
}

func (clnt *BgpDClient) GetObject(obj models.ConfigObj) (models.ConfigObj, bool) {
	return nil, false
}

func (clnt *LACPDClient) GetObject(obj models.ConfigObj) (models.ConfigObj, bool) {
	return nil, false
}

func (clnt *DHCPRELAYDClient) GetObject(obj models.ConfigObj) (models.ConfigObj, bool) {
	return nil, false
}

func (clnt *LocalClient) GetObject(obj models.ConfigObj) (models.ConfigObj, bool) {
	return nil, false
}
