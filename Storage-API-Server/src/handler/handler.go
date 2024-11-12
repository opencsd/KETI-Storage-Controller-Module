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
	fmt.Println("[OpenCSD Storage API Server] /node/info/volume Called\n")
	w.Write([]byte("[OpenCSD Storage API Server] /node/info/volume Called\n"))
}

func StorageVolumeAllocate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Storage API Server] /volume/allocate Called\n")
	w.Write([]byte("[OpenCSD Storage API Server] /volume/allocate Called\n"))
}

func StorageVolumeDeallocate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Storage API Server] /volume/deallocate Called\n")
	w.Write([]byte("[OpenCSD Storage API Server] /volume/deallocate Called\n"))
}

func NodeInfoStorage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Storage API Server] /node/info/storage Called\n")
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(storagestruct.NodeStorageInfo_)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		fmt.Println("Error encoding JSON:", err)
	}
}

func NodeMetricAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Storage API Server] /node/metric/all Called\n")
	var jsonResponse []byte
	response := []storagestruct.NodeMetric{}

	datanum := r.URL.Query().Get("datanum")

	if datanum == "" {
		datanum = "1"
	}

	q := client.Query{
		Command:  "SELECT * FROM node_metric ORDER BY DESC LIMIT " + datanum + " TZ('Asia/Seoul')",
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
				nodeMetric.StorageTotal = parseFloat(value[4])
				nodeMetric.StorageUsed = parseFloat(value[5])
				nodeMetric.StorageUtilization = parseFloat(value[6])
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
	fmt.Println(string(jsonResponse))
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeMetricCpu(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Storage API Server] /node/metric/cpu Called\n")
	var jsonResponse []byte
	response := []storagestruct.CpuMetric{}

	datanum := r.URL.Query().Get("datanum")

	if datanum == "" {
		datanum = "1"
	}

	q := client.Query{
		Command:  "select cpu_total, cpu_usage, cpu_utilization, node_name from node_metric order by desc limit " + datanum + " TZ('Asia/Seoul')",
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
	fmt.Println(string(jsonResponse))
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeMetricPower(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Storage API Server] /node/metric/power Called\n")
	var jsonResponse []byte
	response := []storagestruct.PowerMetric{}

	datanum := r.URL.Query().Get("datanum")

	if datanum == "" {
		datanum = "1"
	}

	q := client.Query{
		Command:  "select power_usage, node_name from node_metric order by desc limit " + datanum + " TZ('Asia/Seoul')",
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
	fmt.Println(string(jsonResponse))
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeMetricMemory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Storage API Server] /node/metric/memory Called\n")
	var jsonResponse []byte
	response := []storagestruct.MemoryMetric{}

	datanum := r.URL.Query().Get("datanum")

	if datanum == "" {
		datanum = "1"
	}

	q := client.Query{
		Command:  "select memory_total, memory_usage, memory_utilization, node_name from node_metric order by desc limit " + datanum + " TZ('Asia/Seoul')",
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
	fmt.Println(string(jsonResponse))
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeMetricNetwork(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Storage API Server] /node/metric/network Called\n")
	var jsonResponse []byte
	response := []storagestruct.NetworkMetric{}

	datanum := r.URL.Query().Get("datanum")

	if datanum == "" {
		datanum = "1"
	}

	q := client.Query{
		Command:  "select network_bandwidth, network_rx_data, network_tx_data, node_name from node_metric order by desc limit " + datanum + " TZ('Asia/Seoul')",
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
	fmt.Println(string(jsonResponse))
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeMetricStorage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Storage API Server] /node/metric/storage Called\n")
	var jsonResponse []byte
	response := []storagestruct.DiskMetric{}

	datanum := r.URL.Query().Get("datanum")

	if datanum == "" {
		datanum = "1"
	}

	q := client.Query{
		Command:  "select disk_total, disk_usage, disk_utilization, node_name from node_metric order by desc limit " + datanum + " TZ('Asia/Seoul')",
		Database: storagestruct.INFLUX_DB,
	}

	if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
		for _, row := range result.Results[0].Series {
			for _, value := range row.Values {
				diskMetric := storagestruct.DiskMetric{}

				diskMetric.Time = fmt.Sprintf("%v", value[0])
				diskMetric.StorageTotal = parseFloat(value[1])
				diskMetric.StorageUsed = parseFloat(value[2])
				diskMetric.StorageUtilization = parseFloat(value[3])
				diskMetric.Name = fmt.Sprintf("%v", value[4])

				response = append(response, diskMetric)
			}

		}
	} else {
		fmt.Println("Error executing query:", err)
	}

	jsonResponse, _ = json.Marshal(response)
	fmt.Println(string(jsonResponse))
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Storage API Server] /storage/metric/all Called\n")
	var jsonResponse []byte
	response := storagestruct.StorageMetric{}

	datanum := r.URL.Query().Get("datanum")

	if datanum == "" {
		datanum = "1"
	}

	for _, csd := range storagestruct.NodeStorageInfo_.CsdList {
		if csd.Status == "READY" {
			measurementName := "csd_metric_" + csd.CsdName
			q := client.Query{
				Command:  "select * from " + measurementName + " order by desc limit " + datanum + " TZ('Asia/Seoul')",
				Database: storagestruct.INFLUX_DB,
			}

			if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
				for _, row := range result.Results[0].Series {
					for _, value := range row.Values {
						csdMetric := storagestruct.CsdMetric{}

						csdMetric.Time = fmt.Sprintf("%v", value[0])
						csdMetric.CpuTotal = parseFloat(value[1])
						csdMetric.CpuUsed = parseFloat(value[2])
						csdMetric.CpuUtilization = parseFloat(value[3])
						csdMetric.MemoryTotal = parseFloat(value[4])
						csdMetric.MemoryUsed = parseFloat(value[5])
						csdMetric.MemoryUtilization = parseFloat(value[6])
						csdMetric.StorageTotal = parseFloat(value[7])
						csdMetric.StorageUsed = parseFloat(value[8])
						csdMetric.StorageUtilization = parseFloat(value[9])
						csdMetric.NetworkRxData = parseFloat(value[10])
						csdMetric.NetworkTxData = parseFloat(value[11])
						csdMetric.NetworkBandwidth = parseFloat(value[12])
						csdMetric.CsdMetricScore = parseFloat(value[13])
						csdMetric.CsdWorkingBlockCount = parseFloat(value[14])
						csdMetric.Status = fmt.Sprintf("%v", value[15])
						csdMetric.Name = csd.CsdName

						response.CsdList = append(response.CsdList, csdMetric)
					}

				}
			} else {
				fmt.Println("Error executing query:", err)
			}
		}
	}

	for _, ssdName := range storagestruct.NodeStorageInfo_.SsdList {
		measurementName := "ssd_metric_" + ssdName
		q := client.Query{
			Command:  "select * from " + measurementName + " order by desc limit " + datanum + " TZ('Asia/Seoul')",
			Database: storagestruct.INFLUX_DB,
		}

		if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
			for _, row := range result.Results[0].Series {
				for _, value := range row.Values {
					ssdMetric := storagestruct.DiskMetric{}

					ssdMetric.Time = fmt.Sprintf("%v", value[0])
					ssdMetric.StorageTotal = parseFloat(value[1])
					ssdMetric.StorageUsed = parseFloat(value[2])
					ssdMetric.StorageUtilization = parseFloat(value[3])
					ssdMetric.Name = ssdName

					response.SsdList = append(response.SsdList, ssdMetric)
				}

			}
		} else {
			fmt.Println("Error executing query:", err)
		}
	}

	jsonResponse, _ = json.Marshal(response)
	fmt.Println(string(jsonResponse))
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricCpu(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Storage API Server] /storage/metric/cpu Called\n")
	var jsonResponse []byte
	response := []storagestruct.CpuMetric{}

	datanum := r.URL.Query().Get("datanum")

	if datanum == "" {
		datanum = "1"
	}

	for _, csd := range storagestruct.NodeStorageInfo_.CsdList {
		measurementName := "csd_metric_" + csd.CsdName
		q := client.Query{
			Command:  "select cpu_total, cpu_usage, cpu_utilization  from " + measurementName + " order by desc limit " + datanum + " TZ('Asia/Seoul')",
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
					cpuMetric.Name = csd.CsdName

					response = append(response, cpuMetric)
				}

			}
		} else {
			fmt.Println("Error executing query:", err)
		}
	}

	jsonResponse, _ = json.Marshal(response)
	fmt.Println(string(jsonResponse))
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricPower(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Storage API Server] /storage/metric/power Called\n")
	var jsonResponse []byte
	response := []storagestruct.PowerMetric{}

	datanum := r.URL.Query().Get("datanum")

	if datanum == "" {
		datanum = "1"
	}

	for _, csd := range storagestruct.NodeStorageInfo_.CsdList {
		measurementName := "csd_metric_" + csd.CsdName
		q := client.Query{
			Command:  "select power_used from " + measurementName + " order by desc limit " + datanum + " TZ('Asia/Seoul')",
			Database: storagestruct.INFLUX_DB,
		}

		if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
			for _, row := range result.Results[0].Series {
				for _, value := range row.Values {
					powerMetric := storagestruct.PowerMetric{}

					powerMetric.Time = fmt.Sprintf("%v", value[0])
					powerMetric.Name = csd.CsdName

					response = append(response, powerMetric)
				}

			}
		} else {
			fmt.Println("Error executing query:", err)
		}
	}

	jsonResponse, _ = json.Marshal(response)
	fmt.Println(string(jsonResponse))
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricMemory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Storage API Server] /storage/metric/memory Called\n")
	var jsonResponse []byte
	response := []storagestruct.MemoryMetric{}

	datanum := r.URL.Query().Get("datanum")

	if datanum == "" {
		datanum = "1"
	}

	for _, csd := range storagestruct.NodeStorageInfo_.CsdList {
		measurementName := "csd_metric_" + csd.CsdName
		q := client.Query{
			Command:  "select memory_total, memory_usage, memory_utilization from " + measurementName + " order by desc limit " + datanum + " TZ('Asia/Seoul')",
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
					memoryMetric.Name = csd.CsdName

					response = append(response, memoryMetric)
				}

			}
		} else {
			fmt.Println("Error executing query:", err)
		}
	}

	jsonResponse, _ = json.Marshal(response)
	fmt.Println(string(jsonResponse))
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricNetwork(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Storage API Server] /storage/metric/network Called\n")
	var jsonResponse []byte
	response := []storagestruct.NetworkMetric{}

	datanum := r.URL.Query().Get("datanum")

	if datanum == "" {
		datanum = "1"
	}

	for _, csd := range storagestruct.NodeStorageInfo_.CsdList {
		measurementName := "csd_metric_" + csd.CsdName
		q := client.Query{
			Command:  "select network_bandwidth,network_rx_data,network_tx_data  from " + measurementName + " order by desc limit " + datanum + " TZ('Asia/Seoul')",
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
					networkMetric.Name = csd.CsdName

					response = append(response, networkMetric)
				}

			}
		} else {
			fmt.Println("Error executing query:", err)
		}
	}

	jsonResponse, _ = json.Marshal(response)
	fmt.Println(string(jsonResponse))
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricDisk(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Storage API Server] /storage/metric/disk Called\n")
	var jsonResponse []byte
	response := []storagestruct.DiskMetric{}

	datanum := r.URL.Query().Get("datanum")

	if datanum == "" {
		datanum = "1"
	}

	for _, csd := range storagestruct.NodeStorageInfo_.CsdList {
		measurementName := "csd_metric_" + csd.CsdName
		q := client.Query{
			Command:  "select disk_total, disk_usage, disk_utilization from " + measurementName + " order by desc limit " + datanum + " TZ('Asia/Seoul')",
			Database: storagestruct.INFLUX_DB,
		}

		if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
			for _, row := range result.Results[0].Series {
				for _, value := range row.Values {
					diskMetric := storagestruct.DiskMetric{}

					diskMetric.Time = fmt.Sprintf("%v", value[0])
					diskMetric.StorageTotal = parseFloat(value[1])
					diskMetric.StorageUsed = parseFloat(value[2])
					diskMetric.StorageUtilization = parseFloat(value[3])
					diskMetric.Name = csd.CsdName

					response = append(response, diskMetric)
				}

			}
		} else {
			fmt.Println("Error executing query:", err)
		}
	}

	for _, ssdName := range storagestruct.NodeStorageInfo_.SsdList {
		measurementName := "ssd_metric_" + ssdName
		q := client.Query{
			Command:  "select disk_total, disk_usage, disk_utilization from " + measurementName + " order by desc limit " + datanum + " TZ('Asia/Seoul')",
			Database: storagestruct.INFLUX_DB,
		}

		if result, err := storagestruct.INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
			for _, row := range result.Results[0].Series {
				for _, value := range row.Values {
					diskMetric := storagestruct.DiskMetric{}

					diskMetric.Time = fmt.Sprintf("%v", value[0])
					diskMetric.StorageTotal = parseFloat(value[1])
					diskMetric.StorageUsed = parseFloat(value[2])
					diskMetric.StorageUtilization = parseFloat(value[3])
					diskMetric.Name = ssdName

					response = append(response, diskMetric)
				}

			}
		} else {
			fmt.Println("Error executing query:", err)
		}
	}

	jsonResponse, _ = json.Marshal(response)
	fmt.Println(string(jsonResponse))
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
