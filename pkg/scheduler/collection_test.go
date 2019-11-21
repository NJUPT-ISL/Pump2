package scheduler

import (
	"testing"
)

func TestAddNodeInfo(t *testing.T) {
	if err := AddNodeInfo("127.0.0.1:5020");err != nil {
		println(err)
	}

}
