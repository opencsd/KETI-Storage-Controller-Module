package main

import (
	"net/http"
	"fmt"

	handler "storage-manager/src/handler"
)

func main() {
	fmt.Println("[OpenCSD Storage Manager] Running...")

	//handler
	http.HandleFunc("/info", handler.StorageVolumeInfo)
	http.HandleFunc("/allocate", handler.StorageVolumeAllocate)

	http.ListenAndServe(":40307", nil)	
}
