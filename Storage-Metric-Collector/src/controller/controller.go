package controller

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	// influxdb v1 client
	client "github.com/influxdata/influxdb/client/v2"
)

func (storageMetricCollector *MetricCollector) HandleConnection(conn net.Conn) {
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	var csdMetric *CsdMetric

	message := string(buffer[:n])
	err = json.Unmarshal([]byte(message), &csdMetric)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	csdID := extractCSDId(csdMetric.IP)
	key := "csd" + csdID

	csdMetric.Name = storageMetricCollector.CsdMetrics[key].Name
	csdMetric.Status = READY

	csdMetric.mutex.Lock()
	defer csdMetric.mutex.Unlock()

	storageMetricCollector.CsdMetrics[key] = csdMetric
}

func (storageMetricCollector *MetricCollector) RunMetricCollector(mode string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			storageMetricCollector.NodeMetric.mutex.Lock()

			storageMetricCollector.updateCpu()
			storageMetricCollector.updateMemory()
			storageMetricCollector.updateNetwork()
			storageMetricCollector.updateStorage()
			storageMetricCollector.updatePower()

			if storageMetricCollector.NodeType == SSD {
				storageMetricCollector.updateSsdMetric()
			}

			storageMetricCollector.saveNodeMetric(mode)

			storageMetricCollector.NodeMetric.mutex.Unlock()
		}
	}
}

func (storageMetricCollector *MetricCollector) updateCpu() {
	file, err := os.Open("/proc/stat")
	if err != nil {
		fmt.Println("cannot open file: ", err)
	} else {
		var cpuID string

		var curJiffies, diffJiffies StJiffies

		_, err = fmt.Fscanf(file, "%5s %d %d %d %d", &cpuID, &curJiffies.User, &curJiffies.Nice, &curJiffies.System, &curJiffies.Idle)
		if err != nil {
			fmt.Println("Error reading data from file:", err)
		}

		diffJiffies.User = curJiffies.User - storageMetricCollector.NodeMetric.Cpu.StJiffies.User
		diffJiffies.Nice = curJiffies.Nice - storageMetricCollector.NodeMetric.Cpu.StJiffies.Nice
		diffJiffies.System = curJiffies.System - storageMetricCollector.NodeMetric.Cpu.StJiffies.System
		diffJiffies.Idle = curJiffies.Idle - storageMetricCollector.NodeMetric.Cpu.StJiffies.Idle

		totalJiffies := diffJiffies.User + diffJiffies.Nice + diffJiffies.System + diffJiffies.Idle

		utilization := 100.0 * (1.0 - float64(diffJiffies.Idle)/float64(totalJiffies))
		storageMetricCollector.NodeMetric.Cpu.Utilization, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", utilization), 64)
		used := float64(storageMetricCollector.NodeMetric.Cpu.Total) * (1.0 - float64(diffJiffies.Idle)/float64(totalJiffies))
		storageMetricCollector.NodeMetric.Cpu.Used, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", used), 64)

		storageMetricCollector.NodeMetric.Cpu.StJiffies = curJiffies
	}
	file.Close()
}

func (storageMetricCollector *MetricCollector) updateMemory() {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		fmt.Println("cannot open file: ", err)
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)

			if len(fields) < 2 {
				continue
			}

			key := fields[0]
			value, err := strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				fmt.Println("Error parsing value:", err)
				continue
			}

			switch key {
			case "MemFree:":
				storageMetricCollector.NodeMetric.Memory.Free = value
			case "Buffers:":
				storageMetricCollector.NodeMetric.Memory.Buffers = value
			case "Cached:":
				storageMetricCollector.NodeMetric.Memory.Cached = value
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}

		storageMetricCollector.NodeMetric.Memory.Used = storageMetricCollector.NodeMetric.Memory.Total - storageMetricCollector.NodeMetric.Memory.Free - storageMetricCollector.NodeMetric.Memory.Buffers - storageMetricCollector.NodeMetric.Memory.Cached
		utilization := float64(storageMetricCollector.NodeMetric.Memory.Used) / float64(storageMetricCollector.NodeMetric.Memory.Total) * 100.0
		storageMetricCollector.NodeMetric.Memory.Utilization, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", utilization), 64)
	}
	file.Close()
}

