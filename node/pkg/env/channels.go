package env

var Events chan (string)
var ElectionChannel_bully chan (*MsgElectionBully)
var ElectionChannel_fl chan (*MsgElectionFL)
var OkChannel chan (*MsgOk)
var CoordChannel chan (*MsgCoordinator)
