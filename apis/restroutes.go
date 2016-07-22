//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package apis

import (
	"config/actions"
	"config/objects"
	"fmt"
	"github.com/gorilla/mux"
	modelObjs "models/objects"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"utils/logging"
)

type ApiRoute struct {
	Name        string           // Unique Identifier to identify this route
	Method      string           // REST Method POST/GET/PATCH....
	Pattern     string           // Endpoint URI
	HandlerFunc http.HandlerFunc // Function reposnsible for executing the request
}

type ApiRoutes []ApiRoute

type ApiMgr struct {
	logger        *logging.Writer
	objectMgr     *objects.ObjectMgr
	actionMgr     *actions.ActionMgr
	dbHdl         *objects.DbHandler
	apiVer        string
	apiBase       string
	apiBaseConfig string
	apiBaseState  string
	apiBaseAction string
	apiBaseEvent  string
	basePath      string
	fullPath      string
	pRestRtr      *mux.Router
	restRoutes    []ApiRoute
	ApiCallStats  ApiCallStats
}

var gApiMgr *ApiMgr

type ApiCallStats struct {
	NumCreateCalls        int32
	NumCreateCallsSuccess int32
	NumDeleteCalls        int32
	NumDeleteCallsSuccess int32
	NumUpdateCalls        int32
	NumUpdateCallsSuccess int32
	NumGetCalls           int32
	NumGetCallsSuccess    int32
	NumActionCalls        int32
	NumActionCallsSuccess int32
}

type LoginResponse struct {
	SessionId uint64 `json: "SessionId"`
}

func InitializeApiMgr(paramsDir string, logger *logging.Writer, dbHdl *objects.DbHandler, objectMgr *objects.ObjectMgr, actionMgr *actions.ActionMgr) *ApiMgr {
	var err error
	mgr := new(ApiMgr)
	mgr.logger = logger
	mgr.dbHdl = dbHdl
	mgr.objectMgr = objectMgr
	mgr.actionMgr = actionMgr
	mgr.apiVer = "v1"
	mgr.apiBase = "/public/" + mgr.apiVer + "/"
	mgr.apiBaseConfig = mgr.apiBase + "config" + "/"
	mgr.apiBaseState = mgr.apiBase + "state" + "/"
	mgr.apiBaseAction = mgr.apiBase + "action" + "/"
	mgr.apiBaseEvent = mgr.apiBase + "events"
	if mgr.fullPath, err = filepath.Abs(paramsDir); err != nil {
		logger.Err(fmt.Sprintln("Unable to get absolute path for %s, error [%s]\n", paramsDir, err))
		return nil
	}
	mgr.basePath, _ = filepath.Split(mgr.fullPath)
	gApiMgr = mgr
	return mgr
}

func (mgr *ApiMgr) InitializeActionRestRoutes() bool {
	var rt ApiRoute
	actionList := mgr.actionMgr.GetAllActions()
	for _, action := range actionList {
		rt = ApiRoute{action + "Action",
			"POST",
			mgr.apiBaseAction + action,
			HandleRestRouteAction,
		}
		mgr.restRoutes = append(mgr.restRoutes, rt)

	}
	return true
}

func (mgr *ApiMgr) InitializeEventRestRoutes() bool {
	var rt ApiRoute
	rt = ApiRoute{"Events",
		"GET",
		mgr.apiBaseEvent,
		HandleRestRouteEvent,
	}
	mgr.restRoutes = append(mgr.restRoutes, rt)

	return true
}

