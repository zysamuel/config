package main

import (
	"asicd/asicdConstDefs"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"models"
	"strconv"
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
	UpdateObject(dbObj models.ConfigObj, obj models.ConfigObj, attrSet []byte, objKey string, dbHdl *sql.DB) bool
}

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

//
//  This method reads the config file and connects to all the clients in the list
//
func (mgr *ConfigMgr) InitializeClientHandles(paramsFile string) bool {
	var clientsList []ClientJson
	mgr.clients = make(map[string]ClientIf)

	bytes, err := ioutil.ReadFile(paramsFile)
	if err != nil {
		logger.Println("Error in reading configuration file", paramsFile)
		return false
	}

	err = json.Unmarshal(bytes, &clientsList)
	if err != nil {
		logger.Println("Error in Unmarshalling Json")
		return false
	}
	for _, client := range clientsList {
		if ClientInterfaces[client.Name] != nil {
			mgr.clients[client.Name] = ClientInterfaces[client.Name]
			mgr.clients[client.Name].Initialize(client.Name, "localhost:"+strconv.Itoa(client.Port))
		}
	}

	return true
}

//
//  This method connects to all the config daemon's cleints
//
func (mgr *ConfigMgr) ConnectToAllClients(clientsUp chan bool) bool {
	unconnectedClients := make([]string, len(mgr.clients))
	idx := 0
	for clntName, client := range mgr.clients {
		unconnectedClients[idx] = clntName
		idx++
		client.ConnectToServer()
		client.IsConnectedToServer()
	}
	waitCount := 0
	for t := range mgr.reconncetTimer.C {
		_ = t
		if waitCount == 0 {
			logger.Println("Looking for clients ", unconnectedClients)
		}
		for i := 0; i < len(unconnectedClients); i++ {
			if waitCount%100 == 0 {
				logger.Println("Waiting to connect to these clients", unconnectedClients[i])
			}
			if mgr.clients[unconnectedClients[i]].IsConnectedToServer() {
				unconnectedClients = append(unconnectedClients[:i], unconnectedClients[i+1:]...)
			} else {
				mgr.clients[unconnectedClients[i]].ConnectToServer()
			}
		}
		if len(unconnectedClients) == 0 {
			mgr.reconncetTimer.Stop()
			break
		}
		waitCount++
	}
	mgr.systemReady = true
	clientsUp <- true
	return true
}

//
// This method is to check if config manager is ready to accept config requests
//
func (mgr *ConfigMgr) IsReady() bool {
	return mgr.systemReady
}

//
// This method is to disconnect from all clients
//
func (mgr *ConfigMgr) disconnectFromAllClients() bool {
	return false
}

//
// This method is to get Port interfaces from Asicd and store in DB for config update on those ports
//
func (mgr *ConfigMgr) StartPortInterfaceThread(clientsUp chan bool) bool {
	logger.Println("Waiting for PortConfig server")
	serverUp := <-clientsUp
	logger.Println("PortConfig server is up? ", serverUp)
	resource := "PortIntfConfig"
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		var resp GetBulkResponse
		var err error
		obj, _ := GetConfigObj(nil, objHdl)
		//currentIndex, objCount := ExtractGetBulkParams(r)
		//resp.CurrentMarker = currentIndex
		currentIndex := int64(asicdConstDefs.MIN_SYS_PORTS)
		objCount := int64(asicdConstDefs.MAX_SYS_PORTS)
		err, resp.ObjCount, resp.NextMarker, resp.MoreExist,
			resp.StateObjects = gMgr.objHdlMap[resource].owner.GetBulkObject(obj, currentIndex, objCount)
		logger.Println("PortIntf thread: ", err, resp.ObjCount, resp.NextMarker, resp.MoreExist)
		if err == nil {
			for i := int64(0); i < resp.ObjCount; i++ {
				portConfig := resp.StateObjects[i].(models.PortIntfConfig)
				logger.Println("PortIntf: ", i, portConfig)
				_, err = portConfig.StoreObjectInDb(mgr.dbHdl)
				if err != nil {
					logger.Println("Failed to store PortIntfConfig in DB ", i, portConfig, err)
				}
			}
		}
	}
	return true
}
