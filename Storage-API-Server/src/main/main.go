package main

import (
	"net/http"
	"fmt"
	"os"
	"log"

	"github.com/influxdata/influxdb/client/v2"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"

	sh "api-server/src/handler"
)

var(
	INFLUX_IP = os.Getenv("INFLUX_IP")
    INFLUX_PORT = os.Getenv("INFLUX_PORT")
    INFLUX_USERNAME = os.Getenv("INFLUX_USERNAME")
    INFLUX_PASSWORD = os.Getenv("INFLUX_PASSWORD")

	MYSQL_IP = os.Getenv("MYSQL_IP")
    MYSQL_PORT = os.Getenv("MYSQL_PORT")
    MYSQL_USERNAME = os.Getenv("MYSQL_USERNAME")
    MYSQL_PASSWORD = os.Getenv("MYSQL_PASSWORD")

	MYSQL_DB = "metric"
	INFLUX_DB = "opencsd_management_platform"
)

func main() {
	fmt.Println("[Storage API Server] Running..")
	var err error;

	//mysql Connection
	sh.Mysql_db, err = sql.Open("mysql", MYSQL_USERNAME+":"+MYSQL_PASSWORD+"@tcp("+MYSQL_IP+":"+MYSQL_PORT+")/"+MYSQL_DB)
	if err != nil {
		log.Fatal(err)
	}
	defer sh.Mysql_db.Close()

	err = sh.Mysql_db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("*** Connected to MySQL Database")

	//influx Connection
	sh.Influx_db, err = client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://" + INFLUX_IP + ":" + INFLUX_PORT, 
		Username: INFLUX_USERNAME,
		Password: INFLUX_PASSWORD,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sh.Influx_db.Close()
	fmt.Println("*** Connected to Influx Database")

	//handler

	//1. ClusterNodeList
	http.HandleFunc("/dashboard/cluster/nodelist", sh.ClusterNodeListHandler)
	//2. NodeStorageList
	http.HandleFunc("/dashboard/node/storagelist", sh.NodeStorageListHandler)
	http.HandleFunc("/storagepage/storage/storagelist", sh.NodeStorageListHandler)
	http.HandleFunc("/diskpage/disk/storagelist", sh.NodeStorageListHandler)
	//3. NodeDiskInfo
	http.HandleFunc("/dashboard/node/diskinfo", sh.NodeDiskInfoHandler)
	//4. NodeStorageInfo
	http.HandleFunc("/dashboard/storage/storageinfo", sh.NodeStorageInfoHandler)
	//5. NodeMetricInfo
	http.HandleFunc("/dashboard/node/metricinfo", sh.NodeMetricInfoHandler)
	//6. StorageInfo
	http.HandleFunc("/storagepage/storage/storageinfo", sh.StorageInfoHandler)
	//7. CSDMetricInfo
	http.HandleFunc("/storagepage/storage/csdmetricinfo", sh.CSDMetricInfoHandler)
	//8. CSDCpuInfo
	http.HandleFunc("/storagepage/storage/csdcpuinfo", sh.CPUInfoHandler)
	//9. CSDMemInfo
	http.HandleFunc("/storagepage/storage/csdmeminfo", sh.MemInfoHandler)
	//10. CSDNetInfo
	http.HandleFunc("/storagepage/storage/csdnetinfo", sh.NetInfoHandler)
	//11. CSDDiskInfo
	http.HandleFunc("/storagepage/storage/csddiskinfo", sh.DiskInfoHandler)

	http.ListenAndServe(":8000", nil)
}
