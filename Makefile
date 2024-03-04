# Define variables
GO_VERSION 				= 1.22.0
PROJECT_NAME 			= Despensa
BINARY_NAME 			= despensa
MIGRATION_NAME    = 

# Define directories
GOPATH 						= $(shell go env GOPATH)
SRC_DIR						= .
BIN_DIR 					= $(SRC_DIR)/bin
ASSETS_DIR 				= $(SRC_DIR)/assets
PUBLIC_DIR 				= $(SRC_DIR)/public
NODE_DIR 					= $(SRC_DIR)/node_modules

# Define tools and flags
SERVER_ADDRESS		= "127.0.0.1"
SERVER_PORT 			= 8080
DB_FILE						= $(BIN_DIR)/$(BINARY_NAME).db

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# #################### BUILD #################### #

public/favicon: assets/favicon
	mkdir -p $(PUBLIC_DIR)/favicon
	cp -f $(ASSETS_DIR)/favicon/* $(PUBLIC_DIR)/favicon
	mv $(PUBLIC_DIR)/favicon/favicon.ico $(PUBLIC_DIR)/favicon.ico
	@sed -i '' 's/"name": ""/"name": "$(PROJECT_NAME)"/g' $(PUBLIC_DIR)/favicon/site.webmanifest 
	@sed -i '' 's/"short-name": ""/"ahort-name": "$(PROJECT_NAME)"/g' $(PUBLIC_DIR)/favicon/site.webmanifest 

#  create public/js folder if not exists
public/js:
	mkdir -p $(PUBLIC_DIR)/js

#  create public/css folder if not exists
public/css:
	mkdir -p $(PUBLIC_DIR)/css

#  copy the htmx.min.js file to dist if not there or changed the source
public/js/htmx.min.js: public/js node_modules/htmx.org/dist/htmx.min.js
	cp -f $(NODE_DIR)/htmx.org/dist/htmx.min.js $(PUBLIC_DIR)/js/htmx.min.js

public/css/style.min.css: public/css assets/css/style.css
	npx tailwindcss -i $(ASSETS_DIR)/css/style.css -o $(PUBLIC_DIR)/css/style.min.css --minify

#  build in development mode (user by make watch)
.PHONY: dev/build
dev/build: public/js/htmx.min.js public/favicon
	go build -o "${BIN_DIR}/${BINARY_NAME}" -ldflags="-X main.build=dev" ${SRC_DIR} 

## dev/decisions:build printing the build optimization decisions
.PHONY: analyse/escape
analyse/escape: public/js/htmx.min.js public/favicon
	go build -gcflags=-m -o "${BIN_DIR}/${BINARY_NAME}" -ldflags="-X main.build=dev" ${SRC_DIR} 2>&1 | grep -vE "inlin|not escape"

## dev/run: run templ, tailwindcss in watch mode and air app reload
.PHONY: dev/run
dev/run: public/css 
	echo "for some reason needs this here otherwise the first command exits" & \
	ADDRESS=$(SERVER_ADDRESS) PORT=$(SERVER_PORT) DATABASE_FILE=$(DB_FILE) air --build.bin="${BIN_DIR}/${BINARY_NAME}"& \
	npx tailwindcss -i $(ASSETS_DIR)/css/style.css -o $(PUBLIC_DIR)/css/style.min.css --watch & \
  trap 'echo " Stopping..."; jobs -r | awk "{print $1}" | xargs kill -SIGTERM; sleep 5; echo "Gracefully exited."' SIGINT
	sleep 3
	templ generate --watch --proxy=http://$(SERVER_ADDRESS):$(SERVER_PORT) -path=$(SRC_DIR)/internal/views

## dev/clean: deletes all the generated resources
.PHONY: dev/clean
dev/clean:
	rm -fr $(PUBLIC_DIR)
	rm -fr $(SRC_DIR)/internal/views/**/*_templ.go


## dev/install: install all golang and npm dependencies
.PHONY: dev/install
dev/install:
	npm install
	go get ./...

## dev/update: updates all golang and npm dependencies
.PHONY: dev/update
dev/update:
	npm update
	go get -u ./...
	go mod tidy

## dev/tools: install tools required to build the project
.PHONY: dev/tools
dev/tools:
	go install golang.org/x/tools/cmd/deadcode@latest
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/boyter/dcd@latest
	go install github.com/alexkohler/prealloc@latest
	go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
	go install github.com/jgautheron/goconst/cmd/goconst@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/kisielk/errcheck@latest
	go install go.uber.org/mock/mockgen@latest
	go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest





## release/build: build in release mode a single binary
.PHONY: release/build
release/build: public/js/htmx.min.js public/css/style.min.css public/favicon tidy audit
	templ generate -path=$(SRC_DIR)/internal/views
	go build -o "${BIN_DIR}/${BINARY_NAME}" -ldflags="-X 'main.version=1.0.0' -X 'main.name=$(PROJECT_NAME)' -w -s" ${SRC_DIR} 

## release/clean: deletes the BIN_DIR
.PHONY: release/clean
release/clean:
	rm -fr $(BIN_DIR)

## clean: deletes the node_modules, release/clean and dev/clean
.PHONY: clean
clean: dev/clean release/clean
	rm -fr $(NODE_DIR)



# ###
# REVIEW
#

## dirty: check for uncommited changes
.PHONY: dirty
dirty:
	git diff --exit-code

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet --all .
	shadow -strict ./...
	prealloc -forloops -set_exit_status -simple -rangeloops ./...
	goconst ./internal/... main.go
	staticcheck -checks=all,-ST1000 ./...
	errcheck ./...
	govulncheck ./...
	go test -race -vet=off ./...


.PHONY: test
test: test/mocks
	go test -v -race ./...

# test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover: test/mocks
	go test -v -race -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## test/mocks: generate mocks for testing
test/mocks: internal/storage/location.go internal/storage/item.go
	mockgen -source=internal/storage/location.go -destination=test/mocks/location_store.go -package=mocks
	mockgen -source=internal/storage/item.go -destination=test/mocks/item_store.go -package=mocks






# push: push changes to the remote Git repository
# .PHONY: push
# push: tidy audit no-dirty
# 	git push




## migrate/create: creates a new migration files with the provided MIGRATION_NAME
.PHONY: migrate/create
migrate/create:
	migrate create -ext sql -dir database/migration/ -seq $(MIGRATION_NAME)

## migrate/up: runs all the up migrations
.PHONY: migrate/up
migrate/up:
	migrate -path database/migration/ -database "sqlite://$(DB_FILE)" -verbose up

## migrate/down: runs all the down migrations
.PHONY: migrate/down
migrate/down:
	migrate -path database/migration/ -database "sqlite://$(DB_FILE)" -verbose down

