package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"models"
	"net/http"
	"path/filepath"
	"time"
	//"strings"
	//"encoding/base64"
	//"fmt"
)

type ConfigMgr struct {
	clients            map[string]ClientIf
	apiVer             string
	apiBase            string
	basePath           string
	fullPath           string
	pRestRtr           *mux.Router
	dbHdl              *sql.DB
	restRoutes         []ApiRoute
	reconncetTimer     *time.Ticker
	objHdlMap          map[string]ConfigObjInfo
	systemReady        bool
	users              []UserData
	sessionId          uint32
	sessionChan        chan uint32
}

//
//  This method reads the model data and creates rest route interfaces.
//
func (mgr *ConfigMgr) InitializeRestRoutes() bool {
	var rt ApiRoute

	for key, _ := range models.ConfigObjectMap {
		rt = ApiRoute{key + "Show",
			"GET",
			mgr.apiBase + key,
			HandleRestRouteShowConfig,
		}
		mgr.restRoutes = append(mgr.restRoutes, rt)
		rt = ApiRoute{key + "Create",
			"POST",
			mgr.apiBase + key,
			HandleRestRouteCreate,
		}
		mgr.restRoutes = append(mgr.restRoutes, rt)
		rt = ApiRoute{key + "Delete",
			"DELETE",
			mgr.apiBase + key + "/" + "{objId}",
			HandleRestRouteDelete,
		}
		mgr.restRoutes = append(mgr.restRoutes, rt)

		rt = ApiRoute{key + "s",
			"GET",
			mgr.apiBase + key + "s/" + "{objId}",
			HandleRestRouteGet,
		}
		rt = ApiRoute{key + "s",
			"GET",
			mgr.apiBase + key + "s",
			HandleRestRouteGet,
		}
		mgr.restRoutes = append(mgr.restRoutes, rt)
		rt = ApiRoute{key + "Update",
			"PATCH",
			mgr.apiBase + key + "/" + "{objId}",
			HandleRestRouteUpdate,
		}
		mgr.restRoutes = append(mgr.restRoutes, rt)

	}
	return true
}

func HandleRestRouteShowConfig(w http.ResponseWriter, r *http.Request) {
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	ShowConfigObject(w, r)
}

func HandleRestRouteCreate(w http.ResponseWriter, r *http.Request) {
/*
	resource := strings.TrimPrefix(r.URL.String(), gMgr.apiBase)
	fmt.Println("Create: ", *r)
	fmt.Println("Resource: ", resource)
	fmt.Println("URL: ", r.URL.String())
	auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)
	fmt.Printf("UserName: %s Password: %s\n", pair[0], pair[1])
	return
*/
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	ConfigObjectCreate(w, r)
}

func HandleRestRouteDelete(w http.ResponseWriter, r *http.Request) {
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	ConfigObjectDelete(w, r)
}

func HandleRestRouteUpdate(w http.ResponseWriter, r *http.Request) {
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	ConfigObjectUpdate(w, r)
}

func HandleRestRouteGet(w http.ResponseWriter, r *http.Request) {
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	ConfigObjectsBulkGet(w, r)
}

//
//  This method creates new rest router interface
//
func (mgr *ConfigMgr) InstantiateRestRtr() *mux.Router {
	mgr.pRestRtr = mux.NewRouter().StrictSlash(true)
	mgr.pRestRtr.PathPrefix("/api-docs/").Handler(http.StripPrefix("/api-docs/",
		http.FileServer(http.Dir(mgr.fullPath+"/docsui"))))

	for _, route := range mgr.restRoutes {
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
// This function would work as a classical constructor for the
// configMgr object
//
func NewConfigMgr(paramsDir string) *ConfigMgr {
	var rc bool
	mgr := new(ConfigMgr)
	var err error
	mgr.apiVer = "v1"
	mgr.apiBase = "/public/" + mgr.apiVer + "/"
	if mgr.fullPath, err = filepath.Abs(paramsDir); err != nil {
		logger.Printf("ERROR: Unable to get absolute path for %s, error [%s]\n", paramsDir, err)
		return nil
	}
	mgr.basePath, _ = filepath.Split(mgr.fullPath)

	objectConfigFile := paramsDir + "/objectconfig.json"
	paramsFile := paramsDir + "/clients.json"

	rc = mgr.InitializeClientHandles(paramsFile)
	if rc == false {
		logger.Println("ERROR: Error in Initializing Client handles")
		return nil
	}
	rc = mgr.InitializeObjectHandles(objectConfigFile)
	if rc == false {
		logger.Println("ERROR: Error in Initializing Object handles")
		return nil
	}
	mgr.InitializeRestRoutes()
	mgr.InstantiateRestRtr()
	mgr.InstantiateDbIf()
	logger.Println("Initialization Done!")
	return mgr
}
