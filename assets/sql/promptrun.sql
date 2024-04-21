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

 Date: 21/04/2024 11:03:43
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
  `seller_id` int NULL DEFAULT NULL COMMENT '被喜欢的提示词的卖家 id',
  `create_time` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 37 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of likes
-- ----------------------------
INSERT INTO `likes` VALUES (11, 4, 1, 1, '2024-04-12 00:53:51');
INSERT INTO `likes` VALUES (31, 4, 3, 1, '2024-04-12 22:12:22');
INSERT INTO `likes` VALUES (32, 4, 7, 2, '2024-04-12 23:09:07');
INSERT INTO `likes` VALUES (33, 4, 17, 1, '2024-04-15 00:20:37');
INSERT INTO `likes` VALUES (34, 4, 18, 1, '2024-04-16 01:19:45');
INSERT INTO `likes` VALUES (36, 4, 20, 1, '2024-04-20 16:55:48');

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `from_user_id` int NULL DEFAULT NULL,
  `to_user_id` int NULL DEFAULT NULL,
  `type` smallint NULL DEFAULT NULL COMMENT '消息类型：1-活动，2-售出，3-审核，4-点赞，5-提现',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `isRead` smallint NULL DEFAULT NULL COMMENT '是否已读：0-未读，1-已读',
  `create_time` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of message
-- ----------------------------

-- ----------------------------
-- Table structure for model
-- ----------------------------
DROP TABLE IF EXISTS `model`;
CREATE TABLE `model`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '模型名称',
  `media_type` int NULL DEFAULT NULL COMMENT '媒体类型，1：文本、2：图片、3：视频',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of model
-- ----------------------------
INSERT INTO `model` VALUES (1, 'GPT', 1);
INSERT INTO `model` VALUES (2, 'Claude 3', 1);
INSERT INTO `model` VALUES (3, 'Midjourney', 2);
INSERT INTO `model` VALUES (4, 'DALL·E', 2);
INSERT INTO `model` VALUES (5, 'Stable Diffusion', 2);

