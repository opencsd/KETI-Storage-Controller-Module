package storagestruct

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	SSD     = "SSD"
	CSD     = "CSD"
	CROSS   = "CROSS"
	UNKNOWN = "UNKNOWN"
)

const (
	READY    = "READY"
	NOTREADY = "NOTREADY"
	BROKEN   = "BROKEN"
	NORMAL   = "NORMAL"
)

var (
	STORAGE_API_SERVER_PORT            = "40306"
	STORAGE_METRIC_COLLECTOR_PORT_HTTP = "40307"
)

var (
	INFLUX_CLIENT   client.HTTPClient
	INFLUX_PORT     = os.Getenv("INFLUXDB_PORT")
	INFLUX_USERNAME = os.Getenv("INFLUXDB_USER")
	INFLUX_PASSWORD = os.Getenv("INFLUXDB_PASSWORD")
	INFLUX_DB       = os.Getenv("INFLUXDB_DB")
)

var NodeStorageInfo_ NodeStorageInfo

type CsdEntry struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SsdEntry struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type NodeStorageInfo struct {
	NodeName   string                `json:"nodeName"`
	CsdList    []CsdEntry            `json:"csdList"`
	SsdList    []SsdEntry            `json:"ssdList"`
	NodeType   string                `json:"nodeType"`
	VolumeInfo map[string]VolumeInfo `json:"volumeInfo"`
}

type VolumeInfo struct {
	VolumeName    string  `json:"volumeName"`
	VolumePath    string  `json:"volumePath"`
	NodeName      string  `json:"nodeName"`
	SizeTotal     float64 `json:"sizeTotal"`
	SizeUsed      float64 `json:"sizeUsed"`
	SizeAvailable float64 `json:"sizeAvailable"`
	Utilization   float64 `json:"instanceUtilization"`
	StorageType   string  `json:"storageType"`
	VolumeType    string  `json:"volumeType"`
}

type NodeMetric struct {
	Time              string  `json:"timestamp"`
	NodeName          string  `json:"name"`
	CpuTotal          float64 `json:"cpuTotal"`
	CpuUsed           float64 `json:"cpuUsed"`
	CpuUtilization    float64 `json:"cpuUtilization"`
	MemoryTotal       float64 `json:"memoryTotal"`
	MemoryUsed        float64 `json:"memoryUsed"`
	MemoryUtilization float64 `json:"memoryUtilization"`
	DiskTotal         float64 `json:"diskTotal"`
	DiskUsed          float64 `json:"diskUsed"`
	DiskUtilization   float64 `json:"diskUtilization"`
	NetworkRxData     float64 `json:"networkRxData"`
	NetworkTxData     float64 `json:"networkTxData"`
	NetworkBandwidth  float64 `json:"networkBandwidth"`
	PowerUsed         float64 `json:"powerUsed"`
}

type CpuMetric struct {
	Time           string  `json:"timestamp"`
	Name           string  `json:"name"`
	CpuTotal       float64 `json:"cpuTotal"`
	CpuUsed        float64 `json:"cpuUsed"`
	CpuUtilization float64 `json:"cpuUtilization"`
}

type PowerMetric struct {
	Time      string  `json:"timestamp"`
	Name      string  `json:"name"`
	PowerUsed float64 `json:"powerUsed"`
}

type MemoryMetric struct {
	Time              string  `json:"timestamp"`
	Name              string  `json:"name"`
	MemoryTotal       float64 `json:"memoryTotal"`
	MemoryUsed        float64 `json:"memoryUsed"`
	MemoryUtilization float64 `json:"memoryUtilization"`
}

type NetworkMetric struct {
	Time             string  `json:"timestamp"`
	Name             string  `json:"name"`
	NetworkRxData    float64 `json:"networkRxData"`
	NetworkTxData    float64 `json:"networkTxData"`
	NetworkBandwidth float64 `json:"networkBandwidth"`
}

type DiskMetric struct {
	Time            string  `json:"timestamp"`
	Name            string  `json:"name"`
	DiskTotal       float64 `json:"diskTotal"`
	DiskUsed        float64 `json:"diskUsed"`
	DiskUtilization float64 `json:"diskUtilization"`
}

type SsdMetric struct {
	Time            string  `json:"timestamp"`
	Id              string  `json:"id"`
	Name            string  `json:"name"`
	DiskTotal       float64 `json:"diskTotal"`
	DiskUsed        float64 `json:"diskUsed"`
	DiskUtilization float64 `json:"diskUtilization"`
	Status          string  `json:"status"`
}

type CsdMetric struct {
	Time                 string  `json:"timestamp"`
	Id                   string  `json:"id"`
	Name                 string  `json:"name"`
	Ip                   string  `json:"ip"`
	CpuTotal             float64 `json:"cpuTotal"`
	CpuUsed              float64 `json:"cpuUsed"`
	CpuUtilization       float64 `json:"cpuUtilization"`
	MemoryTotal          float64 `json:"memoryTotal"`
	MemoryUsed           float64 `json:"memoryUsed"`
	MemoryUtilization    float64 `json:"memoryUtilization"`
	DiskTotal            float64 `json:"diskTotal"`
	DiskUsed             float64 `json:"diskUsed"`
	DiskUtilization      float64 `json:"diskUtilization"`
	NetworkRxData        float64 `json:"networkRxData"`
	NetworkTxData        float64 `json:"networkTxData"`
	NetworkBandwidth     float64 `json:"networkBandwidth"`
	CsdMetricScore       float64 `json:"csdMetricScore"`
	CsdWorkingBlockCount float64 `json:"csdWorkingBlockCount"`
	Status               string  `json:"status"`
}

