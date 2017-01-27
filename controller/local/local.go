package local

import (
	"errors"
	"net/url"
	"os"

	"github.com/ehazlett/circuit/controller"
	"github.com/ehazlett/circuit/ds"
	"github.com/ehazlett/circuit/ds/local"
	"github.com/ehazlett/circuit/ipam"
	"github.com/ehazlett/circuit/lb"
	"github.com/ehazlett/circuit/lb/ipvs"
	"github.com/sirupsen/logrus"
)

const (
	defaultMTU                    = 1500
	defaultContainerInterfaceName = "eth0"
)

var (
	ErrNetworkExists = errors.New("network already exists")
)

type localController struct {
	config *controller.ControllerConfig
	ds     ds.Backend
	ipam   *ipam.IPAM
	lb     lb.LoadBalancer
}

func NewLocalController(c *controller.ControllerConfig) (*localController, error) {
	// TODO: parse DsURI to create ds backend
	u, err := url.Parse(c.DsURI)
	if err != nil {
		return nil, err
	}

	l := &localController{
		config: c,
	}

	switch u.Scheme {
	case "file":
		// TODO: instantiate file backend and set in controller
		if err := os.MkdirAll(u.Path, 0600); err != nil {
			logrus.Fatalf("error initializing state directory: %s", err)
		}

		ls, err := local.NewLocalDS(u.Path)
		if err != nil {
			logrus.Fatalf("error initializing datastore: %s", err)
		}

		l.ds = ls
	case "consul":
		logrus.Debug("configuring state path for consul")
		// TODO: instantiate consul backend and set in controller
	default:
		logrus.Fatalf("unknown datastore uri: %s", c.DsURI)
	}

	// IPAM
	ipm, err := ipam.NewIPAM(l.ds)
	if err != nil {
		logrus.Fatalf("error initializing ipam: %s", err)
	}
	l.ipam = ipm

	// LoadBalancing
	lb, err := ipvs.NewIPVSLB(l.ds)
	if err != nil {
		logrus.Fatalf("error initializing load balancer: %s", err)
	}
	l.lb = lb

	return l, nil
}
