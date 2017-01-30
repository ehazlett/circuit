package ipvs

import (
	"fmt"
	"strings"

	"github.com/ehazlett/circuit/config"
	"github.com/sirupsen/logrus"
)

func (i *ipvsLB) CreateService(svc *config.Service) error {
	ipvs, err := getIPVS()
	if err != nil {
		return fmt.Errorf("error getting ipvs: %s", err)
	}

	protocol, err := getProtocol(svc.Protocol)
	if err != nil {
		return fmt.Errorf("error getting protocol: %s", err)
	}

	host, port, err := getHostPort(svc.Addr)
	if err != nil {
		return fmt.Errorf("error getting host and port: %s", err)
	}

	scheduler := fmt.Sprintf("%s", svc.Scheduler)
	if err := ipvs.AddService(host, port, protocol, scheduler); err != nil {
		if strings.Index(err.Error(), "errorcode is: 17") > -1 {
			logrus.Debugf("service %s appears to be configured", svc.Name)
			return nil
		}

		return fmt.Errorf("error adding ipvs service: %s", err)
	}

	if err := i.ds.SaveService(svc); err != nil {
		return fmt.Errorf("error saving service: %s", err)
	}

	return nil
}
