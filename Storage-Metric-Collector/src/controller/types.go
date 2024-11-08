package controller

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	STORAGE_METRIC_COLLECTOR_PORT string
	STORAGE_METRIC_DB_PORT        string
)

const (
	SSD     = 0
	CSD     = 1
	UNKNOWN = 2
)

type MetricCollector struct {
	NodeName   string
	NodeMetric *NodeMetric
	CsdMetrics map[string]*CsdMetric
	SsdMetrics map[string]*Storage
	NodeType   int
}

func NewMetricCollector() *MetricCollector {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("cannot get hostname:", err)
		hostname = ""
	}

	NodeMetric := NewNodeMetric()
	NodeMetric.InitNodeMetric()

	nodeName := os.Getenv("NODE_NAME")
	if nodeName == "" {
		fmt.Println("NODE_NAME environment variable is not set")
	}

	var labelValue string
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("InClusterConfig error :", err)
	} else {
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println("NewForConfig error :", err)
		} else {
			node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
			if err != nil {
				fmt.Println("Get error :", err)
			} else {
				labelValue = node.Labels["type"]
			}
		}
	}

	var nodeType int
	switch labelValue {
	case "ssd":
		nodeType = SSD
	case "csd":
		nodeType = CSD
	default:
		nodeType = UNKNOWN
	}

	return &MetricCollector{
		NodeName:   hostname,
		NodeMetric: NodeMetric,
		CsdMetrics: make(map[string]*CsdMetric),
		SsdMetrics: make(map[string]*Storage),
		NodeType:   nodeType,
	}
}

func (metricCollector *MetricCollector) InitMetricCollector() {
	cmd := exec.Command("lsblk", "-o", "NAME,SIZE,MOUNTPOINT")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return
	}

	lines := strings.Split(out.String(), "\n")

	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		name := fields[0]
		size := fields[1]

		if strings.HasPrefix(name, "sd") {
			totalSize := convertSizeToMB(size)
			disk := &Storage{
				Total:       totalSize,
				Used:        0,
				Utilization: 0,
			}
			metricCollector.SsdMetrics[name] = disk
		}
	}
}

type Config struct {
	Config    *rest.Config
	Clientset *kubernetes.Clientset
	ClusterIP string
}

func NewConfig() *Config {
	hostConfig, _ := rest.InClusterConfig()
	hostKubeClient := kubernetes.NewForConfigOrDie(hostConfig)

	return &Config{
		Config:    hostConfig,
		Clientset: hostKubeClient,
		ClusterIP: hostConfig.Host,
	}
}

type NodeMetric struct {
	mutex   sync.Mutex
	Cpu     Cpu
	Memory  Memory
	Storage Storage
	Network Network
	Power   Power
}

func NewNodeMetric() *NodeMetric {
	return &NodeMetric{
		Cpu:     NewCpu(),
		Memory:  NewMemory(),
		Storage: NewStorage(),
		Network: NewNetwork(),
		Power:   NewPower(),
	}
}

type Cpu struct {
	Total       int
	Used        float64
	Utilization float64
	StJiffies   StJiffies
}

func NewCpu() Cpu {
	return Cpu{
		Total:       0,
		Used:        0,
		Utilization: 0,
		StJiffies:   NewStJiffies(),
	}
}

type StJiffies struct {
	User   int
	Nice   int
	System int
	Idle   int
}

func NewStJiffies() StJiffies {
	return StJiffies{
		User:   0,
		Nice:   0,
		System: 0,
		Idle:   0,
	}
}

type Memory struct {
	Total       int64
	Used        int64
	Utilization float64
	Free        int64
	Buffers     int64
	Cached      int64
}

func NewMemory() Memory {
	return Memory{
		Total:       0,
		Used:        0,
		Utilization: 0,
		Free:        0,
		Buffers:     0,
		Cached:      0,
	}
}

type Storage struct {
	Total       int64
	Used        int64
	Utilization float64
}

func NewStorage() Storage {
	return Storage{
		Total:       0,
		Used:        0,
		Utilization: 0,
	}
}

type Network struct {
	RxByte    int64
	TxByte    int64
	RxData    int64
	TxData    int64
	Bandwidth int64
}

func NewNetwork() Network {
	return Network{
		RxByte:    0,
		TxByte:    0,
		RxData:    0,
		TxData:    0,
		Bandwidth: 0,
	}
}

type Power struct {
	Energy1 int64
	Energy2 int64
	Used    int64
}

func NewPower() Power {
	return Power{
		Energy1: 0,
		Energy2: 0,
		Used:    0,
	}
}

