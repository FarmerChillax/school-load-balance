package registry

type Registration struct {
	ServiceName ServiceName
	ServiceURL  string
}

type ServiceName string

const (
	LogService = ServiceName("LogService")
)

// type patchEntry struct {
// 	Name ServiceName
// 	URL  string
// }

// type patch struct {
// 	Added   []patchEntry
// 	Removed []patchEntry
// }
