package config

import "time"

type QOSConfig struct {
	Interface string
	Delay     time.Duration
	Rate      int
	Ceiling   int
	Buffer    int
	Cbuffer   int
	Priority  int
}
