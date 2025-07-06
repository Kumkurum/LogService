# Переменные
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc
GO := go
PROTO_DIR := proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT_DIR := internal/transport
SRC_DIR := cmd/main.go
BINARY_NAME := LoggingService
TEST_DIR := ./test/*
TEST_RESULTS_DIR := test_results
COVERAGE_OUT := $(TEST_RESULTS_DIR)/coverage.out
COVERAGE_HTML := $(TEST_RESULTS_DIR)/coverage.html
.PHONY: all clean generate build

all: generate build

# Генерация Go-кода из .proto файлов
generate:
	@echo "Generating protobuf and gRPC code..."
	@mkdir -p $(GO_OUT_DIR)
	protoc \
		--go_opt=paths=source_relative \
		--go_out=$(GO_OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT_DIR) \
		-I=$(PROTO_DIR) \
		$(PROTO_FILES)
	@echo "Code generated in $(GO_OUT_DIR)"

# Сборка основного бинарника
build: generate
	@echo "Building $(BINARY_NAME)..."
	@$(GO) build -o bin/$(BINARY_NAME) $(SRC_DIR)
	@echo "Binary built: bin/$(BINARY_NAME)"

# Запуск тестов
test-coverage: generate
	@echo "Running tests with coverage..."
	@mkdir -p $(TEST_RESULTS_DIR)  # Создаем папку, если её нет
	@$(GO) test -coverprofile=$(COVERAGE_OUT) $(TEST_DIR)
	@$(GO) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

# Проверка форматирования кода
fmt:
	@echo "Checking code formatting..."
	@$(GO) fmt ./...
	@echo "Code is properly formatted"

# Проверка зависимостей
tidy:
	@echo "Tidying up dependencies..."
	@$(GO) mod tidy
	@echo "Dependencies are up to date"

# Запуск приложения
run: build
	@./bin/$(BINARY_NAME)


build-static:
	CGO_ENABLED=0 go build -a -ldflags '-s -w' -o bin/static/$(BINARY_NAME) cmd/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME) cmd/main.go

# Очистка сгенерированных файлов и бинарника
clean:
	@echo "Cleaning generated files and binary..."
	@rm -rf $(GO_OUT_DIR)/* bin
	@rm -rf $(GO_OUT_DIR)/* test_results
	@echo "Clean complete"