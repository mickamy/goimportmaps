APP_NAME = goimportmaps
VERSION ?= dev
BUILD_DIR = bin
GORELEASER ?= go tool goreleaser
VERSION_VARIABLE = github.com/mickamy/goimportmaps/internal/cli/version/version.version

.PHONY: all build install uninstall clean version test fmt

all: build

build:
	@echo "🔨 Building $(APP_NAME)..."
	go build -ldflags "-X $(VERSION_VARIABLE)=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME) ./cmd/goimportmaps

install:
	@echo "📦 Installing $(APP_NAME)..."
	go install -ldflags "-X $(VERSION_VARIABLE)=$(VERSION)" ./cmd/goimportmaps

uninstall:
	@echo "🗑️ Uninstalling $(APP_NAME)..."
	@bin_dir=$$(go env GOBIN); \
	if [ -z "$$bin_dir" ]; then \
		bin_dir=$$(go env GOPATH)/bin; \
	fi; \
	echo "Removing $$bin_dir/$(APP_NAME)"; \
	rm -f $$bin_dir/$(APP_NAME)

clean:
	@echo "🧹 Cleaning up..."
	rm -rf $(BUILD_DIR)

version:
	@echo "🔖 Version: $(VERSION)"

test:
	@echo "🧪 Running tests..."
	go test ./...

fmt:
	@echo "📝 Formatting code..."
	gofmt -w -l .

release:
	@echo "🚀 Running release..."
	$(GORELEASER) release --clean

snapshot:
	@echo "🔍 Running snapshot release (dry run)..."
	$(GORELEASER) release --snapshot --clean
