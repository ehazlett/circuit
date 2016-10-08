package local

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

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
	data, err := json.Marshal(network)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(netPath, 0700); err != nil {
		return err
	}

	if err := ioutil.WriteFile(configPath, data, 0600); err != nil {
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

func (l *localDS) SaveIPAddr(ip, network string) error {
	netPath := l.netPath(network)
	ipConfigPath := filepath.Join(netPath, ipConfigName)
	if _, err := os.Stat(ipConfigPath); err != nil {
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
	data, err := json.Marshal(ips)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(netPath, 0700); err != nil {
		return err
	}
	if err := ioutil.WriteFile(ipConfigPath, data, 0600); err != nil {
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
