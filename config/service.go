package config

type Protocol string

const (
	ProtocolTCP Protocol = "tcp"
	ProtocolUDP Protocol = "udp"
)

type Scheduler string

const (
	SchedulerRR    Scheduler = "rr"
	SchedulerWRR   Scheduler = "wrr"
	SchedulerLC    Scheduler = "lc"
	SchedulerWLC   Scheduler = "wlc"
	SchedulerLBLC  Scheduler = "lblc"
	SchedulerLBLCR Scheduler = "lblcr"
	SchedulerDH    Scheduler = "dh"
	SchedulerSH    Scheduler = "sh"
	SchedulerSED   Scheduler = "sed"
	SchedulerNQ    Scheduler = "nq"
)

type Service struct {
	Name      string    `json:",omitempty"`
	Addr      string    `json:",omitempty"`
	Protocol  Protocol  `json:",omitempty"`
	Scheduler Scheduler `json:",omitempty"`
	Targets   []string  `json:",omitempty"`
}
