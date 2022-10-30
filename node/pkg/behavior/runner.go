// runner.go
package behavior

import (
	bully "distributedelection/node/pkg/behavior/bully"
	fredricksonlynch "distributedelection/node/pkg/behavior/fredricksonlynch"
	. "distributedelection/node/pkg/env"
	net "distributedelection/node/pkg/net"
	smlog "distributedelection/tools/smlog"
	"fmt"
	"os"
)

// Let the selected algorithm be run, based on the ALGORITHM param value in the env.Cfg structure.
// Previous check on that value are already previously done.
func Run() {
	net.InitializeNetMW()
	smlog.Initialize(false, Cfg.TERMINAL_SMLOG_LEVEL)
	if Cfg.ALGORITHM == DE_ALGORITHM_BULLY {
		bully.Run()
	}
	if Cfg.ALGORITHM == DE_ALGORITHM_FREDRICKSONLYNCH {
		fredricksonlynch.Run()
	}
	fmt.Println("Unreachable code while trying to running algorithm %s, something is really odd", Cfg.ALGORITHM)
	os.Exit(1)
}
