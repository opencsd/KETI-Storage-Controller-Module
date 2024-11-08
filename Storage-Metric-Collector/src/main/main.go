package main

import (
	"fmt"
	"net"
	"os"

	"storage-metric-collector/src/controller"
)

func main() {
	controller.STORAGE_METRIC_COLLECTOR_PORT = os.Getenv("STORAGE_METRIC_COLLECTOR_PORT")
	controller.STORAGE_METRIC_DB_PORT = os.Getenv("STORAGE_METRIC_DB_PORT")

	StorageMetricCollector := controller.NewMetricCollector()
	StorageMetricCollector.InitMetricCollector()

	go StorageMetricCollector.RunMetricCollector()

	if StorageMetricCollector.NodeType == controller.CSD {
		go StorageMetricCollector.SaveCsdMetric()
	}

	listener, err := net.Listen("tcp", "0.0.0.0:"+controller.STORAGE_METRIC_COLLECTOR_PORT)
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("[Storage Metric Collector] run on 0.0.0.0:", controller.STORAGE_METRIC_COLLECTOR_PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go StorageMetricCollector.HandleConnection(conn)
	}
}

/*
csd :  2 metric
&{mutex:{state:1 sema:0} IP:10.1.2.2 CpuTotal:4 CpuUsed:0.11999999731779099 CpuUtilization:3.0399999618530273 MemoryTotal:6111708 MemoryUsed:2486808 MemoryUtilization:40.689998626708984 StorageTotal:208562268 StorageUsed:37444232 StorageUtilization:17.950000762939453 NetworkRxData:206 NetworkTxData:716 NetworkBandwidth:1472 CsdMetricScore:74.37 CsdWorkingBlockCount:0}
csd :  3 metric
&{mutex:{state:1 sema:0} IP:10.1.3.2 CpuTotal:4 CpuUsed:0.1599999964237213 CpuUtilization:3.9200000762939453 MemoryTotal:6111708 MemoryUsed:2480804 MemoryUtilization:40.59000015258789 StorageTotal:208562268 StorageUsed:35088192 StorageUtilization:16.81999969482422 NetworkRxData:206 NetworkTxData:713 NetworkBandwidth:1464 CsdMetricScore:74.08 CsdWorkingBlockCount:0}
csd :  4 metric
&{mutex:{state:1 sema:0} IP:10.1.4.2 CpuTotal:4 CpuUsed:0.10000000149011612 CpuUtilization:2.4000000953674316 MemoryTotal:6111708 MemoryUsed:2471596 MemoryUtilization:40.439998626708984 StorageTotal:208562268 StorageUsed:35011760 StorageUtilization:16.790000915527344 NetworkRxData:206 NetworkTxData:715 NetworkBandwidth:1472 CsdMetricScore:74.78 CsdWorkingBlockCount:0}
csd :  5 metric
&{mutex:{state:1 sema:0} IP:10.1.5.2 CpuTotal:4 CpuUsed:0.10000000149011612 CpuUtilization:2.549999952316284 MemoryTotal:6111708 MemoryUsed:2478580 MemoryUtilization:40.54999923706055 StorageTotal:208562268 StorageUsed:34814984 StorageUtilization:16.690000534057617 NetworkRxData:206 NetworkTxData:714 NetworkBandwidth:1472 CsdMetricScore:74.65 CsdWorkingBlockCount:0}
csd :  6 metric
&{mutex:{state:1 sema:0} IP:10.1.6.2 CpuTotal:4 CpuUsed:0.09000000357627869 CpuUtilization:2.1500000953674316 MemoryTotal:6111708 MemoryUsed:2489836 MemoryUtilization:40.7400016784668 StorageTotal:208562268 StorageUsed:34815112 StorageUtilization:16.690000534057617 NetworkRxData:206 NetworkTxData:713 NetworkBandwidth:1464 CsdMetricScore:74.7 CsdWorkingBlockCount:0}
csd :  7 metric
&{mutex:{state:1 sema:0} IP:10.1.7.2 CpuTotal:4 CpuUsed:0.10999999940395355 CpuUtilization:2.6500000953674316 MemoryTotal:6111708 MemoryUsed:2476344 MemoryUtilization:40.52000045776367 StorageTotal:208562268 StorageUsed:35011848 StorageUtilization:16.790000915527344 NetworkRxData:206 NetworkTxData:713 NetworkBandwidth:1464 CsdMetricScore:74.63 CsdWorkingBlockCount:0}
csd :  8 metric
&{mutex:{state:1 sema:0} IP:10.1.8.2 CpuTotal:4 CpuUsed:0.07999999821186066 CpuUtilization:2.0999999046325684 MemoryTotal:6111708 MemoryUsed:2484176 MemoryUtilization:40.650001525878906 StorageTotal:208562268 StorageUsed:35011448 StorageUtilization:16.790000915527344 NetworkRxData:206 NetworkTxData:715 NetworkBandwidth:1472 CsdMetricScore:74.77 CsdWorkingBlockCount:0}
csd :  1 metric
&{mutex:{state:1 sema:0} IP:10.1.1.2 CpuTotal:4 CpuUsed:0.09000000357627869 CpuUtilization:2.259999990463257 MemoryTotal:6111708 MemoryUsed:2477200 MemoryUtilization:40.529998779296875 StorageTotal:208562268 StorageUsed:34810776 StorageUtilization:16.690000534057617 NetworkRxData:206 NetworkTxData:716 NetworkBandwidth:1472 CsdMetricScore:74.78 CsdWorkingBlockCount:0}
*/
