package config

type PeerType string

const (
	ContainerPeer PeerType = "container"
	HostPeer      PeerType = "host"
)

type IPPeer struct {
	IP   string   `json:",omitempty"`
	Pid  int      `json:",omitempty"`
	Type PeerType `json:",omitempty"`
}

type Network struct {
	Name   string             `json:",omitempty"`
	Subnet string             `json:",omitempty"`
	IPs    map[string]*IPPeer `json:",omitempty"`
}
