package common

// Msg 生成标准错误消息：op + target + "失败"
// 例：Msg("创建", "部门") → "创建部门失败"
func Msg(op, target string) string {
	return op + target + "失败"
}
