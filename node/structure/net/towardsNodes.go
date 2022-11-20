// faultTolerance.glob.go
package net

import (
	"context"
	. "distributedelection/node/env"
	. "distributedelection/tools/api"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"

	pbn "distributedelection/node/pb"
	"time"

	pb "distributedelection/node/pb"
	ncl "distributedelection/node/tools/ncl"
	sa "distributedelection/node/tools/structadapter"
	. "distributedelection/tools/formatting"

	"google.golang.org/grpc"
)

type DEANode struct {
	pbn.UnimplementedDistrElectNodeServer
}

func updateLocalCache() *SMNode {
	requested := AskForNodeInfo(CurState.NodeInfo.GetId() + 1)
	if requested.GetId() != 1 {
		NextNode = requested
		smlog.Debug(LOG_UNDEFINED, "NextNode initialized.", NextNode)
	}
	return requested
}

// Unique point of RMI invocation for sending message purposes.
// This is done for fault tolerance, trying to avoid new election when not strictly necessary.
// - an RMI invocation is tried RMI_RETRY_TOLERANCE times, to address temporary failures due to,
//   e.g., temporary node overload
// - FL-specific: if next node in the ring is failed, try with the next one
// - FL-specific: if next node is the election starter and it is failed, stop forwarding its
//   message, to avoid making it turn the ring more than once. The starter will eventually
//   start another election by itself.
func SafeRMI_Ring(msgType MsgType, dest *SMNode, electionMsg *MsgElectionFL, coordMsg *MsgCoordinator) (failedNodeExistence bool) {
	// update local cache the first time a sequential message is sent
	if NextNode.GetId() == 0 {
		dest = updateLocalCache()
	}
	var actualDest *SMNode
	var starter int32
	var rmiErr error
	success := false
	attempts := 0
	for {
		if attempts == 0 {
			actualDest = dest
		} else {
			actualDest = AskForNodeInfo(actualDest.GetId() + 1)
			smlog.Info(LOG_NETWORK, "Trying next node: %v@%v", actualDest.GetId(), actualDest.GetFullAddr())
			attempts = 0
		}
		nextId := actualDest.GetId()
		nextAddr := actualDest.GetFullAddr()

		// if I am the next node, I am the only node into the ring
		if nextAddr == CurState.NodeInfo.GetFullAddr() {
			smlog.InfoU("No other nodes found in the ring.")
			break
		}

		// connect to the next node
		conn := ConnectToNode(nextAddr)
		defer conn.Close()
		nodeClient := pb.NewDistrElectNodeClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Millisecond)
		defer cancel()

		// attempt to send message
		for {
			rmiErr = error(nil)
			attempts++
			ncl.SimulateDelay()
			// Do the actual RMI invocation, based on the message type
			switch msgType {
			case MSG_ELECTION_FL:
				starter = electionMsg.GetStarter()
				netMsg := sa.ToNetElectionFLMsg(electionMsg)
				smlog.Debug(LOG_NETWORK, "Attempt: ELECT %v to %d@%s", netMsg, nextId, nextAddr)
				_, rmiErr = nodeClient.ForwardElectionFL(ctx, netMsg)
				if rmiErr == nil {
					smlog.Info(LOG_MSG_SENT, ColorBlkBckgrGreen+BoldBlack+"ELECT %v to %d@%s"+ColorReset, netMsg, nextId, nextAddr)
				}
				break
			case MSG_COORDINATOR:
				starter = coordMsg.GetStarter()
				netMsg := sa.ToNetCoordinatorMsg(coordMsg)
				smlog.Debug(LOG_NETWORK, "Attempt: COORD %v to %d@%s", netMsg, nextId, nextAddr)
				_, rmiErr = nodeClient.ForwardCoordinator(ctx, netMsg)
				if rmiErr == nil {
					smlog.Info(LOG_MSG_SENT, ColorBlkBckgrGreen+BoldBlack+"COORD %v to %d@%s"+ColorReset, netMsg, nextId, nextAddr)
				}
				break
			default:
				smlog.Fatal(LOG_UNDEFINED, "Invalid/not handled msg type in SafeRMI_Ring")
				break
			}

			if rmiErr == nil {
				break
			}
			// failed node detected
			if attempts >= Cfg.RMI_RETRY_TOLERANCE {
				break
			}
			// ... but we can try again
			smlog.Warn(LOG_NETWORK, "Failed attempt n. %d to contact %d@%v.", attempts, nextId, nextAddr)
			smlog.Debug(LOG_NETWORK, "(%s)", rmiErr)
			smlog.Debug(LOG_NETWORK, "Trying again...", attempts, nextId, nextAddr)
		}
		if rmiErr == nil {
			// no error occurred
			smlog.Trace(LOG_NETWORK, "RMI invoked correctly, exiting from SafeRMI...")
			success = true
			break
		}

		// no more attempts are planned for nextNode: it is for sure failed
		smlog.Error(LOG_NETWORK, "Could not invoke RMI on %d@%v", nextId, nextAddr)
		smlog.Debug(LOG_NETWORK, "(%s)", rmiErr)

		// if starter is failed, stop letting the message to be forwarded a second time
		// into the ring, the election starter will eventually start a new election
		// by itself
		if nextId == starter {
			smlog.Error(LOG_NETWORK, "Election starter failed! Stopping forwarding...")
			break
		}
	}
	return !success
}

