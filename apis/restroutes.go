package apis

import (
	"config/objects"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/gorilla/mux"
	"math/big"
	"models"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"utils/logging"
)

type ApiRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type ApiRoutes []ApiRoute

type ApiMgr struct {
	logger        *logging.Writer
	objectMgr     *objects.ObjectMgr
	dbHdl         *objects.DbHandler
	apiVer        string
	apiBase       string
	apiBaseConfig string
	apiBaseState  string
	apiBaseAction string
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

func InitializeApiMgr(paramsDir string, logger *logging.Writer, dbHdl *objects.DbHandler, objectMgr *objects.ObjectMgr) *ApiMgr {
	var err error
	mgr := new(ApiMgr)
	mgr.logger = logger
	mgr.dbHdl = dbHdl
	mgr.objectMgr = objectMgr
	mgr.apiVer = "v1"
	mgr.apiBase = "/public/" + mgr.apiVer + "/"
	mgr.apiBaseConfig = mgr.apiBase + "config" + "/"
	mgr.apiBaseState = mgr.apiBase + "state" + "/"
	mgr.apiBaseAction = mgr.apiBase + "action" + "/"
	if mgr.fullPath, err = filepath.Abs(paramsDir); err != nil {
		logger.Err(fmt.Sprintln("Unable to get absolute path for %s, error [%s]\n", paramsDir, err))
		return nil
	}
	mgr.basePath, _ = filepath.Split(mgr.fullPath)
	gApiMgr = mgr
	return mgr
}

//
//  This method reads the model data and creates rest route interfaces.
//
func (mgr *ApiMgr) InitializeRestRoutes() bool {
	var rt ApiRoute
	for key, _ := range models.ConfigObjectMap {
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
			rt = ApiRoute{key + "Action",
				"POST",
				mgr.apiBaseAction + key,
				HandleRestRouteAction,
			}
			mgr.restRoutes = append(mgr.restRoutes, rt)
		}
	}
	return true
}

func ConfigMgrCheck(certPath string, keyPath string) error {
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return err
	} else if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return err
	}
	return nil
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			gApiMgr.logger.Err(fmt.Sprintln("Unable to marshal ECDSA private key: %v", err))
			return nil
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func ConfigMgrGenerate(certPath string, keyPath string) error {
	var priv interface{}
	var err error
	validFor := 365 * 24 * time.Hour
	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		gApiMgr.logger.Err(fmt.Sprintln("failed to generate private key: %s\n", err))
		return err
	}

	var notBefore time.Time
	notBefore = time.Now()
	notAfter := notBefore.Add(validFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		gApiMgr.logger.Err(fmt.Sprintln("failed to generate serial number: %s\n", err))
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"SnapRoute"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				template.IPAddresses = append(template.IPAddresses, ipnet.IP)
			}
		}
	}
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		gApiMgr.logger.Err(fmt.Sprintln("Failed to create certificate: %s\n", err))
		return err
	}

	certOut, err := os.Create(certPath)
	if err != nil {
		gApiMgr.logger.Err(fmt.Sprintln("failed to open "+certPath+" for writing: %s\n", err))
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		gApiMgr.logger.Err(fmt.Sprintln("failed to open "+keyPath+" for writing:", err))
		return err
	}
	pem.Encode(keyOut, pemBlockForKey(priv))
	keyOut.Close()
	return nil
}

func HandleRestRouteCreate(w http.ResponseWriter, r *http.Request) {
	/*
		// TODO: this will be uncommented for session authentication
		resource := strings.TrimPrefix(r.URL.String(), gMgr.apiBaseConfig)
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		// When a user logs in - pair[0] contains username and pair[1] contains password
		// when a user logs out -  pair[0] contains username and pair[1] contains session ID
		// All other configurations - pair[0] contains session ID
		switch resource {
		case "Login":
			userName := pair[0]
			password := pair[1]
			if sessionId, ok := LoginUser(userName, password); ok {
				var resp LoginResponse
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				resp.SessionId = sessionId
				js, err := json.Marshal(resp)
				if err != nil {
					http.Error(w, SRErrString(SRRespMarshalErr), http.StatusInternalServerError)
					return
				} else {
					w.Write(js)
				}
				logger.Printf("User %s logged in. Session id %d\n", userName, sessionId)
				return
			} else {
				http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
				logger.Println("Login failed for user ", userName)
				return
			}
		case "Logout":
			userName := pair[0]
			sessionId, _ := strconv.ParseUint(pair[1], 10, 64)
			if ok := LogoutUser(userName, sessionId); ok {
				var resp LoginResponse
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				resp.SessionId = sessionId
				js, err := json.Marshal(resp)
				if err != nil {
					http.Error(w, SRErrString(SRRespMarshalErr), http.StatusInternalServerError)
					return
				} else {
					w.Write(js)
				}
				logger.Printf("Logout: User %s Session %d\n", userName, sessionId)
				return
			} else {
				http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
				logger.Println("Logout failed for user ", userName)
				return
			}
		default:
			sessionId, _ := strconv.ParseUint(pair[0], 10, 64)
			if sessionId != 0 {
				if ok:= AuthenticateSessionId(sessionId); ok == false {
					http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
					return
				}
			}
			if IsLocalObject(r) != true && CheckIfSystemIsReady() != true {
				http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
				return
			}
			ConfigObjectCreate(w, r)
		}
	*/
	ConfigObjectCreate(w, r)
	return
}

