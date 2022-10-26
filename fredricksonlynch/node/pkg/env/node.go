// nde.go
package env

import (
	"fmt"
	. "fredricksonlynch/tools/smlog"
	smlog "fredricksonlynch/tools/smlog"
)

type SMNode struct {
	Id   int32
	Host string
	Port int32
}

//TODO per relazione: specificare che getter/setters non sono necessariamente richiesti per golang
// ref.: https://go.dev/doc/effective_go#Getters

func (msg *SMNode) SetId(id int32) {
	msg.Id = id
}
func (msg *SMNode) SetHost(host string) {
	msg.Host = host
}
func (msg *SMNode) SetPort(portp int32) {
	smlog.Critical(LOG_UNDEFINED, "%d", portp)
	msg.Port = portp
}
func (msg *SMNode) GetId() int32 {
	return msg.Id
}
func (msg *SMNode) GetHost() string {
	return msg.Host
}
func (msg *SMNode) GetPort() int32 {
	return msg.Port
}
func (msg *SMNode) GetFullAddr() string {
	//str, _ := strconv.ParseInt(msg.Port, 10, 32)
	return msg.Host + ":" + fmt.Sprint(msg.Port)
}