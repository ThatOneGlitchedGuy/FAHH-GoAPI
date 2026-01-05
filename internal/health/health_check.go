package health

import (
	"sync"
	"time"
)

type ComponentStatus string

const (
	StatusHealthy   ComponentStatus = "HEALTHY"
	StatusDegraded  ComponentStatus = "DEGRADED"
	StatusUnhealthy ComponentStatus = "UNHEALTHY"
	StatusUnknown   ComponentStatus = "UNKNOWN"
)

type ComponentHealth struct {
	Name        string          `json:"name"`
	Status      ComponentStatus `json:"status"`
	Message     string          `json:"message,omitempty"`
	LastChecked time.Time       `json:"last_checked"`
	ResponseTime time.Duration  `json:"response_time_ms"`
}

type HealthChecker interface {
	Check() *ComponentHealth
}

type HealthCheckResult struct {
	Status        ComponentStatus              `json:"status"`
	Timestamp     time.Time                    `json:"timestamp"`
	Components    map[string]*ComponentHealth  `json:"components"`
	OverallHealth string                       `json:"overall_health"`
}

type HealthManager struct {
	checkers map[string]HealthChecker
	mu       sync.RWMutex
}

func NewHealthManager() *HealthManager {
	return &HealthManager{
		checkers: make(map[string]HealthChecker),
	}
}

func (hm *HealthManager) Register(name string, checker HealthChecker) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hm.checkers[name] = checker
}

func (hm *HealthManager) GetHealth() *HealthCheckResult {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	result := &HealthCheckResult{
		Timestamp:  time.Now(),
		Components: make(map[string]*ComponentHealth),
	}

	unhealthyCount := 0
	degradedCount := 0

	for name, checker := range hm.checkers {
		componentHealth := checker.Check()
		result.Components[name] = componentHealth

		if componentHealth.Status == StatusUnhealthy {
			unhealthyCount++
		} else if componentHealth.Status == StatusDegraded {
			degradedCount++
		}
	}

	if unhealthyCount > 0 {
		result.Status = StatusUnhealthy
		result.OverallHealth = "CRITICAL"
	} else if degradedCount > 0 {
		result.Status = StatusDegraded
		result.OverallHealth = "DEGRADED"
	} else {
		result.Status = StatusHealthy
		result.OverallHealth = "HEALTHY"
	}

	return result
}

type SimpleHealthChecker struct {
	name    string
	healthy bool
	message string
}

func NewSimpleHealthChecker(name string) *SimpleHealthChecker {
	return &SimpleHealthChecker{
		name:    name,
		healthy: true,
		message: "OK",
	}
}

func (shc *SimpleHealthChecker) SetStatus(healthy bool, message string) {
	shc.healthy = healthy
	shc.message = message
}

func (shc *SimpleHealthChecker) Check() *ComponentHealth {
	start := time.Now()
	status := StatusHealthy
	if !shc.healthy {
		status = StatusUnhealthy
	}

	return &ComponentHealth{
		Name:        shc.name,
		Status:      status,
		Message:     shc.message,
		LastChecked: time.Now(),
		ResponseTime: time.Since(start),
	}
}
