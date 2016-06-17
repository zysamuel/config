package actions

import (
	"flag"
	"fmt"
	"testing"
	"utils/logging"
	"config/clients"
	modelObjs "models/objects"
)

//usage: go test or
//       go test -params=<paramsDir> e,g: go test -params=/opt/flexswitch/params/

var paramsDir = flag.String("params", "../../models/actions", "Location of actionJsonFile")
var Logger *logging.Writer
var infoListFile []string
var actionMgr *ActionMgr


func getSystemStatus() modelObjs.SystemStatusState {
	systemStatus := modelObjs.SystemStatusState{}
	return systemStatus
}

func Init() {
	fmt.Println("ActionMgr: Start logger")
	logger, err := logging.NewLogger("actionMgr", "ActionMgr", true)
	if err != nil {
		fmt.Println("Failed to start logger. Nothing will be logged ...")
	}
	Logger = logger
	paramsDirName := *paramsDir
	fmt.Println("paramsDirName:", paramsDirName)
	actionConfigFiles := [...]string{paramsDirName + "/genActionConfig.json"}
	infoListFile = make([]string, 0)
	infoListFile = actionConfigFiles[:]
}

func TestInit(t *testing.T) {
	Init()
	systemSwVersion := modelObjs.SystemSwVersionState{}
	clientMgr := clients.InitializeClientMgr(".././params/clients.json", Logger, getSystemStatus(),  systemSwVersion)
	fmt.Println("call initializeactionMgr")
	v := InitializeActionMgr(infoListFile, Logger, clientMgr)
	fmt.Println("returned v:", v)
	t.Log("For ", infoListFile, " nil, nil",
		"got", v)
	actionMgr = v
}

func TestGetAllActions(t *testing.T) {
	actions := actionMgr.GetAllActions()
	fmt.Println("actions:", actions)
}
