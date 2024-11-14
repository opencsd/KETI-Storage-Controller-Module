package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	storagestruct "opencsd-storage-api-server/src/struct"

	"github.com/influxdata/influxdb/client/v2"
)

func StorageVolumeInfo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Storage API Server] /node/info/volume Called\n"))
}

func StorageVolumeAllocate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Storage API Server] /volume/allocate Called\n"))
}

func StorageVolumeDeallocate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Storage API Server] /volume/deallocate Called\n"))
}

func NodeInfoStorageList(w http.ResponseWriter, r *http.Request) {
	// serverAddress := "http://localhost:" + storagestruct.STORAGE_METRIC_COLLECTOR_PORT_HTTP + "/node/info/storage"

	// for {
	// 	resp, err := http.Get(serverAddress)
	// 	if err != nil {
	// 		fmt.Println("Error sending request:", err)
	// 		time.Sleep(5 * time.Second)
	// 		continue
	// 	}
	// 	defer resp.Body.Close()

	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	var response storagestruct.NodeStorageInfo

	// 	err = json.Unmarshal(body, &response)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	storagestruct.NodeStorageInfo_ = response

	// 	break
	// }

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(storagestruct.NodeStorageInfo_)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		fmt.Println("Error encoding JSON:", err)
	}
}

func NodeMetricAll(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := []storagestruct.NodeMetric{}

	count := r.URL.Query().Get("count")

	if count == "" {
		count = "1"
	}

	q := client.Query{
		Command:  "SELECT * FROM node_metric ORDER BY DESC LIMIT " + count + " TZ('Asia/Seoul')",
		Database: storagestruct.INFLUX_DB,
	}

	if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
		for _, row := range result.Results[0].Series {
			for _, value := range row.Values {
				nodeMetric := storagestruct.NodeMetric{}

				nodeMetric.Time = fmt.Sprintf("%v", value[0])
				nodeMetric.CpuTotal = parseFloat(value[1])
				nodeMetric.CpuUsed = parseFloat(value[2])
				nodeMetric.CpuUtilization = parseFloat(value[3])
				nodeMetric.DiskTotal = parseFloat(value[4])
				nodeMetric.DiskUsed = parseFloat(value[5])
				nodeMetric.DiskUtilization = parseFloat(value[6])
				nodeMetric.MemoryTotal = parseFloat(value[7])
				nodeMetric.MemoryUsed = parseFloat(value[8])
				nodeMetric.MemoryUtilization = parseFloat(value[9])
				nodeMetric.NetworkBandwidth = parseFloat(value[10])
				nodeMetric.NetworkRxData = parseFloat(value[11])
				nodeMetric.NetworkTxData = parseFloat(value[12])
				nodeMetric.NodeName = fmt.Sprintf("%v", value[13])
				nodeMetric.PowerUsed = parseFloat(value[14])

				response = append(response, nodeMetric)
			}
		}
	} else {
		fmt.Println("Error executing query:", err)
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeMetricCpu(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := []storagestruct.CpuMetric{}

	count := r.URL.Query().Get("count")

	if count == "" {
		count = "1"
	}

	q := client.Query{
		Command:  "select cpu_total, cpu_usage, cpu_utilization, node_name from node_metric order by desc limit " + count + " TZ('Asia/Seoul')",
		Database: storagestruct.INFLUX_DB,
	}

	if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
		for _, row := range result.Results[0].Series {
			for _, value := range row.Values {
				cpuMetric := storagestruct.CpuMetric{}

				cpuMetric.Time = fmt.Sprintf("%v", value[0])
				cpuMetric.CpuTotal = parseFloat(value[1])
				cpuMetric.CpuUsed = parseFloat(value[2])
				cpuMetric.CpuUtilization = parseFloat(value[3])
				cpuMetric.Name = fmt.Sprintf("%v", value[4])

				response = append(response, cpuMetric)
			}
		}
	} else {
		fmt.Println("Error executing query:", err)
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeMetricPower(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := []storagestruct.PowerMetric{}

	count := r.URL.Query().Get("count")

	if count == "" {
		count = "1"
	}

	q := client.Query{
		Command:  "select power_usage, node_name from node_metric order by desc limit " + count + " TZ('Asia/Seoul')",
		Database: storagestruct.INFLUX_DB,
	}

	if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
		for _, row := range result.Results[0].Series {
			for _, value := range row.Values {
				powerMetric := storagestruct.PowerMetric{}

				powerMetric.Time = fmt.Sprintf("%v", value[0])
				powerMetric.PowerUsed = parseFloat(value[1])
				powerMetric.Name = fmt.Sprintf("%v", value[2])

				response = append(response, powerMetric)
			}

		}
	} else {
		fmt.Println("Error executing query:", err)
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeMetricMemory(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := []storagestruct.MemoryMetric{}

	count := r.URL.Query().Get("count")

	if count == "" {
		count = "1"
	}

	q := client.Query{
		Command:  "select memory_total, memory_usage, memory_utilization, node_name from node_metric order by desc limit " + count + " TZ('Asia/Seoul')",
		Database: storagestruct.INFLUX_DB,
	}

	if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
		for _, row := range result.Results[0].Series {
			for _, value := range row.Values {
				memoryMetric := storagestruct.MemoryMetric{}

				memoryMetric.Time = fmt.Sprintf("%v", value[0])
				memoryMetric.MemoryTotal = parseFloat(value[1])
				memoryMetric.MemoryUsed = parseFloat(value[2])
				memoryMetric.MemoryUtilization = parseFloat(value[3])
				memoryMetric.Name = fmt.Sprintf("%v", value[4])

				response = append(response, memoryMetric)
			}

		}
	} else {
		fmt.Println("Error executing query:", err)
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeMetricNetwork(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := []storagestruct.NetworkMetric{}

	count := r.URL.Query().Get("count")

	if count == "" {
		count = "1"
	}

	q := client.Query{
		Command:  "select network_bandwidth, network_rx_data, network_tx_data, node_name from node_metric order by desc limit " + count + " TZ('Asia/Seoul')",
		Database: storagestruct.INFLUX_DB,
	}

	if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
		for _, row := range result.Results[0].Series {
			for _, value := range row.Values {
				networkMetric := storagestruct.NetworkMetric{}

				networkMetric.Time = fmt.Sprintf("%v", value[0])
				networkMetric.NetworkRxData = parseFloat(value[1])
				networkMetric.NetworkTxData = parseFloat(value[2])
				networkMetric.NetworkBandwidth = parseFloat(value[3])
				networkMetric.Name = fmt.Sprintf("%v", value[4])

				response = append(response, networkMetric)
			}

		}
	} else {
		fmt.Println("Error executing query:", err)
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeMetricDisk(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := []storagestruct.DiskMetric{}

	count := r.URL.Query().Get("count")

	if count == "" {
		count = "1"
	}

	q := client.Query{
		Command:  "select disk_total, disk_usage, disk_utilization, node_name from node_metric order by desc limit " + count + " TZ('Asia/Seoul')",
		Database: storagestruct.INFLUX_DB,
	}

	if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
		for _, row := range result.Results[0].Series {
			for _, value := range row.Values {
				diskMetric := storagestruct.DiskMetric{}

				diskMetric.Time = fmt.Sprintf("%v", value[0])
				diskMetric.DiskTotal = parseFloat(value[1])
				diskMetric.DiskUsed = parseFloat(value[2])
				diskMetric.DiskUtilization = parseFloat(value[3])
				diskMetric.Name = fmt.Sprintf("%v", value[4])

				response = append(response, diskMetric)
			}

		}
	} else {
		fmt.Println("Error executing query:", err)
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageInfo(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := storagestruct.NewStorageInfoMessage()

	count := r.URL.Query().Get("count")
	targetStorage := r.URL.Query().Get("storage")

	if count == "" {
		count = "1"
	}

	for _, csd := range storagestruct.NodeStorageInfo_.CsdList {
		if targetStorage == "" || targetStorage == csd.Id {
			if csd.Status == storagestruct.READY {
				measurementName := "csd_metric_" + csd.Id
				q := client.Query{
					Command:  "select disk_total, disk_usage, disk_utilization, storage_name, id, ip, metric_score, status from " + measurementName + " order by desc limit " + count + " TZ('Asia/Seoul')",
					Database: storagestruct.INFLUX_DB,
				}

				if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
					csdMetrics := []storagestruct.CsdMetricMin{}

					for _, row := range result.Results[0].Series {
						for _, value := range row.Values {
							csdMetric := storagestruct.CsdMetricMin{}

							csdMetric.Time = fmt.Sprintf("%v", value[0])
							csdMetric.DiskTotal = parseFloat(value[1])
							csdMetric.DiskUsed = parseFloat(value[2])
							csdMetric.DiskUtilization = parseFloat(value[3])
							csdMetric.Name = fmt.Sprintf("%v", value[4])
							csdMetric.Id = fmt.Sprintf("%v", value[5])
							csdMetric.Ip = fmt.Sprintf("%v", value[6])
							csdMetric.CsdMetricScore = parseFloat(value[7])
							csdMetric.Status = fmt.Sprintf("%v", value[8])

							csdMetrics = append(csdMetrics, csdMetric)
						}
					}

					response.CsdList[csd.Id] = csdMetrics
				} else {
					fmt.Println("Error executing query:", err)
				}
			}
		}
	}

	for _, ssd := range storagestruct.NodeStorageInfo_.SsdList {
		if targetStorage == "" || targetStorage == ssd.Id {
			measurementName := "ssd_metric_" + ssd.Id
			q := client.Query{
				Command:  "select * from " + measurementName + " order by desc limit " + count + " TZ('Asia/Seoul')",
				Database: storagestruct.INFLUX_DB,
			}

			if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
				ssdMetrics := []storagestruct.SsdMetric{}
				for _, row := range result.Results[0].Series {
					for _, value := range row.Values {
						ssdMetric := storagestruct.SsdMetric{}

						ssdMetric.Time = fmt.Sprintf("%v", value[0])
						ssdMetric.DiskTotal = parseFloat(value[1])
						ssdMetric.DiskUsed = parseFloat(value[2])
						ssdMetric.DiskUtilization = parseFloat(value[3])
						ssdMetric.Id = fmt.Sprintf("%v", value[4])
						ssdMetric.Status = fmt.Sprintf("%v", value[5])
						ssdMetric.Name = fmt.Sprintf("%v", value[6])

						ssdMetrics = append(ssdMetrics, ssdMetric)
					}
				}
				response.SsdList[ssd.Id] = ssdMetrics
			} else {
				fmt.Println("Error executing query:", err)
			}
		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricAll(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := storagestruct.NewStorageMetricMessage()

	count := r.URL.Query().Get("count")
	targetStorage := r.URL.Query().Get("storage")

	if count == "" {
		count = "1"
	}

	for _, csd := range storagestruct.NodeStorageInfo_.CsdList {
		if targetStorage == "" || targetStorage == csd.Id {
			if csd.Status == storagestruct.READY {
				measurementName := "csd_metric_" + csd.Id
				q := client.Query{
					Command:  "select * from " + measurementName + " order by desc limit " + count + " TZ('Asia/Seoul')",
					Database: storagestruct.INFLUX_DB,
				}

				if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
					csdMetrics := []storagestruct.CsdMetric{}
					for _, row := range result.Results[0].Series {
						for _, value := range row.Values {
							csdMetric := storagestruct.CsdMetric{}

							csdMetric.Time = fmt.Sprintf("%v", value[0])
							csdMetric.CpuTotal = parseFloat(value[1])
							csdMetric.CpuUsed = parseFloat(value[2])
							csdMetric.CpuUtilization = parseFloat(value[3])
							csdMetric.DiskTotal = parseFloat(value[4])
							csdMetric.DiskUsed = parseFloat(value[5])
							csdMetric.DiskUtilization = parseFloat(value[6])
							csdMetric.Id = fmt.Sprintf("%v", value[7])
							csdMetric.Ip = fmt.Sprintf("%v", value[8])
							csdMetric.MemoryTotal = parseFloat(value[9])
							csdMetric.MemoryUsed = parseFloat(value[10])
							csdMetric.MemoryUtilization = parseFloat(value[11])
							csdMetric.CsdMetricScore = parseFloat(value[12])
							csdMetric.NetworkBandwidth = parseFloat(value[13])
							csdMetric.NetworkRxData = parseFloat(value[14])
							csdMetric.NetworkTxData = parseFloat(value[15])
							csdMetric.Status = fmt.Sprintf("%v", value[16])
							csdMetric.Name = fmt.Sprintf("%v", value[17])
							csdMetric.CsdWorkingBlockCount = parseFloat(value[18])

							csdMetrics = append(csdMetrics, csdMetric)
						}
					}
					response.CsdList[csd.Id] = csdMetrics
				} else {
					fmt.Println("Error executing query:", err)
				}
			}
		}
	}

	for _, ssd := range storagestruct.NodeStorageInfo_.SsdList {
		if targetStorage == "" || targetStorage == ssd.Id {
			measurementName := "ssd_metric_" + ssd.Id
			q := client.Query{
				Command:  "select * from " + measurementName + " order by desc limit " + count + " TZ('Asia/Seoul')",
				Database: storagestruct.INFLUX_DB,
			}

			if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
				ssdMetrics := []storagestruct.SsdMetric{}
				fmt.Println(result)
				for _, row := range result.Results[0].Series {
					for _, value := range row.Values {
						ssdMetric := storagestruct.SsdMetric{}

						ssdMetric.Time = fmt.Sprintf("%v", value[0])
						ssdMetric.DiskTotal = parseFloat(value[1])
						ssdMetric.DiskUsed = parseFloat(value[2])
						ssdMetric.DiskUtilization = parseFloat(value[3])
						ssdMetric.Id = fmt.Sprintf("%v", value[4])
						ssdMetric.Status = fmt.Sprintf("%v", value[5])
						ssdMetric.Name = fmt.Sprintf("%v", value[6])

						ssdMetrics = append(ssdMetrics, ssdMetric)
					}

				}
				response.SsdList[ssd.Id] = ssdMetrics
			} else {
				fmt.Println("Error executing query:", err)
			}
		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricCpu(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := map[string][]storagestruct.CpuMetric{}

	count := r.URL.Query().Get("count")
	targetStorage := r.URL.Query().Get("storage")

	if count == "" {
		count = "1"
	}

	for _, csd := range storagestruct.NodeStorageInfo_.CsdList {
		if targetStorage == "" || targetStorage == csd.Id {

			if csd.Status == storagestruct.READY {
				measurementName := "csd_metric_" + csd.Id
				q := client.Query{
					Command:  "select cpu_total, cpu_usage, cpu_utilization, storage_name from " + measurementName + " order by desc limit " + count + " TZ('Asia/Seoul')",
					Database: storagestruct.INFLUX_DB,
				}

				if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
					cpuMetrics := []storagestruct.CpuMetric{}

					for _, row := range result.Results[0].Series {
						for _, value := range row.Values {
							cpuMetric := storagestruct.CpuMetric{}

							cpuMetric.Time = fmt.Sprintf("%v", value[0])
							cpuMetric.CpuTotal = parseFloat(value[1])
							cpuMetric.CpuUsed = parseFloat(value[2])
							cpuMetric.CpuUtilization = parseFloat(value[3])
							cpuMetric.Name = fmt.Sprintf("%v", value[4])

							cpuMetrics = append(cpuMetrics, cpuMetric)
						}
					}
					response[csd.Id] = cpuMetrics
				} else {
					fmt.Println("Error executing query:", err)
				}
			}
		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricPower(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := map[string][]storagestruct.PowerMetric{}

	count := r.URL.Query().Get("count")
	targetStorage := r.URL.Query().Get("storage")

	if count == "" {
		count = "1"
	}

	for _, csd := range storagestruct.NodeStorageInfo_.CsdList {
		if targetStorage == "" || targetStorage == csd.Id {

			if csd.Status == storagestruct.READY {
				measurementName := "csd_metric_" + csd.Id
				q := client.Query{
					Command:  "select power_used, storage_name from " + measurementName + " order by desc limit " + count + " TZ('Asia/Seoul')",
					Database: storagestruct.INFLUX_DB,
				}

				if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
					powerMetrics := []storagestruct.PowerMetric{}

					for _, row := range result.Results[0].Series {
						for _, value := range row.Values {
							powerMetric := storagestruct.PowerMetric{}

							powerMetric.Time = fmt.Sprintf("%v", value[0])
							powerMetric.Name = fmt.Sprintf("%v", value[1])

							powerMetrics = append(powerMetrics, powerMetric)
						}
					}
					response[csd.Id] = powerMetrics
				} else {
					fmt.Println("Error executing query:", err)
				}
			}
		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricMemory(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := map[string][]storagestruct.MemoryMetric{}

	count := r.URL.Query().Get("count")
	targetStorage := r.URL.Query().Get("storage")

	if count == "" {
		count = "1"
	}

	for _, csd := range storagestruct.NodeStorageInfo_.CsdList {
		if targetStorage == "" || targetStorage == csd.Id {

			if csd.Status == storagestruct.READY {
				measurementName := "csd_metric_" + csd.Id
				q := client.Query{
					Command:  "select memory_total, memory_usage, memory_utilization, storage_name from " + measurementName + " order by desc limit " + count + " TZ('Asia/Seoul')",
					Database: storagestruct.INFLUX_DB,
				}

				if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
					memoryMetrics := []storagestruct.MemoryMetric{}

					for _, row := range result.Results[0].Series {
						for _, value := range row.Values {
							memoryMetric := storagestruct.MemoryMetric{}

							memoryMetric.Time = fmt.Sprintf("%v", value[0])
							memoryMetric.MemoryTotal = parseFloat(value[1])
							memoryMetric.MemoryUsed = parseFloat(value[2])
							memoryMetric.MemoryUtilization = parseFloat(value[3])
							memoryMetric.Name = fmt.Sprintf("%v", value[4])

							memoryMetrics = append(memoryMetrics, memoryMetric)
						}
					}
					response[csd.Id] = memoryMetrics
				} else {
					fmt.Println("Error executing query:", err)
				}
			}
		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricNetwork(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := map[string][]storagestruct.NetworkMetric{}

	count := r.URL.Query().Get("count")
	targetStorage := r.URL.Query().Get("storage")

	if count == "" {
		count = "1"
	}

	for _, csd := range storagestruct.NodeStorageInfo_.CsdList {
		if targetStorage == "" || targetStorage == csd.Id {
			if csd.Status == storagestruct.READY {
				measurementName := "csd_metric_" + csd.Id
				q := client.Query{
					Command:  "select network_bandwidth,network_rx_data,network_tx_data, storage_name from " + measurementName + " order by desc limit " + count + " TZ('Asia/Seoul')",
					Database: storagestruct.INFLUX_DB,
				}

				if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
					networkMetrics := []storagestruct.NetworkMetric{}

					for _, row := range result.Results[0].Series {
						for _, value := range row.Values {
							networkMetric := storagestruct.NetworkMetric{}

							networkMetric.Time = fmt.Sprintf("%v", value[0])
							networkMetric.NetworkRxData = parseFloat(value[1])
							networkMetric.NetworkTxData = parseFloat(value[2])
							networkMetric.NetworkBandwidth = parseFloat(value[3])
							networkMetric.Name = fmt.Sprintf("%v", value[4])

							networkMetrics = append(networkMetrics, networkMetric)
						}
					}
					response[csd.Id] = networkMetrics
				} else {
					fmt.Println("Error executing query:", err)
				}
			}
		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricDisk(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := map[string][]storagestruct.DiskMetric{}

	count := r.URL.Query().Get("count")
	targetStorage := r.URL.Query().Get("storage")

	if count == "" {
		count = "1"
	}

	for _, csd := range storagestruct.NodeStorageInfo_.CsdList {
		if targetStorage == "" || targetStorage == csd.Id {
			if csd.Status == storagestruct.READY {
				measurementName := "csd_metric_" + csd.Id
				q := client.Query{
					Command:  "select disk_total, disk_usage, disk_utilization, storage_name from " + measurementName + " order by desc limit " + count + " TZ('Asia/Seoul')",
					Database: storagestruct.INFLUX_DB,
				}

				if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
					diskMetrics := []storagestruct.DiskMetric{}

					for _, row := range result.Results[0].Series {
						for _, value := range row.Values {
							diskMetric := storagestruct.DiskMetric{}

							diskMetric.Time = fmt.Sprintf("%v", value[0])
							diskMetric.DiskTotal = parseFloat(value[1])
							diskMetric.DiskUsed = parseFloat(value[2])
							diskMetric.DiskUtilization = parseFloat(value[3])
							diskMetric.Name = fmt.Sprintf("%v", value[4])

							diskMetrics = append(diskMetrics, diskMetric)
						}
					}
					response[csd.Id] = diskMetrics
				} else {
					fmt.Println("Error executing query:", err)
				}
			}
		}
	}

	for _, ssd := range storagestruct.NodeStorageInfo_.SsdList {
		if targetStorage == "" || targetStorage == ssd.Id {
			measurementName := "ssd_metric_" + ssd.Id
			q := client.Query{
				Command:  "select disk_total, disk_usage, disk_utilization, storage_name from " + measurementName + " order by desc limit " + count + " TZ('Asia/Seoul')",
				Database: storagestruct.INFLUX_DB,
			}

			if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
				diskMetrics := []storagestruct.DiskMetric{}

				for _, row := range result.Results[0].Series {
					for _, value := range row.Values {
						diskMetric := storagestruct.DiskMetric{}

						diskMetric.Time = fmt.Sprintf("%v", value[0])
						diskMetric.DiskTotal = parseFloat(value[1])
						diskMetric.DiskUsed = parseFloat(value[2])
						diskMetric.DiskUtilization = parseFloat(value[3])
						diskMetric.Name = fmt.Sprintf("%v", value[4])

						diskMetrics = append(diskMetrics, diskMetric)
					}
				}
				response[ssd.Id] = diskMetrics
			} else {
				fmt.Println("Error executing query:", err)
			}
		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func parseFloat(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return math.Round(v*100) / 100
	case int:
		v_ := float64(v)
		return math.Round(v_*100) / 100
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return math.Round(f*100) / 100
		}
	case json.Number:
		f, _ := v.Float64()
		return math.Round(f*100) / 100
	default:
		fmt.Printf("Unknown type: %T with value: %v\n", v, v)
	}
	return 0
}
