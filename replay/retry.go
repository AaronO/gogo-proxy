package proxy

import (
	"fmt"
	"time"
)

// Retries calling a function N times, returns last error if all fail
// otherwise nil
func retryWait(f func() error, n int, sleep time.Duration) error {
	var err error

	// Failed
	if n < 1 {
		return fmt.Errorf("retry count must be great than 1")
	}

	for i := 0; i < n; i++ {
		// Test for success
		if err = f(); err == nil {
			return nil
		}
		// Wait before retrying again
		time.Sleep(sleep)
	}

	// We tried but failed every single time
	return err
}

func retry(f func() error, n int) error {
	return retryWait(f, n, time.Duration(0))
}
