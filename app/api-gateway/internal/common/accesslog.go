package common

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	oplogclient "oplog-service/oplogclient"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

var sensitiveAccessLogKeys = map[string]struct{}{
	"password":      {},
	"newpassword":   {},
	"token":         {},
	"authorization": {},
	"admintoken":    {},
}

const (
	opLogBufSize   = 512
	opLogBatchSize = 50
	opLogFlushMs   = 100
)

// OpLogWriter buffers access log entries and batch-writes them via gRPC.
type OpLogWriter struct {
	ch  chan *oplogclient.OpLogRecord
	rpc oplogclient.OpLog
	done chan struct{}
}

// NewOpLogWriter starts a background goroutine to batch-write access logs.
func NewOpLogWriter(rpc oplogclient.OpLog) *OpLogWriter {
	w := &OpLogWriter{
		ch:   make(chan *oplogclient.OpLogRecord, opLogBufSize),
		rpc:  rpc,
		done: make(chan struct{}),
	}
	go w.run()
	return w
}

func (w *OpLogWriter) run() {
	batch := make([]*oplogclient.OpLogRecord, 0, opLogBatchSize)
	ticker := time.NewTicker(time.Duration(opLogFlushMs) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case entry, ok := <-w.ch:
			if !ok {
				w.flush(batch)
				close(w.done)
				return
			}
			batch = append(batch, entry)
			if len(batch) >= opLogBatchSize {
				w.flush(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				w.flush(batch)
				batch = batch[:0]
			}
		}
	}
}

func (w *OpLogWriter) flush(batch []*oplogclient.OpLogRecord) {
	if len(batch) == 0 {
		return
	}
	_, err := w.rpc.BatchCreate(context.Background(), &oplogclient.BatchCreateReq{Logs: batch})
	if err != nil {
		logx.Errorf("[gateway] oplog BatchCreate failed: %v", err)
	}
}

// Stop drains remaining logs and waits for the background goroutine to finish.
func (w *OpLogWriter) Stop() {
	close(w.ch)
	<-w.done
}

// Enqueue adds an entry to the write buffer. Non-blocking; drops on full buffer.
func (w *OpLogWriter) Enqueue(entry *oplogclient.OpLogRecord) {
	select {
	case w.ch <- entry:
	default:
	}
}

// NewAccessLogMiddleware returns a middleware that records access logs.
func NewAccessLogMiddleware(writer *OpLogWriter, accessSecret string) func(http.HandlerFunc) http.HandlerFunc {
	if writer == nil {
		return func(next http.HandlerFunc) http.HandlerFunc {
			return next
		}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r == nil || next == nil {
				return
			}

			startedAt := time.Now()
			recorder := newAccessLogResponseWriter(w)
			var loginUsername string
			if isLoginRoute(r) {
				loginUsername = peekLoginUsername(r)
			}
			next(recorder, r)
			biz := parseAccessLogBizResult(recorder.statusCode, recorder.body.Bytes())
			recorder.flush(accessLogTraceID(r), biz)

			if !shouldRecordAccessLog(r) {
				return
			}

			userId := accessLogUserID(r, accessSecret)
			if userId == "" {
				userId = loginUsername
			}

			entry := &oplogclient.OpLogRecord{
				RequestId:  accessLogTraceID(r),
				UserId:     userId,
				ReIp:       accessLogIP(r),
				ReTime:     startedAt.UnixMilli(),
				ReUa:       strings.TrimSpace(r.UserAgent()),
				ReUrl:      accessLogURL(r),
				ReMethod:   strings.ToUpper(strings.TrimSpace(r.Method)),
				ReContent:  accessLogContent(r),
				Success:    biz.Success,
				BizCode:    biz.BizCode,
				BizMsg:     biz.BizMsg,
				AccessTime: time.Since(startedAt).Milliseconds(),
			}
			writer.Enqueue(entry)
		}
	}
}

// ---- route filtering ----