func (storageMetricCollector *MetricCollector) updateNetwork() {
	statisticsFilePath := "/sys/class/net/eno1/statistics/"

	rxBytesFieldName := statisticsFilePath + "rx_bytes"
	txBytesFieldName := statisticsFilePath + "tx_bytes"

	currentRxBytesStr, err := readStatisticsField(rxBytesFieldName)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentTxBytesStr, err := readStatisticsField(txBytesFieldName)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentRxBytes, _ := strconv.ParseInt(currentRxBytesStr, 10, 64)
	currentTxBytes, _ := strconv.ParseInt(currentTxBytesStr, 10, 64)

	storageMetricCollector.NodeMetric.Network.RxData = currentRxBytes - storageMetricCollector.NodeMetric.Network.RxByte
	storageMetricCollector.NodeMetric.Network.TxData = currentTxBytes - storageMetricCollector.NodeMetric.Network.TxByte

	storageMetricCollector.NodeMetric.Network.Bandwidth = (storageMetricCollector.NodeMetric.Network.RxData + storageMetricCollector.NodeMetric.Network.TxData) / 5 * 8

	storageMetricCollector.NodeMetric.Network.RxByte = currentRxBytes
	storageMetricCollector.NodeMetric.Network.TxByte = currentTxBytes
}

func (storageMetricCollector *MetricCollector) updateStorage() {
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
				storageMetricCollector.NodeMetric.Disk.Used, _ = strconv.ParseInt(fields[2], 10, 64)
				break
			}
		}
	}

	if storageMetricCollector.NodeMetric.Disk.Total > 0 {
		utilization := (float64(storageMetricCollector.NodeMetric.Disk.Used) / float64(storageMetricCollector.NodeMetric.Disk.Total)) * 100
		storageMetricCollector.NodeMetric.Disk.Utilization, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", utilization), 64)
	}
}

func (storageMetricCollector *MetricCollector) updatePower() {
	energyFieldName1 := "/sys/class/powercap/intel-rapl:0/energy_uj"
	energyFieldName2 := "/sys/class/powercap/intel-rapl:1/energy_uj"

	currentEnergyStr1, err := readStatisticsField(energyFieldName1)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentEnergyStr2, err := readStatisticsField(energyFieldName2)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentEnergy1, _ := strconv.ParseInt(currentEnergyStr1, 10, 64)
	currentEnergy2, _ := strconv.ParseInt(currentEnergyStr2, 10, 64)

	energyDiffJ1 := float64(currentEnergy1-storageMetricCollector.NodeMetric.Power.Energy1) / 1e6
	energyDiffJ2 := float64(currentEnergy2-storageMetricCollector.NodeMetric.Power.Energy2) / 1e6

	storageMetricCollector.NodeMetric.Power.Used = int64((energyDiffJ1 + energyDiffJ2) / 1.0)
	storageMetricCollector.NodeMetric.Power.Energy1 = currentEnergy1
	storageMetricCollector.NodeMetric.Power.Energy2 = currentEnergy2
}

func (storageMetricCollector *MetricCollector) updateSsdMetric() {
	for _, ssd := range storageMetricCollector.SsdMetrics {
		mountpoint := fmt.Sprintf("/dev/%s", ssd.Name)
		dfCmd := exec.Command("df", "-BM", "--output=used", mountpoint)
		var dfOut bytes.Buffer
		dfCmd.Stdout = &dfOut
		dfCmd.Run()
		dfLines := strings.Split(dfOut.String(), "\n")
		if len(dfLines) > 1 {
			usedStr := strings.TrimSuffix(dfLines[1], "M")
			usedSize, _ := strconv.ParseInt(usedStr, 10, 64)
			utilization := (float64(usedSize) / float64(ssd.Total)) * 100
			ssd.Used = usedSize
			ssd.Utilization, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", utilization), 64)
		}
	}
}

