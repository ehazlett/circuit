package local

func (c *localController) SetBandwidthLimit(device string, limitBytes, maxBytes int) error {
	// TODO: tc qdisc add dev eth0 handle 1: root htb default 11

	// TODO: tc class add dev eth0 parent 1: classid 1:1 htb rate 1kbps

	// TODO: tc class add dev eth0 parent 1:1 classid 1:11 htb rate 1kbps

	return nil
}
