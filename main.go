package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	port := flag.String("port", "8080", "port to listen on")
	flag.Parse()

	http.HandleFunc("/restart_port", restartPortHandler)
	http.HandleFunc("/start_port", startPortHandler)
	http.HandleFunc("/stop_port", stopPortHandler)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func restartPortHandler(w http.ResponseWriter, r *http.Request) {
	executeCommand(w, r, "restart")
}

func startPortHandler(w http.ResponseWriter, r *http.Request) {
	executeCommand(w, r, "start")
}

func stopPortHandler(w http.ResponseWriter, r *http.Request) {
	executeCommand(w, r, "stop")
}

func executeCommand(w http.ResponseWriter, r *http.Request, command string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query := r.URL.Query()
	port, ok := query["port"]

	if !ok || len(port[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Port not specified"}\n`)
		return
	}

	cmd := exec.Command("/usr/bin/supervisorctl", command, fmt.Sprintf("port%s", port[0]))

	if err := cmd.Run(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Failed to execute command: %s"}\n`, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Execution was successful"}`)
}
