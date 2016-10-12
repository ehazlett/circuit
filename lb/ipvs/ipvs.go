package ipvs

import (
	"fmt"
	"net"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/ehazlett/circuit/config"
	"github.com/tehnerd/gnl2go"
)

type ipvsLB struct {
}

func getIPVS() (*gnl2go.IpvsClient, error) {
	i := new(gnl2go.IpvsClient)
	if err := i.Init(); err != nil {
		return nil, err
	}

	return i, nil
}

func NewIPVSLB() (*ipvsLB, error) {
	l := &ipvsLB{}
	return l, nil
}

func (i *ipvsLB) CreateService(svc *config.Service) error {
	ipvs, err := getIPVS()
	if err != nil {
		return err
	}

	protocol, err := getProtocol(svc.Protocol)
	if err != nil {
		return err
	}

	host, port, err := getHostPort(svc.Addr)
	if err != nil {
		return err
	}

	scheduler := fmt.Sprintf("%s", svc.Scheduler)
	if err := ipvs.AddService(host, port, protocol, scheduler); err != nil {
		return err
	}

	return nil
}

func (i *ipvsLB) RemoveService(svc *config.Service) error {
	ipvs, err := getIPVS()
	if err != nil {
		return err
	}

	protocol, err := getProtocol(svc.Protocol)
	if err != nil {
		return err
	}

	host, port, err := getHostPort(svc.Addr)
	if err != nil {
		return err
	}

	if err := ipvs.DelService(host, port, protocol); err != nil {
		return err
	}

	return nil
}

func (i *ipvsLB) AddTargetsToService(serviceAddr string, p config.Protocol, targets []string) error {
	ipvs, err := getIPVS()
	if err != nil {
		return err
	}

	protocol, err := getProtocol(p)
	if err != nil {
		return err
	}

	host, port, err := getHostPort(serviceAddr)
	if err != nil {
		return err
	}

	for _, target := range targets {
		thost, tport, err := getHostPort(target)
		if err != nil {
			return err
		}

		logrus.Debugf("adding server: %s", target)
		if err := ipvs.AddDestPort(host, port, thost, tport, protocol, 10, gnl2go.IPVS_MASQUERADING); err != nil {
			return err
		}
	}

	return nil
}

func (i *ipvsLB) RemoveTargetsFromService(serviceAddr string, p config.Protocol, targets []string) error {
	ipvs, err := getIPVS()
	if err != nil {
		return err
	}

	protocol, err := getProtocol(p)
	if err != nil {
		return err
	}

	host, port, err := getHostPort(serviceAddr)
	if err != nil {
		return err
	}

	for _, target := range targets {
		thost, tport, err := getHostPort(target)
		if err != nil {
			return err
		}

		logrus.Debugf("removing server: %s", target)
		if err := ipvs.DelDestPort(host, port, thost, tport, protocol); err != nil {
			return err
		}
	}

	return nil
}

func (i *ipvsLB) ClearServices() error {
	ipvs, err := getIPVS()
	if err != nil {
		return err
	}

	return ipvs.Flush()
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
