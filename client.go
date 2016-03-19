package main

import (
	"asicd/asicdConstDefs"
	"encoding/json"
	"io/ioutil"
	"models"
	"strconv"
	"time"
)

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
//  This method connects to all the config daemon's clients
//
func (mgr *ConfigMgr) ConnectToAllClients(clientsUp chan bool) bool {
	unconnectedClients := make([]string, 0)
	mgr.reconncetTimer = time.NewTicker(time.Millisecond * 1000)
	mgr.systemReady = false
	idx := 0
	for clntName, client := range mgr.clients {
		client.ConnectToServer()
		if client.IsConnectedToServer() == false {
			unconnectedClients = append(unconnectedClients, clntName)
			//unconnectedClients[idx] = clntName
			idx++
		}
	}
	waitCount := 0
	if idx > 0 {
		for t := range mgr.reconncetTimer.C {
			_ = t
			if waitCount == 0 {
				logger.Println("Looking for clients ", unconnectedClients)
			}
			for i := 0; i < len(unconnectedClients); i++ {
				if waitCount%100 == 0 {
					logger.Println("Waiting to connect to these clients", unconnectedClients[i])
				}
				if len(unconnectedClients) > i {
					if mgr.clients[unconnectedClients[i]].IsConnectedToServer() {
						unconnectedClients = append(unconnectedClients[:i], unconnectedClients[i+1:]...)
					} else {
						mgr.clients[unconnectedClients[i]].ConnectToServer()
					}
				}
			}
			if len(unconnectedClients) == 0 {
				mgr.reconncetTimer.Stop()
				break
			}
			waitCount++
		}
	}
	logger.Println("Connected to all clients")
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
func (mgr *ConfigMgr) DiscoverSystemObjects(clientsUp chan bool) bool {
	// Wait till confd connects to all the servers, i.e. all the clients are up
	<-clientsUp
	logger.Println("Discovering system objects")

	logger.Println("Discover ports")
	resource := "PortConfig"
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		var objects []models.ConfigObj
		var err error
		_, obj, _ := GetConfigObj(nil, objHdl)
		currentIndex := int64(asicdConstDefs.MIN_SYS_PORTS)
		objCount := int64(asicdConstDefs.MAX_SYS_PORTS)
		err, _, _, _, objects = gMgr.objHdlMap[resource].owner.GetBulkObject(obj, currentIndex, objCount)
		if err == nil {
			for i := 0; i < len(objects); i++ {
				portConfig := objects[i].(models.Port)
				objKey, _ := portConfig.GetKey()
				_, err := portConfig.GetObjectFromDb(objKey, mgr.dbHdl)
				// if we can not find the port in DB then go ahead and store
				if err != nil {
					_, err = portConfig.StoreObjectInDb(mgr.dbHdl)
					if err != nil {
						logger.Println("Failed to store PortConfig in DB ", i, portConfig, err)
					} else {
						uuid, err := StoreUuidToKeyMapInDb(portConfig)
						if err != nil {
							logger.Println("Failed to store uuid map for PortConfig in DB ", portConfig, err)
						} else {
							logger.Println("Stored uuid map for PortConfig in DB ", portConfig, uuid)
						}
					}
				}
			}
		}
	}
	return true
}

func (mgr *ConfigMgr) MonitorSystemStatus() bool {
	return true
}
