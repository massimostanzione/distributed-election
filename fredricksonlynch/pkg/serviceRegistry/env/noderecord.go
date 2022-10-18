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

func (record *NodeRecord) GetFullAddress() string {
	return record.Host + ":" + fmt.Sprint(record.Port)
}
