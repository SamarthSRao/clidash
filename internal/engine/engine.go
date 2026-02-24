package engine

import (
	"clidash/internal/models"
	"math/rand"
	"time"
)

type Engine struct {
	Services []models.Microservice
	State    models.OptimizerState
}

func NewEngine() *Engine {
	return &Engine{
		Services: []models.Microservice{
			{ID: "1", Name: "Payment Gateway", Type: models.Payment, Consistency: models.Strong, IsCritical: true, AutoPilot: true},
			{ID: "2", Name: "Inventory Manager", Type: models.Inventory, Consistency: models.Strong, IsCritical: true, AutoPilot: true},
			{ID: "3", Name: "Product Catalog", Type: models.Catalog, Consistency: models.Eventual, IsCritical: false, AutoPilot: true},
			{ID: "4", Name: "Recommendation Engine", Type: models.Analytics, Consistency: models.Eventual, IsCritical: false, AutoPilot: true},
			{ID: "5", Name: "User Sessions", Type: models.UserCart, Consistency: models.Session, IsCritical: true, AutoPilot: true},
		},
		State: models.OptimizerState{
			Confidence: 0.95,
			Reward:     100.0,
		},
	}
}

func (e *Engine) Update() {
	for i := range e.Services {
		s := &e.Services[i]

		// Simulate traffic fluctuations
		s.RequestsPerSec = rand.Intn(500) + 50
		s.Latency = time.Duration(rand.Intn(20)+10) * time.Millisecond

		if s.AutoPilot {
			e.optimize(s)
		}
	}

	e.State.DecisionsCount++
}

func (e *Engine) optimize(s *models.Microservice) {
	// Simple Logic mimicking the DQN Agent's decisions

	// Rule 1: High Latency Detection
	if s.RequestsPerSec > 400 {
		if !s.IsCritical {
			// Relax non-critical services immediately
			if s.Consistency != models.Eventual {
				s.Consistency = models.Eventual
				e.State.LastDecision = "Relaxed " + s.Name + " due to high load"
				e.State.Reward += 10
				e.State.LatencyReduction += 15.5
			}
		} else {
			// Critical services under load - predictive relaxation of non-criticals
			// to protect these, but keep these Strong if possible unless extreme
			if s.RequestsPerSec > 480 {
				e.State.LastDecision = "Protecting " + s.Name + " by throttling non-criticals"
			}
		}
	}

	// Rule 2: Transaction mode detection
	if s.Type == models.Payment || s.Type == models.Inventory {
		if s.Consistency != models.Strong {
			s.Consistency = models.Strong
			e.State.LastDecision = "Ensuring Perfect mode for " + s.Name
			e.State.Reward += 5
		}
	} else if s.RequestsPerSec < 100 {
		// Low load, can afford higher consistency
		if s.Consistency == models.Eventual {
			s.Consistency = models.Session
			e.State.LastDecision = "Improving consistency for " + s.Name + " (Low Load)"
		}
	}
}
