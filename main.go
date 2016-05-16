Copyright [2016] [SnapRoute Inc]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

	 Unless required by applicable law or agreed to in writing, software
	 distributed under the License is distributed on an "AS IS" BASIS,
	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	 See the License for the specific language governing permissions and
	 limitations under the License.
package main

import (
	"config/server"
	"flag"
	"fmt"
	"net/http"
	"utils/keepalive"
	"utils/logging"
)

func main() {
	fmt.Println("Starting ConfigMgr daemon")
	paramsDir := flag.String("params", "./params", "Directory Location for config files")
	flag.Parse()
	paramsDirName := *paramsDir
	if paramsDirName[len(paramsDirName)-1] != '/' {
		paramsDirName = paramsDirName + "/"
	}

	fmt.Println("ConfigMgr: Start logger")
	logger, err := logging.NewLogger("confd", "ConfigMgr", true)
	if err != nil {
		fmt.Println("Failed to start logger. Nothing will be logged ...")
	}

	configMgr := server.NewConfigMgr(paramsDirName, logger)
	if configMgr == nil {
		logger.Err("Failed to initialize CONF Mgr. Exiting!!!")
		return
	}

	foundConfPort, confPort := server.GetConfigHandlerPort(paramsDirName)
	restRtr := configMgr.ApiMgr.GetRestRtr()

	// Start keepalive routine
	go keepalive.InitKeepAlive("confd", paramsDirName)
	/*
		// TODO: uncomment this section for https server
		certFile := *paramsDir+"/cert.pem"
		keyFile := *paramsDir+"/key.pem"
		err = ConfigMgrCheck(certFile, keyFile)
		if err != nil {
			err = ConfigMgrGenerate(certFile, keyFile)
			if err != nil {
				syslogger.Info("### CONF Mgr Failed to generate certs")
			}
		}
		if foundConfPort {
			http.ListenAndServeTLS(":"+confPort, certFile, keyFile, restRtr)
		} else
			http.ListenAndServeTLS(":8080", certFile, keyFile, restRtr)
		}
	*/
	if foundConfPort {
		http.ListenAndServe(":"+confPort, restRtr)
	} else {
		http.ListenAndServe(":8080", restRtr)
	}
	logger.Info("CONF Mgr. Exiting!!!")
}
