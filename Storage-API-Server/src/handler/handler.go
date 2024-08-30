package storagehandler

import (
	"net/http"
	"log"
	"encoding/json"
	"fmt"
	"bufio"
	"os/exec"

	"github.com/influxdata/influxdb/client/v2"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"

	data "api-server/src/struct"
)

var Mysql_db *sql.DB
var Influx_db client.HTTPClient

var(
	INFLUX_DB = "opencsd_management_platform"
)

//0
func StorageVolumeInfo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Storage API Server] StorageVolumeInfo Completed\n"))
}

func StorageVolumeAllocate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Storage API Server] StorageVolumeAllocate Completed\n"))
}

func StorageVolumeDeallocate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Storage API Server] StorageVolumeDeallocate Completed\n"))
}


//1
func ClusterNodeListHandler(w http.ResponseWriter, r *http.Request) {
	var result map[string][]data.NodeInfo
	var result_to_json []byte
	result = make(map[string][]data.NodeInfo)

	//MySQL Query
	rows, err := Mysql_db.Query("select c.cluster_name, n.node_name, n.node_status from cluster_info c, node_info n where c.cluster_id=n.cluster_id;") 
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cluster_name string
		var node_name string
		var node_status string

		err := rows.Scan(&cluster_name, &node_name, &node_status)
		if err != nil {
			log.Fatal(err)
		}

		node := data.NodeInfo{node_name, node_status}
		result[cluster_name] = append(result[cluster_name], node)
	
		result_to_json, err = json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("[1] : ",string(result_to_json))
	w.Write([]byte(string(result_to_json)+"\n"))
}
//2
func NodeStorageListHandler(w http.ResponseWriter, r *http.Request) {

	var result map[string][]data.StorageInfo
	var result_to_json []byte
	result = make(map[string][]data.StorageInfo)

	//MySQL Query
	rows, err := Mysql_db.Query("select n.node_name, s.storage_name, s.storage_status from node_info n, storage_info s where n.node_id=s.node_id;") 
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var node_name string
		var storage_name string
		var storage_status string

		err := rows.Scan(&node_name, &storage_name, &storage_status)
		if err != nil {
			log.Fatal(err)
		}

		storage := data.StorageInfo{storage_name, storage_status}
		result[node_name] = append(result[node_name], storage)
	
		result_to_json, err = json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("[2] : ",string(result_to_json))
	w.Write([]byte(string(result_to_json)+"\n"))
}

//3
func NodeDiskInfoHandler(w http.ResponseWriter, r *http.Request) {
	nodename := r.URL.Query().Get("nodename")
	datanum := r.URL.Query().Get("datanum")

	var result_to_json []byte
	var node_id string
	var node_info [][]interface{}

	disk_metric := []data.MetricValue{}

	//MySQL Query - Get 'node id'
	rows, err := Mysql_db.Query("select node_id from node_info where node_name = \""+nodename+"\";") 
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&node_id)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Influx Query
	q := client.Query{ 
		Command:  "select storage_total, storage_usage from node"+node_id+"_metric order by desc limit "+datanum,
		Database: INFLUX_DB,
	}

	if response, err := Influx_db.Query(q); err == nil && response.Error() == nil {
		node_series := response.Results[0].Series

		if len(node_series) > 0 {
			node_info = response.Results[0].Series[0].Values

			if len(node_series) > 0 {
				for i:=0 ; i<len(node_info);i++{
					tmp := data.MetricValue{}

					time := node_info[i][0]
					//disk_percent := node_info[i][1]
					disk_total := node_info[i][1]
					disk_usage := node_info[i][2]

					tmp.Time = fmt.Sprintf("%v", time)
					tmp.Capacity, _ = disk_total.(json.Number).Float64()
					tmp.Usage, _ = disk_usage.(json.Number).Float64()
					//tmp.Utilization, _ = disk_percent.(json.Number).Float64()

					disk_metric = append(disk_metric, tmp)
				}
			}
		}
	}

	result_to_json, _ = json.Marshal(disk_metric)
	fmt.Println("[3] : ",string(result_to_json))
	w.Write([]byte(string(result_to_json)+"\n"))
}

