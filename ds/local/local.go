package local

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/ehazlett/circuit/config"
	"github.com/ehazlett/circuit/ds"
)

const (
	networksPath = "networks"
	configName   = "config.json"
	ipConfigName = "ips.json"
)

type localDS struct {
	statePath string
}

func NewLocalDS(statePath string) (*localDS, error) {
	l := &localDS{
		statePath: statePath,
	}

	if err := os.MkdirAll(filepath.Join(l.statePath, networksPath), 0700); err != nil {
		return nil, err
	}

	return l, nil
}

func (l *localDS) netPath(netName string) string {
	return filepath.Join(l.statePath, networksPath, netName)
}

func (l *localDS) addrPath(netName, ip string) string {
	return filepath.Join(l.netPath(netName), ip)
}

func (l *localDS) SaveNetwork(network *config.Network) error {
	netPath := l.netPath(network.Name)
	configPath := filepath.Join(netPath, configName)

	if err := saveData(network, configPath); err != nil {
		return err
	}

	return nil
}

func (l *localDS) GetNetwork(name string) (*config.Network, error) {
	configPath := filepath.Join(l.netPath(name), configName)
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return nil, ds.ErrNetworkDoesNotExist
		} else {
			return nil, err
		}
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var network *config.Network
	if err := json.Unmarshal(data, &network); err != nil {
		return nil, err
	}

	return network, nil
}

func (l *localDS) DeleteNetwork(name string) error {
	netPath := l.netPath(name)
	return os.RemoveAll(netPath)
}

func (l *localDS) SaveIPAddr(ip, network string) error {
	netPath := l.netPath(network)
	ipConfigPath := filepath.Join(netPath, ipConfigName)
	if _, err := os.Stat(ipConfigPath); err != nil {
		// ignore not exists error; it's the first one
		if !os.IsNotExist(err) {
			return err
		}
	}

	ips := []net.IP{}
	currentIPs, err := l.GetNetworkIPs(network)
	if err != nil {
		return err
	}

	ips = append(currentIPs, net.ParseIP(ip))
	if err := saveData(ips, ipConfigPath); err != nil {
		return err
	}

	return nil
}

func (l *localDS) DeleteIPAddr(ip, network string) error {
	// TODO: improve this with a map
	netPath := l.netPath(network)
	ipConfigPath := filepath.Join(netPath, ipConfigName)
	if _, err := os.Stat(ipConfigPath); err != nil {
		// if there are no IPs configured then there is nothing to delete
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}
	ips := []net.IP{}
	currentIPs, err := l.GetNetworkIPs(network)
	if err != nil {
		return err
	}

	for _, i := range currentIPs {
		if i.String() != ip {
			ips = append(ips, i)
		}
	}

	if err := saveData(ips, ipConfigPath); err != nil {
		return err
	}

	return nil
}

func (l *localDS) GetNetworkIPs(name string) ([]net.IP, error) {
	ips := []net.IP{}
	ipConfigPath := filepath.Join(l.netPath(name), ipConfigName)
	if _, err := os.Stat(ipConfigPath); err != nil {
		if os.IsNotExist(err) {
			return ips, nil
		} else {
			return nil, err
		}
	}

	data, err := ioutil.ReadFile(ipConfigPath)
	if err != nil {
		return ips, nil
	}

	if err := json.Unmarshal(data, &ips); err != nil {
		return ips, err
	}

	return ips, nil
}

func (l *localDS) GetNetworks() ([]*config.Network, error) {
	basePath := filepath.Join(l.statePath, networksPath)
	nets, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}
	var networks []*config.Network

	for _, p := range nets {
		n, err := l.GetNetwork(p.Name())
		if err != nil {
			logrus.Warnf("unable to get info for network: %s", p.Name())
		}

		networks = append(networks, n)
	}

	return networks, nil
}

func saveData(d interface{}, fPath string) error {
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	basePath := filepath.Dir(fPath)
	logrus.Debugf("ds: creating base from path: %s base=%s", fPath, basePath)
	if err := os.MkdirAll(basePath, 0700); err != nil {
		return err
	}
	if err := ioutil.WriteFile(fPath, data, 0600); err != nil {
		return err
	}

	return nil
}
