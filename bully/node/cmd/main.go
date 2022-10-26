package main

import (
	. "bully/node/pkg/env"
	. "bully/node/pkg/net"
	. "bully/node/pkg/statemachine"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"

	//"reflect"
	"strconv"
	"strings"
	"syscall"

	"github.com/go-ini/ini"
)

func main() {
	loadConfig()
	Sigchan = make(chan os.Signal, 1)
	signal.Notify(Sigchan, syscall.SIGTSTP)
	go func() {
		sig := <-Sigchan
		log.Printf("SEGNALE!", sig)
		Pause = !Pause
		SwitchServerState(Pause)
		log.Printf("Pause = %v", Pause)
	}()
	StartStateMachine()
}

// Load Config. Priority order: i) default values, ii) INI file, iii) flags
func loadConfig() {
	// Parse flags
	iniPath := flag.String("c", "./../config.ini", "path of a INI file containing environment configuration")
	nodeport := flag.Int("p", 40043, "target port")
	servicereghost := flag.String("sh", "localhost", "host of the service registry, e.g. \"localhost\", 127.0.0.1 or whatever IP address")
	serviceregport := flag.Int64("sp", 40042, "target port of the service registry")
	ncl := flag.String("ncl", "ABSENT", "network congestion level: ABSENT, LIGHT, MEDIUM, SEVERE, CUSTOM. If custom, please specify flags -nclmin and -nclmax")
	nclmin := flag.Int64("nclmin", 0, "minimum network delay, valid only if -ncl CUSTOM is set")
	nclmax := flag.Int64("nclmax", 500, "maximum network delay, valid only if -ncl CUSTOM is set")
	help := flag.Bool("help", false, "show this message")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// 1. set default parameters, to be sure that all parameters are set
	Cfg = DEFAULT_CONFIG_ENV

	// 2 if present, override parameters via INI file
	//   (could not have all the parameters set, or some of them could be invalid)
	iniFile, err := ini.Load(*iniPath)
	if err != nil {
		log.Printf("Could not load %s.\nLoading default parameters...", *iniPath)
	} else {
		iniSections := iniFile.Sections()
		for i := 0; i < len(iniSections); i++ {
			sectName := iniSections[i].Name()
			sect, err := iniFile.GetSection(sectName)
			if err != nil {
				log.Printf("err")
			}
			sect.MapTo(&Cfg)
		}
		key, _ := iniFile.Section("delay-conf").GetKey("NCL_CONGESTION_LEVEL")
		Cfg.NCL_CONGESTION_LEVEL = ToNCL(key.String())
		if Cfg.NCL_CUSTOM_DELAY_MIN > Cfg.NCL_CUSTOM_DELAY_MAX {
			log.Printf("Cannot consider min delay > max delay! Falling back to default value.")
			Cfg.NCL_CUSTOM_DELAY_MIN = 0
			Cfg.NCL_CUSTOM_DELAY_MAX = 500
		}
	}

	// 3. if flags are set, override INI parameters
	if isFlagPassed("ncl") {
		parsed := ToNCL(*ncl)
		if parsed == NCL_CUSTOM {
			if isFlagPassed("nclmin") && isFlagPassed("nclmax") {
				if *nclmin <= *nclmax {
					Cfg.NCL_CONGESTION_LEVEL = NCL_CUSTOM
					Cfg.NCL_CUSTOM_DELAY_MIN = float32(*nclmin)
					Cfg.NCL_CUSTOM_DELAY_MAX = float32(*nclmax)
				} else {
					log.Printf("Cannot consider min delay > max delay! Falling back to default value.")
				}
			} else {
				log.Printf("CUSTOM specified as net congestion level but without -nclmin and -nclmax. Falling back to default parameters.")
			}
		} else {
			Cfg.NCL_CONGESTION_LEVEL = parsed
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

	Me.SetHost(GetOutboundIP())
	Me.SetPort(int32(Cfg.NODE_PORT))

	ServRegAddr = Cfg.SERVREG_HOST + ":" + strconv.FormatInt(Cfg.SERVREG_PORT, 10)
	log.Println(*Cfg)
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	split := strings.Split(localAddr.String(), ":")
	return split[0]
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
