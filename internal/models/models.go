package models

import "time"

type ConsistencyLevel string

const (
	Strong   ConsistencyLevel = "STRONG"
	Eventual ConsistencyLevel = "EVENTUAL"
	Session  ConsistencyLevel = "SESSION"
)

type ServiceType string

const (
	Payment    ServiceType = "PAYMENT"
	Inventory  ServiceType = "INVENTORY"
	Catalog    ServiceType = "CATALOG"
	UserCart   ServiceType = "USER_CART"
	Analytics  ServiceType = "ANALYTICS"
)

type Microservice struct {
	ID               string
	Name             string
	Type             ServiceType
	Consistency      ConsistencyLevel
	Latency          time.Duration
	RequestsPerSec   int
	ErrorRate        float64
	SLACompliance    float64
	IsCritical       bool
	AutoPilot        bool
}

type OptimizerState struct {
	Reward           float64
	Confidence       float64
	DecisionsCount   int
	LastDecision     string
	LatencyReduction float64
}
