// noderecord.go
package env

import (
	"fmt"
)

type NodeRecord struct {

	// TODO documentazione: reportedAsFailed in modo da poter riassegnare stesso numero se il falllimento è temporaneo
	Id               int
	Host             string
	Port             int32
	ReportedAsFailed bool
}

func (record *NodeRecord) GetId() int {
	return record.Id
}
func (record *NodeRecord) GetHost() string {
	return record.Host
}
func (record *NodeRecord) GetPort() int32 {
	return record.Port
}
func (record *NodeRecord) GetFullAddress() string {
	return record.Host + ":" + fmt.Sprint(record.Port)
}
