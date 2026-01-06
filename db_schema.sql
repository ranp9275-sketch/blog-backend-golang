CREATE DATABASE IF NOT EXISTS newblog DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE newblog;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` varchar(191) NOT NULL,
  `email` longtext,
  `name` longtext,
  `password` longtext,
  `avatar` longtext,
  `bio` longtext,
  `role` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for categories
-- ----------------------------
DROP TABLE IF EXISTS `categories`;
CREATE TABLE `categories` (
  `id` varchar(191) NOT NULL,
  `name` longtext,
  `slug` longtext,
  `desc` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for tags
-- ----------------------------
DROP TABLE IF EXISTS `tags`;
CREATE TABLE `tags` (
  `id` varchar(191) NOT NULL,
  `name` longtext,
  `slug` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for articles
-- ----------------------------
DROP TABLE IF EXISTS `articles`;
CREATE TABLE `articles` (
  `id` varchar(191) NOT NULL,
  `title` longtext,
  `slug` longtext,
  `content` longtext,
  `excerpt` longtext,
  `cover_image` longtext,
  `category_id` varchar(191) DEFAULT NULL,
  `status` longtext,
  `views` bigint(20) DEFAULT NULL,
  `author_id` varchar(191) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `published_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_articles_category` (`category_id`),
  KEY `fk_articles_author` (`author_id`),
  CONSTRAINT `fk_articles_author` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_articles_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for comments
-- ----------------------------
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments` (
  `id` varchar(191) NOT NULL,
  `article_id` varchar(191) DEFAULT NULL,
  `author` longtext,
  `email` longtext,
  `content` longtext,
  `status` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_articles_comments` (`article_id`),
  CONSTRAINT `fk_articles_comments` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for article_views
-- ----------------------------
DROP TABLE IF EXISTS `article_views`;
CREATE TABLE `article_views` (
  `id` varchar(191) NOT NULL,
  `article_id` varchar(191) DEFAULT NULL,
  `ip` longtext,
  `user_agent` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_articles_article_views` (`article_id`),
  CONSTRAINT `fk_articles_article_views` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for article_tags
-- ----------------------------
DROP TABLE IF EXISTS `article_tags`;
CREATE TABLE `article_tags` (
  `article_id` varchar(191) NOT NULL,
  `tag_id` varchar(191) NOT NULL,
  PRIMARY KEY (`article_id`,`tag_id`),
  KEY `fk_article_tags_tag` (`tag_id`),
  CONSTRAINT `fk_article_tags_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`),
  CONSTRAINT `fk_article_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for favorites
-- ----------------------------
DROP TABLE IF EXISTS `favorites`;
CREATE TABLE `favorites` (
  `id` varchar(191) NOT NULL,
  `user_id` varchar(191) DEFAULT NULL,
  `article_id` varchar(191) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_favorites_user_id` (`user_id`),
  KEY `idx_favorites_article_id` (`article_id`),
  CONSTRAINT `fk_favorites_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_favorites_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Insert Default Data
-- ----------------------------
-- Insert Default Data
-- ----------------------------
INSERT INTO categories (id, name, slug, `desc`, created_at, updated_at) 
VALUES ('cat-1', 'Category 1', 'category-1', 'Description for Category 1', NOW(), NOW());

INSERT INTO tags (id, name, slug, created_at, updated_at) 
VALUES ('tag-1', 'Tag 1', 'tag-1', NOW(), NOW());

INSERT INTO users (id, email, name, password, role, created_at, updated_at) 
VALUES ('admin-uuid', '4553664@qq.com', 'Administrator', '$2a$10$byOAW35giCZ8.Ytg/WWAueymtqQnZX/Ow3xCsG06HdrgCZG5psu5K', 'admin', NOW(), NOW());

SET FOREIGN_KEY_CHECKS = 1;
