FROM golang:1.20.2-alpine3.16 as builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go main.go
COPY cmd/ cmd/
COPY pkg/ pkg/

RUN go mod vendor && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o /app/grafana-metrics-analyzer

FROM alpine

WORKDIR /

COPY --from=builder /app/grafana-metrics-analyzer .

ENTRYPOINT [ "/grafana-metrics-analyzer" ]