func shouldRecordAccessLog(r *http.Request) bool {
	if r == nil {
		return false
	}

	path := accessLogURL(r)
	if path == "/api/opLog/list" || path == "/api/opLog/loginLogList" {
		return false
	}
	if path == "/api/user/login" || path == "/api/user/userLogin" {
		return true
	}

	switch strings.ToUpper(strings.TrimSpace(r.Method)) {
	case http.MethodPost, http.MethodPut, http.MethodDelete:
		return true
	default:
		return false
	}
}

func isLoginRoute(r *http.Request) bool {
	if r == nil {
		return false
	}
	path := accessLogURL(r)
	return path == "/api/user/login" || path == "/api/user/userLogin" || path == "/api/user/userRegister"
}

// peekLoginUsername reads the request body to extract the username field
// for login/register routes, then restores the body for downstream handlers.
func peekLoginUsername(r *http.Request) string {
	if r == nil || r.Body == nil {
		return ""
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	// Restore body so the handler can still read it
	r.Body = io.NopCloser(bytes.NewReader(body))

	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return ""
	}

	var payload struct {
		Username string `json:"username"`
	}
	if err := json.Unmarshal(trimmed, &payload); err != nil {
		return ""
	}
	return strings.TrimSpace(payload.Username)
}

// ---- request data extraction ----

func accessLogURL(r *http.Request) string {
	if r == nil || r.URL == nil {
		return ""
	}
	return r.URL.Path
}

