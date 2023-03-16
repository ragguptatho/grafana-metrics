FROM golang:1.20.2-alpine3.16

WORKDIR /app

COPY . .

RUN go mod vendor && go build -mod=vendor -o /app/grafana-metrics-analyzer

ENTRYPOINT ["/app/grafana-metrics-analyzer"]
