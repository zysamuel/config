package actions

import (
	"testing"
)

func TestInit(t *testing.T) {
	v := InitializeActionMgr(nil, nil)
	t.Error("For", "nil, nil",
		"expected", nil,
		"got", v)
}
