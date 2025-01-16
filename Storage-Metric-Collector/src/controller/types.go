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

	client "github.com/influxdata/influxdb/client/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	STORAGE_METRIC_COLLECTOR_PORT_TCP  = "40304"
	STORAGE_METRIC_COLLECTOR_PORT_HTTP = "40307"
)

var (
	INFLUX_CLIENT   client.HTTPClient
	INFLUX_PORT     = os.Getenv("INFLUXDB_PORT")
	INFLUX_USERNAME = os.Getenv("INFLUXDB_USER")
	INFLUX_PASSWORD = os.Getenv("INFLUXDB_PASSWORD")
	INFLUX_DB       = os.Getenv("INFLUXDB_DB")
)

const (
	SSD     = "SSD"
	CSD     = "CSD"
	CROSS   = "CROSS"
	UNKNOWN = "UNKOWN"
)

const (
	READY    = "READY"
	NOTREADY = "NOTREADY"
	BROKEN   = "BROKEN"
	NORMAL   = "NORMAL"
)

type MetricCollector struct {
	NodeName   string
	NodeMetric *NodeMetric           `json:"nodeMetric"`
	CsdMetrics map[string]*CsdMetric `json:"csdMetrics"`
	SsdMetrics map[string]*SsdMetric `json:"ssdMetrics"`
	NodeType   string                `json:"nodeType"`
}

func NewMetricCollector(mode string) *MetricCollector {
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

	var nodeType string

	if mode == "off" {
		nodeType = CSD
	} else {
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
		switch labelValue {
		case "ssd":
			nodeType = SSD
		case "csd":
			nodeType = CSD
		default:
			nodeType = UNKNOWN
		}
	}

	return &MetricCollector{
		NodeName:   hostname,
		NodeMetric: NodeMetric,
		CsdMetrics: make(map[string]*CsdMetric),
		SsdMetrics: make(map[string]*SsdMetric),
		NodeType:   nodeType,
	}
}

func (metricCollector *MetricCollector) InitMetricCollector(mode string) {
	csdCount := 0
	if mode == "off" {
		csdCountStr := os.Getenv("CSD_COUNT")
		csdCount, _ = strconv.Atoi(csdCountStr)
	} else {
		{
			file, err := os.Open("/etc/lspci-result.txt")
			if err != nil {
				fmt.Println("lspci-result read error, ", err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()

				parts := strings.Fields(line)
				if len(parts) != 2 {
					continue
				}

				node, value := parts[0], parts[1]

				if node == metricCollector.NodeName {
					csdCount, _ = strconv.Atoi(value)
				}
			}
		}
	}

	if metricCollector.NodeType == SSD {
		cmd := exec.Command("lsblk", "-o", "NAME,SIZE,MOUNTPOINT")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			fmt.Println("lsblk error: ", err)
			return
		}

		lines := strings.Split(out.String(), "\n")

		id := 1
		for _, line := range lines[1:] {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			name := fields[0]
			size := fields[1]

			if strings.HasPrefix(name, "sd") {
				totalSize := convertSizeToGB(size)
				key := "ssd" + strconv.Itoa(id)
				ssdMetric := &SsdMetric{
					Name:        name,
					Total:       totalSize,
					Used:        0,
					Utilization: 0,
					Status:      NORMAL,
				}
				metricCollector.SsdMetrics[key] = ssdMetric
				id++
			}
		}
	} else if metricCollector.NodeType == CSD {
		cmd := exec.Command("lsblk", "-o", "NAME,MODEL")
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error executing lsblk command:", err)
			return
		}

		scanner := bufio.NewScanner(bytes.NewReader(output))
		ngdDevices := []string{}

		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "NAME") {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) > 1 && strings.HasPrefix(fields[1], "NGD") {
				deviceName := fields[0]
				deviceName = deviceName[:len(deviceName)-2]
				ngdDevices = append(ngdDevices, deviceName)
			}
		}

		for _, deviceName := range ngdDevices {
			csdMetric := NewCsdMetric()
			csdMetric.Name = deviceName
			id := strings.TrimPrefix(deviceName, "nvme")
			key := "csd" + id
			metricCollector.CsdMetrics[key] = csdMetric
		}

		for i := 1; i <= csdCount; i++ {
			id := strconv.Itoa(i)
			key := "csd" + id

			if _, exists := metricCollector.CsdMetrics[key]; !exists {
				csdMetric := NewCsdMetric()
				csdMetric.Name = "nvme" + id
				csdMetric.Status = "BROKEN"
				metricCollector.CsdMetrics[key] = csdMetric
			}
		}

	} else {
		fmt.Println("[error] not supported node type: ", metricCollector.NodeType)
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
	mutex  sync.Mutex `json:"-"`
	Cpu    Cpu        `json:"cpu"`
	Memory Memory     `json:"memory"`
	// Disk    Disk       `json:"disk"`
	Network Network `json:"network"`
	Power   Power   `json:"power"`
}

