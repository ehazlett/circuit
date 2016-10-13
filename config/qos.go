package config

import "time"

type QOSConfig struct {
	Interface string        `json:",omitempty"`
	Delay     time.Duration `json:",omitempty"`
	Rate      int           `json:",omitempty"`
	Ceiling   int           `json:",omitempty"`
	Buffer    int           `json:",omitempty"`
	Cbuffer   int           `json:",omitempty"`
	Priority  int           `json:",omitempty"`
}
