package kafka

import (
	"time"
)

// RetryPolicy defines retry behavior
type RetryPolicy struct {
	MaxRetries int
	Backoff    time.Duration
}

// Retry executes fn with retry + backoff
func Retry(policy RetryPolicy, fn func() error) error {
	var err error

	for attempt := 1; attempt <= policy.MaxRetries; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		time.Sleep(policy.Backoff)
	}

	return err
}
