FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY go.work go.work
COPY app/user-service/ app/user-service/
WORKDIR /build/app/user-service
RUN GOWORK=off go mod tidy && GOWORK=off go build -o userservice user.go

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/app/user-service/userservice /app/userservice
COPY --from=builder /build/app/user-service/etc /app/etc
EXPOSE 9101
ENTRYPOINT ["/app/userservice", "-f", "/app/etc/user.yaml"]
