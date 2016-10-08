package local

func (c *localController) DeleteNetwork(name string) error {
	// TODO: stop veth pair

	// TODO: remove tc (tc qdisc del dev <veth> root)

	// TODO: remove veth pair (ip link del veth0)

	// TODO: "release" IPs back to pool

	return nil
}