// Unique point of RMI invocation for sending message purposes.
// This is done for fault tolerance, trying to avoid new election when not strictly necessary.
// - an RMI invocation is tried RMI_RETRY_TOLERANCE times, to address temporary failures due to,
//   e.g., temporary node overload
func SafeRMI(msgType MsgType, dest *SMNode, electMsg *MsgElectionBully, okMsg *MsgOk, coordMsg *MsgCoordinator) (failedNodeExistence bool) {
	var rmiErr error
	success := false
	attempts := 0
	nextId := dest.GetId()
	nextAddr := dest.GetFullAddr()
	// connect to the next node
	conn := ConnectToNode(nextAddr)
	defer conn.Close()
	nodeClient := pb.NewDistrElectNodeClient(conn)

	//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Millisecond)

	// attempt to send message
	for {
		rmiErr = error(nil)
		attempts++
		ncl.SimulateDelay()

		// Do the actual RMI invocation, based on the message type
		switch msgType {
		case MSG_ELECTION_BULLY:
			netMsg := sa.ToNetElectionBullyMsg(electMsg)
			smlog.Debug(LOG_NETWORK, "Attempt: ELECT %v to %d@%s", netMsg, nextId, nextAddr)
			_, rmiErr = nodeClient.ForwardElectionBully(ctx, netMsg)
			if rmiErr == nil {
				smlog.Info(LOG_MSG_SENT, ColorBlkBckgrGreen+BoldBlack+"ELECT %v to %d@%s"+ColorReset, netMsg, nextId, nextAddr)
			}
			break
		case MSG_OK:
			netMsg := sa.ToNetOkMsg(okMsg)
			smlog.Debug(LOG_NETWORK, "Attempt: OK %v to %d@%s", netMsg, nextId, nextAddr)
			_, rmiErr = nodeClient.ForwardOk(ctx, netMsg)
			if rmiErr == nil {
				smlog.Info(LOG_MSG_SENT, ColorBlkBckgrGreen+BoldBlack+"OK %v to %d@%s"+ColorReset, netMsg, nextId, nextAddr)
			}
			break
		case MSG_COORDINATOR:
			netMsg := sa.ToNetCoordinatorMsg(coordMsg)
			smlog.Debug(LOG_NETWORK, "Attempt: COORD %v to %d@%s", netMsg, nextId, nextAddr)
			_, rmiErr = nodeClient.ForwardCoordinator(ctx, netMsg)
			if rmiErr == nil {
				smlog.Info(LOG_MSG_SENT, ColorBlkBckgrGreen+BoldBlack+"COORD %v to %d@%s"+ColorReset, netMsg, nextId, nextAddr)
			}
			break
		default:
			smlog.Fatal(LOG_UNDEFINED, "Invalid/not handled msg type in SafeRMI")
			break
		}
		defer cancel()

		if rmiErr == nil {
			break
		}
		// failed node detected
		if attempts >= Cfg.RMI_RETRY_TOLERANCE {
			break
		}

		// ... but we can try again
		smlog.Warn(LOG_NETWORK, "Failed attempt n. %d to contact %d@%v.", attempts, nextId, nextAddr)
		smlog.Debug(LOG_NETWORK, "(%s)", rmiErr)
		smlog.Warn(LOG_NETWORK, "Trying again...", attempts, nextId, nextAddr)
	}
	if rmiErr == nil {
		// no error occurred
		smlog.Trace(LOG_NETWORK, "RMI invoked correctly, exiting from SafeRMI...")
		success = true
	} else {
		// no more attempts are planned for nextNode: it is for sure failed
		smlog.Error(LOG_NETWORK, "Could not invoke RMI on %d@%v", nextId, nextAddr)
		smlog.Debug(LOG_NETWORK, "(%s)", rmiErr)
	}
	return !success
}

func SafeHB(hb *pbn.Heartbeat, node *SMNode) {
	connN := ConnectToNode(node.GetFullAddr())
	defer connN.Close()
	// New server instance and service registering
	nodoServer := grpc.NewServer()
	pbn.RegisterDistrElectNodeServer(nodoServer, &DEANode{})
	csN := pbn.NewDistrElectNodeClient(connN)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Millisecond)
	defer cancel()
	smlog.Debug(LOG_NETWORK, "Attempt: HB to %d@%s", node.GetId(), node.GetFullAddr())
	ncl.SimulateDelay()
	_, errq := csN.SendHeartBeat(ctx, hb)
	if errq != nil {
		smlog.Error(LOG_NETWORK, "error while contacting %d@%v", node.GetId(), node.GetFullAddr())
		smlog.Debug(LOG_NETWORK, "(%s)", errq)
	} else {
		smlog.Info(LOG_MSG_SENT, "HB to %d@%s", node.GetId(), node.GetFullAddr())
	}
}
