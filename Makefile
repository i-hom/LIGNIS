# Loads environment variables from .env file
ifneq (,$(wildcard ./build/.env))
    include ./build/.env
    export $(shell sed 's/=.*//' ./build/.env)
endif

# Makefile variables
PROJECT = "lignis-api"
GEN_PATH = "./internal/generated"
SWAGGER_TEMPLATE = "./spec/template"
SWAGGER_API = "./swagger/api.yaml"
VERSION = $(shell git rev-parse --abbrev-ref HEAD)

REGISTRY = ""
# Common section
gen-server:
	@echo "Generating server files..."
	@ogen -v --target internal/generated/api --clean swagger/api.yaml

clean-gen:
	@echo "Removing generated files..."
	@rm -rf $(GEN_PATH)

gen-spec: clean-gen gen-server

# Docker section
docker-build:
	@echo "Building docker image..."
	@docker build --ssh default --build-arg TAG=${VERSION} \
		--build-arg BUILD_TIME=${BUILD_TIME} \
		--target builder \
		-t ${PROJECT}:build \
		-f build/Dockerfile .

docker-build-prod:
	@echo "Building docker image..."
	@docker build --ssh default --build-arg TAG=${VERSION} \
		--build-arg BUILD_TIME=${BUILD_TIME} \
		--target production \
		--platform linux/amd64 \
		-t ${PROJECT}:prod \
		-f build/Dockerfile .

docker-push:
	@echo "Pushing docker image..."
	@docker tag ${PROJECT}:prod ${REGISTRY}:latest
	@docker tag ${PROJECT}:prod ${REGISTRY}:${VERSION}
	@docker push ${REGISTRY}:latest
	@docker push ${REGISTRY}:${VERSION}

docker-deploy: docker-build-prod docker-push

docker-network-create:
	@docker network create lignis-api || true

docker-infra: docker-network-create
	@docker-compose -f build/docker_compose_infra.yml up -d

run:
	@clear
	@go run cmd/main.go

