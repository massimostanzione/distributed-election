// logevent.go
package smlog

type LogEvent uint8

const (
	LOG_UNDEFINED LogEvent = iota
	LOG_MSG_SENT
	LOG_MSG_RECV
	LOG_MONITORING
	LOG_NETWORK
	LOG_SERVREG
	LOG_STATEMACHINE
	LOG_ELECTION
)

// String implements Stringer.
func (event LogEvent) String() string {
	switch event {
	case LOG_UNDEFINED:
		return "LOG_UNDEFINED"
	case LOG_MSG_SENT:
		return "LOG_MSG_SENT"
	case LOG_MSG_RECV:
		return "LOG_MSG_RECV"
	case LOG_MONITORING:
		return "LOG_MONITORING"
	case LOG_NETWORK:
		return "LOG_NETWORK"
	case LOG_SERVREG:
		return "LOG_SERVREG"
	case LOG_STATEMACHINE:
		return "LOG_STATEMACHINE"
	case LOG_ELECTION:
		return "LOG_ELECTION"
	default:
		return "<unknown>"
	}
}

// Short returns a five character string to use in
// aligned logging output.
func (event LogEvent) Short() string {
	switch event {
	case LOG_UNDEFINED:
		return "N/D  "
	case LOG_MSG_SENT:
		return "MSENT"
	case LOG_MSG_RECV:
		return "MRECV"
	case LOG_MONITORING:
		return "MONIT"
	case LOG_NETWORK:
		return "NETWK"
	case LOG_SERVREG:
		return "SVREG"
	case LOG_STATEMACHINE:
		return "STMAC"
	case LOG_ELECTION:
		return "ELECT"
	default:
		return "<unknown>"
	}
}
