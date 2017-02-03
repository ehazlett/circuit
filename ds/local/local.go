package local

import (
	"os"
	"path/filepath"
	"sync"
)

const (
	networksPath     = "networks"
	servicesPath     = "services"
	configName       = "config.json"
	peerConfigName   = "peers.json"
	targetConfigName = "targets.json"
)

type localDS struct {
	statePath string
	lock      *sync.Mutex
}

func NewLocalDS(statePath string) (*localDS, error) {
	l := &localDS{
		statePath: statePath,
		lock:      &sync.Mutex{},
	}

	dirs := []string{networksPath, servicesPath}

	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(l.statePath, d), 0700); err != nil {
			return nil, err
		}
	}

	return l, nil
}

func (l *localDS) netPath(netName string) string {
	return filepath.Join(l.statePath, networksPath, netName)
}

func (l *localDS) peerPath(netName string) string {
	return filepath.Join(l.netPath(netName), peerConfigName)
}