-- ----------------------------
-- Table structure for order
-- ----------------------------
DROP TABLE IF EXISTS `order`;
CREATE TABLE `order`  (
  `id` bigint NOT NULL,
  `prompt_id` int NULL DEFAULT NULL COMMENT '购买的提示词 id',
  `seller_id` int NULL DEFAULT NULL COMMENT '提示词所属的卖家 id',
  `buyer_id` int NULL DEFAULT NULL COMMENT '购买提示词的买家 id（等于 user_id）',
  `price` decimal(10, 2) NULL DEFAULT NULL COMMENT '买入的价格',
  `is_rating` int NULL DEFAULT NULL COMMENT '买家个人是否已进行评分，0：否，1：是',
  `rating` float(2, 1) NULL DEFAULT NULL COMMENT '买家个人对该提示词的评分，1.0-5.0',
  `create_time` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of order
-- ----------------------------
INSERT INTO `order` VALUES (1, 1, 1, 4, 4.99, 0, NULL, '2023-06-26 23:47:22');
INSERT INTO `order` VALUES (2, 2, 1, 4, 12.99, 0, NULL, '2023-07-26 23:27:22');
INSERT INTO `order` VALUES (3, 3, 1, 3, 4.99, 0, NULL, '2023-09-26 10:47:22');
INSERT INTO `order` VALUES (4, 4, 1, 4, 4.99, 0, NULL, '2023-09-26 23:47:22');
INSERT INTO `order` VALUES (5, 6, 1, 4, 5.99, 0, NULL, '2023-09-22 23:47:22');
INSERT INTO `order` VALUES (6, 5, 1, 3, 4.99, 0, NULL, '2023-11-26 23:47:22');
INSERT INTO `order` VALUES (7, 7, 1, 4, 4.99, 0, NULL, '2023-12-26 23:47:22');
INSERT INTO `order` VALUES (8, 8, 1, 3, 7.99, 0, NULL, '2024-02-26 23:47:22');
INSERT INTO `order` VALUES (9, 9, 1, 4, 4.99, 0, NULL, '2024-02-26 23:47:22');
INSERT INTO `order` VALUES (10, 10, 1, 4, 2.99, 0, NULL, '2024-04-02 21:37:22');
INSERT INTO `order` VALUES (11, 11, 1, 4, 4.99, 0, NULL, '2024-04-06 23:47:22');
INSERT INTO `order` VALUES (12, 20, 1, 4, 3.89, NULL, NULL, '2024-04-19 23:47:22');
INSERT INTO `order` VALUES (13, 21, 1, 4, 4.99, 1, 5.0, '2024-04-19 23:57:22');

-- ----------------------------
-- Table structure for order_rating
-- ----------------------------
DROP TABLE IF EXISTS `order_rating`;
CREATE TABLE `order_rating`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `order_id` bigint NULL DEFAULT NULL,
  `prompt_id` int NULL DEFAULT NULL,
  `seller_id` int NULL DEFAULT NULL,
  `rating` float(2, 1) NULL DEFAULT NULL,
  `create_time` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of order_rating
-- ----------------------------
INSERT INTO `order_rating` VALUES (1, 13, 21, 1, 5.0, '2024-04-20 16:38:01');

-- ----------------------------
-- Table structure for prompt
-- ----------------------------
DROP TABLE IF EXISTS `prompt`;
CREATE TABLE `prompt`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `seller_id` int NULL DEFAULT NULL COMMENT '卖家 id，逻辑关联到卖家表',
  `model_id` int NULL DEFAULT NULL COMMENT '模型 id，逻辑关联模型表。',
  `category_type` int NULL DEFAULT NULL COMMENT '分类类型，具体枚举看文档',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '提示词标题',
  `intro` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '提示词介绍',
  `input_example` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT 'Prompt 输入示例（文本类模型媒体）',
  `output_example` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT 'Prompt 输出示例（文本类模型媒体）',
  `price` decimal(10, 2) NULL DEFAULT NULL COMMENT '价格',
  `rating` float(2, 1) NULL DEFAULT NULL COMMENT '评分，1.0-5.0',
  `score` double NULL DEFAULT NULL COMMENT '分数，热度排行使用',
  `sell_amount` int NULL DEFAULT NULL COMMENT '销量',
  `browse_amount` int NULL DEFAULT NULL COMMENT '浏览数量',
  `like_amount` int NULL DEFAULT NULL COMMENT '提示词被喜欢数量',
  `publish_status` int NULL DEFAULT NULL COMMENT '上架状态，0: 下架，1：上架',
  `audit_status` int NULL DEFAULT NULL COMMENT '审核状态，0：审核失败，1：审核中，2：审核通过',
  `create_time` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 23 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of prompt
-- ----------------------------
INSERT INTO `prompt` VALUES (1, 1, 1, 1, '电影评价助手', '我是一个电影评价小助手，快来用我写影评吧！', NULL, NULL, 4.99, 4.4, 2, 2, 266, 30, 1, 2, '2024-03-23 15:18:55');
INSERT INTO `prompt` VALUES (2, 1, 2, 2, '小说评价助手', '我是一个小说评价小助手，快来用我写小说评价吧！', NULL, NULL, 4.99, 4.7, 5, 5, 0, 0, 1, 2, '2024-03-20 15:18:55');
INSERT INTO `prompt` VALUES (3, 1, 1, 1, '小说评价助手', '我是一个小说评价小助手，快来用我写小说评价吧！', NULL, NULL, 4.99, 3.2, 2, 5, 22, 1, 1, 2, '2024-03-24 15:18:55');
INSERT INTO `prompt` VALUES (4, 2, 2, 3, '小说评价助手', '我是一个小说评价小助手，快来用我写小说评价吧！', NULL, NULL, 4.99, 3.3, 2, 5, 0, 0, 0, 2, '2024-03-24 15:18:55');
INSERT INTO `prompt` VALUES (5, 1, 1, 1, '小说评价助手', '我是一个小说评价小助手，快来用我写小说评价吧！', NULL, NULL, 4.99, 4.2, 2, 5, 63, 0, 1, 2, '2024-03-24 15:18:55');
INSERT INTO `prompt` VALUES (6, 1, 3, 2, '小说评价助手', '我是一个小说评价小助手，快来用我写小说评价吧！', NULL, NULL, 4.99, 1.4, 2, 5, 1, 0, 1, 1, '2024-03-24 15:18:55');
INSERT INTO `prompt` VALUES (7, 2, 1, 1, '小说评价助手', '我是一个小说评价小助手，快来用我写小说评价吧！', NULL, NULL, 4.99, 0.4, 2, 5, 2, 1, 1, 2, '2024-03-24 15:18:55');
INSERT INTO `prompt` VALUES (8, 1, 5, 1, '小说评价助手', '我是一个小说评价小助手，快来用我写小说评价吧！', NULL, NULL, 4.99, 3.7, 2, 5, 1, 0, 0, 2, '2024-03-24 15:18:55');
INSERT INTO `prompt` VALUES (9, 1, 1, 1, '小说评价助手', '我是一个小说评价小助手，快来用我写小说评价吧！', NULL, NULL, 4.99, 5.0, 2, 5, 0, 0, 1, 2, '2024-03-24 15:18:55');
INSERT INTO `prompt` VALUES (10, 2, 4, 2, '小说评价助手', '我是一个小说评价小助手，快来用我写小说评价吧！', NULL, NULL, 4.99, 4.9, 2, 5, 0, 0, 1, 1, '2024-03-24 15:18:55');
INSERT INTO `prompt` VALUES (11, 1, 1, 1, '小说评价助手', '我是一个小说评价小助手，快来用我写小说评价吧！', NULL, NULL, 4.99, 5.0, 2, 5, 0, 0, 1, 2, '2024-03-24 15:18:55');
INSERT INTO `prompt` VALUES (12, 1, 1, 3, '小说评价助手', '我是一个小说评价小助手，快来用我写小说评价吧！', NULL, NULL, 4.99, 4.1, 2, 5, 0, 0, 1, 0, '2024-03-24 15:18:55');
INSERT INTO `prompt` VALUES (13, 2, 5, 1, '小说评价助手', '我是一个小说评价小助手，快来用我写小说评价吧！', NULL, NULL, 4.99, 5.0, 2, 5, 0, 0, 1, 2, '2024-03-24 15:18:55');
INSERT INTO `prompt` VALUES (14, 1, 1, 1, '小说评价助手', '我是一个小说评价小助手，快来用我写小说评价吧！', NULL, NULL, 4.99, 5.0, 2, 5, 0, 0, 1, 2, '2024-03-24 15:18:55');
INSERT INTO `prompt` VALUES (16, 1, 2, 2, '测试 Prompt001', '这是一个测试 Prompt。', '输入示例测试', '输出示例测试', 12.99, 0.0, 0, 0, 4, 0, 0, 2, '2024-04-14 23:10:18');
INSERT INTO `prompt` VALUES (17, 1, 2, 1, '测试Prompt测试测试测试', '测试测试测试测试测试测试测试测试测试测试测试测试', '输入示例测试', '输出示例测试', 7.99, 0.0, 0, 0, 6, 1, 1, 2, '2024-04-15 00:19:09');
INSERT INTO `prompt` VALUES (18, 1, 5, 1, '背景生成器', '背景生成器测试测试测试', NULL, NULL, 8.00, 0.0, 0, 0, 12, 1, 1, 2, '2024-04-16 01:18:59');
INSERT INTO `prompt` VALUES (19, 1, 2, 1, '测试 Prompt 4/16', '测试测试测试测试测试', '输入示例测试', '输出示例测试', 9.49, 0.0, 0, 0, 30, 0, 0, 2, '2024-04-16 01:36:01');
INSERT INTO `prompt` VALUES (20, 1, 1, 2, '文本写作助手', '我是一个文本写作助手，可以辅助您完成任何文本类相关的写作需求。', '我需要写一部浪漫的爱情电影，故事发生在巴黎。', '我很高兴能够扮演编剧的角色，为您创作一部浪漫的爱情电影。让我们开始吧！\n\n故事背景：爱丽舍宫、卢浮宫、圣母院、塞纳河畔……这些美丽的地方都是巴黎的标志性景点。我们的故事也将发生在这个浪漫的城市里。\n\n角色塑造：男主角是一位年轻的画家，他来自美国，来到巴黎是为了追寻自己的梦想。女主角是一位法国女孩，是一名古董商人，她经常在巴黎的艺术市场上寻找稀有的古董。\n\n故事情节：男主角在巴黎的街头遇到了女主角，两人一见钟情。男主角是一个非常有才华的画家，他在巴黎的街头为游客画肖像，同时也在寻找灵感创作自己的画作。女主角是一个非常有眼光的古董商人，她对艺术品有着独特的鉴赏能力。\n\n两人在一次偶然的相遇中，开始了一段浪漫的爱情故事。男主角为女主角画了一幅肖像，女主角也为男主角提供了一些寻找灵感的建议。两人开始一起游览巴黎的美景，一起探索巴黎的艺术氛围。他们的关系越来越亲密，但是男主角的签证即将过期，他必须回到美国。\n\n女主角为男主角举办了一场告别派对，男主角在派对上展示了自己的画作，这些画作都是他在巴黎创作的。女主角被男主角的才华和深情所打动，她向男主角表白，并且邀请他留在巴黎。男主角想要留下来，但是他还有一些事情需要处理，他必须回到美国。\n\n在回美国的途中，男主角意识到自己真正的梦想是留在巴黎，和女主角在一起。他决定回到巴黎，和女主角开始新的生活。最终，在塞纳河畔，男主角为女主角画了一幅肖像，两人在巴黎的美景中幸福地生活着。\n\n这就是我为您创作的浪漫爱情电影，希望您喜欢！', 3.89, 5.0, 0, 0, 44, 1, 1, 2, '2024-04-18 00:04:46');
INSERT INTO `prompt` VALUES (21, 1, 1, 2, '文本写作助手2', '我是一个文本写作助手，可以辅助您完成任何文本类相关的写作需求。', '我需要写一部浪漫的爱情电影，故事发生在巴黎。', '我很高兴能够扮演编剧的角色，为您创作一部浪漫的爱情电影。\n\n让我们开始吧！\n\n\n**故事背景**：\n\n- 爱丽舍宫、卢浮宫、圣母院、塞纳河畔……这些美丽的地方都是巴黎的标志性景点。我们的故事也将发生在这个浪漫的城市里。\n\n**角色塑造**：\n\n- 男主角是一位年轻的画家，他来自美国，来到巴黎是为了追寻自己的梦想。女主角是一位法国女孩，是一名古董商人，她经常在巴黎的艺术市场上寻找稀有的古董。\n\n**故事情节**：\n\n1. 男主角在巴黎的街头遇到了女主角，两人一见钟情。男主角是一个非常有才华的画家，他在巴黎的街头为游客画肖像，同时也在寻找灵感创作自己的画作。女主角是一个非常有眼光的古董商人，她对艺术品有着独特的鉴赏能力。\n\n2. 两人在一次偶然的相遇中，开始了一段浪漫的爱情故事。男主角为女主角画了一幅肖像，女主角也为男主角提供了一些寻找灵感的建议。两人开始一起游览巴黎的美景，一起探索巴黎的艺术氛围。他们的关系越来越亲密，但是男主角的签证即将过期，他必须回到美国。\n\n3. 女主角为男主角举办了一场告别派对，男主角在派对上展示了自己的画作，这些画作都是他在巴黎创作的。女主角被男主角的才华和深情所打动，她向男主角表白，并且邀请他留在巴黎。男主角想要留下来，但是他还有一些事情需要处理，他必须回到美国。\n\n4. 在回美国的途中，男主角意识到自己真正的梦想是留在巴黎，和女主角在一起。他决定回到巴黎，和女主角开始新的生活。最终，在塞纳河畔，男主角为女主角画了一幅肖像，两人在巴黎的美景中幸福地生活着。\n\n\n这就是我为您创作的浪漫爱情电影，希望您喜欢！\n', 4.99, 5.0, 0, 0, 33, 0, 1, 2, '2024-04-18 00:28:44');
INSERT INTO `prompt` VALUES (22, 2, 1, 1, 'Prompt 测试', '测试测试测试测试测试测试', '测试', '测试', 0.01, 5.0, 0, 0, 21, 0, 1, 2, '2024-04-21 01:05:02');

-- ----------------------------
-- Table structure for prompt_detail
-- ----------------------------
DROP TABLE IF EXISTS `prompt_detail`;
CREATE TABLE `prompt_detail`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `prompt_id` int NULL DEFAULT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT 'Prompt 具体内容',
  `media_type` int NULL DEFAULT NULL,
  `use_suggestion` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '使用建议',
  `create_time` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of prompt_detail
-- ----------------------------
INSERT INTO `prompt_detail` VALUES (1, 16, 'Prompt 内容测试', 1, '使用建议测试', '2024-04-14 23:10:18');
INSERT INTO `prompt_detail` VALUES (2, 17, '测试测试测试测试测试测试', 1, '测试测试测试', '2024-04-15 00:19:09');
INSERT INTO `prompt_detail` VALUES (3, 19, '测试', 1, '测试', '2024-04-16 01:36:02');
INSERT INTO `prompt_detail` VALUES (4, 20, '我需要写一些文本类的文章，请你现在充当一个文学专家，对我提出的要求进行写作，不要生硬，要富有情感，内容有跌宕起伏的感觉。', 1, '暂无。', '2024-04-18 00:04:46');
INSERT INTO `prompt_detail` VALUES (5, 21, '我是一个文本写作助手，可以辅助您完成任何文本类相关的写作需求。', 1, '暂无。', '2024-04-18 00:28:44');
INSERT INTO `prompt_detail` VALUES (6, 22, '测试', 1, '测试', '2024-04-21 01:05:06');

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
) ENGINE = InnoDB AUTO_INCREMENT = 26 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of prompt_img
-- ----------------------------
INSERT INTO `prompt_img` VALUES (1, 1, 'https://run-notes.oss-cn-beijing.aliyuncs.com/notes/202304202145141.png', 0);
INSERT INTO `prompt_img` VALUES (2, 1, 'https://www.educative.io/api/page/5819622709264384/image/download/4750982421413888', 0);
INSERT INTO `prompt_img` VALUES (3, 2, 'https://contentstatic.techgig.com/photo/82278297/5-top-advantages-of-using-golang-programming-language.jpg?35743', 1);
INSERT INTO `prompt_img` VALUES (4, 2, 'https://www.fotolog.com/wp-content/uploads/2020/07/Golang-1-750x456.jpg', 0);
INSERT INTO `prompt_img` VALUES (5, 1, 'https://draft.dev/learn/assets/posts/golang.png', 1);
INSERT INTO `prompt_img` VALUES (6, 3, 'https://draft.dev/learn/assets/posts/golang.png', 1);
INSERT INTO `prompt_img` VALUES (7, 4, 'https://draft.dev/learn/assets/posts/golang.png', 1);
INSERT INTO `prompt_img` VALUES (8, 5, 'https://draft.dev/learn/assets/posts/golang.png', 1);
INSERT INTO `prompt_img` VALUES (9, 6, 'https://draft.dev/learn/assets/posts/golang.png', 1);
INSERT INTO `prompt_img` VALUES (10, 7, 'https://draft.dev/learn/assets/posts/golang.png', 1);
INSERT INTO `prompt_img` VALUES (11, 8, 'https://draft.dev/learn/assets/posts/golang.png', 1);
INSERT INTO `prompt_img` VALUES (12, 9, 'https://draft.dev/learn/assets/posts/golang.png', 1);
INSERT INTO `prompt_img` VALUES (13, 10, 'https://draft.dev/learn/assets/posts/golang.png', 1);
INSERT INTO `prompt_img` VALUES (14, 11, 'https://draft.dev/learn/assets/posts/golang.png', 1);
INSERT INTO `prompt_img` VALUES (15, 12, 'https://draft.dev/learn/assets/posts/golang.png', 1);
INSERT INTO `prompt_img` VALUES (16, 13, 'https://draft.dev/learn/assets/posts/golang.png', 1);
INSERT INTO `prompt_img` VALUES (17, 14, 'https://draft.dev/learn/assets/posts/golang.png', 1);
INSERT INTO `prompt_img` VALUES (18, 17, 'https://promptrun.oss-cn-shanghai.aliyuncs.com/prompt_img/17-2024-04-15_001908.png', 1);
INSERT INTO `prompt_img` VALUES (19, 18, 'https://promptrun.oss-cn-shanghai.aliyuncs.com/prompt_img/18-2024-04-16_011858.png', 1);
INSERT INTO `prompt_img` VALUES (20, 18, 'https://promptrun.oss-cn-shanghai.aliyuncs.com/prompt_img/18-2024-04-16_011859.png', 0);
INSERT INTO `prompt_img` VALUES (21, 18, 'https://promptrun.oss-cn-shanghai.aliyuncs.com/prompt_img/18-2024-04-16_011900.png', 0);
INSERT INTO `prompt_img` VALUES (22, 19, 'https://promptrun.oss-cn-shanghai.aliyuncs.com/prompt_img/19-2024-04-16_013601.png', 1);
INSERT INTO `prompt_img` VALUES (23, 20, 'https://promptrun.oss-cn-shanghai.aliyuncs.com/prompt_img/20-2024-04-18_000446.png', 1);
INSERT INTO `prompt_img` VALUES (24, 21, 'https://promptrun.oss-cn-shanghai.aliyuncs.com/prompt_img/21-2024-04-18_002843.png', 1);
INSERT INTO `prompt_img` VALUES (25, 22, 'https://promptrun.oss-cn-shanghai.aliyuncs.com/prompt_img/22-2024-04-21_010501.png', 1);

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
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `user_id`(`user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 14 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of seller
-- ----------------------------
INSERT INTO `seller` VALUES (1, 4, 4.8, 1, '一位拥有 2 年大模型调试经验的提示词工程师，擅长 GPT、DALL·E 等模型，欢迎购买我的 Prompt。', 0, 3, '2024-02-14 01:09:44');
INSERT INTO `seller` VALUES (2, 3, 4.4, 1, '小陆 Prompts', 0, 1, '2024-03-24 19:19:59');

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `header_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `type` int NULL DEFAULT NULL COMMENT '用户类型，0: 买家，1: 卖家。做个字段冗余，方便查询',
  `create_time` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, '1298911600@qq.com', '$2a$10$J9e59kpzy65S10SwBYwFVeu5UP4ZS4uf9.RI2BuyYRoZ73tYREGs2', '1298911600@qq.com', '', 0, '2024-03-13 23:38:37');
