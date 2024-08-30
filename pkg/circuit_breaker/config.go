package circuit_breaker

import (
	"fmt"
	"strings"
)

// CircuitBreakerConfig is the configuration that all the circuit breakers will use
type CircuitBreakerConfig struct {
	// SampleRate is the time interval (in seconds) at which the data source
	// is polled to determine the number successful and failed requests
	SampleRate uint64 `json:"sample_rate"`

	// ErrorTimeout is the time (in seconds) after which a circuit breaker goes
	// into the half-open state from the open state
	ErrorTimeout uint64 `json:"error_timeout"`

	// FailureThreshold is the % of failed requests in the observability window
	// after which the breaker will go into the open state
	FailureThreshold float64 `json:"failure_threshold"`

	// FailureCount total number of failed requests in the observability window
	FailureCount uint64 `json:"failure_count"`

	// SuccessThreshold is the % of successful requests in the observability window
	// after which a circuit breaker in the half-open state will go into the closed state
	SuccessThreshold uint64 `json:"success_threshold"`

	// ObservabilityWindow is how far back in time (in minutes) the data source is
	// polled when determining the number successful and failed requests
	ObservabilityWindow uint64 `json:"observability_window"`

	// NotificationThresholds These are the error counts after which we will send out notifications.
	NotificationThresholds []uint64 `json:"notification_thresholds"`

	// ConsecutiveFailureThreshold determines when we ultimately disable the endpoint.
	// E.g., after 10 consecutive transitions from half-open → open we should disable it.
	ConsecutiveFailureThreshold uint64 `json:"consecutive_failure_threshold"`
}

func (c *CircuitBreakerConfig) Validate() error {
	var errs strings.Builder

	if c.SampleRate == 0 {
		errs.WriteString("SampleRate must be greater than 0")
		errs.WriteString("; ")
	}

	if c.ErrorTimeout == 0 {
		errs.WriteString("ErrorTimeout must be greater than 0")
		errs.WriteString("; ")
	}

	if c.FailureThreshold < 0 || c.FailureThreshold > 1 {
		errs.WriteString("FailureThreshold must be between 0 and 1")
		errs.WriteString("; ")
	}

	if c.FailureCount == 0 {
		errs.WriteString("FailureCount must be greater than 0")
		errs.WriteString("; ")
	}

	if c.SuccessThreshold == 0 {
		errs.WriteString("SuccessThreshold must be greater than 0")
		errs.WriteString("; ")
	}

	if c.ObservabilityWindow == 0 {
		errs.WriteString("ObservabilityWindow must be greater than 0")
		errs.WriteString("; ")
	}

	// ObservabilityWindow is in minutes and SampleRate is in seconds
	if (c.ObservabilityWindow * 60) <= c.SampleRate {
		errs.WriteString("ObservabilityWindow must be greater than the SampleRate")
		errs.WriteString("; ")
	}

	if len(c.NotificationThresholds) == 0 {
		errs.WriteString("NotificationThresholds must contain at least one threshold")
		errs.WriteString("; ")
	} else {
		for i := 0; i < len(c.NotificationThresholds); i++ {
			if c.NotificationThresholds[i] == 0 {
				errs.WriteString(fmt.Sprintf("Notification thresholds at index [%d] = %d must be greater than 0", i, c.NotificationThresholds[i]))
				errs.WriteString("; ")
			}
		}

		for i := 0; i < len(c.NotificationThresholds)-1; i++ {
			if c.NotificationThresholds[i] >= c.NotificationThresholds[i+1] {
				errs.WriteString("NotificationThresholds should be in ascending order")
				errs.WriteString("; ")
			}
		}
	}

	if c.ConsecutiveFailureThreshold == 0 {
		errs.WriteString("ConsecutiveFailureThreshold must be greater than 0")
		errs.WriteString("; ")
	}

	if errs.Len() > 0 {
		return fmt.Errorf("config validation failed with errors: %s", errs.String())
	}

	return nil
}
