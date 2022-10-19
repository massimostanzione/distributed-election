// netMiddleware
package net

//var cs pb.DistGrepClient
import (
	pb "fredricksonLynch/pb/serviceRegistry"
	. "fredricksonLynch/tools/smlog"
	smlog "fredricksonLynch/tools/smlog"

	//	. "fredricksonLynch/pkg/serviceRegistry/env"
	//"fredricksonLynch/pkg/serviceRegistry/statemachine"
	//"fredricksonLynch/pkg/serviceRegistry"

	"flag"
	//. "fredricksonLynch/tools/smlog"
	//smlog "fredricksonLynch/tools/smlog"
	//"log"
	"net"

	//	"strconv"

	"google.golang.org/grpc"
)

type DGnode struct {
	pb.UnimplementedDistGrepServer
}

//var NONE = &pb.NONE{}
//var cs pb.DistGrepClient

var w *grpc.Server
var lis net.Listener
var serverConn *grpc.ClientConn //server

const RMI_RETRY_TOLERANCE = 3
const LATE_HB_TOLERANCE = 3

func InitializeNetMW() {
	//	portParam := flag.String("p", "40043", "porta")
	flag.Parse()
	/*	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}*/
	//port = *portParam
	//addr = "localHost:" + port
	//port, _ := strconv.ParseInt(*portParam, 10, 32)
	port := "40042"
	//	smsmlog.Critical(LOG_UNDEFINED, "%d", int32(port))
	//	Me.SetPort(int32(port))
	//Me.SetHost("localhost")
	//addr = "localHost:" + port
	// Parsing input arguments
	/*	filepath := flag.String("f", "../../ILIAD_1STBOOK_IT_ALTERED", "source file to be \"fredricksonLynchp-ed\"")
		substr := flag.String("substr", "Achille", "substr to be searched into the source file")
		serverAddr := flag.String("s", "localHost:40042", "server address and port, in the format ADDRESS:PORT")
		highlight := flag.String("hl", "classic", "[classic/asterisks/none] set substr highlighting in the output\nNOTICE: \"classic\" option may be not available on all systems.")
		help := flag.Bool("help", false, "show this message")

		flag.Parse()
		_, exists := HighlightType[*highlight]
		if !exists {
			fmt.Println("\"-hl\" flag not correctly set.\nSee 'fredricksonLynch -help' for allowed values.")
			os.Exit(-1)
		}
		if *help {
			flag.PrintDefaults()
			os.Exit(0)
		}*/

	// Contacting the server
	// il centrale espone il servizio di identificazione dei nodi
	//host := "localhost"
	serverAddr := "localHost:40042" //TODO configurare
	conn := ConnectToNode(serverAddr)
	serverConn = conn
	/*	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		serverConn = conn
		if err != nil {
			smlog.Fatal(LOG_UNDEFINED,"Error while contacting server on %v:\n %v", serverAddr, err)
		}
	*/
	//defer conn.Close()
	//	locCtx = ctx
	// MI METTO IN ASCOLTO: la porta su cui ascolto è
	// la stessa che invierò al gestore dell'anello,
	// in quanto è quella sulla quale sarò contattato
	// Start listening for incoming calls
	//port := "40046"
	/*liss, err := net.Listen("tcp", host+":"+(string)(port))
	lis = liss
	if err != nil {
		smlog.Fatal(LOG_UNDEFINED,"Error while trying to listen to port %v:\n%v", port, err)
	}*/
	/*
		smsmlog.Info(LOG_NETWORK, "Listening on port %v.", port)
		if err := w.Serve(lis); err != nil {
			smlog.Fatal(LOG_UNDEFINED,"Error while trying to serve request: %v", err)
		}
	*/
	liss, err := net.Listen("tcp", serverAddr)
	lis = liss
	if err != nil {
		smlog.Fatal(LOG_UNDEFINED, "Error while trying to listen to port %v:\n%v", port, err)
	}
	smlog.InfoU("--------------------------")

	// New server instance and service registering
	w = grpc.NewServer()
	pb.RegisterDistGrepServer(w, &DGnode{})
	// Serve incoming calls

	// Defining client interface, to be used to invoke the fredricksonLynch service
	//cs = pb.NewDistGrepClient(serverConn)
}

func ConnectToNode(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		smlog.Fatal(LOG_UNDEFINED, "Error while contacting server on %v:\n %v", addr, err)
	}
	return conn
}
func Listen(host string, port string) {

	smlog.InfoU("Listening on port %v...", port)
	// New server instance and service registering
	s := grpc.NewServer()
	pb.RegisterDistGrepServer(s, &DGserver{})
	// Serve incoming calls
	if err := s.Serve(lis); err != nil {
		smlog.Fatal(LOG_UNDEFINED, "Error while trying to serve request: %v", err)
	}

	//	for pause {
	//	}
}
