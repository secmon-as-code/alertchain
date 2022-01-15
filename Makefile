ROOT_DIR := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))

ASSET_DIR=./assets
DIST_DIR=$(ASSET_DIR)/dist
ASSET_JS=$(DIST_DIR)/bundle.js
ASSET_SRC=$(ASSET_DIR)/src/*.tsx
ASSETS=$(ASSET_JS) $(DIST_DIR)/index.html

SRC=*.go pkg/**/*.go

ENT_DIR=./gen/ent
ENT_SRC=$(ENT_DIR)/ent.go
ENT_SCHEMA_DIR=./pkg/domain/schema

all: alertchain

ent: $(ENT_SRC)

$(ENT_SRC): $(ENT_SCHEMA_DIR)/*.go
	ent generate $(ENT_SCHEMA_DIR) --target $(ENT_DIR)

$(CHAIN): $(EXAMPLE_SRC_DIR)/*.go $(SRC) $(ENT_SRC)
	go build -buildmode=plugin -o chain.so $(EXAMPLE_SRC_DIR)

test:
	go test

alertchain: $(ASSETS) $(SRC) $(ENT_SRC)
	go build -o alertchain .
