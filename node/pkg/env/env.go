package env

import (
	"time"
)

// Environment variables specifying all the parameters that define
// the configuration for the execution of the program.
type ConfigEnv struct {
	// the algorithm to be run
	ALGORITHM DEAlgorithm

	// node port to be listening to
	NODE_PORT int

	// information about the service registry
	SERVREG_HOST string
	SERVREG_PORT int64

	// log and output
	TERMINAL_SMLOG_LEVEL string
	VERBOSE              bool

	// monitoring parameters
	MONITORING_TIMEOUT   float32
	MONITORING_TOLERANCE float32

	// bully-specific parameters
	ELECTION_ESPIRY           int
	ELECTION_ESPIRY_TOLERANCE int

	// network management
	RESPONSE_TIME_LIMIT int
	IDLE_WAIT_LIMIT     int

	// fault tolerance
	RMI_RETRY_TOLERANCE int

	// simulation of network delays
	NCL_CONGESTION_LEVEL NetCongestionLevel
	NCL_CUSTOM_DELAY_MIN float32
	NCL_CUSTOM_DELAY_MAX float32
}

var DEFAULT_CONFIG_ENV = &ConfigEnv{
	DE_ALGORITHM_UNDEFINED, // ALGORITHM
	40043,                  // NODE_PORT
	"localhost",            // SERVERG_HOST
	40042,                  // SERVERG_PORT
	"INFO",                 // TERMINAL_LOG_LEVEL
	false,                  // VERBOSE
	1000,                   // HB_TIMEOUT
	500,                    // HB_TOLERANCE
	500,                    // ELECTION_ESPIRY
	10,                     // ELECTION_ESPIRY_TOLERANCE
	1000,                   // RESPONSE_TIME_LIMIT
	3000,                   // IDLE_WAIT_LIMIT
	1,                      // RMI_RETRY_TOLERANCE
	NCL_ABSENT,             // NCL_CONGESTION_LEVEL
	0,                      // NCL_CUSTOM_DELAY_MIN
	500,                    // NCL_CUSTOM_DELAY_MAX
}

var Cfg *ConfigEnv = &ConfigEnv{}

type NodeState struct {
	NodeInfo    *SMNode
	ServRegAddr string
	Coordinator int32
	Participant bool
}

var State *NodeState = &NodeState{}

// Node knowledge
//var Me *SMNode = &SMNode{}
var NextNode *SMNode = &SMNode{}

//var ServRegAddr string

// Election
//var CoordId int32 = -1

// Monitoring (needed to be global)
var Heartbeat chan (*MsgHeartbeat)
var SuccessfulHB = -1

// Bully-specific variables (needed to be global)
var ElectionTimer *time.Timer
var IsElectionStarted = false

// WaitingMaps
type WaitingStruct struct {
	Name    MsgType
	Waiting bool
	Timer   *time.Timer
}

var WaitingMap = map[MsgType]*WaitingStruct{}

// limit servReg requests if network is not changed,
// i.e. if no election has occurred
var DirtyNetList = false
