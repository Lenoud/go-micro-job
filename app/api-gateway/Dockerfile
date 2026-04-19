FROM golang:1.25-alpine AS builder
WORKDIR /build
# 先复制 go.work 和两个模块
COPY go.work go.work
COPY app/user-service/go.mod app/user-service/go.sum user-service/
COPY app/api-gateway/go.mod app/api-gateway/go.sum api-gateway/
WORKDIR /build/api-gateway
RUN go mod tidy
WORKDIR /build
# 复制完整源码
COPY app/user-service/ user-service/
COPY app/api-gateway/ api-gateway/
WORKDIR /build/api-gateway
RUN go build -o apigateway apigateway.go

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/api-gateway/apigateway /app/apigateway
COPY --from=builder /build/api-gateway/etc /app/etc
EXPOSE 9000
ENTRYPOINT ["/app/apigateway", "-f", "/app/etc/apigateway.yaml"]
