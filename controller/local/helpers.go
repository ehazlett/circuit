package local

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/containernetworking/cni/libcni"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

const (
	MaxInterfaceCount = 10
)

func (c *localController) getCniConfig(networkName string, confPath string, containerPid int, ifaceName string) (*libcni.CNIConfig, *libcni.NetworkConfig, *libcni.RuntimeConf, error) {
	cfg, err := c.ds.GetNetwork(networkName)
	if err != nil {
		return nil, nil, nil, err
	}

	f, err := os.Create(filepath.Join(confPath, "01-circuit.conf"))
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	if err := ioutil.WriteFile(filepath.Join(confPath, "01-circuit.conf"), cfg.Bytes, 0644); err != nil {
		return nil, nil, nil, err
	}

	netconf, err := libcni.LoadConf(confPath, networkName)
	if err != nil {
		return nil, nil, nil, err
	}

	cninet := &libcni.CNIConfig{
		Path: c.config.CNIPath,
	}

	rt := &libcni.RuntimeConf{
		ContainerID: fmt.Sprintf("%d", containerPid),
		NetNS:       fmt.Sprintf("/proc/%d/ns/net", containerPid),
		IfName:      ifaceName,
	}

	return cninet, netconf, rt, nil
}

func (c *localController) generateIfaceName(containerPid int) (string, error) {
	originalNs, err := netns.Get()
	if err != nil {
		return "", err

	}
	defer originalNs.Close()

	cntNs, err := netns.GetFromPid(containerPid)
	if err != nil {
		return "", err
	}
	defer cntNs.Close()

	ifaceName := ""
	netns.Set(cntNs)
	for i := 0; i < MaxInterfaceCount; i++ {
		n := fmt.Sprintf("eth%d", i)
		if _, err := netlink.LinkByName(n); err != nil {
			if !strings.Contains(err.Error(), "no such network interface") {
				ifaceName = n
				break
			}
		}
	}
	netns.Set(originalNs)

	if ifaceName == "" {
		return "", fmt.Errorf("unable to generate device name; maximum number of devices reached (%d)", MaxInterfaceCount)
	}

	return ifaceName, nil
}

func (c *localController) getContainerIfaceNames(containerPid int) ([]string, error) {
	originalNs, err := netns.Get()
	if err != nil {
		return nil, err
	}
	defer originalNs.Close()

	cntNs, err := netns.GetFromPid(containerPid)
	if err != nil {
		return nil, err
	}
	defer cntNs.Close()

	ifaces := []string{}
	netns.Set(cntNs)
	for i := 0; i < MaxInterfaceCount; i++ {
		n := fmt.Sprintf("eth%d", i)
		if _, err := netlink.LinkByName(n); err == nil {
			ifaces = append(ifaces, n)
		}
	}
	netns.Set(originalNs)

	return ifaces, nil
}
