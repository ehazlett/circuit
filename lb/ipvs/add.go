package ipvs

import (
	"github.com/sirupsen/logrus"
	"github.com/tehnerd/gnl2go"
)

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
