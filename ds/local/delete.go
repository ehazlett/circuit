package local

import "os"

func (l *localDS) DeleteNetwork(name string) error {
	netPath := l.netPath(name)
	return os.RemoveAll(netPath)
}

func (l *localDS) DeleteNetworkPeer(name string, containerPid int) error {
	peers, err := l.GetNetworkPeers(name)
	if err != nil {
		return err
	}

	for ip, info := range peers {
		if info.ContainerPid == containerPid {
			delete(peers, ip)
		}
	}
	peerPath := l.peerPath(name)
	return l.saveData(peers, peerPath)
}

func (l *localDS) ClearNetworkPeers(name string) error {
	peerPath := l.peerPath(name)
	if _, err := os.Stat(peerPath); err == nil {
		if err := os.Remove(peerPath); err != nil {
			return err
		}
	}

	return nil
}
