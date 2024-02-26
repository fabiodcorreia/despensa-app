# Define variables
GO_VERSION 			= 1.22.0
PROJECT_NAME 		= Despensa
BINARY_NAME 		= despensa
MIGRATION_NAME  =  

# Define directories
GOPATH 					= $(shell go env GOPATH)
SRC_DIR       	= .
BIN_DIR 				= $(SRC_DIR)/bin
ASSETS_DIR 			= $(SRC_DIR)/assets
PUBLIC_DIR 			= $(SRC_DIR)/public
NODE_DIR 				= $(SRC_DIR)/node_modules

# Define tools and flags
SERVER_ADDRESS	= "127.0.0.1"
SERVER_PORT 		= 8080
DB_FILE 				= $(BIN_DIR)/$(BINARY_NAME).db

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
.PHONY: build/dev
build/dev: public/js/htmx.min.js public/favicon
	go build -o "${BIN_DIR}/${BINARY_NAME}" -ldflags="-X main.build=dev" ${SRC_DIR} 



## watch: run templ, tailwindcss in watch mode and air app reload
.PHONY: watch
watch: public/css 
	echo "for some reason needs this here otherwise the first command exits" & \
	ADDRESS=$(SERVER_ADDRESS) PORT=$(SERVER_PORT) DATABASE_FILE=$(DB_FILE) air & \
	npx tailwindcss -i $(ASSETS_DIR)/css/style.css -o $(PUBLIC_DIR)/css/style.min.css --watch & \
  trap 'echo " Stopping..."; jobs -r | awk "{print $1}" | xargs kill -SIGTERM; sleep 5; echo "Gracefully exited."' SIGINT
	sleep 3
	templ generate --watch --proxy=http://$(SERVER_ADDRESS):$(SERVER_PORT) -path=$(SRC_DIR)/internal/views



## build: build in release mode
#  https://gophercoding.com/reduce-go-binary-size/
.PHONY: build
build: public/js/htmx.min.js public/css/style.min.css
	templ generate -path=$(SRC_DIR)/internal/views
	go build -o "${BIN_DIR}/${BINARY_NAME}" -ldflags="-X main.version=1.0.0 -X main.name=$(PROJECT_NAME)-w -s" ${SRC_DIR} 


## clean: delete all the generated resources
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -fr ${BIN_DIR}/*
	rm -fr ${PUBLIC_DIR}/js
	rm -fr ${PUBLIC_DIR}/css
	rm -fr ${PUBLIC_DIR}/favicon
	rm -fr ${PUBLIC_DIR}/favicon.ico



























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
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/boyter/dcd@latest
	go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	npm install -g tailwindcss
	npm install


## push: push changes to the remote Git repository
.PHONY: push
push: tidy audit no-dirty
	git push




## migration/create: creates a new migration files with the provided MIGRATION_NAME
.PHONY: migration/create
migration/create:
	migrate create -ext sql -dir database/migration/ -seq $(MIGRATION_NAME)

## migration/up: runs all the up migrations
.PHONY: migration/up
migration/up:
	migrate -path database/migration/ -database "sqlite://despensa.db" -verbose up

## migration/down: runs all the down migrations
.PHONY: migration/down
migration/down:
	migrate -path database/migration/ -database "sqlite://despensa.db" -verbose down

