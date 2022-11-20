package smlog

import (
	. "distributedelection/node/env"
	. "distributedelection/tools/formatting"
	"fmt"
	"os"
	"strings"

	loggo "github.com/juju/loggo"
)

var curState = ""
var IsServReg bool

var loggoLogger = loggo.GetLogger("")
var verboseFile *os.File

func Initialize(isServRegExec bool, levelParam string) {
	if Cfg.VERBOSE {
		fmt.Println("Verbose log will be placed in distributed-election/docs/verbose folder.")
		str := CurState.NodeInfo.GetFullAddr()
		err := os.MkdirAll("../docs/verbose", 0666)
		if err != nil {
			fmt.Println(err)
		}
		a, err := os.OpenFile("../docs/verbose/verbose"+str+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		verboseFile = a
		if err != nil {
			fmt.Println(err)
		}
		//defer f.Close()
	}
	IsServReg = isServRegExec
	if IsServReg {
		fmt.Println("[ServReg] Time     Lvl   Event  Description")
		fmt.Println("[ServReg] -------- ----- ------ ---------------")
	} else {
		fmt.Println("[Node] Time     Lvl   Prtcp Event  Description")
		fmt.Println("[Node] -------- ----- ----- ------ ---------------")
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
	if Cfg.VERBOSE && (typee == LOG_MSG_SENT || typee == LOG_MSG_RECV) {

		loggo.ReplaceDefaultWriter(NewWriter(verboseFile))
		message = strings.ReplaceAll(message, ColorBlkBckgrYellow, "")
		message = strings.ReplaceAll(message, ColorBlkBckgrGreen, "")
		message = strings.ReplaceAll(message, BoldBlack, "")
		message = strings.ReplaceAll(message, ColorReset, "")

		loggoLogger.Logf(level, " "+typee.Short()+" "+message, args...)
		loggo.ReplaceDefaultWriter(NewSMColorWriter(os.Stderr))
	}
	if IsServReg {
		loggoLogger.Logf(level, " "+typee.Short()+" "+message, args...)
	} else {
		loggoLogger.Logf(level, " "+participantToString()+" "+typee.Short()+" "+message, args...)
	}
}
func init() {
	loggo.ReplaceDefaultWriter(NewSMColorWriter(os.Stderr))
}