//4
func NodeStorageInfoHandler(w http.ResponseWriter, r *http.Request) {
	nodename := r.URL.Query().Get("nodename")

	var result_to_json []byte
	var node_id string
	storage_info := []data.StorageInfoDetailV1{}

	//MySQL Query - Get 'node id'
	rows, err := Mysql_db.Query("select node_id from node_info where node_name = \""+nodename+"\";") 
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&node_id)
		if err != nil {
			log.Fatal(err)
		}
	}

	//MySQL Query
	rows, err = Mysql_db.Query("select s.storage_id, s.storage_name, s.storage_type, s.storage_capacity, s.storage_status, c.cluster_name, n.node_name from cluster_info c, node_info n, storage_info s where c.cluster_id=n.cluster_id and n.node_id=s.node_id and n.node_name = \""+nodename+"\";") 
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var storage_id string
		var storage_name string
		var storage_type string
		var storage_capacity string
		var storage_status string
		var cluster_name string
		var node_name string

		err := rows.Scan(&storage_id, &storage_name, &storage_type, &storage_capacity, &storage_status, &cluster_name, &node_name)

		if err != nil {
			log.Fatal(err)
		}

		tmp := data.StorageInfoDetailV1{storage_name,storage_type,storage_capacity,storage_status,0,"",cluster_name,node_name}

		// Influx Query
		q := client.Query{ 
			Command:  "select score,grade from csd"+storage_id+"_metric order by desc limit 1",
			Database: INFLUX_DB,
		}
		
		if response, err := Influx_db.Query(q); err == nil && response.Error() == nil {
			node_series := response.Results[0].Series

			if len(node_series) > 0 {
				score := response.Results[0].Series[0].Values[0][1]
				grade := response.Results[0].Series[0].Values[0][2]

				tmp.StorageScore, _ = score.(json.Number).Float64()
				tmp.StorageGrade = fmt.Sprintf("%v", grade)
			}
		}

		storage_info = append(storage_info, tmp)
	}

	result_to_json, _ = json.Marshal(storage_info)
	fmt.Println("[4] : ",string(result_to_json))
	w.Write([]byte(string(result_to_json)+"\n"))
}
//5
func NodeMetricInfoHandler(w http.ResponseWriter, r *http.Request) {
	nodename := r.URL.Query().Get("nodename")
	time := r.URL.Query().Get("time")

	var result_to_json []byte
	var node_id string
	var node_info [][]interface{}

	node_metric := []data.NodeMetricInfo{}

	//MySQL Query - Get 'node id'
	rows, err := Mysql_db.Query("select node_id from node_info where node_name = \""+nodename+"\";") 
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&node_id)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Influx Query
	q := client.Query{ 
		Command:  "select cpu_total,cpu_usage,memory_total,memory_usage,storage_total,storage_usage from node"+node_id+"_metric where time > now() - "+time,
		Database: INFLUX_DB,
	}

	if response, err := Influx_db.Query(q); err == nil && response.Error() == nil {
		node_series := response.Results[0].Series

		if len(node_series) > 0 {
			node_info = response.Results[0].Series[0].Values

			if len(node_info) > 0 {
				for i:=0 ; i<len(node_info);i++{
					tmp := data.NodeMetricInfo{}

					time := node_info[i][0]
					//cpu_percent := node_info[i][1]
					cpu_total := node_info[i][1]
					cpu_usage := node_info[i][2]
					//memory_percent := node_info[i][4]
					memory_total := node_info[i][3]
					memory_usage := node_info[i][4]
					//disk_percent := node_info[i][7]
					disk_total := node_info[i][5]
					disk_usage := node_info[i][6]

					tmp.Time = fmt.Sprintf("%v", time)
					tmp.CPUCapacity, _ = cpu_total.(json.Number).Float64()
					tmp.CPUUsage, _ = cpu_usage.(json.Number).Float64()
					//tmp.CPUUtilization, _ = cpu_percent.(json.Number).Float64()
					tmp.MEMCapacity, _ = memory_total.(json.Number).Float64()
					tmp.MEMUsage, _ = memory_usage.(json.Number).Float64()
					//tmp.MEMUtilization, _ = memory_percent.(json.Number).Float64()
					tmp.DiskCapacity, _ = disk_total.(json.Number).Float64()
					tmp.DiskUsage, _ = disk_usage.(json.Number).Float64()
					//tmp.DiskUtilization, _ = disk_percent.(json.Number).Float64()

					node_metric = append(node_metric, tmp)
				}
			}
		}
	}

	result_to_json, _ = json.Marshal(node_metric)
	fmt.Println("[5] : ",string(result_to_json))
	w.Write([]byte(string(result_to_json)+"\n"))
	
}
//6
func StorageInfoHandler(w http.ResponseWriter, r *http.Request) {
	nodename := r.URL.Query().Get("nodename")
	storagename := r.URL.Query().Get("storagename")

	var storage_id string
	storage_info := data.StorageInfoDetailV2{}
	var result_to_json []byte

	//MySQL Query
	rows, err := Mysql_db.Query("select s.storage_id, s.storage_name, s.storage_type, s.storage_capacity,s.storage_status, s.csd_type, s.csd_ip, s.csd_port, c.cluster_name, n.node_name from cluster_info c, node_info n, storage_info s where c.cluster_id=n.cluster_id and n.node_id=s.node_id and n.node_name = \""+nodename+"\" and s.storage_name = \""+storagename+"\";") 
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var storage_name string
		var storage_type string
		var storage_capacity string
		var storage_status string
		var csd_type string
		var csd_ip string
		var csd_port string
		var cluster_name string
		var node_name string

		err := rows.Scan(&storage_id, &storage_name, &storage_type, &storage_capacity, &storage_status, &csd_type, &csd_ip, &csd_port, &cluster_name, &node_name)

		if err != nil {
			log.Fatal(err)
		}

		storage_info = data.StorageInfoDetailV2{storage_name,storage_type,storage_capacity,storage_status,csd_type,csd_ip,csd_port,cluster_name,node_name,0,"",0}
	}

	// Influx Query - working block 수는 추후 추가 예정
	q := client.Query{ 
		//Command:  "select score,grade,working_block from csd"+node_id+"_metric limit 1",
		Command:  "select score,grade from csd"+storage_id+"_metric limit 1",
		Database: INFLUX_DB,
	}
	
	if response, err := Influx_db.Query(q); err == nil && response.Error() == nil {
		node_series := response.Results[0].Series

		if len(node_series) > 0 {
			score := response.Results[0].Series[0].Values[0][1]
			grade := response.Results[0].Series[0].Values[0][2]
			//working_block := response.Results[0].Series[0].Values[0][3]

			storage_info.StorageScore, _ = score.(json.Number).Float64()
			storage_info.StorageGrade = fmt.Sprintf("%v", grade)
			//storage_info.WorkingBlocks, _ = working_block.(json.Number).Float64()
		}
	}

	result_to_json, _ = json.Marshal(storage_info)
	fmt.Println("[6] : ",string(result_to_json))
	w.Write([]byte(string(result_to_json)+"\n"))

}
//7
func CSDMetricInfoHandler(w http.ResponseWriter, r *http.Request) {
	nodename := r.URL.Query().Get("nodename")
	storagename := r.URL.Query().Get("storagename")
	datanum := r.URL.Query().Get("datanum")

	var storage_id string
	var result_to_json []byte
	var storage_info [][]interface{}

	csd_metric := []data.CSDMetricInfo{}

	//MySQL Query - Get 'storage id'
	rows, err := Mysql_db.Query("select storage_id from storage_info s, node_info n where n.node_name = \""+nodename+"\" and s.storage_name = \""+storagename+"\";") 
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&storage_id)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Influx Query
	q := client.Query{ 
		Command:  "select cpu_percent,cpu_total,cpu_usage,memory_percent,memory_total,memory_usage,disk_percent,disk_total,disk_usage from csd"+storage_id+"_metric order by desc limit "+datanum,
		Database: INFLUX_DB,
	}
	
	if response, err := Influx_db.Query(q); err == nil && response.Error() == nil {
		storage_series := response.Results[0].Series

		if len(storage_series) > 0 {
			storage_info = response.Results[0].Series[0].Values

			if len(storage_info) > 0 {
				for i:=0 ; i<len(storage_info);i++{
					tmp := data.CSDMetricInfo{}

					time := storage_info[i][0]
					cpu_percent := storage_info[i][1]
					cpu_total := storage_info[i][2]
					cpu_usage := storage_info[i][3]
					memory_percent := storage_info[i][4]
					memory_total := storage_info[i][5]
					memory_usage := storage_info[i][6]
					disk_percent := storage_info[i][7]
					disk_total := storage_info[i][8]
					disk_usage := storage_info[i][9]

					tmp.Time = fmt.Sprintf("%v", time)
					tmp.CPUCapacity, _ = cpu_total.(json.Number).Float64()
					tmp.CPUUsage, _ = cpu_usage.(json.Number).Float64()
					tmp.CPUUtilization, _ = cpu_percent.(json.Number).Float64()
					tmp.MEMCapacity, _ = memory_total.(json.Number).Float64()
					tmp.MEMUsage, _ = memory_usage.(json.Number).Float64()
					tmp.MEMUtilization, _ = memory_percent.(json.Number).Float64()
					tmp.DiskCapacity, _ = disk_total.(json.Number).Float64()
					tmp.DiskUsage, _ = disk_usage.(json.Number).Float64()
					tmp.DiskUtilization, _ = disk_percent.(json.Number).Float64()

					csd_metric = append(csd_metric, tmp)
				}
			}
		}
	}

	result_to_json, _ = json.Marshal(csd_metric)
	fmt.Println("[7] : ",string(result_to_json))
	w.Write([]byte(string(result_to_json)+"\n"))
}

