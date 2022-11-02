package fredricksonlynch

import (
	"context"
	pb "distributedelection/node/pb"
	. "distributedelection/node/pkg/env"
	. "distributedelection/tools/api"
	. "distributedelection/tools/formatting"

	"time"

	. "distributedelection/node/pkg/net"

	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"

	"google.golang.org/grpc"
)

func RedudantElectionCheck(voter int32, electionMsg *MsgElection) bool {
	for _, i := range electionMsg.GetVoters() {
		if i == voter {
			return true
		}
	}
	return false
}

// Unique point of RMI invocation for sending message purposes.
// This is done for fault tolerance, trying to avoid new election when not strictly necessary.
// - an RMI invocation is tried RMI_RETRY_TOLERANCE times, to address temporary failures due to,
//   e.g., temporary node overload
// - FL-specific: if next node in the ring is failed, try with the next one
// - FL-specific: if next node is the election starter and it is failed, stop forwarding its
//   message, to avoid making it turn the ring more than once. The starter will eventually
//   start another election by itself.
func SafeRMI(msgType MsgType, dest *SMNode, tryNextWhenFailed bool, electionMsg *MsgElection, coordMsg *MsgCoordinator) (failedNodeExistence bool) {

	// update local cache the first time a sequential message is sent
	if NextNode.GetId() == 0 {
		dest = updateLocalCache()
	}

	nextNode := dest
	nextId := nextNode.GetId()
	nextAddr := nextNode.GetFullAddr()

	var rmiErr error
	var starter int32
	time.Sleep(time.Duration(GenerateDelay()) * time.Millisecond)
	for {

		// Connect to the next node and register it as gRPC server for the service
		conn := ConnectToNode(nextAddr)
		defer conn.Close()
		serverNode := grpc.NewServer()
		pb.RegisterDistrElectNodeServer(serverNode, &DEANode{})
		nodeClient := pb.NewDistrElectNodeClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Millisecond)
		defer cancel()

		escape := false
		if !tryNextWhenFailed {
			escape = true
		}
		attempts := 0
		starter = -1
		failedNodeExistence = false
		for {
			rmiErr = nil
			attempts++
			// Do the actual RMI invocation, based on the message type
			switch msgType {
			case MSG_ELECTION:
				starter = electionMsg.GetStarter()
				netMsg := ToNetElectionMsg(electionMsg)
				smlog.Info(LOG_MSG_SENT, ColorBlkBckgrGreen+BoldBlack+"SENDING ELECT %v to %s"+ColorReset, netMsg, nextAddr)
				_, rmiErr = nodeClient.ForwardElection(ctx, netMsg)
				break
			case MSG_COORDINATOR:
				starter = coordMsg.GetStarter()
				netMsg := ToNetCoordinatorMsg(coordMsg)
				smlog.Info(LOG_MSG_RECV, ColorBlkBckgrGreen+BoldBlack+"SENDING COORD %v to %s"+ColorReset, netMsg, nextAddr)
				_, rmiErr = nodeClient.ForwardCoordinator(ctx, netMsg)
				break
			default:
				break
			}

			// Adapt behavior based on the returned error
			if rmiErr != nil {
				if attempts != Cfg.RMI_RETRY_TOLERANCE {
					// failed node detected, but we can try again
					smlog.Warn(LOG_NETWORK, "Failed attempt n. %d to contact %v, trying again...", attempts, nextAddr)
					smlog.Debug(LOG_NETWORK, "(%s)", rmiErr)
				} else {
					// no more attempts are planned for nextNode: it is for sure failed
					smlog.Error(LOG_NETWORK, "Could not invoke RMI on %v", nextAddr)
					smlog.Debug(LOG_NETWORK, "(%s)", rmiErr)
					failedNodeExistence = true

					// if starter is failed, stop letting the message to be forwarded a second time
					// into the ring, the election starter will eventually start a new election
					// by itself
					//escape, tryNextWhenFailed = checkForStarterFailure(nextId, starter)
					if nextId == starter {
						smlog.Error(LOG_NETWORK, "Election starter failed! Stopping forwarding...")
						smlog.Debug(LOG_NETWORK, "(%s)", rmiErr)
						escape = true
						tryNextWhenFailed = false
					}

					if msgType == MSG_HEARTBEAT {
						SuccessfulHB--
					}

					if tryNextWhenFailed {
						nextNode = AskForNodeInfo(nextNode.GetId() + 1)
						smlog.Info(LOG_NETWORK, "Trying next node: %v@%v", nextNode.GetId(), nextNode.GetFullAddr())
						nextId = nextNode.GetId()
						nextAddr = nextNode.GetFullAddr()
						if nextAddr == CurState.NodeInfo.GetFullAddr() {
							smlog.InfoU("Sono rimasto solo io")
							escape = true
						}
					}
				}
			} else {
				// no error occurred
				smlog.Trace(LOG_NETWORK, "RMI invoked correctly, exiting from SafeRMI...")
				escape = true
			}
			if failedNodeExistence || (!failedNodeExistence && escape) {
				break
			}
			// try again with another attempt
		}
		if escape {
			break
		}
		// try again with another node
		attempts = 0
	}
	return failedNodeExistence
}

func updateLocalCache() *SMNode {
	requested := AskForNodeInfo(CurState.NodeInfo.GetId() + 1)
	if requested.GetId() != 1 {
		NextNode = requested
		smlog.Debug(LOG_UNDEFINED, "NextNode initialized.", NextNode)
	}
	return requested
}
