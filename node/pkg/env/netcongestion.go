// netcongestion
package env

type NetCongestionLevel uint8

const (
	NCL_ABSENT NetCongestionLevel = iota
	NCL_LIGHT
	NCL_MEDIUM
	NCL_SEVERE
	NCL_CUSTOM
)
