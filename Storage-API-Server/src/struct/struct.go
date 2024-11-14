package storagestruct

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	STORAGE_API_SERVER_PORT            = os.Getenv("STORAGE_API_SERVER_PORT")
	STORAGE_METRIC_COLLECTOR_PORT_HTTP = os.Getenv("STORAGE_METRIC_COLLECTOR_PORT_HTTP")
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
	NodeName string     `json:"nodeName"`
	CsdList  []CsdEntry `json:"csdList"`
	SsdList  []SsdEntry `json:"ssdList"`
	NodeType string     `json:"nodeType"`
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
