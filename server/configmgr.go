package server

import (
	"asicd/asicdCommonDefs"
	"config/apis"
	"config/clients"
	"config/objects"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"models"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"utils/logging"
)

type ConfigMgr struct {
	logger      *logging.Writer
	dbHdl       *objects.DbHandler
	bringUpTime time.Time
	swVersion   SwVersion
	users       []UserData
	sessionId   uint64
	sessionChan chan uint64
	ApiMgr      *apis.ApiMgr
	clientMgr   *clients.ClientMgr
	objectMgr   *objects.ObjectMgr
}

type Repo struct {
	Name   string `json:Name`
	Sha1   string `json:Sha1`
	Branch string `json:Branch`
	Time   string `json:Time`
}

type Version struct {
	Major string `json:major`
	Minor string `json:minor`
	Patch string `json:patch`
	Build string `json:build`
}

type SwVersion struct {
	SwVersion string
	Repos     []Repo
}

var gConfigMgr *ConfigMgr

type ConfdGlobals struct {
	Name  string `json: "Name"`
	Value string `json: "Value"`
}

// Get the http port on which rest api calls will be received
func GetConfigHandlerPort(paramsDir string) (bool, string) {
	var globals []ConfdGlobals
	var port string

	globalsFile := paramsDir + "/globals.json"
	bytes, err := ioutil.ReadFile(globalsFile)
	if err != nil {
		gConfigMgr.logger.Err(fmt.Sprintln("Error in reading globals file", globalsFile))
		return false, port
	}

	err = json.Unmarshal(bytes, &globals)
	if err != nil {
		gConfigMgr.logger.Err("Failed to Unmarshall Json")
		return false, port
	}
	for _, global := range globals {
		if global.Name == "httpport" {
			port = global.Value
			return true, port
		}
	}
	return false, port
}

//
// This function would work as a classical constructor for the
// configMgr object
//
func NewConfigMgr(paramsDir string, logger *logging.Writer) *ConfigMgr {
	mgr := new(ConfigMgr)
	mgr.logger = logger

	paramsFile := paramsDir + "/clients.json"
	mgr.clientMgr = clients.InitializeClientMgr(paramsFile, logger, GetSystemStatus, GetSystemSwVersion)

	objects.CreateObjectMap()
	objectConfigFiles := [...]string{paramsDir + "/objectconfig.json",
		paramsDir + "/genObjectConfig.json"}
	mgr.objectMgr = objects.InitializeObjectMgr(objectConfigFiles[:], logger, mgr.clientMgr)
	mgr.dbHdl = objects.InstantiateDbIf(logger)

	mgr.ApiMgr = apis.InitializeApiMgr(paramsDir, logger, mgr.dbHdl, mgr.objectMgr)
	mgr.ApiMgr.InitializeRestRoutes()
	mgr.ApiMgr.InstantiateRestRtr()

	mgr.bringUpTime = time.Now()
	logger.Info("Initialization Done!")

	go mgr.ReadSystemSwVersion(paramsDir)
	go mgr.clientMgr.ConnectToAllClients()
	go mgr.DiscoverSystemObjects()
	go mgr.SigHandler()

	// These user management routines are not used right now.
	//go mgr.CreateDefaultUser()
	//go mgr.ReadConfiguredUsersFromDb()
	//go mgr.StartUserSessionHandler()

	gConfigMgr = mgr

	return mgr
}

func (mgr *ConfigMgr) SigHandler() {
	sigChan := make(chan os.Signal, 1)
	signalList := []os.Signal{syscall.SIGHUP}
	signal.Notify(sigChan, signalList...)

	for {
		select {
		case signal := <-sigChan:
			switch signal {
			case syscall.SIGHUP:
				mgr.logger.Info("Exting!!!")
				os.Exit(0)
			default:
			}
		}
	}
}

func GetSystemStatus() models.SystemStatusState {
	systemStatus := models.SystemStatusState{}
	systemStatus.Name, _ = os.Hostname()
	systemStatus.Ready = gConfigMgr.clientMgr.IsReady()
	if systemStatus.Ready == false {
		reason := "Not connected to"
		unconnectedClients := gConfigMgr.clientMgr.GetUnconnectedClients()
		for idx := 0; idx < len(unconnectedClients); idx++ {
			reason = reason + " " + unconnectedClients[idx]
		}
		systemStatus.Reason = reason
	} else {
		systemStatus.Reason = "None"
	}
	systemStatus.UpTime = time.Since(gConfigMgr.bringUpTime).String()
	systemStatus.NumCreateCalls =
		fmt.Sprintf("Total %d Success %d", gConfigMgr.ApiMgr.ApiCallStats.NumCreateCalls, gConfigMgr.ApiMgr.ApiCallStats.NumCreateCallsSuccess)
	systemStatus.NumDeleteCalls =
		fmt.Sprintf("Total %d Success %d", gConfigMgr.ApiMgr.ApiCallStats.NumDeleteCalls, gConfigMgr.ApiMgr.ApiCallStats.NumDeleteCallsSuccess)
	systemStatus.NumUpdateCalls =
		fmt.Sprintf("Total %d Success %d", gConfigMgr.ApiMgr.ApiCallStats.NumUpdateCalls, gConfigMgr.ApiMgr.ApiCallStats.NumUpdateCallsSuccess)
	systemStatus.NumGetCalls =
		fmt.Sprintf("Total %d Success %d", gConfigMgr.ApiMgr.ApiCallStats.NumGetCalls, gConfigMgr.ApiMgr.ApiCallStats.NumGetCallsSuccess)
	systemStatus.NumActionCalls =
		fmt.Sprintf("Total %d Success %d", gConfigMgr.ApiMgr.ApiCallStats.NumActionCalls, gConfigMgr.ApiMgr.ApiCallStats.NumActionCallsSuccess)

	// Read DaemonStates from db
	var daemonState models.DaemonState
	daemonStates, _ := daemonState.GetAllObjFromDb(gConfigMgr.dbHdl)
	systemStatus.FlexDaemons = make([]models.DaemonState, len(daemonStates))
	for idx, daemonState := range daemonStates {
		systemStatus.FlexDaemons[idx] = daemonState.(models.DaemonState)
	}
	return systemStatus
}

