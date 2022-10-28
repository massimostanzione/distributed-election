package env

var Events chan (string)
var ElectionChannel chan (*MsgElection)
var OkChannel chan (*MsgOk)
var CoordChannel chan (*MsgCoordinator)
