package env

import (
	//	pb "fredricksonLynch/cmd/node/intf"
	//"fredricksonLynch/node/pkg/net"

	//"fredricksonLynch/node/pkg/statemachine"
	"os"
	"time"
)

var Me *SMNode = &SMNode{}
var NextNode *SMNode = &SMNode{}
var ServRegAddr string
var CoordId int32 = -1
var Pause = false

type NetCongestionLevel uint8

const (
	NCL_ABSENT NetCongestionLevel = iota
	NCL_LIGHT
	NCL_MEDIUM
	NCL_SEVERE
	NCL_CUSTOM
)

type ConfigEnv struct {
	NODE_PORT            int
	SERVREG_HOST         string
	SERVREG_PORT         int64
	HB_TIMEOUT           float32
	HB_TOLERANCE         float32
	RESPONSE_TIME_LIMIT  int
	IDLE_WAIT_LIMIT      int
	RMI_RETRY_TOLERANCE  int
	LATE_HB_TOLERANCE    int
	NCL_CONGESTION_LEVEL NetCongestionLevel
	NCL_CUSTOM_DELAY_MIN float32
	NCL_CUSTOM_DELAY_MAX float32
}

var DEFAULT_CONFIG_ENV = &ConfigEnv{40043, "localhost", 40042, 1000, 500, 1000, 1000, 3, 3, NCL_ABSENT, 0, 500}

var Cfg *ConfigEnv = &ConfigEnv{}
var SuccessfulHB = -1
var Heartbeat chan (*MsgHeartbeat)
var EventsSend chan (string)
var EventsList chan (string)
var ElectionChannel chan (*MsgElection)
var CoordChannel chan (*MsgCoordinator)
var Sigchan chan (os.Signal)

type WaitingStruct struct {
	Name    MsgType
	Waiting bool
	Timer   *time.Timer
}

var WaitingMap = map[MsgType]*WaitingStruct{}
