# Integration Tests Makefile
# This Makefile provides convenient commands for running integration tests

.PHONY: help test test-short test-verbose test-parallel test-profile test-comment test-post test-reaction test-subscription test-clean test-ci setup-test-env clean-test-env

# Default target
help:
	@echo "Available targets:"
	@echo "  test              - Run all integration tests"
	@echo "  test-short        - Run integration tests in short mode"
	@echo "  test-verbose      - Run integration tests with verbose output"
	@echo "  test-parallel     - Run integration tests in parallel"
	@echo "  test-profile      - Run only profile repository tests"
	@echo "  test-comment      - Run only comment repository tests"
	@echo "  test-post         - Run only post repository tests"
	@echo "  test-reaction     - Run only reaction repository tests"
	@echo "  test-subscription - Run only subscription repository tests"
	@echo "  test-clean        - Run tests with clean environment"
	@echo "  test-ci           - Run tests in CI mode"
	@echo "  setup-test-env    - Setup test environment only"
	@echo "  clean-test-env    - Clean test environment"
	@echo "  test-benchmark    - Run benchmark tests"
	@echo "  test-coverage     - Run tests with coverage report"

# Go test configuration
GO_TEST_CMD = go test
GO_TEST_FLAGS = -v
GO_TEST_TIMEOUT = 10m
GO_TEST_PARALLEL = 4

# Test package
TEST_PACKAGE = .

# Environment variables
export GOMAXPROCS ?= 4
export CGO_ENABLED = 1

# Run all integration tests
test:
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -timeout $(GO_TEST_TIMEOUT) $(TEST_PACKAGE)

# Run tests in short mode (skips long-running tests)
test-short:
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -short -timeout 5m $(TEST_PACKAGE)

# Run tests with verbose output
test-verbose:
	$(GO_TEST_CMD) -v -timeout $(GO_TEST_TIMEOUT) $(TEST_PACKAGE)

# Run tests in parallel
test-parallel:
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -parallel $(GO_TEST_PARALLEL) -timeout $(GO_TEST_TIMEOUT) $(TEST_PACKAGE)

# Run specific test suites
test-profile:
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -run TestProfileRepositoryTestSuite -timeout 5m $(TEST_PACKAGE)

test-comment:
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -run TestCommentRepositoryTestSuite -timeout 5m $(TEST_PACKAGE)

test-post:
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -run TestPostRepositoryTestSuite -timeout 5m $(TEST_PACKAGE)

test-reaction:
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -run TestReactionRepositoryTestSuite -timeout 5m $(TEST_PACKAGE)

test-subscription:
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -run TestSubscriptionRepositoryTestSuite -timeout 5m $(TEST_PACKAGE)

# Run tests with clean environment
test-clean:
	docker system prune -f
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -timeout $(GO_TEST_TIMEOUT) $(TEST_PACKAGE)

# Run tests in CI mode
test-ci:
	@echo "Running integration tests in CI mode..."
	CI=true $(GO_TEST_CMD) -v -timeout 15m -parallel 2 $(TEST_PACKAGE)

# Setup test environment only (useful for debugging)
setup-test-env:
	@echo "Setting up test environment..."
	$(GO_TEST_CMD) -run TestMain -timeout 5m $(TEST_PACKAGE) || true

# Clean test environment
clean-test-env:
	@echo "Cleaning test environment..."
	docker container prune -f
	docker volume prune -f
	docker network prune -f

# Run benchmark tests
test-benchmark:
	$(GO_TEST_CMD) -bench=. -benchmem -timeout 15m $(TEST_PACKAGE)

# Run tests with coverage report
test-coverage:
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -timeout $(GO_TEST_TIMEOUT) -coverprofile=coverage.out $(TEST_PACKAGE)
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -race -timeout $(GO_TEST_TIMEOUT) $(TEST_PACKAGE)

# Run tests with memory profiling
test-memprofile:
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -memprofile=mem.prof -timeout $(GO_TEST_TIMEOUT) $(TEST_PACKAGE)
	go tool pprof mem.prof

