
package main

import (
)

// SR error codes
cons (
	SRFail                                         = 0
	SRSuccess                                      = 1
	SRSystemNotReady                               = 2
)

// SR error strings
var ErrString = map[int]string {
	SRFail:                                        "Configuration failed"
	SRSuccess:                                     "Configuration applied successfully"
	SRSystemNotReady:                              "System not ready"
}

//Given a code reurn error string
func SRErrString(errCode int) string {
	return ErrString[errCode]
}
