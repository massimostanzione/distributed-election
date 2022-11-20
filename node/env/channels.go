// Channels involved in the execution.
package env

var ElectionChannel_bully chan (*MsgElectionBully)
var ElectionChannel_fl chan (*MsgElectionFL)
var OkChannel chan (*MsgOk)
var CoordChannel chan (*MsgCoordinator)

// Monitoring (needed to be global)
var Heartbeat chan (*MsgHeartbeat)