INSERT INTO `user` VALUES (2, '1298900@qq.com', '$2a$10$J9e59kpzy65S10SwBYwFVeu5UP4ZS4uf9.RI2BuyYRoZ73tYREGs2', '1298900@qq.com', '', 0, '2024-03-14 13:03:05');
INSERT INTO `user` VALUES (3, '123@qq.com', '$2a$10$J9e59kpzy65S10SwBYwFVeu5UP4ZS4uf9.RI2BuyYRoZ73tYREGs2', '123@qq.com', '', 1, '2024-03-15 22:46:13');
INSERT INTO `user` VALUES (4, 'aarynlu@qq.com', '$2a$10$2bzwRNmYxugHleznHNruBeaQvwyiD.sk4n6KL29cPkfm7z3O3ZCEW', 'AarynLu', 'https://promptrun.oss-cn-shanghai.aliyuncs.com/header_img/4-2024-04-06_005519.jpeg', 1, '2024-03-23 15:15:39');

-- ----------------------------
-- Table structure for wallet
-- ----------------------------
DROP TABLE IF EXISTS `wallet`;
CREATE TABLE `wallet`  (
  `user_id` int NOT NULL,
  `wallet_income` decimal(10, 2) NULL DEFAULT 0.00 COMMENT '钱包总收入额',
  `wallet_outcome` decimal(10, 2) NULL DEFAULT 0.00 COMMENT '钱包总支出额',
  `balance` decimal(10, 2) NULL DEFAULT 0.00 COMMENT '钱包余额',
  `create_time` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of wallet
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;
