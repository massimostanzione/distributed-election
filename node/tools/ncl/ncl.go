// netcongestion
package ncl

import (
	. "distributedelection/node/env"
	smlog "distributedelection/tools/smlog"
	"math/rand"
	"time"
)

// Network congestion simulator.
// Further details in netMW.

type NetCongestionLevel uint8

const (
	NCL_ABSENT NetCongestionLevel = iota
	NCL_LIGHT
	NCL_MEDIUM
	NCL_SEVERE
	NCL_CUSTOM
)

func ToNCL(input string) NetCongestionLevel {
	switch input {
	case "ABSENT":
		return NCL_ABSENT
	case "LIGHT":
		return NCL_LIGHT
	case "MEDIUM":
		return NCL_MEDIUM
	case "SEVERE":
		return NCL_SEVERE
	case "CUSTOM":
		return NCL_CUSTOM
	default:
		smlog.Warn(smlog.LOG_UNDEFINED, "Cannot parse %s as NCL level. Switching to NCL_ABSENT.")
		return NCL_ABSENT
	}
}
func SimulateDelay() {
	time.Sleep(time.Duration(generateDelay()) * time.Millisecond)
}
func generateDelay() int32 {
	var min float32
	var max float32
	switch ToNCL(Cfg.NCL_CONGESTION_LEVEL) {
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
	default:
		// unreachable, ref. ToNCL()
	}
	ret := (rand.Float32() * (max - min)) + min
	return int32(ret)
}
