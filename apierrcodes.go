
package main

import (
)

// SR error codes
const (
	SRFail                                         = 0
	SRSuccess                                      = 1
	SRSystemNotReady                               = 2
	SRRespMarshalErr                               = 3
	SRNotFound                                     = 4
	SRIdStoreFail                                  = 5
	SRIdDeleteFail                                 = 6
	SRServerError                                  = 7
	SRObjHdlError                                  = 8
	SRObjMapError                                  = 9
)

// SR error strings
var ErrString = map[int]string {
	SRFail:                                        "Configuration failed.",
	SRSuccess:                                     "Configuration applied successfully.",
	SRSystemNotReady:                              "System not ready.",
	SRRespMarshalErr:                              "Configuration applied successfully. However, failed to marshal response.",
	SRNotFound:                                    "Failed to find entry.",
	SRIdStoreFail:                                  "Failed to store Id in DB. However, configuration has been applied.",
	SRIdDeleteFail:                                "Failed to delete Id from DB. However, configuration has been removed.",
	SRServerError:                                 "Backend server failed to apply configuration.",
	SRObjHdlError:                                 "Failed to get object handle.",
	SRObjMapError:                                 "Failed to get object map.",
}

//Given a code reurn error string
func SRErrString(errCode int) string {
	return ErrString[errCode]
}
