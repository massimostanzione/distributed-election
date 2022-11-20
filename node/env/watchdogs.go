// watchdogs.go
package env

import "time"

// Watchdogs: detect if a message within an election gets lost,
//            e.g.: a node that is processing a message fails
//		      during the processing
type Watchdog struct {
	Name    MsgType
	Waiting bool
	Timer   *time.Timer
}

var Watchdogs = map[MsgType]*Watchdog{}

func SetWatchdog(msgType MsgType, active bool) {
	Watchdogs[msgType].Waiting = active
	if active {
		Watchdogs[msgType].Timer.Reset(time.Duration(Cfg.IDLE_WAIT_LIMIT) * time.Millisecond)
	} else {
		Watchdogs[msgType].Timer.Stop()
	}
}
