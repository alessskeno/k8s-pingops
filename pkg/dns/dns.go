// pkg/dns/dns.go
package dns

import (
	"net"
)

func ResolveDNS(hostname string) ([]string, error) {
	ips, err := net.LookupHost(hostname)
	if err != nil {
		return nil, err
	}
	return ips, nil
}
