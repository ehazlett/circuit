package ipam

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/ehazlett/circuit/ds"
)

type IPAM struct {
	ds ds.Backend
}

func NewIPAM(b ds.Backend) (*IPAM, error) {
	return &IPAM{
		ds: b,
	}, nil
}

func (i *IPAM) AllocateIP(subnet *net.IPNet, networkName string) (net.IP, error) {
	logrus.Debugf("allocating IP for subnet: %v", subnet)
	// TODO: allocate IP from pool
	o := subnet.IP.To4()
	// add new source; default is deterministic
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	min := 2   // start at 2
	max := 254 // no greater than 254
	d := min + r.Intn(max-min)

	if len(o) < 3 {
		return nil, fmt.Errorf("error allocating ip: %v", subnet)
	}
	ip := net.IPv4(o[0], o[1], o[2], byte(d))

	// save to ds
	if err := i.ds.SaveIPAddr(ip.String(), networkName); err != nil {
		return ip, err
	}

	return ip, nil
}

func (i *IPAM) ReleaseIP(ip net.IP, networkName string) error {
	// TODO: release IP back to pool
	if err := i.ds.DeleteIPAddr(ip.String(), networkName); err != nil {
		return err
	}

	return nil
}
