package storagestruct

type ClusterNodeInfo struct {
	ClusterName  string   `json:"clustername"`
	NodeList     []NodeInfo
}

type NodeInfo struct {
	NodeName     string `json:"nodename"`
	Status 		 string `json:"status"`
}

type NodeStorageInfo struct {
	NodeName     string `json:"nodename"`
	StorageList  []StorageInfo
}

type StorageInfo struct {
	StorageName     string `json:"storagename"`
	Status   		string `json:"status"`
}

type StorageInfoDetailV1 struct {
	StorageName     string 	`json:"storagename"`
	StorageType		string 	`json:"storagetype"`
	Capacity		string 	`json:"capacity"`
	Status   		string 	`json:"status"`
	StorageScore	float64	`json:"storagescore"`
	StorageGrade	string	`json:"storagegrade"`
	ClusterBelongTo	string	`json:"clusterbelongto"`
	NodeBelongTo	string	`json:"nodebelongto"`
}

type StorageInfoDetailV2 struct {
	StorageName     string 	`json:"storagename"`
	StorageType		string 	`json:"storagetype"`
	Capacity		string 	`json:"capacity"`
	Status   		string 	`json:"status"`
	CSDType			string 	`json:"csdtype"`
	IP				string 	`json:"ip"`
	Port			string	`json:"port"`
	ClusterBelongTo	string	`json:"clusterbelongto"`
	NodeBelongTo	string	`json:"nodebelongto"`
	StorageScore	float64	`json:"storagescore"`
	StorageGrade	string	`json:"storagegrade"`
	WorkingBlocks	float64	`json:"workingblocks"`
}

type NodeDiskInfo struct {
	Capacity     		float64 		`json:"capacity"`
	AvailableCapacity   float64 		`json:"availablecapacity"`
	Utilization			float64   		`json:"utilization"`
}

type MetricInfo struct {
	MetricName     		string 			`json:"metricname"`
	MetricValueList   	[]MetricValue 	`json:"metricvaluelist"`
}

type MetricValue struct {
	Time		string		`json:"time"`
	Capacity	float64		`json:"capacity"`
	Usage		float64		`json:"usage"`
	Utilization	float64		`json:"utilization"`
}

type NetMetricValue struct {
	Time		string		`json:"time"`
	Bandwidth	float64		`json:"bandwidth"`
	RXByte		float64		`json:"rxbyte"`
	TXByte		float64		`json:"txbyte"`
}

type NodeMetricInfo struct {
	Time			string		`json:"time"`
	CPUCapacity		float64		`json:"cpucapacity"`
	CPUUsage		float64		`json:"cpuusage"`
	CPUUtilization	float64		`json:"cputotal"`
	MEMCapacity		float64		`json:"memcapacity"`
	MEMUsage		float64		`json:"memusage"`
	MEMUtilization	float64		`json:"memtotal"`
	DiskCapacity	float64		`json:"diskcapacity"`
	DiskUsage		float64		`json:"diskusage"`
	DiskUtilization	float64		`json:"disktotal"`
}

type CSDMetricInfo struct {
	Time			string		`json:"time"`
	CPUCapacity		float64		`json:"cpucapacity"`
	CPUUsage		float64		`json:"cpuusage"`
	CPUUtilization	float64		`json:"cputotal"`
	MEMCapacity		float64		`json:"memcapacity"`
	MEMUsage		float64		`json:"memusage"`
	MEMUtilization	float64		`json:"memtotal"`
	DiskCapacity	float64		`json:"diskcapacity"`
	DiskUsage		float64		`json:"diskusage"`
	DiskUtilization	float64		`json:"disktotal"`
}