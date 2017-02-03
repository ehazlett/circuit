package local

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/containernetworking/cni/libcni"
	"github.com/ehazlett/circuit/config"
	"github.com/ehazlett/circuit/ds"
	"github.com/sirupsen/logrus"
)

func (l *localDS) GetNetwork(name string) (*libcni.NetworkConfig, error) {
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

	var network *libcni.NetworkConfig
	if err := json.Unmarshal(data, &network); err != nil {
		return nil, err
	}
	return network, nil
}

func (l *localDS) GetNetworks() ([]*libcni.NetworkConfig, error) {
	basePath := filepath.Join(l.statePath, networksPath)
	nets, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	var networks []*libcni.NetworkConfig
	for _, p := range nets {
		n, err := l.GetNetwork(p.Name())
		if err != nil {
			logrus.Warnf("unable to get info for network: %s", p.Name())
		}

		networks = append(networks, n)
	}

	return networks, nil
}

func (l *localDS) GetNetworkPeer(name string, containerPid int) (*config.PeerInfo, error) {
	peers, err := l.GetNetworkPeers(name)
	if err != nil {
		return nil, err
	}

	for _, info := range peers {
		if info.ContainerPid == containerPid {
			return info, nil
		}
	}

	return nil, ds.ErrNetworkPeerDoesNotExist
}

func (l *localDS) GetNetworkPeers(name string) (map[string]*config.PeerInfo, error) {
	peerPath := l.peerPath(name)
	if _, err := os.Stat(peerPath); err != nil {
		if os.IsNotExist(err) {
			return map[string]*config.PeerInfo{}, nil
		} else {
			return nil, err
		}
	}
	data, err := ioutil.ReadFile(peerPath)
	if err != nil {
		return nil, err
	}

	peers := map[string]*config.PeerInfo{}
	if err := json.Unmarshal(data, &peers); err != nil {
		return nil, err
	}

	return peers, nil
}
