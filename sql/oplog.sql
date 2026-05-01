-- micro_job 数据库初始化：oplog-service 所需的 b_op_log 表

CREATE DATABASE IF NOT EXISTS `micro_job` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `micro_job`;

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
  `re_time` bigint DEFAULT '0' COMMENT '请求时间戳',
  `access_time` bigint DEFAULT '0' COMMENT '访问耗时',
  PRIMARY KEY (`id`),
  KEY `idx_request_id` (`request_id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_re_time` (`re_time`),
  KEY `idx_re_url_time` (`re_url`(191),`re_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';
