package env

import (
	//	pb "fredricksonLynch/cmd/node/intf"
	//"fredricksonLynch/pkg/node/net"

	//"fredricksonLynch/pkg/node/statemachine"
	"os"
	"time"
)

var Me *SMNode = &SMNode{}
var CoordId int32 = -1
var Pause = false

const SIMULATION_DELAY = 0 * time.Second

// TODO file di configurazione con i parametri
const HB_TIMEOUT = 1000 * time.Millisecond
const HB_TOLERANCE = HB_TIMEOUT / 2 //250 * time.Millisecond

var Heartbeat chan (*MsgHeartbeat)
var Events chan (string)
var ElectionChannel chan (*MsgElection)
var CoordChannel chan (*MsgCoordinator)
var Sigchan chan (os.Signal)
