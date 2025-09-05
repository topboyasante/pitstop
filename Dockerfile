ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Install swag for generating API documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

# Generate swagger docs
RUN swag init -g cmd/server/main.go -o docs/v1

RUN go build -v -o /run-app ./cmd/server


FROM debian:bookworm

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /run-app /usr/local/bin/
COPY --from=builder /usr/src/app/docs /docs
CMD ["run-app"]