func GetSystemSwVersion() models.SystemSwVersionState {
	systemSwVersion := models.SystemSwVersionState{}
	systemSwVersion.FlexswitchVersion = gConfigMgr.swVersion.SwVersion
	numRepos := len(gConfigMgr.swVersion.Repos)
	systemSwVersion.Repos = make([]models.RepoInfo, numRepos)
	for i := 0; i < numRepos; i++ {
		systemSwVersion.Repos[i].Name = gConfigMgr.swVersion.Repos[i].Name
		systemSwVersion.Repos[i].Sha1 = gConfigMgr.swVersion.Repos[i].Sha1
		systemSwVersion.Repos[i].Branch = gConfigMgr.swVersion.Repos[i].Branch
		systemSwVersion.Repos[i].Time = gConfigMgr.swVersion.Repos[i].Time
	}
	return systemSwVersion
}

func (mgr *ConfigMgr) DiscoverPorts() error {
	mgr.logger.Debug("Discovering ports")
	asicdConnectionCheckTimer := time.NewTicker(time.Millisecond * 1000)
	i := 0
	for t := range asicdConnectionCheckTimer.C {
		_ = t
		if mgr.clientMgr.IsConnectedClient("asicd") {
			asicdConnectionCheckTimer.Stop()
			break
		} else {
			if i%100 == 0 {
				mgr.logger.Info("Not connected to asicd yet to get all ports")
			}
		}
		i++
	}
	// Get ports present on this system and store in DB for user to update port parameters
	resource := "Port"
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		var objs []models.ConfigObj
		var err error
		_, obj, _ := objects.GetConfigObj(nil, objHdl)
		currentIndex := int64(asicdCommonDefs.MIN_SYS_PORTS)
		objCount := int64(asicdCommonDefs.MAX_SYS_PORTS)
		err, _, _, _, objs = mgr.objectMgr.ObjHdlMap[resource].Owner.GetBulkObject(obj, mgr.dbHdl.DBUtil, currentIndex, objCount)
		if err == nil {
			for i := 0; i < len(objs); i++ {
				portConfig := (*objs[i].(*models.Port))
				_, err := portConfig.GetObjectFromDb(portConfig.GetKey(), mgr.dbHdl)
				// if we can not find the port in DB then go ahead and store
				if err != nil {
					err = portConfig.StoreObjectInDb(mgr.dbHdl)
					if err != nil {
						mgr.logger.Err(fmt.Sprintln("Failed to store Port in DB ", i, portConfig, err))
					} else {
						_, err := mgr.dbHdl.StoreUUIDToObjKeyMap(portConfig.GetKey())
						if err != nil {
							mgr.logger.Err(fmt.Sprintln("Failed to store uuid map for Port in DB ", portConfig, err))
						}
					}
				}
			}
		}
	}
	mgr.logger.Debug("Ports discovered")
	return nil
}

//
// This method is to get system objects and store in DB
//
func (mgr *ConfigMgr) DiscoverSystemObjects() error {
	mgr.logger.Info("Discover system objects")
	mgr.DiscoverPorts()
	return nil
}

func (mgr *ConfigMgr) ReadSystemSwVersion(paramsDir string) error {
	var version Version
	infoDir := strings.TrimSuffix(paramsDir, "params/")
	pkgInfoFile := infoDir + "pkgInfo.json"
	bytes, err := ioutil.ReadFile(pkgInfoFile)
	if err != nil {
		mgr.logger.Err(fmt.Sprintln("Error in reading configuration file", pkgInfoFile))
		return err
	}

	err = json.Unmarshal(bytes, &version)
	if err != nil {
		mgr.logger.Err("Error in Unmarshalling pkgInfo Json")
		return err
	}
	mgr.swVersion.SwVersion = version.Major + "." + version.Minor + "." + version.Patch + "." + version.Build

	buildInfoFile := infoDir + "buildInfo.json"
	bytes, err = ioutil.ReadFile(buildInfoFile)
	if err != nil {
		mgr.logger.Err(fmt.Sprintln("Error in reading configuration file", buildInfoFile))
		return err
	}

	err = json.Unmarshal(bytes, &mgr.swVersion.Repos)
	if err != nil {
		mgr.logger.Err("Error in Unmarshalling buildInfo Json")
		return err
	}
	return nil
}
