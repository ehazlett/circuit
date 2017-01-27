package local

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

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