func HandleRestRouteDeleteForId(w http.ResponseWriter, r *http.Request) {
	/*
		// TODO: this will be uncommented for session authentication
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		sessionId, _ := strconv.ParseUint(pair[0], 10, 64)
		if ok:= AuthenticateSessionId(sessionId); ok == false {
			http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
			return
		}
		if IsLocalObject(r) != true && CheckIfSystemIsReady() != true {
			http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
			return
		}
		ConfigObjectDeleteForId(w, r)
	*/
	ConfigObjectDeleteForId(w, r)
	return
}

func HandleRestRouteDelete(w http.ResponseWriter, r *http.Request) {
	/*
		// TODO: this will be uncommented for session authentication
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		sessionId, _ := strconv.ParseUint(pair[0], 10, 64)
		if ok:= AuthenticateSessionId(sessionId); ok == false {
			http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
			return
		}
		if IsLocalObject(r) != true && CheckIfSystemIsReady() != true {
			http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
			return
		}
		ConfigObjectDelete(w, r)
	*/
	ConfigObjectDelete(w, r)
	return
}

func HandleRestRouteUpdateForId(w http.ResponseWriter, r *http.Request) {
	/*
		// TODO: this will be uncommented for session authentication
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		sessionId, _ := strconv.ParseUint(pair[0], 10, 64)
		if ok:= AuthenticateSessionId(sessionId); ok == false {
			http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
			return
		}
		if IsLocalObject(r) != true && CheckIfSystemIsReady() != true {
			http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
			return
		}
		ConfigObjectUpdateForId(w, r)
	*/
	ConfigObjectUpdateForId(w, r)
	return
}

func HandleRestRouteUpdate(w http.ResponseWriter, r *http.Request) {
	/*
		// TODO: this will be uncommented for session authentication
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		sessionId, _ := strconv.ParseUint(pair[0], 10, 64)
		if ok:= AuthenticateSessionId(sessionId); ok == false {
			http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
			return
		}
		if IsLocalObject(r) != true && CheckIfSystemIsReady() != true {
			http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
			return
		}
		ConfigObjectUpdate(w, r)
	*/
	ConfigObjectUpdate(w, r)
	return
}

func HandleRestRouteGetConfigForId(w http.ResponseWriter, r *http.Request) {
	/*
		// TODO: this will be uncommented for session authentication
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		sessionId, _ := strconv.ParseUint(pair[1], 10, 64)
		if ok:= AuthenticateSessionId(sessionId); ok == false {
			http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
			return
		}
		if IsLocalObject(r) != true && CheckIfSystemIsReady() != true {
			http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
			return
		}
		GetOneConfigObjectForId(w, r)
	*/
	GetOneConfigObjectForId(w, r)
}

func HandleRestRouteGetConfig(w http.ResponseWriter, r *http.Request) {
	/*
		// TODO: this will be uncommented for session authentication
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		sessionId, _ := strconv.ParseUint(pair[1], 10, 64)
		if ok:= AuthenticateSessionId(sessionId); ok == false {
			http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
			return
		}
		if IsLocalObject(r) != true && CheckIfSystemIsReady() != true {
			http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
			return
		}
		GetOneConfigObject(w, r)
	*/
	GetOneConfigObject(w, r)
}

func HandleRestRouteGetStateForId(w http.ResponseWriter, r *http.Request) {
	/*
		// TODO: this will be uncommented for session authentication
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		sessionId, _ := strconv.ParseUint(pair[1], 10, 64)
		if ok:= AuthenticateSessionId(sessionId); ok == false {
			http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
			return
		}
		if IsLocalObject(r) != true && CheckIfSystemIsReady() != true {
			http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
			return
		}
		GetOneStateObjectForId(w, r)
	*/
	GetOneStateObjectForId(w, r)
}

func HandleRestRouteGetState(w http.ResponseWriter, r *http.Request) {
	/*
		// TODO: this will be uncommented for session authentication
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		sessionId, _ := strconv.ParseUint(pair[1], 10, 64)
		if ok:= AuthenticateSessionId(sessionId); ok == false {
			http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
			return
		}
		if IsLocalObject(r) != true && CheckIfSystemIsReady() != true {
			http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
			return
		}
		GetOneStateObject(w, r)
	*/
	GetOneStateObject(w, r)
}

func HandleRestRouteBulkGetConfig(w http.ResponseWriter, r *http.Request) {
	/*
		// TODO: this will be uncommented for session authentication
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		sessionId, _ := strconv.ParseUint(pair[0], 10, 64)
		if ok:= AuthenticateSessionId(sessionId); ok == false {
			http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
			return
		}
		if IsLocalObject(r) != true && CheckIfSystemIsReady() != true {
			http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
			return
		}
		BulkGetConfigObjects(w, r)
	*/
	BulkGetConfigObjects(w, r)
	return
}

func HandleRestRouteBulkGetState(w http.ResponseWriter, r *http.Request) {
	/*
		// TODO: this will be uncommented for session authentication
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		sessionId, _ := strconv.ParseUint(pair[0], 10, 64)
		if ok:= AuthenticateSessionId(sessionId); ok == false {
			http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
			return
		}
		if IsLocalObject(r) != true && CheckIfSystemIsReady() != true {
			http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
			return
		}
		BulkGetStateObjects(w, r)
	*/
	BulkGetStateObjects(w, r)
	return
}

func HandleRestRouteAction(w http.ResponseWriter, r *http.Request) {
	/*
		// TODO: this will be uncommented for session authentication
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		sessionId, _ := strconv.ParseUint(pair[0], 10, 64)
		if ok:= AuthenticateSessionId(sessionId); ok == false {
			http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
			return
		}
		if IsLocalObject(r) != true && CheckIfSystemIsReady() != true {
			http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
			return
		}
		ExecuteActionObject(w, r)
	*/
	ExecuteActionObject(w, r)
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
