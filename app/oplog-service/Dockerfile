FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY go.work go.work
COPY app/shared/ app/shared/
COPY app/oplog-service/ app/oplog-service/
WORKDIR /build/app/oplog-service
RUN GOWORK=off go mod tidy && GOWORK=off go build -o oplogservice oplog.go

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/app/oplog-service/oplogservice /app/oplogservice
COPY --from=builder /build/app/oplog-service/etc /app/etc
EXPOSE 9103
ENTRYPOINT ["/app/oplogservice", "-f", "/app/etc/oplog.yaml"]
