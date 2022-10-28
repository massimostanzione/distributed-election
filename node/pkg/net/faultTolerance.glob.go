// faultTolerance.glob.go
package net

import (
	. "distributedelection/node/pkg/env"

	"math/rand"
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
		return NCL_ABSENT
	}
}
func GenerateDelay() int32 {
	var min float32
	var max float32
	switch Cfg.NCL_CONGESTION_LEVEL {
	case NCL_ABSENT:
		min = 0
		max = 0
	case NCL_LIGHT:
		min = 0
		max = .2 * Cfg.HB_TIMEOUT
	case NCL_MEDIUM:
		min = .3 * Cfg.HB_TIMEOUT
		max = .5 * Cfg.HB_TIMEOUT
	case NCL_SEVERE:
		min = .5 * Cfg.HB_TIMEOUT
		max = 1.5 * Cfg.HB_TIMEOUT
	case NCL_CUSTOM:
		min = Cfg.NCL_CUSTOM_DELAY_MIN
		max = Cfg.NCL_CUSTOM_DELAY_MAX

	}
	ret := (rand.Float32() * (max - min)) + min
	return int32(ret)
}
