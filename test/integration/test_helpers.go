package integration

import (
	"log"
	"os"
	"time"
)

// TimeGenerator provides methods for generating timestamps in tests
type TimeGenerator struct {
	baseTime time.Time
	counter  int64
}

// NewTimeGenerator creates a new time generator starting from a base time
func NewTimeGenerator() *TimeGenerator {
	return &TimeGenerator{
		baseTime: time.Now().Add(-24 * time.Hour), // Start from yesterday to avoid conflicts
		counter:  0,
	}
}

// Next returns the next timestamp, incrementing by 1 millisecond each time
func (tg *TimeGenerator) Next() time.Time {
	tg.counter++
	return tg.baseTime.Add(time.Duration(tg.counter) * time.Millisecond)
}

// RetryOperation retries an operation with exponential backoff
func RetryOperation(operation func() error, maxRetries int, baseDelay time.Duration) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := operation(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if i < maxRetries-1 {
			delay := baseDelay * time.Duration(1<<uint(i)) // Exponential backoff
			if delay > 5*time.Second {
				delay = 5 * time.Second // Cap at 5 seconds
			}
			time.Sleep(delay)
		}
	}
	return lastErr
}

// CreateTestData creates test data structures for common test scenarios
type TestDataBuilder struct {
	timeGen *TimeGenerator
}

// NewTestDataBuilder creates a new test data builder
func NewTestDataBuilder() *TestDataBuilder {
	return &TestDataBuilder{
		timeGen: NewTimeGenerator(),
	}
}

// UserID generates a test user ID
func (tdb *TestDataBuilder) UserID(suffix string) string {
	return "user_" + suffix
}

// Email generates a test email
func (tdb *TestDataBuilder) Email(username string) string {
	return username + "@example.com"
}

// TestProfile represents test profile data
type TestProfile struct {
	UserID    string
	Email     string
	FirstName string
	LastName  string
}

// CreateTestProfile creates a test profile with generated data
func (tdb *TestDataBuilder) CreateTestProfile(suffix string) TestProfile {
	return TestProfile{
		UserID:    tdb.UserID(suffix),
		Email:     tdb.Email("user" + suffix),
		FirstName: "FirstName" + suffix,
		LastName:  "LastName" + suffix,
	}
}

// Helper function to check if running in CI environment
func isCI() bool {
	return os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" || os.Getenv("GITLAB_CI") != ""
}

// Helper function to check if running in short mode
func isShortMode() bool {
	// Cannot use testing.Short() in init function, so check flag manually
	for _, arg := range os.Args {
		if arg == "-test.short" || arg == "--test.short" {
			return true
		}
	}
	return false
}

// LogTestEnvironment logs information about test environment
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Integration test environment:")
	log.Printf("  - CI: %v", isCI())
	log.Printf("  - Short mode: %v", isShortMode())
	log.Printf("  - Go version: %s", os.Getenv("GOVERSION"))

	if isShortMode() {
		log.Println("  - Running in short mode - some tests may be skipped")
	}

	if isCI() {
		log.Println("  - Running in CI environment - using optimized settings")
	}
}
