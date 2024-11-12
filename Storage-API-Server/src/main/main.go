package main

import (
	"fmt"
	"log"
	"net/http"

	handler "opencsd-storage-api-server/src/handler"
	storagestruct "opencsd-storage-api-server/src/struct"

	"github.com/influxdata/influxdb/client/v2"
)

func main() {
	//influx Connection
	var err error
	storagestruct.INFLUX_CLIENT, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:" + storagestruct.INFLUX_PORT,
		Username: storagestruct.INFLUX_USERNAME,
		Password: storagestruct.INFLUX_PASSWORD,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer storagestruct.INFLUX_CLIENT.Close()

	fmt.Println("[OpenCSD Storage API Server] Connected to Influx Database")
	fmt.Println("[OpenCSD Storage API Server] Running..")
	fmt.Println("[OpenCSD Storage API Server] run on 0.0.0.0:", storagestruct.STORAGE_API_SERVER_PORT)

	storagestruct.NodeStorageInfo_.InitNodeStorageInfo()

	//0. VolumeAllocation with Gluesys
	http.HandleFunc("/volume/allocate", handler.StorageVolumeAllocate)
	http.HandleFunc("/volume/deallocate", handler.StorageVolumeDeallocate)

	http.HandleFunc("/node/info/storage", handler.NodeInfoStorage)
	http.HandleFunc("/node/info/volume", handler.StorageVolumeInfo)

	http.HandleFunc("/node/metric/all", handler.NodeMetricAll)
	http.HandleFunc("/node/metric/cpu", handler.NodeMetricCpu)
	http.HandleFunc("/node/metric/power", handler.NodeMetricPower)
	http.HandleFunc("/node/metric/memory", handler.NodeMetricMemory)
	http.HandleFunc("/node/metric/network", handler.NodeMetricNetwork)
	http.HandleFunc("/node/metric/storage", handler.NodeMetricStorage)

	http.HandleFunc("/storage/metric/all", handler.StorageMetricAll)
	http.HandleFunc("/storage/metric/cpu", handler.StorageMetricCpu)
	http.HandleFunc("/storage/metric/power", handler.StorageMetricPower)
	http.HandleFunc("/storage/metric/memory", handler.StorageMetricMemory)
	http.HandleFunc("/storage/metric/network", handler.StorageMetricNetwork)
	http.HandleFunc("/storage/metric/disk", handler.StorageMetricDisk)

	http.ListenAndServe(":"+storagestruct.STORAGE_API_SERVER_PORT, nil)
}
