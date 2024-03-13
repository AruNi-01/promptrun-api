/*
 Navicat Premium Data Transfer

 Source Server         : local_mysql
 Source Server Type    : MySQL
 Source Server Version : 80033 (8.0.33)
 Source Host           : localhost:3306
 Source Schema         : promptrun

 Target Server Type    : MySQL
 Target Server Version : 80033 (8.0.33)
 File Encoding         : 65001

 Date: 12/03/2024 22:31:07
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for likes
-- ----------------------------
DROP TABLE IF EXISTS `likes`;
CREATE TABLE `likes`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NULL DEFAULT NULL COMMENT '触发喜欢的用户 id',
  `prompt_id` int NULL DEFAULT NULL COMMENT '被喜欢的提示词 id',
  `seller_user_id` int NULL DEFAULT NULL COMMENT '被喜欢的提示词的卖家用户 id',
  `create_time` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for prompt
-- ----------------------------
DROP TABLE IF EXISTS `prompt`;
CREATE TABLE `prompt`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `seller_id` int NULL DEFAULT NULL COMMENT '卖家 id，逻辑关联到卖家表',
  `model_id` int NULL DEFAULT NULL COMMENT '模型 id，逻辑关联模型表。1：GPT、2：Claude3、3：Midjourney、4：DALL·E、5：Stable Diffusion、6：Sora',
  `media_type` int NULL DEFAULT NULL COMMENT '媒体类型，0：文本、1：图片、2：视频',
  `category_id` int NULL DEFAULT NULL COMMENT '分类 id，逻辑关联分类表',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '提示词标题',
  `intro` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '提示词介绍',
  `price` decimal(10, 2) NULL DEFAULT NULL COMMENT '价格',
  `rating` float(2, 1) NULL DEFAULT NULL COMMENT '评分，1.0-5.0',
  `score` double NULL DEFAULT NULL COMMENT '分数，热度排行使用',
  `sell_amount` int NULL DEFAULT NULL COMMENT '销量',
  `browse_amount` int NULL DEFAULT NULL COMMENT '浏览数量',
  `like_amount` int NULL DEFAULT NULL COMMENT '提示词被喜欢数量',
  `publish_status` int NULL DEFAULT NULL COMMENT '上架状态，0: 下架，1：上架',
  `audit_status` int NULL DEFAULT NULL COMMENT '审核状态，0：未审核，1：审核中，2：审核通过',
  `create_time` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for prompt_img
-- ----------------------------
DROP TABLE IF EXISTS `prompt_img`;
CREATE TABLE `prompt_img`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `prompt_id` int NULL DEFAULT NULL COMMENT '提示词 id，逻辑关联到提示词表',
  `img_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `is_master` int NULL DEFAULT NULL COMMENT '是否主图，0：非主图，1：主图',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for seller
-- ----------------------------
DROP TABLE IF EXISTS `seller`;
CREATE TABLE `seller`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL COMMENT '逻辑外键，关联到用户 id',
  `rating` float(2, 1) NULL DEFAULT NULL COMMENT '卖家总评分，1.0-5.0',
  `status` int NULL DEFAULT NULL COMMENT '卖家状态，0: 禁用，1: 启用',
  `intro` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '卖家简介',
  `sell_amount` int NULL DEFAULT NULL COMMENT '销量',
  `like_amount` int NULL DEFAULT NULL COMMENT '提示词被喜欢数量',
  `create_time` datetime NULL DEFAULT NULL COMMENT '成为卖家时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `password` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `salt` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '加密盐，每个用户随机生成',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `header_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `type` int NULL DEFAULT NULL COMMENT '用户类型，0: 买家，1: 卖家。做个字段冗余，方便查询',
  `create_time` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
