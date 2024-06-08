# Loads environment variables from .env file
ifneq (,$(wildcard ./build/.env))
    include ./build/.env
    export $(shell sed 's/=.*//' ./build/.env)
endif

# Makefile variables
PROJECT = "lignis-api"
GEN_PATH = "./internal/generated"
SWAGGER_API = "./swagger/api.yaml"

REGISTRY = ""
# Common section
gen-server:
	@echo "Generating server files..."
	@ogen -v --target internal/generated/api --clean swagger/api.yaml

clean-gen:
	@echo "Removing generated files..."
	@rm -rf $(GEN_PATH)

gen-spec: clean-gen gen-server

deploy:
	@go build -o build/main cmd/main.go
	@tar czf - build | ssh root@209.38.204.242 "tar xzf - && docker-compose -f build/docker_compose.yaml up --build -d"

docker-compose-infra:
	@docker-compose -f build/docker_compose.yaml up -d

run:
	@clear
	@go build -o build/main cmd/main.go
	echo "Compilation completed"
	@docker-compose -f build/docker_compose.yaml up --build -d

