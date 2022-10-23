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

const SIMULATION_DELAY = 0 * time.Second

// TODO file di configurazione con i parametri
/*const HB_TIMEOUT = 1000 * time.Millisecond
const HB_TOLERANCE = HB_TIMEOUT / 2 //250 * time.Millisecond

var Heartbeat chan (*MsgHeartbeat)*/
const ELECTION_ESPIRY = 1 * time.Second
const ELECTION_ESPIRY_TOLERANCE = 500 * time.Millisecond

var Events chan (string)
var ElectionChannel chan (*MsgElection)
var OkChannel chan (*MsgOk)
var CoordChannel chan (*MsgCoordinator)
var Sigchan chan (os.Signal)
