-- go_job 数据库：oplog-service 所需的 b_op_log 表（共用单体数据库）
USE `go_job`;

CREATE TABLE IF NOT EXISTS `b_op_log` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `request_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '请求ID',
  `user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '操作用户ID',
  `re_ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '请求IP',
  `re_ua` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT 'User-Agent',
  `re_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '请求URL',
  `re_method` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '请求方法',
  `re_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '请求参数(脱敏后)',
  `success` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '1' COMMENT '是否成功',
  `biz_code` bigint DEFAULT '0' COMMENT '业务响应码',
  `biz_msg` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '业务结果摘要',
  `request_time` bigint DEFAULT '0' COMMENT '请求时间戳',
  `response_ms` bigint DEFAULT '0' COMMENT '响应耗时(ms)',
  PRIMARY KEY (`id`),
  KEY `idx_request_id` (`request_id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_request_time` (`request_time`),
  KEY `idx_re_url_request_time` (`re_url`(191),`request_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';
