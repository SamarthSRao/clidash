package api

type TelemetryUpdate struct {
	ServiceID      string  `json:"service_id"`
	Operation      string  `json:"operation"`
	LatencyMS      float64 `json:"latency_ms"`
	RequestsPerSec int     `json:"rps"`
}

type PolicyUpdate struct {
	ServiceID   string `json:"service_id"`
	Consistency string `json:"consistency"` // STRONG, EVENTUAL, SESSION
}

type GlobalState struct {
	Services     map[string]TelemetryUpdate `json:"services"`
	LastDecision string                     `json:"last_decision"`
	Reward       float64                    `json:"reward"`
}
