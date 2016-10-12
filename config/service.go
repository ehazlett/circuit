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
	Addr      string
	Protocol  Protocol
	Scheduler Scheduler
}
