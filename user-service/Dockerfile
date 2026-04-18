FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY go.work go.work
COPY user-service/ user-service/
WORKDIR /build/user-service
RUN go mod tidy && go build -o userservice user.go

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/user-service/userservice /app/userservice
COPY --from=builder /build/user-service/etc /app/etc
EXPOSE 9101
ENTRYPOINT ["/app/userservice", "-f", "/app/etc/user.yaml"]
