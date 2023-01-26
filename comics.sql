/*
 Navicat Premium Data Transfer

 Source Server         : mydb
 Source Server Type    : MySQL
 Source Server Version : 50735
 Source Host           : localhost:3306
 Source Schema         : comics

 Target Server Type    : MySQL
 Target Server Version : 50735
 File Encoding         : 65001

 Date: 27/01/2023 07:42:38
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for source_chapter
-- ----------------------------
DROP TABLE IF EXISTS `source_chapter`;
CREATE TABLE `source_chapter`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `source` tinyint(1) NOT NULL DEFAULT 1 COMMENT '采集源 1:快看 2:腾讯',
  `source_id` int(11) NOT NULL COMMENT '源漫画id',
  `source_chapter_id` int(11) NOT NULL COMMENT '源章节id',
  `source_uri` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '源uri',
  `cover` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '封面',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '标题',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `source_chapter_id`(`source`, `source_id`, `source_chapter_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '采集-漫画章节' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of source_chapter
-- ----------------------------

-- ----------------------------
-- Table structure for source_comic
-- ----------------------------
DROP TABLE IF EXISTS `source_comic`;
CREATE TABLE `source_comic`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `source` tinyint(1) NOT NULL DEFAULT 1 COMMENT '采集源 1:快看 2:腾讯',
  `source_id` int(11) NOT NULL COMMENT '源漫画id',
  `source_uri` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '源uri',
  `cover` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '封面',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '标题',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `source_id`(`source`, `source_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '采集-漫画' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of source_comic
-- ----------------------------

-- ----------------------------
-- Table structure for source_image
-- ----------------------------
DROP TABLE IF EXISTS `source_image`;
CREATE TABLE `source_image`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `source` tinyint(1) NULL DEFAULT NULL COMMENT '采集源 1:快看 2:腾讯',
  `source_id` int(11) NOT NULL,
  `source_chapter_id` int(11) NOT NULL,
  `images` json NOT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of source_image
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;
