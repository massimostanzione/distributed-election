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

const SIMULATION_DELAY = 0 * time.Second

// TODO file di configurazione con i parametri
const HB_TIMEOUT = 1000 * time.Millisecond
const HB_TOLERANCE = HB_TIMEOUT * 2 //250 * time.Millisecond
const RESPONSE_TIME_LIMIT = 1000 * time.Millisecond
const IDLE_WAIT_LIMIT = 5000 * time.Millisecond

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
