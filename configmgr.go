package main

import (
	"encoding/json"
	"database/sql"
	"github.com/gorilla/mux"
	"models"
	"net/http"
	"path/filepath"
	"time"
	"os"
	"net"
	"strconv"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"strings"
	"encoding/base64"
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
	sessionId          uint64
	sessionChan        chan uint64
}

type LoginResponse struct {
	SessionId     uint64     `json: "SessionId"`
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
			logger.Printf("Unable to marshal ECDSA private key: %v", err)
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
		logger.Printf("failed to generate private key: %s\n", err)
		return err
	}

	var notBefore time.Time
	notBefore = time.Now()
	notAfter := notBefore.Add(validFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		logger.Printf("failed to generate serial number: %s\n", err)
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
			if  ipnet.IP.To4() != nil {
				template.IPAddresses = append(template.IPAddresses, ipnet.IP)
			}
		}
	}
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		logger.Printf("Failed to create certificate: %s\n", err)
		return err
	}

	certOut, err := os.Create(certPath)
	if err != nil {
		logger.Printf("failed to open "+certPath+" for writing: %s\n", err)
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logger.Println("failed to open "+keyPath+" for writing:", err)
		return err
	}
	pem.Encode(keyOut, pemBlockForKey(priv))
	keyOut.Close()
	return nil
}

func HandleRestRouteShowConfig(w http.ResponseWriter, r *http.Request) {
	auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)
	sessionId, _ := strconv.ParseUint(pair[1], 10, 64)
	if ok:= AuthenticateSessionId(sessionId); ok == false {
		http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
		return
	}
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	ShowConfigObject(w, r, sessionId)
}

func HandleRestRouteCreate(w http.ResponseWriter, r *http.Request) {
	resource := strings.TrimPrefix(r.URL.String(), gMgr.apiBase)
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
		if CheckIfSystemIsReady(w) != true {
			http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
			return
		}
		ConfigObjectCreate(w, r, sessionId)
	}
	return
}

func HandleRestRouteDelete(w http.ResponseWriter, r *http.Request) {
	auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)
	sessionId, _ := strconv.ParseUint(pair[0], 10, 64)
	if ok:= AuthenticateSessionId(sessionId); ok == false {
		http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
		return
	}
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	ConfigObjectDelete(w, r, sessionId)
	return
}

func HandleRestRouteUpdate(w http.ResponseWriter, r *http.Request) {
	auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)
	sessionId, _ := strconv.ParseUint(pair[0], 10, 64)
	if ok:= AuthenticateSessionId(sessionId); ok == false {
		http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
		return
	}
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	ConfigObjectUpdate(w, r, sessionId)
	return
}

func HandleRestRouteGet(w http.ResponseWriter, r *http.Request) {
	auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)
	sessionId, _ := strconv.ParseUint(pair[0], 10, 64)
	if ok:= AuthenticateSessionId(sessionId); ok == false {
		http.Error(w, SRErrString(SRAuthFailed), http.StatusUnauthorized)
		return
	}
	if CheckIfSystemIsReady(w) != true {
		http.Error(w, SRErrString(SRSystemNotReady), http.StatusServiceUnavailable)
		return
	}
	ConfigObjectsBulkGet(w, r, sessionId)
	return
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
