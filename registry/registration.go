package registry

type Registration struct {
	ServiceName      ServiceName   //服务名
	ServiceURL       string        //服务地址
	RequiredServices []ServiceName //依赖服务名slice
	ServiceUpdateURl string        //回调地址, 服务中心调用该地址将可用的依赖服务   patch 对象返回给依赖方
}

type ServiceName string

const (
	LogService     = ServiceName("LogService")
	GradingService = ServiceName("GradingService")
)

type patchEntry struct {
	Name ServiceName
	URL  string
}

type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
