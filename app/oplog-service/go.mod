module oplog-service

go 1.25.7

require (
	github.com/zeromicro/go-zero v1.10.1
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
	micro-shared v0.0.0
)

replace micro-shared => ../shared
