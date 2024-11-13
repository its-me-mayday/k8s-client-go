FROM golang:1.23 as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o k8s-go-client .

FROM scratch
COPY --from=builder /app/k8s-go-client /k8s-go-client

ENTRYPOINT ["/k8s-go-client"]