//8
func CPUInfoHandler(w http.ResponseWriter, r *http.Request) {
	nodename := r.URL.Query().Get("nodename")
	storagename := r.URL.Query().Get("storagename")
	time := r.URL.Query().Get("time")

	var storage_id string
	var result_to_json []byte
	var storage_info [][]interface{}

	cpu_metric := []data.MetricValue{}

	//MySQL Query - Get 'storage id'
	rows, err := Mysql_db.Query("select storage_id from storage_info s, node_info n where n.node_name = \""+nodename+"\" and s.storage_name = \""+storagename+"\";") 
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&storage_id)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Influx Query
	q := client.Query{ 
		Command:  "select cpu_percent,cpu_total,cpu_usage from csd"+storage_id+"_metric where time > now() - "+time,
		Database: INFLUX_DB,
	}
	
	if response, err := Influx_db.Query(q); err == nil && response.Error() == nil {
		storage_series := response.Results[0].Series

		if len(storage_series) > 0 {
			storage_info = response.Results[0].Series[0].Values

			if len(storage_info) > 0 {
				for i:=0 ; i<len(storage_info);i++{
					tmp := data.MetricValue{}

					time := storage_info[i][0]
					cpu_percent := storage_info[i][1]
					cpu_total := storage_info[i][2]
					cpu_usage := storage_info[i][3]

					tmp.Time = fmt.Sprintf("%v", time)
					tmp.Capacity, _ = cpu_total.(json.Number).Float64()
					tmp.Usage, _ = cpu_usage.(json.Number).Float64()
					tmp.Utilization, _ = cpu_percent.(json.Number).Float64()

					cpu_metric = append(cpu_metric, tmp)
				}
			}
		}
	}

	result_to_json, _ = json.Marshal(cpu_metric)
	fmt.Println("[8] : ",string(result_to_json))
	w.Write([]byte(string(result_to_json)+"\n"))
}
//9
func MemInfoHandler(w http.ResponseWriter, r *http.Request) {
	nodename := r.URL.Query().Get("nodename")
	storagename := r.URL.Query().Get("storagename")
	time := r.URL.Query().Get("time")

	var storage_id string
	var result_to_json []byte
	var storage_info [][]interface{}

	mem_metric := []data.MetricValue{}

	//MySQL Query - Get 'storage id'
	rows, err := Mysql_db.Query("select storage_id from storage_info s, node_info n where n.node_name = \""+nodename+"\" and s.storage_name = \""+storagename+"\";") 
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&storage_id)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Influx Query
	q := client.Query{ 
		Command:  "select memory_percent,memory_total,memory_usage from csd"+storage_id+"_metric where time > now() - "+time,
		Database: INFLUX_DB,
	}
	
	if response, err := Influx_db.Query(q); err == nil && response.Error() == nil {
		storage_series := response.Results[0].Series

		if len(storage_series) > 0 {
			storage_info = response.Results[0].Series[0].Values

			if len(storage_info) > 0 {
				for i:=0 ; i<len(storage_info);i++{
					tmp := data.MetricValue{}

					time := storage_info[i][0]
					memory_percent := storage_info[i][1]
					memory_total := storage_info[i][2]
					memory_usage := storage_info[i][3]

					tmp.Time = fmt.Sprintf("%v", time)
					tmp.Capacity, _ = memory_total.(json.Number).Float64()
					tmp.Usage, _ = memory_usage.(json.Number).Float64()
					tmp.Utilization, _ = memory_percent.(json.Number).Float64()

					mem_metric = append(mem_metric, tmp)
				}
			}
		}
	}

	result_to_json, _ = json.Marshal(mem_metric)
	fmt.Println("[9] : ",string(result_to_json))
	w.Write([]byte(string(result_to_json)+"\n"))
}
//10
func NetInfoHandler(w http.ResponseWriter, r *http.Request) {
	nodename := r.URL.Query().Get("nodename")
	storagename := r.URL.Query().Get("storagename")
	time := r.URL.Query().Get("time")

	var storage_id string
	var result_to_json []byte
	var storage_info [][]interface{}

	net_metric := []data.NetMetricValue{}

	//MySQL Query - Get 'storage id'
	rows, err := Mysql_db.Query("select storage_id from storage_info s, node_info n where n.node_name = \""+nodename+"\" and s.storage_name = \""+storagename+"\";") 
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&storage_id)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Influx Query
	q := client.Query{ 
		Command:  "select network_bandwidth,network_rx_byte,network_tx_byte from csd"+storage_id+"_metric where time > now() - "+time,
		Database: INFLUX_DB,
	}
	
	if response, err := Influx_db.Query(q); err == nil && response.Error() == nil {
		storage_series := response.Results[0].Series

		if len(storage_series) > 0 {
			storage_info = response.Results[0].Series[0].Values

			if len(storage_info) > 0 {
				for i:=0 ; i<len(storage_info);i++{
					tmp := data.NetMetricValue{}

					time := storage_info[i][0]
					network_bandwidth := storage_info[i][1]
					network_rx_byte := storage_info[i][2]
					network_tx_byte := storage_info[i][3]

					tmp.Time = fmt.Sprintf("%v", time)
					tmp.Bandwidth, _ = network_bandwidth.(json.Number).Float64()
					tmp.RXByte, _ = network_rx_byte.(json.Number).Float64()
					tmp.TXByte, _ = network_tx_byte.(json.Number).Float64()

					net_metric = append(net_metric, tmp)
				}
			}
		}
	}

	result_to_json, _ = json.Marshal(net_metric)
	fmt.Println("[10] : ",string(result_to_json))
	w.Write([]byte(string(result_to_json)+"\n"))
}
//11
func DiskInfoHandler(w http.ResponseWriter, r *http.Request) {
	nodename := r.URL.Query().Get("nodename")
	storagename := r.URL.Query().Get("storagename")
	time := r.URL.Query().Get("time")

	var storage_id string
	var result_to_json []byte
	var storage_info [][]interface{}

	disk_metric := []data.MetricValue{}

	//MySQL Query - Get 'storage id'
	rows, err := Mysql_db.Query("select storage_id from storage_info s, node_info n where n.node_name = \""+nodename+"\" and s.storage_name = \""+storagename+"\";") 
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&storage_id)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Influx Query
	q := client.Query{ 
		Command:  "select disk_percent,disk_total,disk_usage from csd"+storage_id+"_metric where time > now() - "+time,
		Database: INFLUX_DB,
	}
	
	if response, err := Influx_db.Query(q); err == nil && response.Error() == nil {
		storage_series := response.Results[0].Series

		if len(storage_series) > 0 {
			storage_info = response.Results[0].Series[0].Values

			if len(storage_info) > 0 {
				for i:=0 ; i<len(storage_info);i++{
					tmp := data.MetricValue{}

					time := storage_info[i][0]
					disk_percent := storage_info[i][1]
					disk_total := storage_info[i][2]
					disk_usage := storage_info[i][3]

					tmp.Time = fmt.Sprintf("%v", time)
					tmp.Capacity, _ = disk_total.(json.Number).Float64()
					tmp.Usage, _ = disk_usage.(json.Number).Float64()
					tmp.Utilization, _ = disk_percent.(json.Number).Float64()

					disk_metric = append(disk_metric, tmp)
				}
			}
		}
	}

	result_to_json, _ = json.Marshal(disk_metric)
	fmt.Println("[11] : ",string(result_to_json))
	w.Write([]byte(string(result_to_json)+"\n"))
}

func CmdExec(cmdStr string) error{
	cmd := exec.Command("bash", "-c", cmdStr)
	stdoutReader, _ := cmd.StdoutPipe()
	stdoutScanner := bufio.NewScanner(stdoutReader)
	go func() {
		for stdoutScanner.Scan() {
			fmt.Println(stdoutScanner.Text())
		}
	}()
	stderrReader, _ := cmd.StderrPipe()
	stderrScanner := bufio.NewScanner(stderrReader)
	go func() {
		for stderrScanner.Scan() {
			fmt.Println(stderrScanner.Text())
		}
	}()
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error : %v \n", err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error: %v \n", err)
	}

	return nil
}