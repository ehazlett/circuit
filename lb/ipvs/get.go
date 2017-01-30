package ipvs

import (
	"fmt"
	"net"
	"strconv"

	"github.com/ehazlett/circuit/config"
	"github.com/tehnerd/gnl2go"
)

func (i *ipvsLB) GetServices() ([]*config.Service, error) {
	return i.ds.GetServices()
}

func (i *ipvsLB) GetService(name string) (*config.Service, error) {
	svcs, err := i.GetServices()
	if err != nil {
		return nil, err
	}

	for _, svc := range svcs {
		if svc.Name == name {
			return svc, nil
		}
	}

	return nil, nil
}

func getProtocol(p config.Protocol) (uint16, error) {
	var protocol uint16
	switch p {
	case config.ProtocolTCP:
		protocol = uint16(gnl2go.ToProtoNum("tcp"))
	case config.ProtocolUDP:
		protocol = uint16(gnl2go.ToProtoNum("udp"))
	default:
		return 0, fmt.Errorf("unknown protocol: %s", p)
	}

	return protocol, nil
}

func getHostPort(addr string) (string, uint16, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, fmt.Errorf("error parsing service addr: %s", err)
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, fmt.Errorf("error converting port: %s", err)
	}

	return host, uint16(p), nil
}
