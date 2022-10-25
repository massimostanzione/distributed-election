package env

import (
	//	pb "bully/cmd/node/intf"
	//"bully/pkg/node/net"

	//"bully/pkg/node/statemachine"
	"os"
	"time"
)

var Me *SMNode = &SMNode{}
var ServRegAddr string
var CoordId int32 = -1
var Pause = false
var ElectionTimer *time.Timer
var IsElectionStarted = false

const SIMULATION_DELAY = 0 * time.Second

// TODO file di configurazione con i parametri
/*const HB_TIMEOUT = 1000 * time.Millisecond
const HB_TOLERANCE = HB_TIMEOUT / 2 //250 * time.Millisecond

var Heartbeat chan (*MsgHeartbeat)*/

const HB_TIMEOUT = 1000             // ms
const HB_TOLERANCE = HB_TIMEOUT * 2 //250 * time.Millisecond
const ELECTION_ESPIRY = 500 * time.Millisecond
const ELECTION_ESPIRY_TOLERANCE = 10 * time.Millisecond

const RESPONSE_TIME_LIMIT = 1000 * time.Millisecond
const IDLE_WAIT_LIMIT = 5000 * time.Millisecond

var SuccessfulHB = -1
var Events chan (string)
var ElectionChannel chan (*MsgElection)
var OkChannel chan (*MsgOk)
var CoordChannel chan (*MsgCoordinator)
var Heartbeat chan (*MsgHeartbeat)
var Sigchan chan (os.Signal)
