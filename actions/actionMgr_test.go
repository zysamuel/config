package actions

import (
	"config/objects"
	"flag"
	"fmt"
	modelActions "models/actions"
	"net/http"
	"testing"
	"utils/logging"
)

//usage: go test or
//       go test -params=<paramsDir> e,g: go test -params=/opt/flexswitch/params/

var paramsDir = flag.String("params", "../../models/actions", "Location of actionJsonFile")
var Logger *logging.Writer
var infoListFile []string
var actionMgr *ActionMgr

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
	//	systemSwVersion := modelObjs.SystemSwVersionState{}
	//	clientMgr := clients.InitializeClientMgr(".././params/clients.json", Logger, getSystemStatus(),  systemSwVersion)
	fmt.Println("call initializeactionMgr")
	dbHdl := objects.InstantiateDbIf(Logger)
	v := InitializeActionMgr(paramsDir, infoListFile, Logger, dbHdl, nil, nil)
	fmt.Println("returned v:", v)
	t.Log("For ", infoListFile, " nil, nil",
		"got", v)
	actionMgr = v
}

func TestGetAllActions(t *testing.T) {
	actions := actionMgr.GetAllActions()
	fmt.Println("actions:", actions)
}

func TestGetApplyActionObj(t *testing.T) {
	var r *http.Request
	if actionobjHdl, ok := modelActions.ActionMap["ApplyConfig"]; ok {
		fmt.Println("actionObjhdl:", actionobjHdl)
		if body, actionobj, err := GetActionObj(r, actionobjHdl); err == nil {
			fmt.Println("body:", body, " actionobj:", actionobj)
			err := ExecutePerformAction(actionobj)
			fmt.Println("err:", err)
		}
	}
}
func TestGetSaveActionObj(t *testing.T) {
	var r *http.Request
	if actionobjHdl, ok := modelActions.ActionMap["SaveConfig"]; ok {
		fmt.Println("actionObjhdl:", actionobjHdl)
		if body, actionobj, err := GetActionObj(r, actionobjHdl); err == nil {
			fmt.Println("body:", body, " actionobj:", actionobj)
			err := ExecutePerformAction(actionobj)
			fmt.Println("err:", err)
		}
	}
}
