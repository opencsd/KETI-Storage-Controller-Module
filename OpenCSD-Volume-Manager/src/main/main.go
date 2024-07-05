package main

import (
	"net/http"
	"fmt"

	handler "volume-manager/src/handler"
)

func main() {
	fmt.Println("[OpenCSD Volume Manager] Running..")

	//handler
	http.HandleFunc("/info", handler.StorageNodeInfo)
	http.HandleFunc("/allocate", handler.StorageVolumeRequest)

	http.ListenAndServe(":40306", nil)
}
