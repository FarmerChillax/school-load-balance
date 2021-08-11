package network

import "time"

type Timeout time.Duration

type Segment struct {
	Template           string
	SegmentType        rune
	StartPort, EndPort int
}

type SegmentResult struct {
	Host  string
	Ports []int
}

type Address struct {
	Host    string
	Port    int
	Status  bool
	Timeout Timeout
	ssl     bool
}
