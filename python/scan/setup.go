package scan

type Address struct {
	Host   string
	Port   int
	status bool
}

type Addrs []Address
