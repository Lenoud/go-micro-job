FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY go.work go.work
COPY app/shared/ app/shared/
COPY app/department-service/ app/department-service/
WORKDIR /build/app/department-service
RUN GOWORK=off go mod tidy && GOWORK=off go build -o departmentservice department.go

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/app/department-service/departmentservice /app/departmentservice
COPY --from=builder /build/app/department-service/etc /app/etc
EXPOSE 9102
ENTRYPOINT ["/app/departmentservice", "-f", "/app/etc/department.yaml"]
