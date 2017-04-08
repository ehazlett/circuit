package local

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/types"
)

func getLocalDS() (*localDS, error) {
	d, err := ioutil.TempDir("", "circuit-test")
	if err != nil {
		return nil, err
	}

	return NewLocalDS(d)
}

func TestLocalSaveNetwork(t *testing.T) {
	l, err := getLocalDS()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(l.statePath)

	n := &libcni.NetworkConfig{
		Network: &types.NetConf{
			Name: "testing",
		},
	}

	if err := l.SaveNetwork(n); err != nil {
		t.Fatal(err)
	}

	network, err := l.GetNetwork(n.Network.Name)
	if err != nil {
		t.Fatal(err)
	}

	if network.Network.Name != n.Network.Name {
		t.Fatalf("expected network name %q; received %q", n.Network.Name, network.Network.Name)
	}
}
