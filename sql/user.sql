-- micro_job 数据库初始化：仅 user-service 所需的 b_user 表
-- 数据库名称：micro_job

CREATE DATABASE IF NOT EXISTS `micro_job` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `micro_job`;

DROP TABLE IF EXISTS `b_user`;
CREATE TABLE `b_user` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '登录用户名',
  `password` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '密码(MD5+盐)',
  `nickname` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '昵称',
  `mobile` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '手机号',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '邮箱',
  `role` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '1' COMMENT '角色: 1=求职者, 2=HR, 3=管理员',
  `status` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '状态: 0=正常, 1=禁用',
  `token` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '登录令牌',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `push_email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '推送邮箱',
  `push_switch` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '推送开关: 0=关, 1=开',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

INSERT INTO `b_user` VALUES
(1,'admin','6658ed9eaa4108e30a9d04a008f2c2e9','系统管理员','13800138000','admin@maka.com','3','0','','2026-04-15 16:49:35','admin@maka.com','1'),
(2,'skyrisai','6658ed9eaa4108e30a9d04a008f2c2e9','玛咖HR-刘经理','13911112222','hr_maka@163.com','2','0','','2026-04-15 16:49:35','liubiao351719672@gmail.com','1'),
(3,'351719672@qq.com','6658ed9eaa4108e30a9d04a008f2c2e9','张求职','13733334444','351719672@qq.com','1','0','','2026-04-15 16:49:35','351719672@qq.com','1'),
(4,'skyrisai2','6658ed9eaa4108e30a9d04a008f2c2e9','屋里(技术面试官)','13655556666','hr_dev@maka.com','2','0','','2026-04-15 16:49:35','','0'),
(5,'liubiao351719672@gmail.com','6658ed9eaa4108e30a9d04a008f2c2e9','李候选','13577778888','liubiao351719672@gmail.com','1','0','','2026-04-15 16:49:35','','0');
