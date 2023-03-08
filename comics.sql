/*
 Navicat Premium Data Transfer

 Source Server         : loc
 Source Server Type    : MySQL
 Source Server Version : 50737
 Source Host           : 127.0.0.1:3306
 Source Schema         : comics

 Target Server Type    : MySQL
 Target Server Version : 50737
 File Encoding         : 65001

 Date: 13/02/2023 14:40:37
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for admin_menu
-- ----------------------------
DROP TABLE IF EXISTS `admin_menu`;
CREATE TABLE `admin_menu`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) NOT NULL DEFAULT 0,
  `order` int(11) NOT NULL DEFAULT 0,
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `uri` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `permission` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_menu
-- ----------------------------
INSERT INTO `admin_menu` VALUES (1, 0, 1, 'Dashboard', 'fa-bar-chart', '/', NULL, NULL, NULL);
INSERT INTO `admin_menu` VALUES (2, 0, 2, 'Admin', 'fa-tasks', '', NULL, NULL, NULL);
INSERT INTO `admin_menu` VALUES (3, 2, 3, 'Users', 'fa-users', 'auth/users', NULL, NULL, NULL);
INSERT INTO `admin_menu` VALUES (4, 2, 4, 'Roles', 'fa-user', 'auth/roles', NULL, NULL, NULL);
INSERT INTO `admin_menu` VALUES (5, 2, 5, 'Permission', 'fa-ban', 'auth/permissions', NULL, NULL, NULL);
INSERT INTO `admin_menu` VALUES (6, 2, 6, 'Menu', 'fa-bars', 'auth/menu', NULL, NULL, NULL);
INSERT INTO `admin_menu` VALUES (7, 2, 7, 'Operation log', 'fa-history', 'auth/logs', NULL, NULL, NULL);
INSERT INTO `admin_menu` VALUES (8, 0, 0, '爬虫管理', 'fa-bug', '/paw', NULL, '2023-02-13 09:57:14', '2023-02-13 09:57:14');

-- ----------------------------
-- Table structure for admin_operation_log
-- ----------------------------
DROP TABLE IF EXISTS `admin_operation_log`;
CREATE TABLE `admin_operation_log`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `ip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `input` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `admin_operation_log_user_id_index`(`user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 29 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_operation_log
-- ----------------------------
INSERT INTO `admin_operation_log` VALUES (1, 1, 'admin', 'GET', '172.17.0.1', '[]', '2023-02-13 09:32:03', '2023-02-13 09:32:03');
INSERT INTO `admin_operation_log` VALUES (2, 1, 'admin', 'GET', '172.17.0.1', '[]', '2023-02-13 09:35:37', '2023-02-13 09:35:37');
INSERT INTO `admin_operation_log` VALUES (3, 1, 'admin', 'GET', '172.17.0.1', '[]', '2023-02-13 09:36:14', '2023-02-13 09:36:14');
INSERT INTO `admin_operation_log` VALUES (4, 1, 'admin', 'GET', '172.17.0.1', '[]', '2023-02-13 09:36:29', '2023-02-13 09:36:29');
INSERT INTO `admin_operation_log` VALUES (5, 1, 'admin/auth/users', 'GET', '172.17.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-02-13 09:38:59', '2023-02-13 09:38:59');
INSERT INTO `admin_operation_log` VALUES (6, 1, 'admin', 'GET', '172.17.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-02-13 09:39:02', '2023-02-13 09:39:02');
INSERT INTO `admin_operation_log` VALUES (7, 1, 'admin', 'GET', '172.17.0.1', '[]', '2023-02-13 09:52:02', '2023-02-13 09:52:02');
INSERT INTO `admin_operation_log` VALUES (8, 1, 'admin/auth/users', 'GET', '172.17.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-02-13 09:52:08', '2023-02-13 09:52:08');
INSERT INTO `admin_operation_log` VALUES (9, 1, 'admin', 'GET', '172.17.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-02-13 09:52:11', '2023-02-13 09:52:11');
INSERT INTO `admin_operation_log` VALUES (10, 1, 'admin', 'GET', '172.17.0.1', '[]', '2023-02-13 09:56:48', '2023-02-13 09:56:48');
INSERT INTO `admin_operation_log` VALUES (11, 1, 'admin/auth/menu', 'GET', '172.17.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-02-13 09:56:51', '2023-02-13 09:56:51');
INSERT INTO `admin_operation_log` VALUES (12, 1, 'admin/auth/menu', 'POST', '172.17.0.1', '{\"parent_id\":\"0\",\"title\":\"\\u722c\\u866b\\u7ba1\\u7406\",\"icon\":\"fa-bug\",\"uri\":\"\\/paw\",\"roles\":[\"1\",null],\"permission\":null,\"_token\":\"yjukxyDm3Fsbb09y9evsgtgGj47qXyo05gxhtSYe\"}', '2023-02-13 09:57:14', '2023-02-13 09:57:14');
INSERT INTO `admin_operation_log` VALUES (13, 1, 'admin/auth/menu', 'GET', '172.17.0.1', '[]', '2023-02-13 09:57:14', '2023-02-13 09:57:14');
INSERT INTO `admin_operation_log` VALUES (14, 1, 'admin/auth/menu', 'GET', '172.17.0.1', '[]', '2023-02-13 09:57:16', '2023-02-13 09:57:16');
INSERT INTO `admin_operation_log` VALUES (15, 1, 'admin/paw', 'GET', '172.17.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-02-13 09:57:18', '2023-02-13 09:57:18');
INSERT INTO `admin_operation_log` VALUES (16, 1, 'admin/paw', 'GET', '172.17.0.1', '[]', '2023-02-13 10:05:30', '2023-02-13 10:05:30');
INSERT INTO `admin_operation_log` VALUES (17, 1, 'admin/paw', 'GET', '172.17.0.1', '[]', '2023-02-13 10:08:29', '2023-02-13 10:08:29');
INSERT INTO `admin_operation_log` VALUES (18, 1, 'admin/paw', 'GET', '172.17.0.1', '[]', '2023-02-13 10:09:29', '2023-02-13 10:09:29');
INSERT INTO `admin_operation_log` VALUES (19, 1, 'admin/paw', 'GET', '172.17.0.1', '[]', '2023-02-13 10:09:37', '2023-02-13 10:09:37');
INSERT INTO `admin_operation_log` VALUES (20, 1, 'admin/paw', 'GET', '172.17.0.1', '[]', '2023-02-13 10:10:11', '2023-02-13 10:10:11');
INSERT INTO `admin_operation_log` VALUES (21, 1, 'admin/paw', 'GET', '172.17.0.1', '[]', '2023-02-13 10:11:09', '2023-02-13 10:11:09');
INSERT INTO `admin_operation_log` VALUES (22, 1, 'admin/paw', 'GET', '172.17.0.1', '[]', '2023-02-13 10:11:38', '2023-02-13 10:11:38');
INSERT INTO `admin_operation_log` VALUES (23, 1, 'admin/paw', 'GET', '172.17.0.1', '[]', '2023-02-13 10:11:46', '2023-02-13 10:11:46');
INSERT INTO `admin_operation_log` VALUES (24, 1, 'admin/paw', 'GET', '172.17.0.1', '[]', '2023-02-13 10:12:00', '2023-02-13 10:12:00');
INSERT INTO `admin_operation_log` VALUES (25, 1, 'admin/paw', 'GET', '172.17.0.1', '[]', '2023-02-13 10:12:13', '2023-02-13 10:12:13');
INSERT INTO `admin_operation_log` VALUES (26, 1, 'admin/paw', 'GET', '172.17.0.1', '[]', '2023-02-13 10:12:42', '2023-02-13 10:12:42');
INSERT INTO `admin_operation_log` VALUES (27, 1, 'admin/paw', 'GET', '172.17.0.1', '[]', '2023-02-13 10:21:18', '2023-02-13 10:21:18');
INSERT INTO `admin_operation_log` VALUES (28, 1, 'admin/paw', 'GET', '172.17.0.1', '[]', '2023-02-13 10:22:07', '2023-02-13 10:22:07');

-- ----------------------------
-- Table structure for admin_permissions
-- ----------------------------
DROP TABLE IF EXISTS `admin_permissions`;
CREATE TABLE `admin_permissions`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `slug` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `http_method` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `http_path` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `admin_permissions_name_unique`(`name`) USING BTREE,
  UNIQUE INDEX `admin_permissions_slug_unique`(`slug`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_permissions
-- ----------------------------
INSERT INTO `admin_permissions` VALUES (1, 'All permission', '*', '', '*', NULL, NULL);
INSERT INTO `admin_permissions` VALUES (2, 'Dashboard', 'dashboard', 'GET', '/', NULL, NULL);
INSERT INTO `admin_permissions` VALUES (3, 'Login', 'auth.login', '', '/auth/login\r\n/auth/logout', NULL, NULL);
INSERT INTO `admin_permissions` VALUES (4, 'User setting', 'auth.setting', 'GET,PUT', '/auth/setting', NULL, NULL);
INSERT INTO `admin_permissions` VALUES (5, 'Auth management', 'auth.management', '', '/auth/roles\r\n/auth/permissions\r\n/auth/menu\r\n/auth/logs', NULL, NULL);

-- ----------------------------
-- Table structure for admin_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `admin_role_menu`;
CREATE TABLE `admin_role_menu`  (
  `role_id` int(11) NOT NULL,
  `menu_id` int(11) NOT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  INDEX `admin_role_menu_role_id_menu_id_index`(`role_id`, `menu_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_role_menu
-- ----------------------------
INSERT INTO `admin_role_menu` VALUES (1, 2, NULL, NULL);
INSERT INTO `admin_role_menu` VALUES (1, 8, NULL, NULL);

-- ----------------------------
-- Table structure for admin_role_permissions
-- ----------------------------
DROP TABLE IF EXISTS `admin_role_permissions`;
CREATE TABLE `admin_role_permissions`  (
  `role_id` int(11) NOT NULL,
  `permission_id` int(11) NOT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  INDEX `admin_role_permissions_role_id_permission_id_index`(`role_id`, `permission_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_role_permissions
-- ----------------------------
INSERT INTO `admin_role_permissions` VALUES (1, 1, NULL, NULL);

-- ----------------------------
-- Table structure for admin_role_users
-- ----------------------------
DROP TABLE IF EXISTS `admin_role_users`;
CREATE TABLE `admin_role_users`  (
  `role_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  INDEX `admin_role_users_role_id_user_id_index`(`role_id`, `user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_role_users
-- ----------------------------
INSERT INTO `admin_role_users` VALUES (1, 1, NULL, NULL);

-- ----------------------------
-- Table structure for admin_roles
-- ----------------------------
DROP TABLE IF EXISTS `admin_roles`;
CREATE TABLE `admin_roles`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `slug` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `admin_roles_name_unique`(`name`) USING BTREE,
  UNIQUE INDEX `admin_roles_slug_unique`(`slug`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_roles
-- ----------------------------
INSERT INTO `admin_roles` VALUES (1, 'Administrator', 'administrator', '2023-02-13 09:31:24', '2023-02-13 09:31:24');

-- ----------------------------
-- Table structure for admin_user_permissions
-- ----------------------------
DROP TABLE IF EXISTS `admin_user_permissions`;
CREATE TABLE `admin_user_permissions`  (
  `user_id` int(11) NOT NULL,
  `permission_id` int(11) NOT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  INDEX `admin_user_permissions_user_id_permission_id_index`(`user_id`, `permission_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for admin_users
-- ----------------------------
DROP TABLE IF EXISTS `admin_users`;
CREATE TABLE `admin_users`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(190) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `remember_token` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `admin_users_username_unique`(`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_users
-- ----------------------------
INSERT INTO `admin_users` VALUES (1, 'admin', '$2y$10$OFhcLbWA/AJZYadqUg5ClONU6MhhLIJdzNuCKHo4cA9qU05gDVvGq', 'Administrator', NULL, '0Lmx9L3E37NIYVBvgw45uAqCJzXNAW0L3Z377EJ5zWE3cPW9AoZhYnPQZCpY', '2023-02-13 09:31:24', '2023-02-13 09:31:24');

-- ----------------------------
-- Table structure for fail_info
-- ----------------------------
DROP TABLE IF EXISTS `fail_info`;
CREATE TABLE `fail_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `source` tinyint(1) NOT NULL DEFAULT 1 COMMENT '采集源 1:快看 2:腾讯',
  `type` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0漫画列表 1漫画 2章节 3图片',
  `err` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '错误关键词',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '地址',
  `info` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '失败信息记录',
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `type`(`type`) USING BTREE,
  INDEX `err`(`err`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for failed_jobs
-- ----------------------------
DROP TABLE IF EXISTS `failed_jobs`;
CREATE TABLE `failed_jobs`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `connection` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `queue` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `payload` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `exception` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `failed_at` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for migrations
-- ----------------------------
DROP TABLE IF EXISTS `migrations`;
CREATE TABLE `migrations`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `migration` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `batch` int(11) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of migrations
-- ----------------------------
INSERT INTO `migrations` VALUES (1, '2014_10_12_000000_create_users_table', 1);
INSERT INTO `migrations` VALUES (2, '2014_10_12_100000_create_password_resets_table', 1);
INSERT INTO `migrations` VALUES (3, '2019_08_19_000000_create_failed_jobs_table', 1);
INSERT INTO `migrations` VALUES (4, '2016_01_04_173148_create_admin_tables', 2);

-- ----------------------------
-- Table structure for password_resets
-- ----------------------------
DROP TABLE IF EXISTS `password_resets`;
CREATE TABLE `password_resets`  (
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  INDEX `password_resets_email_index`(`email`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for source_chapter
-- ----------------------------
DROP TABLE IF EXISTS `source_chapter`;
CREATE TABLE `source_chapter` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `comic_id` int(11) NOT NULL,
  `source` tinyint(1) NOT NULL DEFAULT '1' COMMENT '采集源 1:快看 2:腾讯',
  `source_chapter_id` int(11) NOT NULL COMMENT '源章节id',
  `source_url` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '源url',
  `cover` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '封面',
  `title` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标题',
  `sort` int(11) NOT NULL DEFAULT '0',
  `is_free` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0免费 1收费',
  `source_data` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0未审核 1通过',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `source` (`source`,`comic_id`,`source_url`) USING BTREE,
  KEY `comic_id` (`comic_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='采集-漫画章节';

-- ----------------------------
-- Table structure for source_comic
-- ----------------------------
DROP TABLE IF EXISTS `source_comic`;
CREATE TABLE `source_comic` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `source` tinyint(1) NOT NULL DEFAULT '1' COMMENT '采集源 1:快看 2:腾讯',
  `source_id` int(11) NOT NULL COMMENT '源漫画id',
  `source_url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '源url',
  `cover` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '封面',
  `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标题',
  `author` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '作者',
  `label` json NOT NULL COMMENT '标签',
  `category` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类',
  `region` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '地区',
  `chapter_count` int(11) NOT NULL DEFAULT '0' COMMENT '章节数量',
  `like` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '喜欢',
  `popularity` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0' COMMENT '人气热度',
  `is_free` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0免费 1收费',
  `is_finish` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0连载 1完结',
  `description` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '描述',
  `source_data` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '源数据',
  `chapter_pick` int(11) NOT NULL DEFAULT '0' COMMENT '章节拨片',
  `retry` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0默认 1重抓',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0未审核 1通过',
  `last_chapter_update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最新章节更新时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `source_id` (`source`,`source_id`) USING BTREE,
  UNIQUE KEY `source_uri` (`source_url`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='采集-漫画';

-- ----------------------------
-- Table structure for source_image
-- ----------------------------
DROP TABLE IF EXISTS `source_image`;
CREATE TABLE `source_image` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `source` tinyint(1) DEFAULT '0' COMMENT '采集源 1:快看 2:腾讯',
  `comic_id` int(11) NOT NULL DEFAULT '0',
  `chapter_id` int(11) NOT NULL,
  `images` json NOT NULL,
  `source_data` json NOT NULL,
  `state` tinyint(1) NOT NULL DEFAULT '0' COMMENT '资源获取:0未开始 1已完成',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `chapter_id` (`chapter_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `email_verified_at` timestamp(0) NULL DEFAULT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `remember_token` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `users_email_unique`(`email`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

DROP TABLE IF EXISTS `comic`;
CREATE TABLE `comic` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `source_comic_id` int(11) NOT NULL COMMENT '采集-漫画id',
  `cover` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '封面',
  `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标题',
  `author` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '作者',
  `label` json NOT NULL COMMENT '标签',
  `category` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类',
  `region` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '地区',
  `chapter_count` int(11) NOT NULL DEFAULT '0' COMMENT '章节数量',
  `like` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '喜欢',
  `popularity` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0' COMMENT '人气热度',
  `is_free` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0免费 1收费',
  `is_finish` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0连载 1完结',
  `description` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '描述',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `comic_id` (`source_comic_id`) USING BTREE,
  KEY `title` (`title`) USING HASH,
  KEY `category` (`category`) USING HASH,
  KEY `region` (`region`) USING HASH,
  KEY `is_free` (`is_free`) USING BTREE,
  KEY `is_finish` (`is_finish`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='漫画';

DROP TABLE IF EXISTS `chapter`;
CREATE TABLE `chapter` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `comic_id` int(11) NOT NULL COMMENT '漫画id',
  `source_comic_id` int(11) NOT NULL COMMENT '采集-漫画id',
  `source_chapter_id` int(11) NOT NULL COMMENT '采集-章节id',
  `source_image_id` int(11) NOT NULL COMMENT '采集-图片id',
  `cover` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '封面',
  `title` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标题',
  `sort` int(11) NOT NULL DEFAULT '0',
  `is_free` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0免费 1收费',
  `images` json NOT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `source_chapter_id` (`source_chapter_id`),
  UNIQUE KEY `source_image_id` (`source_image_id`),
  KEY `comic_id` (`comic_id`) USING BTREE,
  KEY `is_free` (`is_free`) USING BTREE,
  KEY `title` (`title`) USING HASH
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='章节';

SET FOREIGN_KEY_CHECKS = 1;
