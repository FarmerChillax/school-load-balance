package network

import "time"

type Host string
type Port int
type Timeout time.Duration

type Segment struct {
	Template           string
	SegmentType        rune
	StartPort, EndPort Port
}

type SegmentResult struct {
	Host  Host
	Ports []Port
}

type Address struct {
	Host    Host
	Port    Port
	Status  bool
	Timeout Timeout
	ssl     bool
}
