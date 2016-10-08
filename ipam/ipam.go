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
	backend ds.Backend
}

func NewIPAM(b ds.Backend) (*IPAM, error) {
	return &IPAM{
		backend: b,
	}, nil
}

func (i *IPAM) AllocateIP(subnet *net.IPNet) (net.IP, error) {
	logrus.Debugf("allocating IP for subnet: %v", subnet)
	// TODO: allocate IP from pool
	o := subnet.IP.To4()
	// add new source; default is deterministic
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	d := r.Intn(254)
	if d == 0 {
		d++
	}

	if len(o) < 3 {
		return nil, fmt.Errorf("error allocating ip: %v", subnet)
	}
	ip := net.IPv4(o[0], o[1], o[2], byte(d))
	return ip, nil
}

func (i *IPAM) ReleaseIP(ip net.IP) error {
	// TODO: release IP back to pool
	return nil
}
