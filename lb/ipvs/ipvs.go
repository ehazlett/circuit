package ipvs

import (
	"fmt"
	"net"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/ehazlett/circuit/config"
	"github.com/ehazlett/circuit/ds"
	"github.com/tehnerd/gnl2go"
)

type ipvsLB struct {
	ds ds.Backend
}

func getIPVS() (*gnl2go.IpvsClient, error) {
	i := new(gnl2go.IpvsClient)
	if err := i.Init(); err != nil {
		return nil, err
	}

	return i, nil
}

func NewIPVSLB(b ds.Backend) (*ipvsLB, error) {
	l := &ipvsLB{
		ds: b,
	}
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

	if err := i.ds.SaveService(svc); err != nil {
		return err
	}

	return nil
}

func (i *ipvsLB) RemoveService(serviceName string) error {
	ipvs, err := getIPVS()
	if err != nil {
		return err
	}

	svc, err := i.ds.GetService(serviceName)
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

	if err := i.ds.DeleteService(svc.Name); err != nil {
		return err
	}

	return nil
}

func (i *ipvsLB) AddTargetsToService(serviceName string, targets []string) error {
	ipvs, err := getIPVS()
	if err != nil {
		return err
	}

	svc, err := i.ds.GetService(serviceName)
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

	for _, target := range targets {
		thost, tport, err := getHostPort(target)
		if err != nil {
			return err
		}

		logrus.Debugf("adding server: %s", target)
		if err := ipvs.AddDestPort(host, port, thost, tport, protocol, 10, gnl2go.IPVS_MASQUERADING); err != nil {
			return err
		}
		if err := i.ds.AddTargetToService(serviceName, target); err != nil {
			return err
		}
	}

	return nil
}

func (i *ipvsLB) RemoveTargetsFromService(serviceName string, targets []string) error {
	ipvs, err := getIPVS()
	if err != nil {
		return err
	}

	svc, err := i.ds.GetService(serviceName)
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

	for _, target := range targets {
		thost, tport, err := getHostPort(target)
		if err != nil {
			return err
		}

		logrus.Debugf("removing server: %s", target)
		if err := ipvs.DelDestPort(host, port, thost, tport, protocol); err != nil {
			return err
		}

		if err := i.ds.RemoveTargetFromService(serviceName, target); err != nil {
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

func (i *ipvsLB) GetServices() ([]*config.Service, error) {
	return i.ds.GetServices()
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
