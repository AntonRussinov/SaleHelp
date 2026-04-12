FROM golang:1.24.4-alpine3.21 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/salehelp/main.go

FROM alpine:3.21
ARG git_commit=default
ARG version=0.0.0
ARG env=prod
ENV GIT_COMMIT=$git_commit
ENV VERSION=$version
ENV ENV=$env
LABEL GIT_COMMIT=$git_commit \
      VERSION=$version \
      ENV=$env
RUN apk update && apk --no-cache add tzdata ca-certificates && apk upgrade
RUN apk --no-cache add poppler-utils wv
WORKDIR /app
COPY --from=build-stage /app/main /app/main
COPY docs/ ./docs/
COPY configs/ ./configs/
EXPOSE 8080
CMD ["./main", "-c", "configs/config.prod.yml"]