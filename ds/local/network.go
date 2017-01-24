package local

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/ehazlett/circuit/config"
	"github.com/ehazlett/circuit/ds"
)

func (l *localDS) netPath(netName string) string {
	return filepath.Join(l.statePath, networksPath, netName)
}

func (l *localDS) addrPath(netName, ip string) string {
	return filepath.Join(l.netPath(netName), ip)
}

func (l *localDS) SaveNetwork(network *config.Network) error {
	netPath := l.netPath(network.Name)
	configPath := filepath.Join(netPath, configName)

	logrus.Debugf("saving network: %+v", network)
	if err := l.saveData(network, configPath); err != nil {
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

func (l *localDS) SaveIPAddr(ip, networkName string, containerPid int, peerType config.PeerType) error {
	network, err := l.GetNetwork(networkName)
	if err != nil {
		return err
	}
	logrus.Debugf("network: %+v", network)
	if network.Peers == nil {
		network.Peers = map[string]*config.IPPeer{}
	}

	network.Peers[ip] = &config.IPPeer{
		IP:   ip,
		Pid:  containerPid,
		Type: peerType,
	}

	if err := l.SaveNetwork(network); err != nil {
		return err
	}

	return nil
}

func (l *localDS) DeleteIPAddr(ip, networkName string) error {
	network, err := l.GetNetwork(networkName)
	if err != nil {
		return err
	}

	if _, ok := network.Peers[ip]; ok {
		delete(network.Peers, ip)
	}

	if err := l.SaveNetwork(network); err != nil {
		return err
	}

	return nil
}

func (l *localDS) GetNetworkPeers(name string) (map[string]*config.IPPeer, error) {
	network, err := l.GetNetwork(name)
	if err != nil {
		return nil, err
	}

	return network.Peers, nil
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
