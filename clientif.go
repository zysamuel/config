package main

import (
	//"asicdServices"
	"database/sql"
	"models"
	//	"strconv"
	//	"utils/commonDefs"
	//	"utils/ipcutils"
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

/*type RIBDIntClient struct {
	ipcutils.IPCClientBase
	ClientHdl *ribd.RIBDServicesClient
}

func (clnt *RIBDIntClient) Initialize(name string, address string) {
	clnt.Address = address
	return
}
func (clnt *RIBDIntClient) ConnectToServer() bool {

	clnt.TTransport, clnt.PtrProtocolFactory, _ = ipcutils.CreateIPCHandles(clnt.Address)
	if clnt.TTransport != nil && clnt.PtrProtocolFactory != nil {
		clnt.ClientHdl = ribd.NewRIBDServicesClientFactory(clnt.TTransport, clnt.PtrProtocolFactory)
		if clnt.ClientHdl != nil {
			clnt.IsConnected = true
		} else {
			clnt.IsConnected = false
		}
	}
	return true
}

func (clnt *RIBDIntClient) GetObject(obj models.ConfigObj) (models.ConfigObj, bool) {

	switch obj.(type) {

	default:
		break
	}
	return nil, false
}

func (clnt *RIBDIntClient) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {
	logger.Println("### Get Bulk request called with", currMarker, count)
	switch obj.(type) {
	case models.IPv4Route:
		if clnt.ClientHdl != nil {
			var ret_obj models.IPv4Route
			routesInfo, _ := clnt.ClientHdl.GetBulkRoutes(ribdInt.Int(currMarker), ribdInt.Int(count))
			if routesInfo.Count != 0 {
				objCount = int64(routesInfo.Count)
				more = bool(routesInfo.More)
				nextMarker = int64(routesInfo.EndIdx)
				for i := 0; i < int(routesInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}
					ret_obj.DestinationNw = routesInfo.RouteList[i].DestNetIp
					ret_obj.NextHopIp = routesInfo.RouteList[i].NextHopIp
					ret_obj.Cost = uint32(routesInfo.RouteList[i].Metric)
					ret_obj.Protocol = routesInfo.RouteList[i].RoutePrototypeString //strconv.Itoa(int(routesInfo.RouteList[i].Prototype))
					ret_obj.OutgoingInterface = strconv.Itoa(int(routesInfo.RouteList[i].IfIndex))
					if routesInfo.RouteList[i].NextHopIfType == commonDefs.L2RefTypeVlan {
						ret_obj.OutgoingIntfType = "VLAN"
					} else if routesInfo.RouteList[i].NextHopIfType == commonDefs.L2RefTypePort {
						ret_obj.OutgoingIntfType = "PHY"
					} else if routesInfo.RouteList[i].NextHopIfType == commonDefs.IfTypeNull {
						ret_obj.OutgoingIntfType = "NULL"
					} else if routesInfo.RouteList[i].NextHopIfType == commonDefs.IfTypeLoopback {
						ret_obj.OutgoingIntfType = "Lpbk"
					}
					objs = append(objs, ret_obj)
				}
			}
		}
		break
	case models.IPv4RouteState:
		if clnt.ClientHdl != nil {
			var ret_obj models.IPv4RouteState
			routesInfo, _ := clnt.ClientHdl.GetBulkRoutes(ribd.Int(currMarker), ribd.Int(count))
			if routesInfo.Count != 0 {
				objCount = int64(routesInfo.Count)
				more = bool(routesInfo.More)
				nextMarker = int64(routesInfo.EndIdx)
				for i := 0; i < int(routesInfo.Count); i++ {
					if len(objs) == 0 {
						objs = make([]models.ConfigObj, 0)
					}
					ret_obj.DestinationNw = routesInfo.RouteList[i].DestNetIp
					ret_obj.RouteCreatedTime = routesInfo.RouteList[i].RouteCreated
					ret_obj.RouteUpdatedTime = routesInfo.RouteList[i].RouteUpdated
					//ret_obj.PolicyList = make([]string,0)
					  //      routePolicyListInfo := ""
					    //    if routesInfo.RouteList[i].PolicyList != nil {
					      //    for k,v := range routesInfo.RouteList[i].PolicyList {
						    //    routePolicyListInfo = k+":"
					          //  for vv:=0;vv<len(v);vv++ {
					           //   routePolicyListInfo = routePolicyListInfo + v[vv]+","
					            //}
					            //ret_obj.PolicyList = append(ret_obj.PolicyList,routePolicyListInfo)
					         // }
					        //}
					ret_obj.PolicyList = make([]string, 0)
					for j := 0; j < len(routesInfo.RouteList[i].PolicyList); j++ {
						ret_obj.PolicyList = append(ret_obj.PolicyList, routesInfo.RouteList[i].PolicyList[j])
					}
					objs = append(objs, ret_obj)
				}
			}
		}
		break
	case models.IPv4EventState:
		if clnt.ClientHdl != nil {
			var ret_obj models.IPv4EventState
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

	case models.PolicyConditionState:
		logger.Println("PolicyConditionState")
		if clnt.ClientHdl != nil {
			var ret_obj models.PolicyConditionState
			getBulkInfo, _ := clnt.ClientHdl.GetBulkPolicyConditionState(ribdInt.Int(currMarker), ribdInt.Int(count))
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
	case models.PolicyActionState:
		logger.Println("PolicyActionState")
		if clnt.ClientHdl != nil {
			var ret_obj models.PolicyActionState
			getBulkInfo, _ := clnt.ClientHdl.GetBulkPolicyActionState(ribdInt.Int(currMarker), ribdInt.Int(count))
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
			getBulkInfo, _ := clnt.ClientHdl.GetBulkPolicyStmtState(ribdInt.Int(currMarker), ribdInt.Int(count))
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
					ret_obj.PolicyList = make([]string, 0)
					for j := 0; j < len(getBulkInfo.PolicyStmtStateList[i].PolicyList); j++ {
						ret_obj.PolicyList = append(ret_obj.PolicyList, getBulkInfo.PolicyStmtStateList[i].PolicyList[j])
					}
					objs = append(objs, ret_obj)
				}
			}
		}
		break
	case models.PolicyDefinitionState:
		if clnt.ClientHdl != nil {
			var ret_obj models.PolicyDefinitionState
			getBulkInfo, _ := clnt.ClientHdl.GetBulkPolicyDefinitionState(ribdInt.Int(currMarker), ribdInt.Int(count))
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

func (clnt *RIBDIntClient) CreateObject(obj models.ConfigObj, dbHdl *sql.DB) (int64, bool) {
    var err error
	switch obj.(type) {
/*	case models.IPv4Route:
		v4Route := obj.(models.IPv4Route)
		outIntf, _ := strconv.Atoi(v4Route.OutgoingInterface)
		var outIntfType ribd.Int
		if v4Route.OutgoingIntfType == "VLAN" {
			outIntfType = commonDefs.L2RefTypeVlan
		} else if v4Route.OutgoingIntfType == "PHY" {
			outIntfType = commonDefs.L2RefTypePort
		} else if v4Route.OutgoingIntfType == "NULL" {
			outIntfType = commonDefs.IfTypeNull
		}
		//proto, _ := strconv.Atoi(v4Route.Protocol)
		if clnt.ClientHdl != nil {
			_, err = clnt.ClientHdl.CreateV4Route(
				v4Route.DestinationNw, //ribd.Int(binary.BigEndian.Uint32(net.ParseIP(v4Route.DestinationNw).To4())),
				v4Route.NetworkMask,   //ribd.Int(prefixLen),
				ribd.Int(v4Route.Cost),
				v4Route.NextHopIp, //ribd.Int(binary.BigEndian.Uint32(net.ParseIP(v4Route.NextHopIp).To4())),
				outIntfType,
				ribd.Int(outIntf),
				v4Route.Protocol)
			//ribd.Int(proto))
		}
		if err != nil {
			return int64(0), false
		}
		objId, _ := v4Route.StoreObjectInDb(dbHdl)
		return objId, true
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
			_, err = clnt.ClientHdl.CreatePolicyPrefixSet(&cfg)
		}
		if err != nil {
			return int64(0), false
		}
		objId, _ := inCfg.StoreObjectInDb(dbHdl)
		return objId, true*/
/*	case models.PolicyConditionConfig:
		logger.Println("PolicyConditionConfig")
		inCfg := obj.(models.PolicyConditionConfig)
		var cfg ribd.PolicyConditionConfig
		cfg.Name = inCfg.Name
		cfg.ConditionType = inCfg.ConditionType
		switch inCfg.ConditionType {
		case "MatchProtocol":
			logger.Println("MatchProtocol ", inCfg.MatchProtocol)
			cfg.MatchProtocolConditionInfo = inCfg.MatchProtocol
			//dstIpMatchPrefixconditionCfg.Prefix = &cfgIpPrefix
			//cfg.MatchDstIpPrefixConditionInfo = &dstIpMatchPrefixconditionCfg
			break
		case "MatchDstIpPrefix":
			logger.Println("MatchDstIpPrefix")
			inConditionCfg := models.PolicyDstIpMatchPrefixSetCondition{}
			inConditionCfg.Prefix.IpPrefix = inCfg.IpPrefix
			logger.Println("inCfg.MatchDstIpConditionIpPrefix = ", inCfg.IpPrefix)
			logger.Println("inConditionCfg.Prefix.IpPrefix = ", inConditionCfg.Prefix.IpPrefix)
			inConditionCfg.Prefix.MaskLengthRange = inCfg.MaskLengthRange
			var cfgIpPrefix ribd.PolicyPrefix
			var dstIpMatchPrefixconditionCfg ribd.PolicyDstIpMatchPrefixSetCondition
			if len(inConditionCfg.PrefixSet) > 0 && len(inConditionCfg.Prefix.IpPrefix) > 0 {
				logger.Println("cannot set both prefix set name and a prefix")
				return int64(0), true
			}
			dstIpMatchPrefixconditionCfg.PrefixSet = inConditionCfg.PrefixSet
			cfgIpPrefix.IpPrefix = inConditionCfg.Prefix.IpPrefix
			cfgIpPrefix.MasklengthRange = inConditionCfg.Prefix.MaskLengthRange
			logger.Println("cfgIpPrefix.IpPrefix = ", cfgIpPrefix.IpPrefix)
			dstIpMatchPrefixconditionCfg.Prefix = &cfgIpPrefix
			cfg.MatchDstIpPrefixConditionInfo = &dstIpMatchPrefixconditionCfg
			break
		default:
			logger.Println("Invalid condition type")
			return int64(0), true
		}
		if clnt.ClientHdl != nil {
			_, err = clnt.ClientHdl.CreatePolicyCondition(&cfg)
		}
		if err != nil {
			return int64(0), false
		}
		objId, _ := inCfg.StoreObjectInDb(dbHdl)
		return objId, true
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
			cfg.RedistributeAction = inCfg.RedistributeAction
			cfg.RedistributeTargetProtocol = inCfg.RedistributeTargetProtocol
			break
		case "NetworkStatementAdvertise":
			logger.Println("NetworkStatementAdvertise")
			cfg.NetworkStatementTargetProtocol = inCfg.NetworkStatementTargetProtocol
			break
		case "SetAdminDistance":
			logger.Println("SetSdminDistance to inCfg.SetAdminDistanceValue")
			cfg.SetAdminDistanceValue = ribd.Int(inCfg.SetAdminDistanceValue)
			break
		}
		if clnt.ClientHdl != nil {
			_, err = clnt.ClientHdl.CreatePolicyAction(&cfg)
		}
		if err != nil {
			return int64(0), false
		}
		objId, _ := inCfg.StoreObjectInDb(dbHdl)
		return objId, true
	case models.PolicyStmtConfig:
		logger.Println("PolicyStmtConfig")
		var i int
		inCfg := obj.(models.PolicyStmtConfig)
		var cfg ribdInt.PolicyStmtConfig
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
			_, err = clnt.ClientHdl.CreatePolicyStatement(&cfg)
		}
		if err != nil {
			return int64(0), false
		}
		objId, _ := inCfg.StoreObjectInDb(dbHdl)
		return objId, true
	case models.PolicyDefinitionConfig:
		logger.Println("PolicyDefinitionConfig")
		inCfg := obj.(models.PolicyDefinitionConfig)
		var cfg ribdInt.PolicyDefinitionConfig
		cfg.Name = inCfg.Name
		cfg.Precedence = ribdInt.Int(inCfg.Precedence)
		cfg.MatchType = inCfg.MatchType
		logger.Println("Number of statements = ", len(inCfg.StatementList))
		policyDefinitionStatements := make([]ribdInt.PolicyDefinitionStmtPrecedence, len(inCfg.StatementList))
		cfg.PolicyDefinitionStatements = make([]*ribdInt.PolicyDefinitionStmtPrecedence, 0)
		var i int
		for k, v := range inCfg.StatementList {
			logger.Println("k= ", k, " v= ", v)
			/*if v == nil {
				logger.Println("Interface nil at key ", k)
				continue
			}*/
/*inCfgStatementIf := v.(map[string]interface{}) //models.PolicyDefinitionStmtPrecedence)
policyDefinitionStatements[i] = ribd.PolicyDefinitionStmtPrecedence{Precedence: ribd.Int(inCfgStatementIf["Precedence"].(float64)), Statement: inCfgStatementIf["Statement"].(string)}*/
/*policyDefinitionStatements[i] = ribdInt.PolicyDefinitionStmtPrecedence{Precedence: ribdInt.Int(v.Precedence), Statement: v.Statement}
			cfg.PolicyDefinitionStatements = append(cfg.PolicyDefinitionStatements, &policyDefinitionStatements[i])
			i++
		}
		if clnt.ClientHdl != nil {
			_, err = clnt.ClientHdl.CreatePolicyDefinition(&cfg)
		}
		if err != nil {
			return int64(0), false
		}
		objId, _ := inCfg.StoreObjectInDb(dbHdl)
		return objId, true
		break
	default:
		break
	}
	return int64(0), true
}

func (clnt *RIBDIntClient) DeleteObject(obj models.ConfigObj, objKey string, dbHdl *sql.DB) bool {
	logger.Println("### Delete Object is called in RIBClient. ObjectKey: ", objKey, obj)
	switch obj.(type) {
/*	case models.IPv4Route:
		v4Route := obj.(models.IPv4Route)
		logger.Println("### DeleteV4Route is called in RIBClient. ", v4Route.DestinationNw, v4Route.NetworkMask, v4Route.OutgoingInterface)
		if clnt.ClientHdl != nil {
			_, err := clnt.ClientHdl.DeleteV4Route(
				v4Route.DestinationNw, //ribd.Int(binary.BigEndian.Uint32(net.ParseIP(v4Route.DestinationNw).To4())),
				v4Route.NetworkMask,   //ribd.Int(prefixLen),
				v4Route.Protocol,
				v4Route.NextHopIp)
			if err != nil {
				return false
			}
		}
		v4Route.DeleteObjectFromDb(objKey, dbHdl)
		break
	case models.PolicyConditionConfig:
		logger.Println("PolicyConditionConfig")
		inCfg := obj.(models.PolicyConditionConfig)
		var cfg ribd.PolicyConditionConfig
		cfg.Name = inCfg.Name
		if clnt.ClientHdl != nil {
			_, err := clnt.ClientHdl.DeletePolicyCondition(&cfg)
			if err != nil {
				return false
			}
		}
		inCfg.DeleteObjectFromDb(objKey, dbHdl)
		break
	case models.PolicyActionConfig:
		logger.Println("PolicyActionConfig")
		inCfg := obj.(models.PolicyActionConfig)
		var cfg ribd.PolicyActionConfig
		cfg.Name = inCfg.Name
		if clnt.ClientHdl != nil {
			_, err := clnt.ClientHdl.DeletePolicyAction(&cfg)
			if err != nil {
				return false
			}
		}
		inCfg.DeleteObjectFromDb(objKey, dbHdl)
		break
	case models.PolicyStmtConfig:
		logger.Println("PolicyStmtConfig")
		inCfg := obj.(models.PolicyStmtConfig)
		var cfg ribdInt.PolicyStmtConfig
		cfg.Name = inCfg.Name
		if clnt.ClientHdl != nil {
			_, err := clnt.ClientHdl.DeletePolicyStatement(&cfg)
			if err != nil {
				return false
			}
		}
		inCfg.DeleteObjectFromDb(objKey, dbHdl)
		break
	case models.PolicyDefinitionConfig:
		logger.Println("PolicyDefinition")
		inCfg := obj.(models.PolicyDefinitionConfig)
		var cfg ribdInt.PolicyDefinitionConfig
		cfg.Name = inCfg.Name
		if clnt.ClientHdl != nil {
			_, err := clnt.ClientHdl.DeletePolicyDefinition(&cfg)
			if err != nil {
				return false
			}
		}
		inCfg.DeleteObjectFromDb(objKey, dbHdl)
		break

		//default:
		//	logger.Println("OBJECT Type is ", obj.(type))
	}

	return true
}

func (clnt *RIBDIntClient) UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []bool, objKey string, dbHdl *sql.DB) bool {
	logger.Println("### Update Object is called in RIBClient. ", objKey, dbObj, obj, attrSet)
	return true
}
*/
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

func (clnt *ASICDClient) GetObject(obj models.ConfigObj) (models.ConfigObj, bool) {
	return nil, false
}
*/

func (clnt *LACPDClient) GetObject(obj models.ConfigObj) (models.ConfigObj, bool) {
	return nil, false
}

func (clnt *DHCPRELAYDClient) GetObject(obj models.ConfigObj) (models.ConfigObj, bool) {
	return nil, false
}

func (clnt *LocalClient) GetObject(obj models.ConfigObj) (models.ConfigObj, bool) {
	return nil, false
}
