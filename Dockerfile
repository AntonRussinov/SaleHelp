# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.22.1-alpine3.19 as build-stage

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . ./

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/salehelp/main.go

# Deploy the application binary into a lean image
FROM alpine:3.18

ARG git_commit=default
ARG version=0.0.0
ARG env=dev

ENV GIT_COMMIT=$git_commit
ENV VERSION=$version
ENV ENV=$env

LABEL GIT_COMMIT=$git_commit \
      VERSION=$version \
      ENV=$env

RUN apk update && apk --no-cache add tzdata && apk upgrade
RUN apk --no-cache add poppler-utils wv

WORKDIR /app

COPY --from=build-stage /app/main /app/main
COPY docs/ ./docs/
COPY configs/ ./configs/
COPY assets/ ./assets/

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
