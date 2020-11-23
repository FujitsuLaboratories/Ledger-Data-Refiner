PACKAGES = $(shell go list ./...)
OUTPUT ?= ./ledgerdata-refiner
BUILD_TAGS ?= ledgerdata-refiner

DOCKER_DATABASE ?= fujitsu/refinerdb
DOCKER_DATABASE_FILE ?= DOCKER/postgreSQL/Dockerfile
DOCKER_REFINER ?= fujitsu/refiner
DOCKER_REFINER_FILE ?= DOCKER/app/Dockerfile

CGO_ENABLED ?= 0

#################################################################
#																#
#						Build ledgerdata refiner 				#
#																#
#################################################################
build:
	CGO_ENABLED=$(CGO_ENABLED) go build -tags '$(BUILD_TAGS)' -o $(OUTPUT) .
.PHONY: build

install:
	CGO_ENABLED=$(CGO_ENABLED) go install $(BUILD_FLAGS) -tags $(BUILD_TAGS) .
.PHONY: install

clean:
	rm -rf $(BUILD_TAGS)
.PHONY: clean

vendor:
	CGO_ENABLED=$(CGO_ENABLED) go mod vendor
.PHONY: vendor

#################################################################
#																#
#						Build DOCKER images 					#
#																#
#################################################################
docker_all:
	docker build -t $(DOCKER_DATABASE) -f $(DOCKER_DATABASE_FILE) .
	docker build -t $(DOCKER_REFINER) -f $(DOCKER_REFINER_FILE) .
.PHONY: docker_all

docker_refinerdb:
	docker build -t $(DOCKER_DATABASE) -f $(DOCKER_DATABASE_FILE) .
.PHONY: docker_refinerdb

docker_refiner:
	docker build -t $(DOCKER_REFINER) -f $(DOCKER_REFINER_FILE) .
.PHONY: docker_refiner