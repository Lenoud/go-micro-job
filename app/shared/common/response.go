package common

import "time"

const (
	CodeSuccess      int64 = 200
	CodeParam        int64 = 400
	CodeUnauthorized int64 = 401
	CodeForbidden    int64 = 403
	CodeNotFound     int64 = 404
	CodeServer       int64 = 500
)

func CurrentTimeMillis() int64 {
	return time.Now().UnixMilli()
}
