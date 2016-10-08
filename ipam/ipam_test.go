package ipam

import (
	"io/ioutil"
	"net"
	"testing"

	"github.com/ehazlett/circuit/ds"
	"github.com/ehazlett/circuit/ds/local"
)

func getLocalDS() (ds.Backend, error) {
	d, err := ioutil.TempDir("", "circuit-test")
	if err != nil {
		return nil, err
	}

	ds, err := local.NewLocalDS(d)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func TestAllocateIP(t *testing.T) {
	d, err := getLocalDS()
	if err != nil {
		t.Fatal(err)
	}

	ipm, err := NewIPAM(d)
	if err != nil {
		t.Fatal(err)
	}

	subnetIP := net.ParseIP("10.254.10.0")
	subnetMask := net.IPv4Mask(255, 255, 0, 0)
	s := &net.IPNet{
		IP:   subnetIP,
		Mask: subnetMask,
	}
	i, err := ipm.AllocateIP(s)
	if err != nil {
		t.Fatal(err)
	}

	if !s.Contains(i) {
		t.Fatalf("expected ip %s to be in subnet %s", i.String(), s.String())
	}
}