func (nodeMetric *NodeMetric) InitNodeMetric() {
	nodeMetric.mutex.Lock()
	defer nodeMetric.mutex.Unlock()

	{
		cmd := exec.Command("grep", "-c", "processor", "/proc/cpuinfo")
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error: Command execution failed:", err)
		} else {
			coreCountStr := strings.TrimSpace(string(output))
			coreCount, err := strconv.Atoi(coreCountStr)
			if err != nil {
				fmt.Println("Error: Failed to convert core count to integer:", err)
			} else {
				nodeMetric.Cpu.Total = coreCount
			}
		}
	}

	{
		file, err := os.Open("/proc/stat")
		if err != nil {
			fmt.Println("cannot open file: ", err)
		} else {
			var cpuID string
			_, err = fmt.Fscanf(file, "%5s %d %d %d %d", &cpuID, &nodeMetric.Cpu.StJiffies.User, &nodeMetric.Cpu.StJiffies.Nice, &nodeMetric.Cpu.StJiffies.System, &nodeMetric.Cpu.StJiffies.Idle)
			if err != nil {
				fmt.Println("Error reading data from file:", err)
			}
		}
		file.Close()
	}

	{
		file, err := os.Open("/proc/meminfo")
		if err != nil {
			fmt.Println("cannot open file: ", err)
		} else {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()

				if strings.HasPrefix(line, "MemTotal:") {
					fields := strings.Fields(line)
					if len(fields) >= 2 {
						nodeMetric.Memory.Total, err = strconv.ParseInt(fields[1], 10, 64)
						if err != nil {
							fmt.Println("Error parsing memory value:", err)
						}
					}
					break
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading file:", err)
			}
		}
		file.Close()
	}

	{
		statisticsFilePath := "/sys/class/net/eno1/statistics/"

		rxBytesFieldName := statisticsFilePath + "rx_bytes"
		txBytesFieldName := statisticsFilePath + "tx_bytes"

		rxBytes, err := readStatisticsField(rxBytesFieldName)
		if err != nil {
			fmt.Println(err)
			return
		}

		txBytes, err := readStatisticsField(txBytesFieldName)
		if err != nil {
			fmt.Println(err)
			return
		}

		nodeMetric.Network.RxByte, _ = strconv.ParseInt(rxBytes, 10, 64)
		nodeMetric.Network.TxByte, _ = strconv.ParseInt(txBytes, 10, 64)
	}

	{
		cmd := exec.Command("df", "-k", "--total")
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error executing command:", err)
			return
		}

		scanner := bufio.NewScanner(bytes.NewReader(output))
		scanner.Scan()

		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "total") {
				fields := strings.Fields(line)
				if len(fields) >= 3 {
					nodeMetric.Storage.Total, _ = strconv.ParseInt(fields[1], 10, 64)
					break
				}
			}
		}
	}
}

type CsdMetric struct {
	mutex                sync.Mutex
	IP                   string  `json:"ip"`
	CpuTotal             int     `json:"cpuTotal"`
	CpuUsed              float64 `json:"cpuUsed"`
	CpuUtilization       float64 `json:"cpuUtilization"`
	MemoryTotal          int64   `json:"memoryTotal"`
	MemoryUsed           int64   `json:"memoryUsed"`
	MemoryUtilization    float64 `json:"memoryUtilization"`
	StorageTotal         int64   `json:"storageTotal"`
	StorageUsed          int64   `json:"storageUsed"`
	StorageUtilization   float64 `json:"storageUtilization"`
	NetworkRxData        int64   `json:"networkRxData"`
	NetworkTxData        int64   `json:"networkTxData"`
	NetworkBandwidth     int64   `json:"networkBandwidth"`
	CsdMetricScore       float64 `json:"csdMetricScore"`
	CsdWorkingBlockCount int64   `json:"csdWorkingBlockCount"`
	// Power		  int     `json:"powerUsage"`
}

func NewCsdMetric() *CsdMetric {
	return &CsdMetric{
		IP:                   "",
		CpuTotal:             0,
		CpuUsed:              0,
		CpuUtilization:       0,
		MemoryTotal:          0,
		MemoryUsed:           0,
		MemoryUtilization:    0,
		StorageTotal:         0,
		StorageUsed:          0,
		StorageUtilization:   0,
		NetworkRxData:        0,
		NetworkTxData:        0,
		NetworkBandwidth:     0,
		CsdMetricScore:       0,
		CsdWorkingBlockCount: 0,
	}
}

func readStatisticsField(fieldName string) (string, error) {
	data, err := os.ReadFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %v", fieldName, err)
	}

	value := strings.TrimSpace(string(data))
	return value, nil
}

func extractCSDId(addr string) string {
	parts := strings.Split(addr, ".")
	if len(parts) > 0 {
		id := parts[2]
		return id
	}
	return ""
}

func convertSizeToMB(sizeStr string) int64 {
	unit := sizeStr[len(sizeStr)-1:]
	sizeValue, _ := strconv.ParseFloat(sizeStr[:len(sizeStr)-1], 64)
	switch unit {
	case "T":
		return int64(sizeValue * 1024 * 1024)
	case "G":
		return int64(sizeValue * 1024)
	case "M":
		return int64(sizeValue)
	default:
		return 0
	}
}
