// runner.go
package behavior

import (
	bully "distributedelection/node/pkg/behavior/bully"
	fredricksonlynch "distributedelection/node/pkg/behavior/fredricksonlynch"
	. "distributedelection/node/pkg/env"
	"fmt"
	"os"
)

// Let the selected algorithm be run, based on the ALGORITHM param value in the env.Cfg structure.
// Previous check on that value are already previously done.
func Run() {
	if Cfg.ALGORITHM == DE_ALGORITHM_BULLY {
		bully.Run()
	}
	if Cfg.ALGORITHM == DE_ALGORITHM_FREDRICKSONLYNCH {
		fredricksonlynch.Run()
	}
	fmt.Println("Unreachable code while trying to running algorithm %s, something is really odd", Cfg.ALGORITHM)
	os.Exit(1)
}
