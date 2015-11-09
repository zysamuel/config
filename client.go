package main

import (
	"encoding/json"
	"io/ioutil"
	"models"
)

type ClientIf interface {
	Initialize(name string, address string)
	ConnectToServer() bool
	IsConnectedToServer() bool
	CreateObject(models.ConfigObj) bool
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
		logger.Println("Error in reading configuration file")
		return false
	}

	err = json.Unmarshal(bytes, &clientsList)
	if err != nil {
		logger.Println("Error in Unmarshalling Json")
		return false
	}

	for _, client := range clientsList {
		logger.Println("#### Client name is ", client.Name)
		mgr.clients[client.Name] = ClientInterfaces[client.Name]
		mgr.clients[client.Name].Initialize(client.Name, "localhost:9090")
		mgr.clients[client.Name].ConnectToServer()
		logger.Println("Initialization of Client: ", client.Name)
	}
	return true
}

//
//  This method connects to all the config daemon's cleints
//
func (mgr *ConfigMgr) ConnectToAllClients() bool {
	logger.Println("connect to all client", mgr.clients)
	for _, clnt := range mgr.clients {
		logger.Println("Trying to connect to client", clnt)
	}
	return false
}

//
// This method is to check if config manager is ready to accept config requests
//
func (mgr *ConfigMgr) IsReady() bool {
	return false
}

//
// This method is to disconnect from all clients
//
func (mgr *ConfigMgr) disconnectFromAllClients() bool {
	return false
}
