package ipvs

import "github.com/sirupsen/logrus"

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