# Run tests with CPU profiling
test-cpuprofile:
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -cpuprofile=cpu.prof -timeout $(GO_TEST_TIMEOUT) $(TEST_PACKAGE)
	go tool pprof cpu.prof

# Run tests with specific test name pattern
test-pattern:
	@if [ -z "$(PATTERN)" ]; then echo "Usage: make test-pattern PATTERN=TestName"; exit 1; fi
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -run $(PATTERN) -timeout 5m $(TEST_PACKAGE)

# Run tests with specific timeout
test-timeout:
	@if [ -z "$(TIMEOUT)" ]; then echo "Usage: make test-timeout TIMEOUT=5m"; exit 1; fi
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -timeout $(TIMEOUT) $(TEST_PACKAGE)

# Run tests with JSON output (useful for CI/CD)
test-json:
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -json -timeout $(GO_TEST_TIMEOUT) $(TEST_PACKAGE)

# Run tests with count (useful for flaky test detection)
test-count:
	@if [ -z "$(COUNT)" ]; then echo "Usage: make test-count COUNT=10"; exit 1; fi
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -count $(COUNT) -timeout $(GO_TEST_TIMEOUT) $(TEST_PACKAGE)

# Run tests with specific parallel setting
test-parallel-custom:
	@if [ -z "$(PARALLEL)" ]; then echo "Usage: make test-parallel-custom PARALLEL=8"; exit 1; fi
	$(GO_TEST_CMD) $(GO_TEST_FLAGS) -parallel $(PARALLEL) -timeout $(GO_TEST_TIMEOUT) $(TEST_PACKAGE)

# Debug tests (run with dlv debugger)
test-debug:
	dlv test $(TEST_PACKAGE) -- -test.v

# Check test dependencies
check-deps:
	@echo "Checking test dependencies..."
	@command -v docker >/dev/null 2>&1 || { echo "Docker is required but not installed. Aborting." >&2; exit 1; }
	@docker --version
	@go version
	@echo "All dependencies are available"

# Lint test files
lint-tests:
	golangci-lint run $(TEST_PACKAGE)/...

# Format test files
format-tests:
	go fmt $(TEST_PACKAGE)/...

# Run all quality checks
quality: format-tests lint-tests test-race test-coverage

# Fast test run (optimized for development)
test-fast:
	$(GO_TEST_CMD) -short -parallel $(GO_TEST_PARALLEL) -timeout 3m $(TEST_PACKAGE)

# Comprehensive test run (for CI/CD)
test-comprehensive: check-deps test-race test-coverage

# Print test configuration
show-config:
	@echo "Test Configuration:"
	@echo "  - Command: $(GO_TEST_CMD)"
	@echo "  - Flags: $(GO_TEST_FLAGS)"
	@echo "  - Timeout: $(GO_TEST_TIMEOUT)"
	@echo "  - Parallel: $(GO_TEST_PARALLEL)"
	@echo "  - Package: $(TEST_PACKAGE)"
	@echo "  - GOMAXPROCS: $(GOMAXPROCS)"
	@echo "  - CGO_ENABLED: $(CGO_ENABLED)"

# Generate test report
test-report:
	$(GO_TEST_CMD) -v -json -timeout $(GO_TEST_TIMEOUT) $(TEST_PACKAGE) | tee test-results.json
	@echo "Test report generated: test-results.json"

# Stress test (run tests multiple times to detect flaky tests)
test-stress:
	@echo "Running stress test (10 iterations)..."
	for i in {1..10}; do \
		echo "=== Iteration $$i ==="; \
		$(GO_TEST_CMD) $(GO_TEST_FLAGS) -timeout 5m $(TEST_PACKAGE) || exit 1; \
	done
	@echo "Stress test completed successfully"

# Clean up generated files
clean:
	rm -f coverage.out coverage.html mem.prof cpu.prof test-results.json

# Show help by default
all: help
