ROOT_DIR := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))

ASSET_DIR=./assets
DIST_DIR=$(ASSET_DIR)/dist
ASSET_JS=$(DIST_DIR)/bundle.js
ASSET_SRC=$(ASSET_DIR)/src/*.tsx
ASSETS=$(ASSET_JS) $(DIST_DIR)/index.html

SRC=*.go \
	pkg/**/*.go \
	types/*.go
ENT_DIR=./pkg/infra/ent
ENT_SRC=$(ENT_DIR)/ent.go
ENT_SCHEMA_DIR=./pkg/infra/schema

CHAIN=chain.so
EXAMPLE_SRC_DIR=./examples/simple

all: alertchain

ent: $(ENT_SRC)

chain: $(CHAIN)

docker:
	docker run -p 127.0.0.1:3306:3306 -e MYSQL_ROOT_PASSWORD -e MYSQL_DATABASE mysql

$(ASSET_JS): $(ASSET_SRC)
	cd $(ASSET_DIR) && npm i && cd $(ROOT_DIR)

$(ENT_SRC): $(ENT_SCHEMA_DIR)/*.go
	ent generate $(ENT_SCHEMA_DIR) --target $(ENT_DIR)

dev: $(CHAIN)
	go run ./cmd/alertchain/ serve -d "root:${MYSQL_ROOT_PASSWORD}@tcp(localhost:3306)/${MYSQL_DATABASE}" -c $(CHAIN)

test: $(SRC) $(ENT_SRC)
	go test ./...

$(CHAIN): $(EXAMPLE_SRC_DIR)/*.go $(SRC) $(ENT_SRC)
	go build -buildmode=plugin -o chain.so $(EXAMPLE_SRC_DIR)

alertchain: $(ASSETS) $(SRC) $(ENT_SRC)
	go build -o alertchain ./cmd/alertchain
