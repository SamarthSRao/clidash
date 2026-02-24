package main

import (
	"clidash/internal/engine"
	"clidash/pkg/api"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type OptimizerServer struct {
	mu     sync.RWMutex
	state  api.GlobalState
	engine *engine.Engine
}

func NewOptimizerServer() *OptimizerServer {
	return &OptimizerServer{
		state: api.GlobalState{
			Services: make(map[string]api.TelemetryUpdate),
		},
		engine: engine.NewEngine(),
	}
}

func (s *OptimizerServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	var update api.TelemetryUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	s.mu.Lock()
	s.state.Services[update.ServiceID] = update
	s.engine.Update() // Trigger RL Logic
	s.state.LastDecision = s.engine.State.LastDecision
	s.state.Reward = s.engine.State.Reward
	s.mu.Unlock()

	// In a real Jaeger-style push, the agent would poll or we'd stream.
	// For now, we just acknowledge.
	w.WriteHeader(200)
}

func (s *OptimizerServer) handleState(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	json.NewEncoder(w).Encode(s.state)
}

func main() {
	srv := NewOptimizerServer()

	http.HandleFunc("/metrics", srv.handleMetrics)
	http.HandleFunc("/state", srv.handleState)

	fmt.Println("CLIDash Optimizer (The Brain) starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
