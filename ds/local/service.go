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

func (l *localDS) servicePath(serviceName string) string {
	return filepath.Join(l.statePath, servicesPath, serviceName)
}
func (l *localDS) SaveService(s *config.Service) error {
	servicePath := l.servicePath(s.Name)
	configPath := filepath.Join(servicePath, configName)

	if err := l.saveData(s, configPath); err != nil {
		return err
	}

	return nil
}

func (l *localDS) DeleteService(name string) error {
	servicePath := l.servicePath(name)
	return os.RemoveAll(servicePath)
}

func (l *localDS) GetService(name string) (*config.Service, error) {
	configPath := filepath.Join(l.servicePath(name), configName)
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return nil, ds.ErrServiceDoesNotExist
		} else {
			return nil, err
		}
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var service *config.Service
	if err := json.Unmarshal(data, &service); err != nil {
		return nil, err
	}

	targets, err := l.GetServiceTargets(name)
	if err != nil {
		return nil, err
	}

	service.Targets = targets

	return service, nil
}

func (l *localDS) GetServiceTargets(serviceName string) ([]string, error) {
	targets := []string{}
	servicePath := l.servicePath(serviceName)
	targetConfigPath := filepath.Join(servicePath, targetConfigName)
	if _, err := os.Stat(targetConfigPath); err != nil {
		// ignore not exists error; it's the first one
		if !os.IsNotExist(err) {
			return targets, err
		}
	}

	data, err := ioutil.ReadFile(targetConfigPath)
	if err != nil {
		return targets, nil
	}

	if err := json.Unmarshal(data, &targets); err != nil {
		return targets, err
	}

	return targets, nil
}

func (l *localDS) AddTargetToService(serviceName, target string) error {
	servicePath := l.servicePath(serviceName)
	targetConfigPath := filepath.Join(servicePath, targetConfigName)
	if _, err := os.Stat(targetConfigPath); err != nil {
		// ignore not exists error; it's the first one
		if !os.IsNotExist(err) {
			return err
		}
	}

	targets := []string{}
	current, err := l.GetServiceTargets(serviceName)
	if err != nil {
		return err
	}

	targets = append(current, target)
	if err := l.saveData(targets, targetConfigPath); err != nil {
		return err
	}

	return nil
}

func (l *localDS) RemoveTargetFromService(serviceName, target string) error {
	servicePath := l.servicePath(serviceName)
	targetConfigPath := filepath.Join(servicePath, ipConfigName)
	if _, err := os.Stat(targetConfigPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}
	targets := []string{}
	current, err := l.GetServiceTargets(serviceName)
	if err != nil {
		return err
	}

	for _, t := range current {
		if t != target {
			targets = append(targets, t)
		}
	}

	if err := l.saveData(targets, targetConfigPath); err != nil {
		return err
	}

	return nil
}

func (l *localDS) GetServices() ([]*config.Service, error) {
	basePath := filepath.Join(l.statePath, servicesPath)
	svcs, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}
	var services []*config.Service

	for _, s := range svcs {
		svc, err := l.GetService(s.Name())
		if err != nil {
			logrus.Warnf("unable to get info for service: %s", s.Name())
		}

		services = append(services, svc)
	}

	return services, nil
}
