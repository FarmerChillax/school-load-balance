package scan

type Host string
type Port int

type Address struct {
	Host   Host
	Port   Port
	status bool
}

type OpenPorts map[Host][]Port
