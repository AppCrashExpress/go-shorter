# syntax=docker/dockerfile:1

FROM golang:1.18-alpine AS build

RUN apk add --no-cache protobuf

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .
COPY ./src/. ./src/
COPY ./main.go .
RUN go mod download && \
    protoc --go_out=. --go_opt=paths=source_relative \
      --go-grpc_out=. --go-grpc_opt=paths=source_relative \
      ./src/api/shortener.proto

RUN CGO_ENABLED=0 go build -o /app/service

FROM gcr.io/distroless/base-debian10:debug

WORKDIR /

COPY --from=build /app/service /app/service

USER nonroot:nonroot

ENTRYPOINT ["/app/service"]

