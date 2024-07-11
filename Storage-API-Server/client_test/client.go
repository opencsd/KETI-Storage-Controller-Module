package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
)

var (
	HOST = "10.0.4.85:30008"
    NODENAME = "storage-node"
    CSDNAME = "csd4"
    TIMEPERIOD = "1m"
)

func requestURL(id string, url string) {
    fmt.Println("[ ID  : ", id , "]")
    fmt.Println("[ URL : ", url , "]")
    req, _ := http.NewRequest("GET", url, nil)

    client := &http.Client{}
    res, err := client.Do(req)

    if err != nil {
        panic(err)
    }

    defer res.Body.Close()
    thisBody, err := ioutil.ReadAll(res.Body)

    if err != nil {
        panic(err)
    }

    fmt.Println(string(thisBody))
}

func main() {
    requestURL("ClusterNodeList", "http://"+HOST+"/dashboard/cluster/nodelist")
    requestURL("NodeStorageList", "http://"+HOST+"/dashboard/node/storagelist")
    requestURL("NodeStorageList", "http://"+HOST+"/storagepage/storage/storagelist")
    requestURL("NodeStorageList", "http://"+HOST+"/diskpage/disk/storagelist")
    requestURL("NodeDiskInfo", "http://"+HOST+"/dashboard/node/diskinfo?nodename="+NODENAME+"&datanum=1")
    requestURL("NodeStorageInfo", "http://"+HOST+"/dashboard/storage/storageinfo?nodename="+NODENAME)
    requestURL("NodeMetricInfo", "http://"+HOST+"/dashboard/node/metricinfo?nodename="+NODENAME+"&time="+TIMEPERIOD)
    requestURL("StorageInfo", "http://"+HOST+"/storagepage/storage/storageinfo?storagename="+CSDNAME+"&nodename="+NODENAME)
    requestURL("CSDMetricInfo", "http://"+HOST+"/storagepage/storage/csdmetricinfo?storagename="+CSDNAME+"&nodename="+NODENAME+"&datanum=1")
    requestURL("CSDCpuInfo", "http://"+HOST+"/storagepage/storage/csdcpuinfo?storagename="+CSDNAME+"&nodename="+NODENAME+"&time="+TIMEPERIOD)
    requestURL("CSDMemInfo", "http://"+HOST+"/storagepage/storage/csdmeminfo?storagename="+CSDNAME+"&nodename="+NODENAME+"&time="+TIMEPERIOD)
    requestURL("CSDNetInfo", "http://"+HOST+"/storagepage/storage/csdnetinfo?storagename="+CSDNAME+"&nodename="+NODENAME+"&time="+TIMEPERIOD)
    requestURL("CSDDiskInfo", "http://"+HOST+"/storagepage/storage/csddiskinfo?storagename="+CSDNAME+"&nodename="+NODENAME+"&time="+TIMEPERIOD)
}

