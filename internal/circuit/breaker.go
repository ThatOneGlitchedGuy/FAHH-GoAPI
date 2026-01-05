package circuit

import (
	"errors"
	"sync"
	"time"
)

type State string

const (
	StateClosed      State = "CLOSED"
	StateOpen        State = "OPEN"
	StateHalfOpen    State = "HALF_OPEN"
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type CircuitBreaker struct {
	maxFailures      int
	successThreshold int
	timeout          time.Duration
	state            State
	failures         int
	successes        int
	lastFailTime     time.Time
	mu               sync.RWMutex
}

func NewCircuitBreaker(maxFailures, successThreshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:      maxFailures,
		successThreshold: successThreshold,
		timeout:          timeout,
		state:            StateClosed,
	}
}

func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateOpen {
		if time.Since(cb.lastFailTime) > cb.timeout {
			cb.state = StateHalfOpen
			cb.successes = 0
		} else {
			return ErrCircuitOpen
		}
	}

	err := fn()

	if err != nil {
		cb.failures++
		cb.lastFailTime = time.Now()

		if cb.state == StateHalfOpen {
			cb.state = StateOpen
			cb.failures = 0
		} else if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
		}

		return err
	}

	cb.failures = 0

	if cb.state == StateHalfOpen {
		cb.successes++
		if cb.successes >= cb.successThreshold {
			cb.state = StateClosed
			cb.successes = 0
		}
	}

	return nil
}

func (cb *CircuitBreaker) GetState() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failures = 0
	cb.successes = 0
}