func accessLogIP(r *http.Request) string {
	if r == nil {
		return ""
	}

	for _, header := range []string{"X-Forwarded-For", "X-Real-IP"} {
		if raw := strings.TrimSpace(r.Header.Get(header)); raw != "" {
			if header == "X-Forwarded-For" {
				parts := strings.Split(raw, ",")
				if len(parts) > 0 {
					return strings.TrimSpace(parts[0])
				}
			}
			return raw
		}
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

func accessLogUserID(r *http.Request, accessSecret string) string {
	if r == nil {
		return ""
	}

	// go-zero JWT middleware already puts claims into context
	if userID, ok := claimStringFromCtx(r.Context(), "userId"); ok {
		return userID
	}
	if accessSecret == "" {
		return ""
	}

	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(strings.TrimSpace(parts[0]), "Bearer") {
		return ""
	}

	tokenString := strings.TrimSpace(parts[1])
	if tokenString == "" {
		return ""
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	if err != nil || token == nil || !token.Valid {
		return ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}

	rawUserID, ok := claims["userId"]
	if !ok {
		return ""
	}

	switch v := rawUserID.(type) {
	case string:
		return strings.TrimSpace(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return ""
	}
}

func claimStringFromCtx(ctx context.Context, key string) (string, bool) {
	if ctx == nil {
		return "", false
	}
	switch v := ctx.Value(key).(type) {
	case string:
		v = strings.TrimSpace(v)
		return v, v != ""
	case float64:
		return strconv.FormatInt(int64(v), 10), true
	case int:
		return strconv.Itoa(v), true
	case int64:
		return strconv.FormatInt(v, 10), true
	default:
		return "", false
	}
}

func accessLogContent(r *http.Request) string {
	if r == nil {
		return ""
	}

	if r.Method == http.MethodGet {
		return sanitizeValues(r.URL.Query())
	}

	contentType := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
	switch {
	case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
		if err := r.ParseForm(); err != nil {
			return ""
		}
		return sanitizeValues(r.PostForm)
	case strings.HasPrefix(contentType, "multipart/form-data"):
		return "multipart/form-data"
	default:
		return sanitizeValues(r.URL.Query())
	}
}

func sanitizeValues(values map[string][]string) string {
	if len(values) == 0 {
		return ""
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		lowerKey := strings.ToLower(strings.TrimSpace(key))
		val := strings.Join(values[key], ",")
		if _, sensitive := sensitiveAccessLogKeys[lowerKey]; sensitive {
			val = "[REDACTED]"
		}
		parts = append(parts, key+"="+val)
	}
	return strings.Join(parts, "&")
}

func accessLogTraceID(r *http.Request) string {
	if r == nil {
		return ""
	}
	return GetRequestID(r.Context())
}

// ---- response interception ----

type accessLogBizResult struct {
	Success string
	BizCode int64
	BizMsg  string
}

func parseAccessLogBizResult(statusCode int, body []byte) accessLogBizResult {
	result := accessLogBizResult{
		Success: accessLogSuccess(statusCode),
		BizCode: defaultBizCode(statusCode),
	}

	trimmed := bytes.TrimSpace(body)
	if len(trimmed) > 0 && trimmed[0] == '{' {
		var payload struct {
			Code int64  `json:"code"`
			Msg  string `json:"msg"`
		}
		if err := json.Unmarshal(trimmed, &payload); err == nil {
			result.BizCode = payload.Code
			result.BizMsg = truncateBizMsg(payload.Msg)
			if payload.Code == 0 || payload.Code == 200 {
				result.Success = "1"
			} else {
				result.Success = "0"
			}
		}
	}

	if statusCode >= http.StatusBadRequest && result.BizCode == 200 {
		result.Success = "0"
		result.BizCode = int64(statusCode)
	}
	if statusCode >= http.StatusBadRequest && result.BizMsg == "" {
		result.BizMsg = truncateBizMsg(http.StatusText(statusCode))
	}

	return result
}

func accessLogSuccess(statusCode int) string {
	if statusCode >= http.StatusBadRequest {
		return "0"
	}
	return "1"
}

func defaultBizCode(statusCode int) int64 {
	if statusCode >= http.StatusBadRequest {
		return int64(statusCode)
	}
	return 200
}

func truncateBizMsg(msg string) string {
	msg = strings.TrimSpace(msg)
	runes := []rune(msg)
	if len(runes) <= 200 {
		return msg
	}
	return string(runes[:200])
}

type accessLogResponseWriter struct {
	target      http.ResponseWriter
	header      http.Header
	statusCode  int
	body        bytes.Buffer
	wroteHeader bool
}

func newAccessLogResponseWriter(w http.ResponseWriter) *accessLogResponseWriter {
	headers := make(http.Header)
	if w != nil {
		for key, values := range w.Header() {
			copied := append([]string(nil), values...)
			headers[key] = copied
		}
	}
	return &accessLogResponseWriter{
		target:     w,
		header:     headers,
		statusCode: http.StatusOK,
	}
}

func (w *accessLogResponseWriter) Header() http.Header {
	return w.header
}

func (w *accessLogResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.wroteHeader = true
}

func (w *accessLogResponseWriter) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.body.Write(data)
}

func (w *accessLogResponseWriter) flush(requestID string, biz accessLogBizResult) {
	if w.target == nil {
		return
	}

	body := injectTraceIntoBody(w.body.Bytes(), requestID)
	targetHeader := w.target.Header()
	for key := range targetHeader {
		delete(targetHeader, key)
	}
	for key, values := range w.header {
		targetHeader[key] = append([]string(nil), values...)
	}
	if requestID != "" {
		targetHeader.Set(requestIDHeader, requestID)
	}
	if len(body) > 0 {
		targetHeader.Set("Content-Length", strconv.Itoa(len(body)))
	}
	w.target.WriteHeader(w.statusCode)
	if len(body) > 0 {
		_, _ = w.target.Write(body)
	}
}

func injectTraceIntoBody(body []byte, requestID string) []byte {
	if requestID == "" || len(bytes.TrimSpace(body)) == 0 {
		return body
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return body
	}
	if _, ok := payload["trace"]; !ok || strings.TrimSpace(toString(payload["trace"])) == "" {
		payload["trace"] = requestID
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return body
	}
	return encoded
}

func toString(v interface{}) string {
	switch value := v.(type) {
	case string:
		return value
	default:
		return ""
	}
}
