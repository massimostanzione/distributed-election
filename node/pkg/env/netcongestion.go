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
