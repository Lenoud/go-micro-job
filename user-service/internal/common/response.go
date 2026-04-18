package common

import (
	"encoding/json"
	"user-service/user"
)

func apiResp(code int64, msg string, data interface{}) *user.ApiResponse {
	var dataStr string
	if data != nil {
		b, _ := json.Marshal(data)
		dataStr = string(b)
	}
	return &user.ApiResponse{
		Code: code,
		Msg:  msg,
		Data: dataStr,
	}
}

func Success(data interface{}) *user.ApiResponse {
	return apiResp(200, "success", data)
}

func SuccessMsg(msg string, data interface{}) *user.ApiResponse {
	return apiResp(200, msg, data)
}

func Fail(msg string) *user.ApiResponse {
	return apiResp(-1, msg, nil)
}

func SuccessPage(list interface{}, total, page, pageSize int64) *user.ApiResponse {
	return apiResp(200, "success", map[string]interface{}{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func SplitIDs(ids string) []string {
	if ids == "" {
		return nil
	}
	parts := make([]string, 0)
	for _, id := range splitComma(ids) {
		id = trimSpace(id)
		if id != "" {
			parts = append(parts, id)
		}
	}
	return parts
}

func splitComma(s string) []string {
	result := make([]string, 0)
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

func GenerateToken(username string) string {
	return Md5(username + Salt + "token")
}
