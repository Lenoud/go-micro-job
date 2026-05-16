FROM golang:1.25-alpine AS builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /build
COPY go.work go.work
COPY app/user-service/go.mod app/user-service/go.sum app/user-service/
COPY app/department-service/go.mod app/department-service/go.sum app/department-service/
COPY app/api-gateway/go.mod app/api-gateway/go.sum app/api-gateway/
WORKDIR /build/app/api-gateway
RUN go mod download
WORKDIR /build
COPY app/user-service/ app/user-service/
COPY app/department-service/ app/department-service/
COPY app/api-gateway/ app/api-gateway/
WORKDIR /build/app/api-gateway
RUN go build -o apigateway apigateway.go

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/app/api-gateway/apigateway /app/apigateway
COPY --from=builder /build/app/api-gateway/etc /app/etc
EXPOSE 9000
ENTRYPOINT ["/app/apigateway", "-f", "/app/etc/apigateway.yaml"]
