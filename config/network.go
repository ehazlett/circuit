package config

type IPPeer struct {
	IP  string `json:",omitempty"`
	Pid int    `json:",omitempty"`
}

type Network struct {
	Name   string             `json:",omitempty"`
	Subnet string             `json:",omitempty"`
	IPs    map[string]*IPPeer `json:",omitempty"`
}
