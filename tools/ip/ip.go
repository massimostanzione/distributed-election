// Utility function for IP retrieving.
package ip

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Printf("Could not retrieve IP address: %s", err)
		os.Exit(1)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	split := strings.Split(localAddr.String(), ":")
	return split[0]
}
