package smlog

import (
	"os"

	loggo "github.com/juju/loggo"
)

var curState = ""

func SetState(state string) {
	curState = state
}

//TODO come event in altro file?
type LogEvent uint8

const (
	LOG_UNDEFINED LogEvent = iota
	LOG_MSG_SENT
	LOG_MSG_RECV
	LOG_HB
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
	case LOG_HB:
		return "LOG_HB"
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
	case LOG_HB:
		return "HB   "
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

var loggoLogger = loggo.GetLogger("")

func Trace(typee LogEvent, message string, args ...interface{}) {
	sendToLoggo(loggo.TRACE, typee, message, args...)
}
func Debug(typee LogEvent, message string, args ...interface{}) {
	sendToLoggo(loggo.DEBUG, typee, message, args...)
}
func Info(typee LogEvent, message string, args ...interface{}) {
	sendToLoggo(loggo.INFO, typee, message, args...)
}
func InfoU(message string, args ...interface{}) {
	sendToLoggo(loggo.INFO, LOG_UNDEFINED, message, args...)
}
func Warn(typee LogEvent, message string, args ...interface{}) {
	sendToLoggo(loggo.WARNING, typee, message, args...)
}
func Error(typee LogEvent, message string, args ...interface{}) {
	sendToLoggo(loggo.ERROR, typee, message, args...)
}
func Critical(typee LogEvent, message string, args ...interface{}) {
	sendToLoggo(loggo.CRITICAL, typee, message, args...)
}
func Fatal(typee LogEvent, message string, args ...interface{}) {
	sendToLoggo(loggo.CRITICAL, typee, message, args...)
	os.Exit(1)
}

/*
func Trace(message string, args ...interface{}) {
	Trace(LOG_UNDEFINED, message, args)
}
func Debug(typee LogEvent, message string, args ...interface{}) {
	Debug(LOG_UNDEFINED, message, args)
}

func Info(typee LogEvent, message string, args ...interface{}) {
	Info(LOG_UNDEFINED, message, args)
}

func Warn(typee LogEvent, message string, args ...interface{}) {
	Warn(LOG_UNDEFINED, message, args)
}

func Error(typee LogEvent, message string, args ...interface{}) {
	Error(LOG_UNDEFINED, message, args)
}

func Critical(typee LogEvent, message string, args ...interface{}) {
	Critical(LOG_UNDEFINED, message, args)
}

func Fatal(typee LogEvent, message string, args ...interface{}) {
	Fatal(LOG_UNDEFINED, message, args)
}
*/
func sendToLoggo(level loggo.Level, typee LogEvent, message string, args ...interface{}) {
	loggoLogger.Logf(level, " "+curState+" "+typee.Short()+" "+message, args...)

}
func init() {
	loggoLogger.SetLogLevel(loggo.TRACE)
	loggo.ReplaceDefaultWriter(NewSMColorWriter(os.Stderr))
}

func SecondCritical(message string) {
	//second.Criticalf("wwwwwwwww" + message)

}
