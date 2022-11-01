// faultTolerance.glob.go
package net

import (
	. "distributedelection/node/pkg/env"

	"math/rand"
	"time"
)

func GenerateDelay() int32 {
	var min float32
	var max float32
	switch Cfg.NCL_CONGESTION_LEVEL {
	case NCL_ABSENT:
		min = 0
		max = 0
	case NCL_LIGHT:
		min = 0
		max = .2 * Cfg.MONITORING_TIMEOUT
	case NCL_MEDIUM:
		min = .3 * Cfg.MONITORING_TIMEOUT
		max = .5 * Cfg.MONITORING_TIMEOUT
	case NCL_SEVERE:
		min = .5 * Cfg.MONITORING_TIMEOUT
		max = 1.5 * Cfg.MONITORING_TIMEOUT
	case NCL_CUSTOM:
		min = Cfg.NCL_CUSTOM_DELAY_MIN
		max = Cfg.NCL_CUSTOM_DELAY_MAX

	}
	ret := (rand.Float32() * (max - min)) + min
	return int32(ret)
}
func SetWatchdog(msgType MsgType, active bool) {
	Watchdogs[msgType].Waiting = active
	if active {
		Watchdogs[msgType].Timer.Reset(time.Duration(Cfg.IDLE_WAIT_LIMIT) * time.Millisecond)
	} else {
		Watchdogs[msgType].Timer.Stop()
	}
}