func (storageMetricCollector *MetricCollector) saveNodeMetric(mode string) {
	if mode == "off" {
		return
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  INFLUX_DB,
		Precision: "s",
	})
	if err != nil {
		fmt.Println("DB NewBatchPoints Error:", err)
		return
	}

	var INFLUXDB_NODE_MEASUREMENT = "node_metric"

	fields := map[string]interface{}{
		"node_name": storageMetricCollector.NodeName,

		"cpu_total":       storageMetricCollector.NodeMetric.Cpu.Total,
		"cpu_usage":       storageMetricCollector.NodeMetric.Cpu.Used,
		"cpu_utilization": storageMetricCollector.NodeMetric.Cpu.Utilization,

		"memory_total":       storageMetricCollector.NodeMetric.Memory.Total,
		"memory_usage":       storageMetricCollector.NodeMetric.Memory.Used,
		"memory_utilization": storageMetricCollector.NodeMetric.Memory.Utilization,

		"disk_total":       storageMetricCollector.NodeMetric.Disk.Total,
		"disk_usage":       storageMetricCollector.NodeMetric.Disk.Used,
		"disk_utilization": storageMetricCollector.NodeMetric.Disk.Utilization,

		"network_bandwidth": storageMetricCollector.NodeMetric.Network.Bandwidth,
		"network_rx_data":   storageMetricCollector.NodeMetric.Network.RxData,
		"network_tx_data":   storageMetricCollector.NodeMetric.Network.TxData,

		"power_usage": storageMetricCollector.NodeMetric.Power.Used,
	}

	pt, err := client.NewPoint(INFLUXDB_NODE_MEASUREMENT, nil, fields, time.Now())
	if err != nil {
		fmt.Println("DB NewPoint Error:", err)
		return
	}
	bp.AddPoint(pt)

	err = INFLUX_CLIENT.Write(bp)
	if err != nil {
		fmt.Println("DB Write Error:", err)
		return
	}

	for key, metric := range storageMetricCollector.SsdMetrics {
		// fmt.Println("ssd : ", key, "metric")
		// fmt.Printf("%+v\n", metric)

		var INFLUXDB_SSD_MEASUREMENT = "ssd_metric_" + key

		fields := map[string]interface{}{
			"id":           key,
			"storage_name": metric.Name,

			"disk_total":       metric.Total,
			"disk_usage":       metric.Used,
			"disk_utilization": metric.Utilization,
			"status":           metric.Status,
		}

		pt, err := client.NewPoint(INFLUXDB_SSD_MEASUREMENT, nil, fields, time.Now())
		if err != nil {
			fmt.Println("DB NewPoint Error:", err)
			return
		}
		bp.AddPoint(pt)

		err = INFLUX_CLIENT.Write(bp)
		if err != nil {
			fmt.Println("DB Write Error:", err)
			return
		}
	}
}

