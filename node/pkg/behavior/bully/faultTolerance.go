package bully

import (
	"context"
	pb "distributedelection/node/pb"

	. "distributedelection/node/pkg/env"
	. "distributedelection/node/pkg/net"
	. "distributedelection/tools/api"
	. "distributedelection/tools/formatting"
	"time"

	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"

	"google.golang.org/grpc"
)

// Unique point of RMI invocation for sending message purposes.
// This is done for fault tolerance, trying to avoid new election when not strictly necessary.
// - an RMI invocation is tried RMI_RETRY_TOLERANCE times, to address temporary failures due to,
//   e.g., temporary node overload
func SafeRMI(msgType MsgType, dest *SMNode, electMsg *MsgElection, okMsg *MsgOk, coordMsg *MsgCoordinator) (failedNodeExistence bool) { //opt ...interface{}) {
	attempts := 0
	nextNode := dest
	nextAddr := nextNode.GetFullAddr()
	failedNodeExistence = false
	var rmiErr error
	time.Sleep(time.Duration(GenerateDelay()) * time.Millisecond)
	for {
		rmiErr = nil

		// Connect to the next node and register it as gRPC server for the service
		conn := ConnectToNode(nextAddr)
		defer conn.Close()
		serverNode := grpc.NewServer()
		pb.RegisterDistrElectNodeServer(serverNode, &DEANode{})
		nodeClient := pb.NewDistrElectNodeClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Millisecond)
		defer cancel()

		attempts++
		// Do the actual RMI invocation, based on the message type
		switch msgType {
		case MSG_ELECTION:
			netMsg := ToNetElectionMsg(electMsg)
			smlog.Warn(LOG_MSG_SENT, ColorBlkBckgrGreen+BoldBlack+"SENDING ELECT %v to %s"+ColorReset, netMsg, nextAddr)
			_, rmiErr = nodeClient.ForwardElection(ctx, netMsg)
			break
		case MSG_OK:
			netMsg := ToNetOkMsg(okMsg)
			smlog.Warn(LOG_MSG_RECV, ColorBlkBckgrGreen+BoldBlack+"SENDING OK %v to %s"+ColorReset, netMsg, nextAddr)
			_, rmiErr = nodeClient.ForwardOk(ctx, netMsg)
			break
		case MSG_COORDINATOR:
			netMsg := ToNetCoordinatorMsg(coordMsg)
			smlog.Warn(LOG_MSG_RECV, ColorBlkBckgrGreen+BoldBlack+"SENDING COORD %v to %s"+ColorReset, netMsg, nextAddr)
			_, rmiErr = nodeClient.ForwardCoordinator(ctx, netMsg)
			break
		default:
			smlog.Fatal(LOG_UNDEFINED, "Unreachable code, msgType not parsable")
			break
		}

		// Adapt behavior based on the returned error
		if rmiErr != nil {
			if attempts != Cfg.RMI_RETRY_TOLERANCE {
				// failed node detected, but we can try again
				smlog.Warn(LOG_NETWORK, "Failed attempt n. %d to contact %v, trying again...", attempts, nextAddr)
				smlog.Debug(LOG_NETWORK, "(%s)", rmiErr)
			} else {
				// no more attempts are planned for nextNode
				smlog.Error(LOG_NETWORK, "Could not invoke RMI on %v", nextAddr)
				smlog.Debug(LOG_NETWORK, "(%s)", rmiErr)
				failedNodeExistence = true
				break
			}
		} else {
			// no error occurred
			smlog.Trace(LOG_NETWORK, "RMI invoked correctly, exiting from SafeRMI...")
			break
		}
		// try again with another attempt
	}
	return failedNodeExistence
}
