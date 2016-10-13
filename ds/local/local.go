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
	ipConfigName     = "ips.json"
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
