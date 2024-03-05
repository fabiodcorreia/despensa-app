# Build stage with NodeJS for TailwindCSS and DaisyUI
FROM node:21 AS web-build-stage 

WORKDIR /app

COPY package*.json tailwind.config.js Makefile .
COPY assets/ assets
COPY internal/views/ internal/views

RUN npm install && make docker/node

# Build stange with Go for the application
FROM golang:1.22 AS build-stage

WORKDIR /app

COPY --from=web-build-stage /app/public ./public
COPY database/ database
COPY internal/ internal
COPY Makefile go.mod go.sum main.go LICENSE .
RUN make docker/tools
RUN make test/mocks
RUN go get ./...
RUN make docker/go

# # Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /app

ENV DATABASE_FILE=/home/nonroot/data/despensa.db
ENV SERVER_PORT=8080

EXPOSE 8080

USER nonroot:nonroot

COPY --from=build-stage /app/bin/despensa /home/nonroot/despensa

VOLUME /home/nonroot/data

ENTRYPOINT ["/home/nonroot/despensa"]
