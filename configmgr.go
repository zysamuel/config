package main
import ("encoding/json"
		  "github.com/gorilla/mux"
		  "net/http"
		  "models"
		  "io/ioutil")

type ConfigMgr struct {
    clients  [] Client
    pRestRtr *mux.Router 
	 restRoutes [] ApiRoute
}

//
//  This method reads the model data and creates rest route interfaces. 
//
func (mgr *ConfigMgr) InitializeRestRoutes() bool{
	 var rt ApiRoute
	 for key, _:= range models.ConfigObjectMap {
		  rt = ApiRoute {key+"Show",
					          "GET",
								 "/" + key,
								 ShowConfigObject,
					 			} 
		  mgr.restRoutes = append (mgr.restRoutes, rt)
		  rt = ApiRoute {key+"Create",
					          "POST",
								 "/"+key,
								 ConfigObjectCreate,
					 			} 
		  mgr.restRoutes = append (mgr.restRoutes, rt)

	 }
	 return true 
}

//
//  This method creates new rest router interface
//
func (mgr *ConfigMgr) InstantiateRestRtr() *mux.Router {
	mgr.pRestRtr = mux.NewRouter().StrictSlash(true)
	for _, route := range mgr.restRoutes{
		var handler http.Handler
		handler = Logger(route.HandlerFunc, route.Name)
		mgr.pRestRtr.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}
	return mgr.pRestRtr 
}

func (mgr *ConfigMgr) GetRestRtr() *mux.Router {
	 return mgr.pRestRtr
}
//
//  This method connects to all the config daemon's cleints
//
func (mgr *ConfigMgr) ConnectToAllClients () bool {
    logger.Println("connect to all client", mgr.clients)
	 for _,clnt := range mgr.clients{		  
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
func (mgr *ConfigMgr) disconnectFromAllClients () bool {
    return false
}

//
// This function would work as a classical constructor for the 
// configMgr object
//
func NewConfigMgr ( paramsFile string)  *ConfigMgr {
    mgr :=  new (ConfigMgr)
	 var clientsList [] Client
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
	 mgr.InitializeRestRoutes()
	 mgr.InstantiateRestRtr()
    logger.Println("Initialization Done!" , mgr.clients)
    return mgr
}
