package env

import (
	//	pb "bully/cmd/node/intf"
	//"bully/node/pkg/net"

	//"bully/pkg/node/statemachine"
	"os"
	"time"
)

type NetCongestionLevel uint8

const (
	NCL_ABSENT NetCongestionLevel = iota
	NCL_LIGHT
	NCL_MEDIUM
	NCL_SEVERE
	NCL_CUSTOM
)

type ConfigEnv struct {
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

var DEFAULT_CONFIG_ENV = &ConfigEnv{40043, "localhost", 40042, "Info", false, 1000, 500, 500, 10, 1000, 1000, 3, 3, NCL_ABSENT, 0, 500}

var Cfg *ConfigEnv = &ConfigEnv{}
var Me *SMNode = &SMNode{}
var ServRegAddr string
var CoordId int32 = -1
var Pause = false
var ElectionTimer *time.Timer
var IsElectionStarted = false

var SuccessfulHB = -1
var Events chan (string)
var ElectionChannel chan (*MsgElection)
var OkChannel chan (*MsgOk)
var CoordChannel chan (*MsgCoordinator)
var Heartbeat chan (*MsgHeartbeat)
var Sigchan chan (os.Signal)