type StorageMetricMessage struct {
	SsdList map[string][]SsdMetric `json:"ssdList"`
	CsdList map[string][]CsdMetric `json:"csdList"`
}

func NewStorageMetricMessage() StorageMetricMessage {
	return StorageMetricMessage{
		SsdList: make(map[string][]SsdMetric),
		CsdList: make(map[string][]CsdMetric),
	}
}

type CsdMetricMin struct {
	Time            string  `json:"timestamp"`
	Id              string  `json:"id"`
	Name            string  `json:"name"`
	Ip              string  `json:"ip"`
	DiskTotal       float64 `json:"diskTotal"`
	DiskUsed        float64 `json:"diskUsed"`
	DiskUtilization float64 `json:"diskUtilization"`
	CsdMetricScore  float64 `json:"csdMetricScore"`
	Status          string  `json:"status"`
}

type StorageInfoMessage struct {
	SsdList map[string][]SsdMetric    `json:"ssdList"`
	CsdList map[string][]CsdMetricMin `json:"csdList"`
}

func NewStorageInfoMessage() StorageInfoMessage {
	return StorageInfoMessage{
		SsdList: make(map[string][]SsdMetric),
		CsdList: make(map[string][]CsdMetricMin),
	}
}

func (nodeStorageInfo *NodeStorageInfo) InitNodeStorageInfo() {
	serverAddress := "http://localhost:" + STORAGE_METRIC_COLLECTOR_PORT_HTTP + "/node/info/storage"

	for {
		resp, err := http.Get(serverAddress)
		if err != nil {
			fmt.Println("Error sending request:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var response NodeStorageInfo

		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatal(err)
		}

		*nodeStorageInfo = response

		fmt.Printf("Decoded response: %+v\n", response)
		break
	}

	nodeStorageInfo.VolumeInfo = getNodeVolumeInfo()
}

func getNodeVolumeInfo() map[string]VolumeInfo {
	volumeInfo := map[string]VolumeInfo{}

	lvmVolumes, err := getLvmVolumes()
	if err != nil {
		fmt.Println("Error collecting LVM volumes:", err)
	} else {
		fmt.Println("LVM Volumes:", lvmVolumes)
	}

	glusterVolumes, err := getGlusterVolumes()
	if err != nil {
		fmt.Println("Error collecting GlusterFS volumes:", err)
	} else {
		fmt.Println("GlusterFS Volumes:", glusterVolumes)
	}

	for k, v := range lvmVolumes {
		volumeInfo[k] = v
	}

	for k, v := range glusterVolumes {
		volumeInfo[k] = v
	}

	return volumeInfo
}

func getLvmVolumes() (map[string]VolumeInfo, error) {
	volumeInfo := map[string]VolumeInfo{}

	cmd := exec.Command("lvs", "--units", "g", "--noheadings", "-o", "lv_name,vg_name,lv_size")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get LVM volumes: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		volumeName := fields[0]
		volumeGroup := fields[1]

		// Get usage using `df`
		filesystem := fmt.Sprintf("/dev/mapper/%s-%s", volumeGroup, volumeName)
		total, used, available, utilization, mountPath, err := df(filesystem)
		if err != nil {
			fmt.Printf("Failed to get usage for %s: %v\n", filesystem, err)
			continue
		}

		volumeInfo[volumeName] = VolumeInfo{
			VolumeName:    volumeName,
			VolumePath:    mountPath,
			NodeName:      NodeStorageInfo_.NodeName,
			SizeTotal:     total,
			SizeUsed:      used,
			SizeAvailable: available,
			Utilization:   utilization,
			VolumeType:    "lvm",
		}
	}

	return volumeInfo, nil
}

func getGlusterVolumes() (map[string]VolumeInfo, error) {
	volumeInfo := map[string]VolumeInfo{}

	// Command to list GlusterFS volumes
	cmd := exec.Command("gluster", "volume", "info")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get GlusterFS volumes: %v", err)
	}

	volumeName := ""
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Volume Name:") {
			volumeName = strings.TrimSpace(strings.TrimPrefix(line, "Volume Name:"))
		}

		if strings.HasPrefix(line, "Status:") && volumeName != "" {
			// Assume volume is mounted at /mnt/<volumeName>
			filesystem := fmt.Sprintf("/mnt/%s", volumeName)
			total, used, available, utilization, mountPath, err := df(filesystem)
			fmt.Println(filesystem)
			if err != nil {
				fmt.Printf("Failed to get usage for %s: %v\n", filesystem, err)
				continue
			}

			volumeInfo[volumeName] = VolumeInfo{
				VolumeName:    volumeName,
				VolumePath:    mountPath,
				SizeTotal:     total,
				SizeUsed:      used,
				SizeAvailable: available,
				Utilization:   utilization,
				VolumeType:    "gluster",
			}

			volumeName = ""
		}
	}

	return volumeInfo, nil
}

func df(filesystem string) (total float64, used float64, available float64, utilization float64, path string, err error) {
	cmd := exec.Command("df", "-BG", filesystem)
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, 0, 0, "", fmt.Errorf("failed to get volume usage: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0, 0, 0, 0, "", fmt.Errorf("invalid df output")
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 5 {
		return 0, 0, 0, 0, "", fmt.Errorf("unexpected df output format")
	}

	total = parseSize(fields[1])
	used = parseSize(fields[2])
	available = parseSize(fields[3])
	utilization = (used / total) * 100
	roundedUtilization := math.Round(utilization*100) / 100
	path = fields[5]

	return total, used, available, roundedUtilization, path, nil
}

func parseSize(target string) float64 {
	resultStr := strings.TrimSuffix(target, "G")
	result, _ := strconv.ParseFloat(resultStr, 64)

	return math.Round(result*100) / 100
}
