package main
import ("encoding/json"
		  "io/ioutil")

type ConfigMgr struct {
    clients  [] Client
}
//
//  This method connects to all the config daemon's cleints
//
func (mgr ConfigMgr) ConnectToAllClients () bool {
    logger.Println("connect to all client", mgr.clients)
	 for _,clnt := range mgr.clients{		  
        logger.Println("Trying to connect to client", clnt)
    }
    return false
}

//
// This method is to check if config manager is ready to accept config requests
//
func (mgr ConfigMgr) IsReady() bool {
    return false
}

//
// This method is to disconnect from all clients
//
func (mgr ConfigMgr) disconnectFromAllClients () bool {
    return false
}

//
// This function would work as a classical constructor for the 
// configMgr object
//
func NewConfigMgr ( paramsFile string)  *ConfigMgr {
    mgr :=  new (ConfigMgr)
	 var clientsList [] Client
	 //bytes, err := ioutil.ReadFile("./params/clients.json")
	 bytes, err := ioutil.ReadFile(paramsFile)
	 if err != nil {
		  logger.Println("Error in reading configuration file")
		  return nil
	 }

	 err = json.Unmarshal(bytes, &clientsList)
	 if err != nil {
		  logger.Println("Error in Unmarshalling Json")
		  return nil
	 }
    mgr.clients =  clientsList
	 for idx, client := range(mgr.clients) {
		  mgr.clients[idx].Intf = ClientInterfaces[client.Name]
		  mgr.clients[idx].Intf.Initialize(client.Name, "localhost:9090")
		  mgr.clients[idx].Intf.ConnectToServer()
		  logger.Println("Initialization of Client: " , client.Name,client.Intf) 
	 }

    logger.Println("Initialization Done!" , mgr.clients)
    return mgr
}
