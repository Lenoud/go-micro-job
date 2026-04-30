-- micro_job 数据库初始化：department-service 所需的 b_department 表

CREATE DATABASE IF NOT EXISTS `micro_job` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `micro_job`;

DROP TABLE IF EXISTS `b_department`;
CREATE TABLE `b_department` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '部门名称',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '部门描述',
  `parent_id` int DEFAULT '0' COMMENT '上级部门ID, 0=顶级部门',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';

INSERT INTO `b_department` (`id`, `title`, `description`, `parent_id`, `create_time`) VALUES
(1,'开发部门','软件开发',0,'2026-04-15 16:50:28'),
(2,'人事部门','技术人员招聘',1,'2026-04-15 16:54:36');
