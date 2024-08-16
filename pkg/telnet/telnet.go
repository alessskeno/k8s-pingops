// pkg/telnet/telnet.go
package telnet

import (
	"fmt"
	"net"
	"time"
)

// CheckTelnetConnection checks if a Telnet connection to the host and port is possible
func CheckTelnetConnection(host string, port int) error {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
