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

 Date: 31/01/2023 15:12:16
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for source_chapter
-- ----------------------------
DROP TABLE IF EXISTS `source_chapter`;
CREATE TABLE `source_chapter`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `comic_id` int(11) NOT NULL,
  `source` tinyint(1) NOT NULL DEFAULT 1 COMMENT '采集源 1:快看 2:腾讯',
  `source_chapter_id` int(11) NOT NULL COMMENT '源章节id',
  `source_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '源url',
  `cover` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '封面',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标题',
  `sort` int(11) NOT NULL DEFAULT 0,
  `is_free` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0免费 1收费',
  `source_data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `source_chapter_id`(`source`, `source_chapter_id`) USING BTREE,
  UNIQUE INDEX `source_uri`(`source_url`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '采集-漫画章节' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for source_comic
-- ----------------------------
DROP TABLE IF EXISTS `source_comic`;
CREATE TABLE `source_comic`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `source` tinyint(1) NOT NULL DEFAULT 1 COMMENT '采集源 1:快看 2:腾讯',
  `source_id` int(11) NOT NULL COMMENT '源漫画id',
  `source_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '源url',
  `cover` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '封面',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标题',
  `author` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '作者',
  `category` json NOT NULL COMMENT '分类',
  `chapter_count` int(11) NOT NULL DEFAULT 0 COMMENT '章节数量',
  `like_count` int(11) NOT NULL DEFAULT 0 COMMENT '喜欢',
  `popularity` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '人气热度',
  `is_free` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0免费 1收费',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '描述',
  `source_data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '源数据',
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `source_id`(`source`, `source_id`) USING BTREE,
  UNIQUE INDEX `source_uri`(`source_url`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '采集-漫画' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for source_image
-- ----------------------------
DROP TABLE IF EXISTS `source_image`;
CREATE TABLE `source_image`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `chapter_id` int(11) NOT NULL,
  `images` json NOT NULL,
  `source_data` json NOT NULL,
  `state` tinyint(1) NOT NULL DEFAULT 0 COMMENT '资源获取:0未开始 1已完成',
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