//
//  This method reads the model data and creates rest route interfaces.
//
func (mgr *ApiMgr) InitializeRestRoutes() bool {
	var rt ApiRoute
	for key, _ := range modelObjs.ConfigObjectMap {
		objInfo := mgr.objectMgr.ObjHdlMap[key]
		if objInfo.Access == "w" || objInfo.Access == "rw" {
			rt = ApiRoute{key + "Create",
				"POST",
				mgr.apiBaseConfig + key,
				HandleRestRouteCreate,
			}
			mgr.restRoutes = append(mgr.restRoutes, rt)
			rt = ApiRoute{key + "Delete",
				"DELETE",
				mgr.apiBaseConfig + key + "/" + "{objId}",
				HandleRestRouteDeleteForId,
			}
			mgr.restRoutes = append(mgr.restRoutes, rt)
			rt = ApiRoute{key + "Delete",
				"DELETE",
				mgr.apiBaseConfig + key,
				HandleRestRouteDelete,
			}
			mgr.restRoutes = append(mgr.restRoutes, rt)
			rt = ApiRoute{key + "Update",
				"PATCH",
				mgr.apiBaseConfig + key + "/" + "{objId}",
				HandleRestRouteUpdateForId,
			}
			mgr.restRoutes = append(mgr.restRoutes, rt)
			rt = ApiRoute{key + "Update",
				"PATCH",
				mgr.apiBaseConfig + key,
				HandleRestRouteUpdate,
			}
			mgr.restRoutes = append(mgr.restRoutes, rt)
			rt = ApiRoute{key + "Get",
				"GET",
				mgr.apiBaseConfig + key + "/" + "{objId}",
				HandleRestRouteGetConfigForId,
			}
			mgr.restRoutes = append(mgr.restRoutes, rt)
			rt = ApiRoute{key + "Get",
				"GET",
				mgr.apiBaseConfig + key,
				HandleRestRouteGetConfig,
			}
			mgr.restRoutes = append(mgr.restRoutes, rt)
			rt = ApiRoute{key + "s",
				"GET",
				mgr.apiBaseConfig + key + "s",
				HandleRestRouteBulkGetConfig,
			}
			mgr.restRoutes = append(mgr.restRoutes, rt)
		} else if objInfo.Access == "r" {
			key = strings.TrimSuffix(key, "State")
			rt = ApiRoute{key + "Show",
				"GET",
				mgr.apiBaseState + key + "/" + "{objId}",
				HandleRestRouteGetStateForId,
			}
			mgr.restRoutes = append(mgr.restRoutes, rt)
			rt = ApiRoute{key + "Show",
				"GET",
				mgr.apiBaseState + key,
				HandleRestRouteGetState,
			}
			mgr.restRoutes = append(mgr.restRoutes, rt)
			rt = ApiRoute{key + "s",
				"GET",
				mgr.apiBaseState + key + "s",
				HandleRestRouteBulkGetState,
			}
			mgr.restRoutes = append(mgr.restRoutes, rt)
		} else if objInfo.Access == "x" {
		}
	}
	return true
}

func HandleRestRouteCreate(w http.ResponseWriter, r *http.Request) {
	ConfigObjectCreate(w, r)
	return
}

func HandleRestRouteDeleteForId(w http.ResponseWriter, r *http.Request) {
	ConfigObjectDeleteForId(w, r)
	return
}

func HandleRestRouteDelete(w http.ResponseWriter, r *http.Request) {
	ConfigObjectDelete(w, r)
	return
}

func HandleRestRouteUpdateForId(w http.ResponseWriter, r *http.Request) {
	ConfigObjectUpdateForId(w, r)
	return
}

func HandleRestRouteUpdate(w http.ResponseWriter, r *http.Request) {
	ConfigObjectUpdate(w, r)
	return
}

func HandleRestRouteGetConfigForId(w http.ResponseWriter, r *http.Request) {
	GetOneConfigObjectForId(w, r)
}

func HandleRestRouteGetConfig(w http.ResponseWriter, r *http.Request) {
	GetOneConfigObject(w, r)
}

func HandleRestRouteGetStateForId(w http.ResponseWriter, r *http.Request) {
	GetOneStateObjectForId(w, r)
}

func HandleRestRouteGetState(w http.ResponseWriter, r *http.Request) {
	GetOneStateObject(w, r)
}

func HandleRestRouteBulkGetConfig(w http.ResponseWriter, r *http.Request) {
	BulkGetConfigObjects(w, r)
	return
}

func HandleRestRouteBulkGetState(w http.ResponseWriter, r *http.Request) {
	BulkGetStateObjects(w, r)
	return
}

func HandleRestRouteAction(w http.ResponseWriter, r *http.Request) {
	ExecuteActionObject(w, r)
	return
}

func HandleRestRouteEvent(w http.ResponseWriter, r *http.Request) {
	ExecuteEventObject(w, r)
	return
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		gApiMgr.logger.Info(fmt.Sprintln("%s\t%s\t%s\t%s\n",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start)))
	})
}

//
//  This method creates new rest router interface
//
func (mgr *ApiMgr) InstantiateRestRtr() *mux.Router {
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

func (mgr *ApiMgr) GetRestRtr() *mux.Router {
	return mgr.pRestRtr
}
