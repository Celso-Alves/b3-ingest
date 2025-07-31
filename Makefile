GO = go
# Variáveis do Projeto
# Nome do executável da aplicação Go
APP_NAME := b3-ingest
# Diretório principal do código Go
GO_SRC := ./
# Diretório onde os arquivos CSV de dados estão localizados
CSV_DATA_DIR := ./data
# Diretório de saída do binário
BIN_DIR := ./cmd
BIN_PATH := $(BIN_DIR)/$(APP_NAME)

# Variáveis Docker
DOCKER_COMPOSE_FILE := docker-compose.yml
DOCKER_APP_SERVICE := app
DOCKER_DB_SERVICE := postgres

# Variáveis de Configuração (pode ser sobrescrito via linha de comando)
# Exemplo: make build GOOS=linux GOARCH=amd64
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GO_BUILD_FLAGS := -ldflags="-s -w" # Reduz o tamanho do binário e remove informações de depuração

# Alvos principais

.PHONY: all build run ingest test coverage docker-build docker-up docker-down docker-ingest clean

all: build # Alvo padrão, constrói a aplicação

# Constrói o executável Go para o sistema operacional e arquitetura atuais
build:
	@echo "Building Go application..."
	@mkdir -p $(BIN_DIR)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GO_BUILD_FLAGS) -o $(BIN_PATH) $(GO_SRC)
	@echo "Build complete: $(BIN_PATH)"

# Executa a aplicação Go localmente
run: build
	@echo "Running Go application..."
	@$(BIN_PATH)

# Executa o processo de ingestão de dados localmente
ingest: build
	@echo "Running data ingestion locally..."
	@$(BIN_PATH) -load

# Executa todos os testes unitários
test:
	@echo "Running unit tests..."
	@go test ./... -v

# Gera relatório de cobertura dos testes
coverage:
	@echo "Checking test coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -func=coverage.out

# Executa o servidor HTTP localmente
serve: build
	@echo "Running HTTP server locally..."
	@$(BIN_PATH) -serve

# Executa o modo download localmente@echo "  lint            : Runs golangci-lint for static code analysis."@echo "  lint            : Runs golangci-lint for static code analysis."
download: build
	@echo "Running download mode locally..."
	@$(BIN_PATH) -download

# Gerenciamento Docker

# Constrói as imagens Docker definidas no docker-compose.yml
docker-build:
	@echo "Building Docker images..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) build

# Sobe os serviços Docker (banco de dados e aplicação)
docker-up: docker-build
	@echo "Starting Docker services..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d --remove-orphans
	@echo "Docker services are up. Access API at http://localhost:8080"

# Executa o processo de ingestão de dados dentro do container Docker
# Garante que o banco de dados esteja pronto antes de iniciar a ingestão
docker-ingest: docker-build
	@echo "Ensuring database is ready before ingestion..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d $(DOCKER_DB_SERVICE)
	@echo "Waiting for database to be fully ready (might take a few seconds)..."
	@sleep 10 # Pequena pausa para o DB iniciar completamente, ajuste se necessário
	@echo "Running data ingestion inside Docker container..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) run --rm $(DOCKER_APP_SERVICE) ./b3-ingest -load
	@echo "Docker ingestion complete."

# Executa o servidor HTTP dentro do container Docker
docker-serve: docker-build
	@echo "Starting HTTP server in Docker..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) run --rm -p 8000:8000 $(DOCKER_APP_SERVICE) ./b3-ingest -serve

# Entra no shell do container da aplicação
shell:
	@echo "Abrindo shell no container da aplicação..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) run --rm $(DOCKER_APP_SERVICE) sh

# Derruba os serviços Docker
docker-down:
	@echo "Stopping and removing Docker services..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down -v # -v remove volumes anonimos

# Limpa os arquivos gerados pelo build
clean:
	@echo "Cleaning up build artifacts..."
	@rm -f $(BIN_PATH)
	@echo "Clean complete."

# Ajuda para os alvos disponíveis
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all             : Builds the Go application (default)."
	@echo "  build           : Compiles the Go application binary."
	@echo "  run             : Runs the compiled Go application locally."
	@echo "  ingest          : Runs the data ingestion process locally."
	@echo "  serve           : Runs the HTTP server locally."
	@echo "  download        : Runs the download mode locally."
	@echo "  test            : Runs all unit tests."
	@echo "  coverage        : Shows test coverage report."
	@echo "  docker-build    : Builds the Docker images."
	@echo "  docker-up       : Starts the Docker services (database and application) in detached mode."
	@echo "  docker-down     : Stops and removes the Docker services and associated volumes."
	@echo "  docker-ingest   : Runs the data ingestion process inside the Docker application container."
	@echo "  clean           : Removes compiled binaries and other temporary files."
	@echo "  help            : Displays this help message."
	@echo ""
	@echo "Variables (can be overridden):"
	@echo "  GOOS            : Target OS (e.g., linux, windows, darwin). Default: $(GOOS)"
	@echo "  GOARCH          : Target architecture (e.g., amd64, arm64). Default: $(GOARCH)"

