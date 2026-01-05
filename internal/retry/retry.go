package retry

import (
	"math"
	"math/rand"
	"time"
)

type RetryPolicy struct {
	maxAttempts      int
	initialDelay     time.Duration
	maxDelay         time.Duration
	backoffMultiplier float64
	jitter           bool
}

func NewRetryPolicy(maxAttempts int, initialDelay, maxDelay time.Duration) *RetryPolicy {
	return &RetryPolicy{
		maxAttempts:       maxAttempts,
		initialDelay:      initialDelay,
		maxDelay:          maxDelay,
		backoffMultiplier: 2.0,
		jitter:            true,
	}
}

func (rp *RetryPolicy) WithBackoffMultiplier(multiplier float64) *RetryPolicy {
	rp.backoffMultiplier = multiplier
	return rp
}

func (rp *RetryPolicy) WithJitter(jitter bool) *RetryPolicy {
	rp.jitter = jitter
	return rp
}

func (rp *RetryPolicy) Execute(fn func() error) error {
	var lastErr error

	for attempt := 0; attempt < rp.maxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		if attempt < rp.maxAttempts-1 {
			delay := rp.calculateDelay(attempt)
			time.Sleep(delay)
		}
	}

	return lastErr
}

func (rp *RetryPolicy) calculateDelay(attempt int) time.Duration {
	exponentialDelay := time.Duration(
		float64(rp.initialDelay) * math.Pow(rp.backoffMultiplier, float64(attempt)),
	)

	if exponentialDelay > rp.maxDelay {
		exponentialDelay = rp.maxDelay
	}

	if rp.jitter {
		jitterValue := time.Duration(rand.Int63n(int64(exponentialDelay)))
		return jitterValue
	}

	return exponentialDelay
}

type RetryableFunc func() error

type AsyncRetry struct {
	policy *RetryPolicy
	result chan error
}

func NewAsyncRetry(policy *RetryPolicy) *AsyncRetry {
	return &AsyncRetry{
		policy: policy,
		result: make(chan error, 1),
	}
}

func (ar *AsyncRetry) Execute(fn RetryableFunc) chan error {
	go func() {
		ar.result <- ar.policy.Execute(fn)
	}()
	return ar.result
}
