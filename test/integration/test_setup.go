package integration

import (
	"context"
	"testing"
	"time"
)

// TestRunner provides methods for running tests efficiently
type TestRunner struct {
	timeGen     *TimeGenerator
	dataBuilder *TestDataBuilder
}

// NewTestRunner creates a new optimized test runner
func NewTestRunner() *TestRunner {
	return &TestRunner{
		timeGen:     NewTimeGenerator(),
		dataBuilder: NewTestDataBuilder(),
	}
}

// RunConcurrently runs multiple operations concurrently and waits for completion
func (otr *TestRunner) RunConcurrently(t *testing.T, operations ...func() error) {
	t.Helper()

	done := make(chan error, len(operations))

	for _, op := range operations {
		go func(operation func() error) {
			done <- operation()
		}(op)
	}

	// Wait for all operations to complete
	for i := 0; i < len(operations); i++ {
		if err := <-done; err != nil {
			t.Errorf("Concurrent operation %d failed: %v", i+1, err)
		}
	}
}

// SkipIfShort skips the test if running in short mode
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
}

// WaitForEventualConsistency waits for eventual consistency in the database
func WaitForEventualConsistency(ctx context.Context, checkFn func() (bool, error)) error {
	return RetryOperation(func() error {
		consistent, err := checkFn()
		if err != nil {
			return err
		}
		if !consistent {
			return &TemporaryError{Message: "not yet consistent"}
		}
		return nil
	}, 10, 50*time.Millisecond)
}

// TemporaryError represents a temporary error that should be retried
type TemporaryError struct {
	Message string
}

func (e *TemporaryError) Error() string {
	return e.Message
}
