package ipvs

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ehazlett/circuit/ds"
	"github.com/sirupsen/logrus"
	"github.com/tehnerd/gnl2go"
)

type ipvsLB struct {
	ds ds.Backend
}

var (
	kernelModules = []string{
		"ip_vs",
		"ip_vs_rr",
		"ip_vs_wrr",
		"ip_vs_lc",
		"ip_vs_wlc",
		"ip_vs_lblc",
		"ip_vs_lblcr",
	}
)

func loadModules() error {
	logrus.Debug("checking kernel modules")
	modules, err := exec.Command("lsmod").CombinedOutput()
	if err != nil {
		return err
	}

	if strings.Index(string(modules), "ip_vs") > -1 {
		logrus.Debug("modules are loaded")
		// already loaded
		return nil
	}

	logrus.Debug("kernel modules not found; attempting to load")
	for _, mod := range kernelModules {
		if out, err := exec.Command("modprobe", "-va", mod).CombinedOutput(); err != nil {
			return fmt.Errorf("error loading ipvs module %s: %s err=%s", mod, strings.TrimSpace(string(out)), err)
		}
	}

	return nil
}

func getIPVS() (*gnl2go.IpvsClient, error) {

	if err := loadModules(); err != nil {
		return nil, err
	}

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

func (i *ipvsLB) ClearServices() error {
	ipvs, err := getIPVS()
	if err != nil {
		return err
	}

	return ipvs.Flush()
}
