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
	http.HandleFunc("/directory/create", handler.StorageDirectoryCreate) // ?path=
	http.HandleFunc("/directory/delete", handler.StorageDirectoryDelete)

	http.HandleFunc("/node/info/storage-list", handler.NodeInfoStorageList)
	http.HandleFunc("/node/info/storage", handler.StorageInfo) // ?storage=&count=
	http.HandleFunc("/node/info/volume", handler.StorageVolumeInfo)

	http.HandleFunc("/node/metric/all", handler.NodeMetricAll)         // ?count=
	http.HandleFunc("/node/metric/cpu", handler.NodeMetricCpu)         // ?count=
	http.HandleFunc("/node/metric/power", handler.NodeMetricPower)     // ?count=
	http.HandleFunc("/node/metric/memory", handler.NodeMetricMemory)   // ?count=
	http.HandleFunc("/node/metric/network", handler.NodeMetricNetwork) // ?count=
	http.HandleFunc("/node/metric/disk", handler.NodeMetricDisk)       // ?count=

	http.HandleFunc("/storage/info", handler.StorageInfo)                    // ?storage=&count=
	http.HandleFunc("/storage/metric/all", handler.StorageMetricAll)         // ?storage=&count=
	http.HandleFunc("/storage/metric/cpu", handler.StorageMetricCpu)         // ?storage=&count=
	http.HandleFunc("/storage/metric/power", handler.StorageMetricPower)     // ?storage=&count=
	http.HandleFunc("/storage/metric/memory", handler.StorageMetricMemory)   // ?storage=&count=
	http.HandleFunc("/storage/metric/network", handler.StorageMetricNetwork) // ?storage=&count=
	http.HandleFunc("/storage/metric/disk", handler.StorageMetricDisk)       // ?storage=&count=

	http.ListenAndServe(":"+storagestruct.STORAGE_API_SERVER_PORT, nil)
}
