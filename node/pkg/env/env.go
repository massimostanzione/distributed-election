package env

import (
	//	pb "distributedelection/cmd/node/intf"
	//"distributedelection/node/pkg/net"

	//"distributedelection/node/pkg/statemachine"
	"fmt"
	"os"
	"time"
)

var Me *SMNode = &SMNode{}
var NextNode *SMNode = &SMNode{}
var ServRegAddr string
var CoordId int32 = -1
var Pause = false

var Heartbeat chan (*MsgHeartbeat)

// bully
var ElectionTimer *time.Timer
var IsElectionStarted = false

type DEAlgorithm uint8

const (
	DE_ALGORITHM_UNDEFINED DEAlgorithm = iota
	DE_ALGORITHM_BULLY
	DE_ALGORITHM_FREDRICKSONLYNCH
)

func ParseDEAlgorithm(input string) DEAlgorithm {
	bullyFlags := []string{"BULLY", "b", "DE_ALGORITHM_BULLY"}
	flFlags := []string{"FREDRICKSONLYNCH", "fl", "DE_ALGORITHM_FREDRICKSONLYNCH"}
	if contains(bullyFlags, input) {
		fmt.Print("YOOssssssssssssOOOOOOO")
		return DE_ALGORITHM_BULLY
	}
	if contains(flFlags, input) {
		fmt.Print("YOOOOOOOOO")
		return DE_ALGORITHM_FREDRICKSONLYNCH
	}
	fmt.Print("qqqqqqqqqq")
	return DE_ALGORITHM_UNDEFINED
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

type NetCongestionLevel uint8

const (
	NCL_ABSENT NetCongestionLevel = iota
	NCL_LIGHT
	NCL_MEDIUM
	NCL_SEVERE
	NCL_CUSTOM
)

type ConfigEnv struct {
	ALGORITHM                 DEAlgorithm
	NODE_PORT                 int
	SERVREG_HOST              string
	SERVREG_PORT              int64
	TERMINAL_SMLOG_LEVEL      string
	VERBOSE                   bool
	HB_TIMEOUT                float32
	HB_TOLERANCE              float32
	ELECTION_ESPIRY           int
	ELECTION_ESPIRY_TOLERANCE int
	RESPONSE_TIME_LIMIT       int
	IDLE_WAIT_LIMIT           int
	RMI_RETRY_TOLERANCE       int
	LATE_HB_TOLERANCE         int
	NCL_CONGESTION_LEVEL      NetCongestionLevel
	NCL_CUSTOM_DELAY_MIN      float32
	NCL_CUSTOM_DELAY_MAX      float32
}

var DEFAULT_CONFIG_ENV = &ConfigEnv{DE_ALGORITHM_UNDEFINED, 40043, "localhost", 40042, "INFO", false, 1000, 500, 500, 10, 1000, 1000, 1, 3, NCL_ABSENT, 0, 500}

var Cfg *ConfigEnv = &ConfigEnv{}
var SuccessfulHB = -1

var Sigchan chan (os.Signal)
