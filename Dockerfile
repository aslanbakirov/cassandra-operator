# build stage
FROM golang:1.9.4 AS builder
WORKDIR /go/src/github.com/aslanbekirov/cassandra-operator
ADD . $WORKDIR
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o operator cmd/main.go

# final stage
FROM alpine
WORKDIR /app
COPY --from=builder /go/src/github.com/aslanbekirov/cassandra-operator/operator .
CMD ["./operator"]
