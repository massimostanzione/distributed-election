// Entry point for node application.
package main

import (
	. "distributedelection/node/env"
	runner "distributedelection/node/structure/behavior"
	ncl "distributedelection/node/tools/ncl"
	. "distributedelection/tools/api"
	ip "distributedelection/tools/ip"
	"flag"
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

func main() {
	fmt.Println("+------------------------------------------------+")
	fmt.Println("|         DISTRIBUTED ELECTION ALGORITHMS        |")
	fmt.Println("|github.com/massimostanzione/distributed-election|")
	fmt.Println("+------------------------------------------------+")
	fmt.Println("|                      Node                      |")
	fmt.Println("+------------------------------------------------+")
	fmt.Println("")
	fmt.Println("Loading configuration environment...")
	loadConfig()
	setNodeKnowledge()
	fmt.Println("... done. Starting...")
	runner.Run()
}

// Load Config. Priority order: i) default values, ii) INI file, iii) flags
func loadConfig() {
	// Parse flags
	// (notice: only some of the config parameteres are settable via flags, for practical use)
	algorithm := flag.String("a", "UNDEFINED", "distributed election algorithm to be run, accepted values: [BULLY b FREDRICKSONLYNCH fl]")
	iniPath := flag.String("c", "./../configs/config.ini", "path of a INI file containing environment configuration")
	nodeport := flag.Int("p", 40043, "target port")
	servicereghost := flag.String("sh", "0.0.0.0", "host of the service registry, e.g. \"localhost\", 127.0.0.1 or whatever IP address")
	serviceregport := flag.Int64("sp", 40042, "target port of the service registry")
	nclParam := flag.String("ncl", "ABSENT", "network congestion level: ABSENT, LIGHT, MEDIUM, SEVERE, CUSTOM. If custom, please specify flags -nclmin and -nclmax")
	nclmin := flag.Int64("nclmin", 0, "minimum network delay, valid only if -ncl CUSTOM is set")
	nclmax := flag.Int64("nclmax", 500, "maximum network delay, valid only if -ncl CUSTOM is set")
	logLevel := flag.String("l", "INFO", "log level: CRITICAL, ERROR, WARNING, INFO, DEBUG, TRACE")
	verbose := flag.Bool("v", false, "verbose")
	help := flag.Bool("help", false, "show this message")
	flag.Parse()

	// 1. set default parameters, to be sure that all parameters are set
	Cfg = DEFAULT_CONFIG_ENV

	// 2. if present, override parameters via INI file
	//   (could not have all the parameters set, or some of them could be invalid)
	iniFile, err := ini.Load(*iniPath)
	if err != nil {
		fmt.Printf("[ERROR] Could not load %s.\n%s\nLoading default parameters...\n", *iniPath, err)
	} else {
		iniSections := iniFile.Sections()
		for i := 0; i < len(iniSections); i++ {
			sectName := iniSections[i].Name()
			sect, err := iniFile.GetSection(sectName)
			if err != nil {
				fmt.Println("[ERROR] Cannot get INI section: %s", err)
			}
			sect.MapTo(&Cfg)
		}
		key, _ := iniFile.Section("algorithm").GetKey("ALGORITHM")
		Cfg.ALGORITHM = ParseDEAlgorithm(key.String())
		key, _ = iniFile.Section("delay-conf").GetKey("NCL_CONGESTION_LEVEL")
		Cfg.NCL_CONGESTION_LEVEL = key.String()
		if Cfg.NCL_CUSTOM_DELAY_MIN > Cfg.NCL_CUSTOM_DELAY_MAX {
			fmt.Println("[ERROR] Cannot consider min delay > max delay! Falling back to default value.")
			Cfg.NCL_CUSTOM_DELAY_MIN = 0
			Cfg.NCL_CUSTOM_DELAY_MAX = 500
		}
	}

	// 3. if flags are set, override INI parameters
	//    (notice: only some of them are settable, for practical use)
	algo := ParseDEAlgorithm(*algorithm)
	if algo == DE_ALGORITHM_UNDEFINED && Cfg.ALGORITHM == DE_ALGORITHM_UNDEFINED {
		fmt.Print("Please specify a distributed election algorithm to be run providing an INI file (-c flag) or via -a flag.\n",
			"Admissible values for -a flag:\n",
			"\tBULLY\n\tb\n\tFREDRICKSONLYNCH\n\tfl\n")
		os.Exit(1)
	}
	if isFlagPassed("a") {
		Cfg.ALGORITHM = algo
	}
	if *help {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if isFlagPassed("ncl") {
		parsed := ncl.ToNCL(*nclParam)
		if parsed == ncl.NCL_CUSTOM {
			if isFlagPassed("nclmin") && isFlagPassed("nclmax") {
				if *nclmin <= *nclmax {
					Cfg.NCL_CONGESTION_LEVEL = "CUSTOM"
					Cfg.NCL_CUSTOM_DELAY_MIN = float32(*nclmin)
					Cfg.NCL_CUSTOM_DELAY_MAX = float32(*nclmax)
				} else {
					fmt.Println("[ERROR] Cannot consider min delay > max delay! Falling back to default value.")
				}
			} else {
				fmt.Println("[ERROR] CUSTOM specified as net congestion level but without -nclmin and -nclmax. Falling back to default parameters.")
			}
		} else {
			Cfg.NCL_CONGESTION_LEVEL = *nclParam
		}
	}
	if isFlagPassed("sh") {
		Cfg.SERVREG_HOST = *servicereghost
	}
	if isFlagPassed("sp") {
		Cfg.SERVREG_PORT = *serviceregport
	}
	if isFlagPassed("p") {
		Cfg.NODE_PORT = *nodeport
	}
	if isFlagPassed("l") {
		// conversion is done natively by loggo
		Cfg.TERMINAL_SMLOG_LEVEL = *logLevel
	}
	if isFlagPassed("v") {
		Cfg.VERBOSE = *verbose
	}
}

func setNodeKnowledge() {
	CurState.NodeInfo = &SMNode{}
	CurState.NodeInfo.SetHost(ip.GetOutboundIP())
	CurState.NodeInfo.SetPort(int32(Cfg.NODE_PORT))
	CurState.DirtyNetCache = false
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
