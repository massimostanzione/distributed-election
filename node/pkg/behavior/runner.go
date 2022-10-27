// runner.go
package behavior

import (
	b "distributedelection/node/pkg/behavior/bully"
	fl "distributedelection/node/pkg/behavior/fredricksonlynch"
	. "distributedelection/node/pkg/env"
	"fmt"
	"os"
)

func Run() {
	if Cfg.ALGORITHM == DE_ALGORITHM_BULLY {
		b.Run()
	}
	if Cfg.ALGORITHM == DE_ALGORITHM_FREDRICKSONLYNCH {
		fl.Run()
	}
	fmt.Println("unreachable code, wtf %s", Cfg.ALGORITHM)
	os.Exit(0)
}
