// Channels involved in the execution.
package env

// FL - Inbound messages
var MsgOrderIn chan (MsgType)
var ElectChIn chan (*MsgElectionFL)
var CoordChIn chan (*MsgCoordinator)

// FL - Outbound messages
var MsgOrderOut chan (MsgType)
var ElectChOut chan (*MsgElectionFL)
var CoordChOut chan (*MsgCoordinator)

// Bully-specific
var ElectionChannel_bully chan (*MsgElectionBully)
var OkChannel chan (*MsgOk)

// Monitoring (needed to be global)
var Heartbeat chan (*MsgHeartbeat)
