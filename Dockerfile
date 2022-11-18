FROM golang:1.19.3-alpine as builder
WORKDIR /app
COPY . /app

ENV BUILD_TAG 1.0.0
ENV GO111MODULE on
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go mod vendor
RUN go build -o crudapi main.go

# stage2.1: rebuild
FROM alpine
WORKDIR /app
COPY --from=builder /app/crudapi /app/crudapi.go
CMD ["./crudapi.go"]


