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
ENT_SCHEMA=$(ENT_DIR)/schema/*.go

all: alertchain

$(ASSET_JS): $(ASSET_SRC)
	cd $(ASSET_DIR) && npm i && cd $(ROOT_DIR)

$(ENT_SRC): $(ENT_SCHEMA)
	go generate $(ENT_DIR)

test: $(SRC) $(ENT_SRC)
	go test ./...

example: ./examples/chain/*.go
	go build -buildmode=plugin -o chain.so ./examples/chain

alertchain: $(ASSETS) $(SRC) $(ENT_SRC)
	go build -o alertchain ./cmd/alertchain
