// Outbound RMI management, towards service registry.
package net

import (
	"context"
	. "distributedelection/node/env"
	pbsr "distributedelection/serviceregistry/pb"
	. "distributedelection/tools/api"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"time"
)

type DGservreg struct {
	pbsr.UnimplementedDistrElectServRegServer
}

func AskForJoining() *SMNode {
	CurState.DirtyNetCache = true
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Second)
	defer cancel()
	smlog.Info(LOG_SERVREG, "Asking for joining the ring...")
	node, err := cs.JoinNetwork(ctx, &pbsr.NodeAddr{Host: CurState.NodeInfo.GetHost(), Port: CurState.NodeInfo.GetPort()})
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while executing fredricksonlynch:\n%v", err)
	}
	return ToSMNode(node)
}

func AskForNodeInfo(id int32) *SMNode {
	smlog.Info(LOG_SERVREG, "Asking for info about node n. %d", id)
	if !CurState.DirtyNetCache {
		i := int(id)
		for {
			i = i % len(NetCache)
			if i == 0 {
				i = len(NetCache)
			}
			return NetCache[i-1]
		}
	}
	smlog.Debug(LOG_SERVREG, "NetCache is dirty - now asking servReg for info about node n. %d", id)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Second)
	defer cancel()
	ret, err := cs.GetNode(ctx, &pbsr.NodeId{Id: int32(id)})
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while executing GetNode:\n%v", err)
		return nil
	}
	return &SMNode{Id: ret.GetId(), Host: ret.GetHost(), Port: ret.GetPort()}
}

// For monitoring use only
func AskForAllNodes() []*SMNode {
	var ret []*SMNode
	smlog.Debug(LOG_SERVREG, "Asking for info about all nodes")
	if CurState.DirtyNetCache {
		ret = updateNetCache()
		NetCache = ret
	} else {
		ret = NetCache
	}
	CurState.DirtyNetCache = false
	return ret
}

func AskForNodesWithGreaterIds(baseId int32) []*SMNode {
	smlog.Trace(LOG_SERVREG, "Asking for info about nodes s.t. id > %d", baseId)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Second)
	//	locCtx = ctx
	defer cancel()

	ret, err := cs.GetAllNodesWithIdGreaterThan(ctx, &pbsr.NodeId{Id: int32(baseId)})
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while calling GetAllNodesWithIdGreaterThan:\n%v", err)
		return nil
	}
	var array []*SMNode
	for _, node := range ret.GetList() {
		array = append(array, ToSMNode(node))
	}
	return array
}

func updateNetCache() []*SMNode {
	var ret []*SMNode
	smlog.Debug(LOG_SERVREG, "Election has occurred, so net could have changed. Asking to ServReg...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Second)
	defer cancel()
	allNodesList, err := cs.GetAllNodes(ctx, new(pbsr.EMPTY_SR))
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while executing GetAllNodes:\n%v", err)
	}
	for _, node := range allNodesList.List {
		ret = append(ret, ToSMNode(node))
	}
	return ret
}
