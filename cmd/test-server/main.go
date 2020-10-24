package main

import (
	"net/http"

	"github.com/calebamiles/example-figlet-service/service"
)

func main() {
	// Don't use Cadence backend
	http.HandleFunc("/figlet", service.HandleFigletizeTextDirect)
	http.HandleFunc("/healthz", service.HandleGetHealthz)

	http.ListenAndServe(":8091", nil)
}
