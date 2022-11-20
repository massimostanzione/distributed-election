// runner.go
package behavior

import (
	. "distributedelection/node/env"
	bully "distributedelection/node/structure/behavior/bully"
	fredricksonlynch "distributedelection/node/structure/behavior/fredricksonlynch"
	net "distributedelection/node/structure/net"
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
	fmt.Printf("Unreachable code while trying to run algorithm %s, something is really odd\n", Cfg.ALGORITHM)
	os.Exit(1)
}
