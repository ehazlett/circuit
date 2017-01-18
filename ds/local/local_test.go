package local

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/ehazlett/circuit/config"
)

func getLocalDS() (*localDS, error) {
	d, err := ioutil.TempDir("", "circuit-test")
	if err != nil {
		return nil, err
	}

	return &localDS{
		statePath: d,
	}, nil
}

func TestLocalSaveNetwork(t *testing.T) {
	l, err := getLocalDS()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(l.statePath)

	n := &config.Network{
		Name:   "testing",
		Subnet: "10.254.0.0",
	}

	if err := l.SaveNetwork(n); err != nil {
		t.Fatal(err)
	}

	network, err := l.GetNetwork(n.Name)
	if err != nil {
		t.Fatal(err)
	}

	if network.Name != n.Name {
		t.Fatalf("expected network name %q; received %q", n.Name, network.Name)
	}
}

func TestLocalSaveIPAddr(t *testing.T) {
	l, err := getLocalDS()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(l.statePath)

	n := &config.Network{
		Name:           "testing",
		Subnet:         "10.254.0.0",
		BandwidthBytes: 2048,
	}

	if err := l.SaveNetwork(n); err != nil {
		t.Fatal(err)
	}

	network, err := l.GetNetwork(n.Name)
	if err != nil {
		t.Fatal(err)
	}

	if network.Name != n.Name {
		t.Fatalf("expected network name %q; received %q", n.Name, network.Name)
	}

	testIP := "10.254.0.1"

	if err := l.SaveIPAddr(testIP, network.Name); err != nil {
		t.Fatal(err)
	}

	ips, err := l.GetNetworkIPs(network.Name)
	if err != nil {
		t.Fatal(err)
	}

	if len(ips) != 1 {
		t.Fatalf("expected 1 network; received %d", len(ips))
	}

	if ips[0].String() != testIP {
		t.Fatalf("expected ip %q; received %q", testIP, ips[0].String())
	}
}
