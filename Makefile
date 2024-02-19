# Define variables
GO_VERSION = 1.22.0
PROJECT_NAME = despensa-app

# Define directories
GOPATH = $(shell go env GOPATH)
SRC_DIR = .
BIN_DIR = ./bin
BINARY_NAME = despensa

# Define tools and flags




## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ####################
# BUILD
# ####################

# copy the htmx.min.js file to dist if not there or changed the source
public/dist/js/htmx.min.js: node_modules/htmx.org/dist/htmx.min.js
	cp -f ./node_modules/htmx.org/dist/htmx.min.js ./public/dist/js/htmx.min.js

# build in development mode (used by air / make watch)
.PHONY: build/dev
build/dev: public/dist/js/htmx.min.js 
	go build -o "${BIN_DIR}/${BINARY_NAME}" -ldflags="-X main.build=dev" ${SRC_DIR} 

## watch: run templ genrate in watch mode and start the air app reload
.PHONY: watch
watch:
	templ generate --watch & \
	npx tailwindcss -i ./public/css/style.css -o ./public/dist/css/style.css --watch & \
 	trap 'kill 0' SIGINT; \
	air

## build: build in release mode
.PHONY: build
build: public/dist/js/htmx.min.js
	templ generate
	npx tailwindcss -i ./public/css/style.css -o ./public/dist/css/style.css --minify
	go build -o "${BIN_DIR}/${BINARY_NAME}" -ldflags="-X main.build=1.0.0" ${SRC_DIR} 

## run: run the application after make a release build
.PHONY: run
run: build
	"${BIN_DIR}/${BINARY_NAME}"

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
	go vet ./...
	# go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	# go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	# go test -race -buildvcs -vet=off ./...


## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out





## install-tools: install tools required to build the project
.PHONY: install-tools
install-tools:
	@echo "Installing tools..."
	# go get -u github.com/golang-migrate/migrate/v4/...
	# go get -u github.com/machielw/templ/v3/...
	# go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest



## push: push changes to the remote Git repository
.PHONY: push
push: tidy audit no-dirty
	git push

## clean: delete all the generated resources
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf ${BIN_DIR}/*
	rm -fr despensa.db


MIGRATION_NAME=

.PHONY: migration/create
migration/create:
	migrate create -ext sql -dir database/migration/ -seq $(MIGRATION_NAME)

migration/up:
	migrate -path database/migration/ -database "sqlite://despensa.db" -verbose up

migration/down:
	migrate -path database/migration/ -database "sqlite://despensa.db" -verbose down

migration/version:
	migrate 
