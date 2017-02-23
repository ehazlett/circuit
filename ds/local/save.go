package local

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/containernetworking/cni/libcni"
	"github.com/ehazlett/circuit/config"
	"github.com/sirupsen/logrus"
)

func (l *localDS) SaveNetwork(cfg *libcni.NetworkConfig) error {
	netPath := l.netPath(cfg.Network.Name)
	configPath := filepath.Join(netPath, configName)

	if err := l.saveData(cfg, configPath); err != nil {
		return err
	}

	return nil
}

func (l *localDS) SaveNetworkPeer(name string, containerPid int, ip string, ifaceName string) error {
	peers, err := l.GetNetworkPeers(name)
	if err != nil {
		return err
	}

	peers[ip] = &config.PeerInfo{
		NetworkName:  name,
		ContainerPid: containerPid,
		IP:           ip,
		IfaceName:    ifaceName,
	}
	peerPath := l.peerPath(name)
	return l.saveData(peers, peerPath)
}

func (l *localDS) saveData(d interface{}, fPath string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

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
