package controller

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func (storageMetricCollector *MetricCollector) HandleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 4096) // 4096바이트 버퍼 생성
	n, err := conn.Read(buffer)  // 데이터 읽기
	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	var csdMetric *CsdMetric

	message := string(buffer[:n])
	fmt.Printf("Received JSON Data: %s\n", message)
	err = json.Unmarshal([]byte(message), &csdMetric)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	csdID := extractCSDId(csdMetric.IP)
	csdName := "nvme" + csdID

	if csdMetric, exists := storageMetricCollector.CsdMetrics[csdName]; exists {
		csdMetric.mutex.Lock()
		defer csdMetric.mutex.Unlock()
	}

	storageMetricCollector.CsdMetrics[csdName] = csdMetric
}

func (storageMetricCollector *MetricCollector) RunMetricCollector() {
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

			storageMetricCollector.saveNodeMetric()

			storageMetricCollector.NodeMetric.mutex.Unlock()
		}
	}
}

func (storageMetricCollector *MetricCollector) updateCpu() {
	file, err := os.Open("/proc/stat" /*"/metric/proc/stat"*/)
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
	file, err := os.Open("/proc/meminfo" /*"/metric/proc/meminfo"*/)
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
	statisticsFilePath := "/sys/class/net/eno1/statistics/" //"/metric/sys/class/net/eno1/statistics/"

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
				storageMetricCollector.NodeMetric.Storage.Used, _ = strconv.ParseInt(fields[2], 10, 64)
				break
			}
		}
	}

	if storageMetricCollector.NodeMetric.Storage.Total > 0 {
		utilization := (float64(storageMetricCollector.NodeMetric.Storage.Used) / float64(storageMetricCollector.NodeMetric.Storage.Total)) * 100
		storageMetricCollector.NodeMetric.Storage.Utilization, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", utilization), 64)
	}
}

func (storageMetricCollector *MetricCollector) updatePower() {
	powerFilePath1 := "intel-rapl:0"
	powerFilePath2 := "intel-rapl:1"
	energyFieldName1 := "/sys/class/powercap/" + powerFilePath1 + "/energy_uj"
	energyFieldName2 := "/sys/class/powercap/" + powerFilePath2 + "/energy_uj"

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

	storageMetricCollector.NodeMetric.Power.Used = (currentEnergy1 - storageMetricCollector.NodeMetric.Power.Energy1) + (currentEnergy2 - storageMetricCollector.NodeMetric.Power.Energy2)

	storageMetricCollector.NodeMetric.Power.Energy1 = currentEnergy1
	storageMetricCollector.NodeMetric.Power.Energy2 = currentEnergy2
}

func (storageMetricCollector *MetricCollector) updateSsdMetric() {
	for key, ssd := range storageMetricCollector.SsdMetrics {
		mountpoint := fmt.Sprintf("/dev/%s", key)
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

func (storageMetricCollector *MetricCollector) saveNodeMetric() {
	//influxdb에 저장
	fmt.Println("save node metric")

	fmt.Printf("%+v\n", storageMetricCollector.NodeMetric)

	for key, metric := range storageMetricCollector.SsdMetrics {
		fmt.Println("ssd : ", key, "metric")
		fmt.Printf("%+v\n", metric)
	}
}

func (storageMetricCollector *MetricCollector) SaveCsdMetric() {
	//influxdb에 저장
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("save csd metric")

			for key, metric := range storageMetricCollector.CsdMetrics {
				metric.mutex.Lock()

				fmt.Println("csd : ", key, "metric")
				fmt.Printf("%+v\n", metric)

				metric.mutex.Unlock()
			}
		}
	}

}
