package discover

import (
	"fmt"
	"math/rand"
	"time"
)

func MakeRandSegment() int {
	rand.Seed(time.Now().Unix())
	i := rand.Intn(255) + 1
	fmt.Println(i)
	return i
}

func checkPort(start, end int) error {
	if start > end {
		return fmt.Errorf("end must big than start")
	}
	if start <= 0 || end > 65535 {
		return fmt.Errorf("scan port out of range")
	}
	return nil
}

func checkProtocol(protocol string) error {
	if len(protocol) == 0 {
		return nil
	}
	if protocol == "http" || protocol == "https" {
		return nil
	}
	return fmt.Errorf("protocol error")
}
