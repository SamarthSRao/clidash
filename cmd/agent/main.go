package main

import (
	"bytes"
	"clidash/pkg/api"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	serviceID := os.Getenv("SERVICE_ID")
	if serviceID == "" {
		serviceID = "inventory-service"
	}

	optimizerURL := "http://localhost:8080/metrics"

	fmt.Printf("Starting CLIDash Agent for %s...\n", serviceID)

	// In a real Jaeger sidecar, this would listen on UDP.
	// For this demo, we simulate a microservice loop.
	go func() {
		for {
			data := api.TelemetryUpdate{
				ServiceID:      serviceID,
				Operation:      "DB_READ",
				LatencyMS:      float64(rand.Intn(50) + 10),
				RequestsPerSec: rand.Intn(100) + 50,
			}

			payload, _ := json.Marshal(data)
			resp, err := http.Post(optimizerURL, "application/json", bytes.NewBuffer(payload))
			if err == nil {
				resp.Body.Close()
			}

			time.Sleep(2 * time.Second)
		}
	}()

	// Simulate listening for Policy Updates from the Optimizer
	fmt.Println("Listening for policy updates on :5775...")
	http.HandleFunc("/policy", func(w http.ResponseWriter, r *http.Request) {
		var policy api.PolicyUpdate
		if err := json.NewDecoder(r.Body).Decode(&policy); err == nil {
			fmt.Printf(">>> POLICY UPDATE RECEIVED: %s is now %s\n", policy.ServiceID, policy.Consistency)
		}
	})

	log.Fatal(http.ListenAndServe(":5775", nil))
}
