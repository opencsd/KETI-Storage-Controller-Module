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
	SSD     = 0
	CSD     = 1
	CROSS   = 2
	UNKNOWN = 3
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
	CsdId  string `json:"csd_id"`
	Status string `json:"status"`
}

type NodeStorageInfo struct {
	NodeName string     `json:"node_name"`
	CsdList  []CsdEntry `json:"csd_list"`
	SsdList  []string   `json:"ssd_list"`
	NodeType string     `json:"node_type"`
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
	Time               string  `json:"timestamp"`
	NodeName           string  `json:"name"`
	CpuTotal           float64 `json:"cpuTotal"`
	CpuUsed            float64 `json:"cpuUsed"`
	CpuUtilization     float64 `json:"cpuUtilization"`
	MemoryTotal        float64 `json:"memoryTotal"`
	MemoryUsed         float64 `json:"memoryUsed"`
	MemoryUtilization  float64 `json:"memoryUtilization"`
	StorageTotal       float64 `json:"storageTotal"`
	StorageUsed        float64 `json:"storageUsed"`
	StorageUtilization float64 `json:"storageUtilization"`
	NetworkRxData      float64 `json:"networkRxData"`
	NetworkTxData      float64 `json:"networkTxData"`
	NetworkBandwidth   float64 `json:"networkBandwidth"`
	PowerUsed          float64 `json:"powerUsed"`
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
	Time               string  `json:"timestamp"`
	Name               string  `json:"name"`
	StorageTotal       float64 `json:"storageTotal"`
	StorageUsed        float64 `json:"storageUsed"`
	StorageUtilization float64 `json:"storageUtilization"`
}

type SsdMetric struct {
	Time               string  `json:"timestamp"`
	Id                 string  `json:"id"`
	Name               string  `json:"name"`
	StorageTotal       float64 `json:"storageTotal"`
	StorageUsed        float64 `json:"storageUsed"`
	StorageUtilization float64 `json:"storageUtilization"`
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
	StorageTotal         float64 `json:"storageTotal"`
	StorageUsed          float64 `json:"storageUsed"`
	StorageUtilization   float64 `json:"storageUtilization"`
	NetworkRxData        float64 `json:"networkRxData"`
	NetworkTxData        float64 `json:"networkTxData"`
	NetworkBandwidth     float64 `json:"networkBandwidth"`
	CsdMetricScore       float64 `json:"csdMetricScore"`
	CsdWorkingBlockCount float64 `json:"csdWorkingBlockCount"`
	Status               string  `json:"status"`
}

type StorageMetric struct {
	SsdList []SsdMetric `json:"ssd_list"`
	CsdList []CsdMetric `json:"csd_list"`
}
