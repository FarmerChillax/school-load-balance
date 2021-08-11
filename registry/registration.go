package registry

// 服务注册结构体
type Registration struct {
	ServiceName ServiceName
	ServiceURL  string
}

type ServiceName string

const (
	LogService   = ServiceName("LogService")
	RedisService = ServiceName("RedisService")
)

// type patchEntry struct {
// 	Name ServiceName
// 	URL  string
// }

// type patch struct {
// 	Added   []patchEntry
// 	Removed []patchEntry
// }
