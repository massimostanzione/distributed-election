package smlog

import (
	"os"

	. "distributedelection/node/pkg/env"
	"fmt"

	loggo "github.com/juju/loggo"
)

var curState = ""
var IsServReg bool

func SetStateSMLogger(state string) {
	curState = state
}

var loggoLogger = loggo.GetLogger("")

func Initialize(isServRegExec bool, levelParam string) {
	IsServReg = isServRegExec
	if IsServReg {
		fmt.Println("[ServReg] Time     Lvl   Event  Description")
		fmt.Println("[ServReg] -------- ----- ------ ---------------")
	} else {
		fmt.Println("[SM] Time     Lvl   Prtcp Event  Description")
		fmt.Println("[SM] -------- ----- ----- ------ ---------------")
	}
	level, _ := loggo.ParseLevel(levelParam)
	loggoLogger.SetLogLevel(level)
	Info(LOG_UNDEFINED, "Starting...")
	Info(LOG_UNDEFINED, "Type CTRL+C to terminate")
	Info(LOG_UNDEFINED, "------------------------")
}

func participantToString() string {
	if CurState.Participant {
		return "yes  "
	}
	return "no   "
}

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
func sendToLoggo(level loggo.Level, typee LogEvent, message string, args ...interface{}) {

	if IsServReg {
		loggoLogger.Logf(level, " "+typee.Short()+" "+message, args...)
	} else {
		loggoLogger.Logf(level, " "+participantToString()+" "+typee.Short()+" "+message, args...)
	}

}
func init() {
	loggo.ReplaceDefaultWriter(NewSMColorWriter(os.Stderr))
}
