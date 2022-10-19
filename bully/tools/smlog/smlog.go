package smlog

import (
	"os"

	"fmt"

	loggo "github.com/juju/loggo"
)

var curState = ""
var IsServReg bool

func SetStateSMLogger(state string) {
	curState = state
}

var loggoLogger = loggo.GetLogger("")

func InitLogger(isServRegExec bool) {
	IsServReg = isServRegExec
	if IsServReg {
		fmt.Println("[ServReg] Time     Lvl   Event  Description")
		fmt.Println("[ServReg] -------- ----- ------ ---------------")
	} else {
		fmt.Println("[SM] Time     Lvl   State Event  Description")
		fmt.Println("[SM] -------- ----- ----- ------ ---------------")
	}

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

	if IsServReg {
		loggoLogger.Logf(level, " "+curState+" "+typee.Short()+" "+message, args...)
	} else {
		loggoLogger.Logf(level, " "+typee.Short()+" "+message, args...)
	}

}
func init() {
	loggoLogger.SetLogLevel(loggo.TRACE)
	loggo.ReplaceDefaultWriter(NewSMColorWriter(os.Stderr))
}

func SecondCritical(message string) {
	//second.Criticalf("wwwwwwwww" + message)

}
