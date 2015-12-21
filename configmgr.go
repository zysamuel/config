package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"models"
	"net/http"
	"path/filepath"
	"time"
)

type ConfigMgr struct {
	clients        map[string]ClientIf
	basePath       string
	fullPath       string
	pRestRtr       *mux.Router
	dbHdl          *sql.DB
	restRoutes     []ApiRoute
	reconncetTimer *time.Ticker
	objHdlMap      map[string]ConfigObjInfo
}

//
//  This method reads the model data and creates rest route interfaces.
//
func (mgr *ConfigMgr) InitializeRestRoutes() bool {
	var rt ApiRoute

	for key, _ := range models.ConfigObjectMap {
		rt = ApiRoute{key + "Show",
			"GET",
			"/" + key,
			ShowConfigObject,
		}
		mgr.restRoutes = append(mgr.restRoutes, rt)
		rt = ApiRoute{key + "Create",
			"POST",
			"/" + key,
			ConfigObjectCreate,
		}
		mgr.restRoutes = append(mgr.restRoutes, rt)
		rt = ApiRoute{key + "Delete",
			"DELETE",
			"/" + key + "/" + "{objId}",
			ConfigObjectDelete,
		}
		mgr.restRoutes = append(mgr.restRoutes, rt)

		rt = ApiRoute{key + "s",
			"GET",
			"/" + key + "s/" + "{objId}",
			ConfigObjectsBulkGet,
		}
		rt = ApiRoute{key + "s",
			"GET",
			"/" + key + "s",
			ConfigObjectsBulkGet,
		}
		mgr.restRoutes = append(mgr.restRoutes, rt)
		rt = ApiRoute{key + "Update",
			"PATCH",
			"/" + key + "/" + "{objId}",
			ConfigObjectUpdate,
		}
		mgr.restRoutes = append(mgr.restRoutes, rt)

	}
	return true
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
	mgr.reconncetTimer = time.NewTicker(time.Millisecond * 1000)
	mgr.InitializeRestRoutes()
	mgr.InstantiateRestRtr()
	mgr.InstantiateDbIf()
	logger.Println("Initialization Done!")
	return mgr
}
