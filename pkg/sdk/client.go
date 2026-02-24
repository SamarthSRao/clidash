package sdk

import (
	"bytes"
	"clidash/pkg/api"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	serviceID    string
	agentURL     string
	mu           sync.RWMutex
	currentLevel string
}

// NewClient creates a new SDK client that communicates with the local sidecar agent.
func NewClient(serviceID string, agentAddr string) *Client {
	c := &Client{
		serviceID:    serviceID,
		agentURL:     fmt.Sprintf("http://%s", agentAddr),
		currentLevel: "STRONG", // Default
	}

	// In a real implementation, the Agent would push to the SDK.
	// Here we simulate a background pull for the latest policy.
	go c.pollPolicy()

	return c
}

func (c *Client) pollPolicy() {
	for {
		// In a real sidecar, this might be a long-poll or gRPC stream
		time.Sleep(5 * time.Second)
	}
}

// RecordOperation sends telemetry to the sidecar agent (Jaeger-style push)
func (c *Client) RecordOperation(op string, latencyMs float64) {
	data := api.TelemetryUpdate{
		ServiceID:      c.serviceID,
		Operation:      op,
		LatencyMS:      latencyMs,
		RequestsPerSec: 1, // Simplified for SDK
	}

	payload, _ := json.Marshal(data)
	// Sending to local sidecar (very low latency)
	http.Post(c.agentURL+"/metrics", "application/json", bytes.NewBuffer(payload))
}

// GetConsistency returns the currently mandated consistency level from the RL Optimizer
func (c *Client) GetConsistency() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentLevel
}