func (storageMetricCollector *MetricCollector) SaveCsdMetric(mode string) {
	if mode == "off" {
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for key, metric := range storageMetricCollector.CsdMetrics {
				if metric.Status == READY {
					metric.mutex.Lock()

					// fmt.Println("csd : ", key)
					// fmt.Printf("%+v\n", metric)

					bp, err := client.NewBatchPoints(client.BatchPointsConfig{
						Database:  INFLUX_DB,
						Precision: "s",
					})
					if err != nil {
						fmt.Println("DB NewBatchPoints Error:", err)
						metric.mutex.Unlock()
						break
					}

					var INFLUXDB_CSD_MEASUREMENT = "csd_metric_" + key

					fields := map[string]interface{}{
						"id":           key,
						"storage_name": metric.Name,
						"ip":           metric.IP,

						"cpu_total":       metric.CpuTotal,
						"cpu_usage":       metric.CpuUsed,
						"cpu_utilization": metric.CpuUtilization,

						"memory_total":       metric.MemoryTotal,
						"memory_usage":       metric.MemoryUsed,
						"memory_utilization": metric.MemoryUtilization,

						"disk_total":       metric.DiskTotal,
						"disk_usage":       metric.DiskUsed,
						"disk_utilization": metric.DiskUtilization,

						"network_bandwidth": metric.NetworkBandwidth,
						"network_rx_data":   metric.NetworkRxData,
						"network_tx_data":   metric.NetworkTxData,

						"metric_score":        metric.CsdMetricScore,
						"working_block_count": metric.CsdWorkingBlockCount,
						"status":              metric.Status,
					}

					pt, err := client.NewPoint(INFLUXDB_CSD_MEASUREMENT, nil, fields, time.Now())
					if err != nil {
						fmt.Println("DB NewPoint Error:", err)
						metric.mutex.Unlock()
						break
					}
					bp.AddPoint(pt)

					err = INFLUX_CLIENT.Write(bp)
					if err != nil {
						fmt.Println("DB Write Error:", err)
						metric.mutex.Unlock()
						break
					}

					metric.mutex.Unlock()
				}
			}
		}
	}
}

func (storageMetricCollector *MetricCollector) HandleNodeInfoStorage(w http.ResponseWriter, r *http.Request) {
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

	response := struct {
		NodeName string     `json:"nodeName"`
		CsdList  []CsdEntry `json:"csdList"`
		SsdList  []SsdEntry `json:"ssdList"`
		NodeType string     `json:"nodeType"`
	}{
		NodeName: storageMetricCollector.NodeName,
		NodeType: storageMetricCollector.NodeType,
	}

	for key, metric := range storageMetricCollector.CsdMetrics {
		response.CsdList = append(response.CsdList, CsdEntry{Id: key, Name: metric.Name, Status: metric.Status})
	}

	for key, metric := range storageMetricCollector.SsdMetrics {
		response.SsdList = append(response.SsdList, SsdEntry{Id: key, Name: metric.Name, Status: metric.Status})
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		fmt.Println("Error encoding JSON:", err)
	}

	fmt.Printf("HandleNodeInfoStorage called %+v\n", response)
}

func (storageMetricCollector *MetricCollector) HandleNodeMetric(w http.ResponseWriter, r *http.Request) {
	response := NodeMetric{}

	response.Cpu = storageMetricCollector.NodeMetric.Cpu
	response.Memory = storageMetricCollector.NodeMetric.Memory
	response.Disk = storageMetricCollector.NodeMetric.Disk
	response.Network = storageMetricCollector.NodeMetric.Network
	response.Power = storageMetricCollector.NodeMetric.Power

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		fmt.Println("Error encoding JSON:", err)
	}

	fmt.Printf("HandleNodeMetric called %+v\n", response)
}

func (storageMetricCollector *MetricCollector) HandleStorageMetric(w http.ResponseWriter, r *http.Request) {
	response := struct {
		CsdMetrics map[string]CsdMetric `json:"csdMetrics"`
		SsdMetrics map[string]SsdMetric `json:"ssdMetrics"`
	}{
		CsdMetrics: make(map[string]CsdMetric),
		SsdMetrics: make(map[string]SsdMetric),
	}

	for key, metric := range storageMetricCollector.CsdMetrics {
		response.CsdMetrics[key] = *metric
	}

	for key, metric := range storageMetricCollector.SsdMetrics {
		response.SsdMetrics[key] = *metric
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		fmt.Println("Error encoding JSON:", err)
	}

	fmt.Printf("HandleStorageMetric called %+v\n", response)
}