func NewNodeMetric() *NodeMetric {
	return &NodeMetric{
		Cpu:    NewCpu(),
		Memory: NewMemory(),
		// Disk:    NewDisk(),
		Network: NewNetwork(),
		Power:   NewPower(),
	}
}

type Cpu struct {
	Total       int     `json:"total"`
	Used        float64 `json:"used"`
	Utilization float64 `json:"utilization"`
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
	Total       float64 `json:"total"`
	Used        float64 `json:"used"`
	Utilization float64 `json:"utilization"`
	Free        float64
	Buffers     float64
	Cached      float64
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

type Disk struct {
	Name        string  `json:"name"`
	Total       float64 `json:"total"`
	Used        float64 `json:"used"`
	Utilization float64 `json:"utilization"`
}

func NewDisk() Disk {
	return Disk{
		Name:        "",
		Total:       0,
		Used:        0,
		Utilization: 0,
	}
}

type Network struct {
	RxByte    int64
	TxByte    int64
	RxData    int64 `json:"rxData"`
	TxData    int64 `json:"txData"`
	Bandwidth int64 `json:"bandwidth"`
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
	Used    int64 `json:"used"`
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
		cmd := exec.Command("grep", "-c", "processor", "/host/proc/cpuinfo")
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
		file, err := os.Open("/host/proc/stat")
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
		file, err := os.Open("/host/proc/meminfo")
		if err != nil {
			fmt.Println("cannot open file: ", err)
		} else {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()

				if strings.HasPrefix(line, "MemTotal:") {
					fields := strings.Fields(line)
					if len(fields) >= 2 {
						memTotalKB, err := strconv.ParseFloat(fields[1], 64)
						if err != nil {
							fmt.Println("Error parsing memory value:", err)
						}
						nodeMetric.Memory.Total = memTotalKB / (1024 * 1024)
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
		statisticsFilePath := "/host/sys/class/net/eno1/statistics/"

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

	// {
	// 	cmd := exec.Command("df", "-k", "--total")
	// 	output, err := cmd.Output()
	// 	if err != nil {
	// 		fmt.Println("Error executing command:", err)
	// 		return
	// 	}

	// 	scanner := bufio.NewScanner(bytes.NewReader(output))
	// 	scanner.Scan()

	// 	for scanner.Scan() {
	// 		line := scanner.Text()

	// 		if strings.Contains(line, "total") {
	// 			fields := strings.Fields(line)
	// 			if len(fields) >= 3 {
	// 				nodeMetric.Disk.Total, _ = strconv.ParseFloat(fields[1], 64)
	// 				break
	// 			}
	// 		}
	// 	}
	// }
}

type SsdMetric struct {
	Name        string  `json:"name"`
	Total       float64 `json:"total"`
	Used        float64 `json:"used"`
	Utilization float64 `json:"utilization"`
	Status      string  `json:"status"`
}

func NewSsdMetric() SsdMetric {
	return SsdMetric{
		Name:        "",
		Total:       0,
		Used:        0,
		Utilization: 0,
	}
}

type CsdMetric struct {
	mutex                sync.Mutex `json:"-"`
	IP                   string     `json:"ip"`
	Name                 string     `json:"name"`
	CpuTotal             int        `json:"cpuTotal"`
	CpuUsed              float64    `json:"cpuUsed"`
	CpuUtilization       float64    `json:"cpuUtilization"`
	MemoryTotal          float64    `json:"memoryTotal"`
	MemoryUsed           float64    `json:"memoryUsed"`
	MemoryUtilization    float64    `json:"memoryUtilization"`
	DiskTotal            float64    `json:"diskTotal"`
	DiskUsed             float64    `json:"diskUsed"`
	DiskUtilization      float64    `json:"diskUtilization"`
	NetworkRxData        int64      `json:"networkRxData"`
	NetworkTxData        int64      `json:"networkTxData"`
	NetworkBandwidth     int64      `json:"networkBandwidth"`
	CsdMetricScore       float64    `json:"csdMetricScore"`
	CsdWorkingBlockCount int64      `json:"csdWorkingBlockCount"`
	Status               string     `json:"status"`
}

func NewCsdMetric() *CsdMetric {
	return &CsdMetric{
		IP:                   "",
		Name:                 "",
		CpuTotal:             0,
		CpuUsed:              0,
		CpuUtilization:       0,
		MemoryTotal:          0,
		MemoryUsed:           0,
		MemoryUtilization:    0,
		DiskTotal:            0,
		DiskUsed:             0,
		DiskUtilization:      0,
		NetworkRxData:        0,
		NetworkTxData:        0,
		NetworkBandwidth:     0,
		CsdMetricScore:       0,
		CsdWorkingBlockCount: 0,
		Status:               NOTREADY,
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

func convertSizeToGB(sizeStr string) float64 {
	unit := sizeStr[len(sizeStr)-1:]
	sizeValue, _ := strconv.ParseFloat(sizeStr[:len(sizeStr)-1], 64)
	switch unit {
	case "T":
		return float64(sizeValue * 1024)
	case "G":
		return float64(sizeValue)
	case "M":
		return float64(sizeValue / 1024)
	default:
		return 0
	}
}
