package config

import "github.com/containernetworking/cni/libcni"

type PeerInfo struct {
	NetworkName  string `json:"network_name,omitempty"`
	ContainerPid int    `json:"container_pid,omitempty"`
	IP           string `json:"ip,omitempty"`
	IfaceName    string `json:"iface_name,omitempty"`
}

type Network struct {
	Name   string                `json:",omitempty"`
	Config *libcni.NetworkConfig `json:"config"`
}
