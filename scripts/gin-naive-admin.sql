-- MySQL dump 10.13  Distrib 8.0.45, for Linux (aarch64)
--
-- Host: localhost    Database: gin-naive-admin
-- ------------------------------------------------------
-- Server version	8.0.45

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `sys_config`
--

DROP TABLE IF EXISTS `sys_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_config` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  `config_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `config_value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `create_by` bigint unsigned DEFAULT '0' COMMENT '创建人用户ID',
  `update_by` bigint unsigned DEFAULT '0' COMMENT '更新人用户ID',
  `delete_by` bigint unsigned DEFAULT '0' COMMENT '删除人用户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_sys_config_config_key` (`config_key`),
  KEY `idx_sys_config_delete_time` (`delete_time`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_config`
--

LOCK TABLES `sys_config` WRITE;
/*!40000 ALTER TABLE `sys_config` DISABLE KEYS */;
INSERT INTO `sys_config` VALUES (1,'2026-03-24 21:22:28.466','2026-03-26 20:43:19.319',NULL,'title','Gin Naive Admin','站点标题',0,0,0),(2,'2026-03-24 21:22:28.466','2026-03-28 21:41:36.309',NULL,'copyright','Copyright © 2026 Gin Naive Admin','页脚版权',0,0,0),(3,'2026-03-26 20:58:01.559','2026-03-26 20:58:54.878',NULL,'phone','027 88888888','联系电话',0,0,0);
/*!40000 ALTER TABLE `sys_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_department`
--

DROP TABLE IF EXISTS `sys_department`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_department` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '父部门ID',
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '部门名称',
  `code` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '部门编码',
  `sort` bigint unsigned DEFAULT '0' COMMENT '排序',
  `status` bigint unsigned DEFAULT '1' COMMENT '状态（1:启用 0:禁用）',
  `remark` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  `create_by` bigint unsigned DEFAULT '0' COMMENT '创建人用户ID',
  `update_by` bigint unsigned DEFAULT '0' COMMENT '更新人用户ID',
  `delete_by` bigint unsigned DEFAULT '0' COMMENT '删除人用户ID',
  PRIMARY KEY (`id`),
  KEY `idx_sys_department_delete_time` (`delete_time`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_department`
--

LOCK TABLES `sys_department` WRITE;
/*!40000 ALTER TABLE `sys_department` DISABLE KEYS */;
INSERT INTO `sys_department` VALUES (1,'2026-03-21 20:37:46.000','2026-03-21 12:39:09.323',NULL,0,'天成科技集团','A01',1,1,'',0,0,0),(2,'2026-03-23 20:58:18.750','2026-03-23 20:58:18.750',NULL,1,'天成投资集团','B01',1,1,'',0,0,0),(3,'2026-03-23 20:58:56.366','2026-03-23 20:58:56.366',NULL,1,'天成软件科技','B02',2,1,'',0,0,0);
/*!40000 ALTER TABLE `sys_department` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_dict_data`
--

DROP TABLE IF EXISTS `sys_dict_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dict_data` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  `type_code` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典类型编码',
  `label` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典标签',
  `value` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典值',
  `status` bigint unsigned DEFAULT '1' COMMENT '状态（1:启用 0:禁用）',
  `remark` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  `sort` bigint unsigned DEFAULT '0' COMMENT '排序',
  `create_by` bigint unsigned DEFAULT '0' COMMENT '创建人用户ID',
  `update_by` bigint unsigned DEFAULT '0' COMMENT '更新人用户ID',
  `delete_by` bigint unsigned DEFAULT '0' COMMENT '删除人用户ID',
  PRIMARY KEY (`id`),
  KEY `idx_sys_dict_data_delete_time` (`delete_time`),
  KEY `idx_sys_dict_data_type_code` (`type_code`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dict_data`
--

LOCK TABLES `sys_dict_data` WRITE;
/*!40000 ALTER TABLE `sys_dict_data` DISABLE KEYS */;
INSERT INTO `sys_dict_data` VALUES (3,'2026-03-21 01:07:12.545','2026-03-21 01:07:20.979',NULL,'sex','男','1',1,'',1,0,0,0),(4,'2026-03-21 01:07:28.278','2026-03-21 01:07:28.278',NULL,'sex','女','2',1,'',2,0,0,0),(5,'2026-03-29 16:48:54.181','2026-03-29 16:48:54.181',NULL,'status','启用','1',1,'',1,0,0,0),(6,'2026-03-29 16:49:07.786','2026-03-29 16:49:07.786',NULL,'status','禁用','2',1,'',2,0,0,0);
/*!40000 ALTER TABLE `sys_dict_data` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_dict_type`
--

DROP TABLE IF EXISTS `sys_dict_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dict_type` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  `type_code` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典类型编码',
  `type_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典类型名称',
  `status` bigint unsigned DEFAULT '1' COMMENT '状态（1:启用 0:禁用）',
  `remark` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  `sort` bigint unsigned DEFAULT '0' COMMENT '排序',
  `create_by` bigint unsigned DEFAULT '0' COMMENT '创建人用户ID',
  `update_by` bigint unsigned DEFAULT '0' COMMENT '更新人用户ID',
  `delete_by` bigint unsigned DEFAULT '0' COMMENT '删除人用户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_sys_dict_type_type_code` (`type_code`),
  KEY `idx_sys_dict_type_delete_time` (`delete_time`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dict_type`
--

LOCK TABLES `sys_dict_type` WRITE;
/*!40000 ALTER TABLE `sys_dict_type` DISABLE KEYS */;
INSERT INTO `sys_dict_type` VALUES (20,'2026-03-21 01:06:03.565','2026-03-21 01:15:21.717',NULL,'sex','用户性别',1,'用户性别',1,0,0,0),(21,'2026-03-29 16:48:43.789','2026-03-29 16:48:43.789',NULL,'status','状态',1,'状态',2,0,0,0);
/*!40000 ALTER TABLE `sys_dict_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_job_level`
--

DROP TABLE IF EXISTS `sys_job_level`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_job_level` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  `level_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '职务级别名称',
  `level` bigint unsigned NOT NULL COMMENT '职务级别数值（越小级别越高）',
  `create_by` bigint unsigned DEFAULT '0' COMMENT '创建人用户ID',
  `update_by` bigint unsigned DEFAULT '0' COMMENT '更新人用户ID',
  `delete_by` bigint unsigned DEFAULT '0' COMMENT '删除人用户ID',
  PRIMARY KEY (`id`),
  KEY `idx_sys_job_level_delete_time` (`delete_time`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_job_level`
--

LOCK TABLES `sys_job_level` WRITE;
/*!40000 ALTER TABLE `sys_job_level` DISABLE KEYS */;
INSERT INTO `sys_job_level` VALUES (1,'2026-03-24 20:16:06.671','2026-03-24 20:20:33.503',NULL,'董事长',1,0,0,0),(2,'2026-03-24 20:16:23.790','2026-03-24 20:16:23.790','2026-03-24 20:21:03.996','财务总监',2,0,0,0),(3,'2026-03-24 20:21:11.347','2026-03-24 20:21:11.347',NULL,'财务',2,0,0,0);
/*!40000 ALTER TABLE `sys_job_level` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_login_log`
--

DROP TABLE IF EXISTS `sys_login_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_login_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '用户ID',
  `account` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '账号',
  `ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'IP',
  `status` bigint DEFAULT NULL COMMENT '状态 1成功 2失败',
  `msg` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '说明',
  `create_by` bigint unsigned DEFAULT '0' COMMENT '创建人用户ID',
  `update_by` bigint unsigned DEFAULT '0' COMMENT '更新人用户ID',
  `delete_by` bigint unsigned DEFAULT '0' COMMENT '删除人用户ID',
  PRIMARY KEY (`id`),
  KEY `idx_sys_login_log_delete_time` (`delete_time`),
  KEY `idx_sys_login_log_user_id` (`user_id`),
  KEY `idx_sys_login_log_account` (`account`),
  KEY `idx_sys_login_log_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_login_log`
--

LOCK TABLES `sys_login_log` WRITE;
/*!40000 ALTER TABLE `sys_login_log` DISABLE KEYS */;
INSERT INTO `sys_login_log` VALUES (1,'2026-03-24 21:22:46.900','2026-03-24 21:22:46.900',NULL,1,'admin','127.0.0.1',1,'登录成功',0,0,0),(2,'2026-03-24 22:20:24.699','2026-03-24 22:20:24.699',NULL,0,'admin','127.0.0.1',2,'验证码错误',0,0,0),(3,'2026-03-24 22:22:04.555','2026-03-24 22:22:04.555',NULL,1,'admin','127.0.0.1',1,'登录成功',0,0,0),(4,'2026-03-25 00:10:54.992','2026-03-25 00:10:54.992',NULL,2,'admin_test','127.0.0.1',2,'密码错误',0,0,0),(5,'2026-03-25 00:11:09.285','2026-03-25 00:11:09.285',NULL,1,'admin','127.0.0.1',1,'登录成功',0,0,0),(6,'2026-03-25 00:11:37.241','2026-03-25 00:11:37.241',NULL,2,'admin_test','127.0.0.1',1,'登录成功',0,0,0),(7,'2026-03-25 00:12:28.019','2026-03-25 00:12:28.019',NULL,1,'admin','127.0.0.1',1,'登录成功',0,0,0),(8,'2026-03-25 00:13:42.286','2026-03-25 00:13:42.286',NULL,2,'admin_test','127.0.0.1',1,'登录成功',0,0,0),(9,'2026-03-25 23:22:35.447','2026-03-25 23:22:35.447',NULL,0,'admin','127.0.0.1',2,'验证码错误',0,0,0),(10,'2026-03-25 23:22:41.012','2026-03-25 23:22:41.012',NULL,1,'admin','127.0.0.1',1,'登录成功',0,0,0),(11,'2026-03-26 00:21:05.556','2026-03-26 00:21:05.556',NULL,1,'admin','127.0.0.1',1,'登录成功',0,0,0),(12,'2026-03-26 20:58:31.609','2026-03-26 20:58:31.609',NULL,2,'admin_test','127.0.0.1',1,'登录成功',0,0,0),(13,'2026-03-28 15:57:31.944','2026-03-28 15:57:31.944',NULL,0,'admin','127.0.0.1',2,'验证码错误',0,0,0),(14,'2026-03-28 15:57:45.097','2026-03-28 15:57:45.097',NULL,1,'admin','127.0.0.1',1,'登录成功',0,0,0),(15,'2026-03-28 21:35:52.084','2026-03-28 21:35:52.084',NULL,0,'admin','127.0.0.1',2,'验证码错误',0,0,0),(16,'2026-03-28 21:36:03.449','2026-03-28 21:36:03.449',NULL,0,'admin','127.0.0.1',2,'验证码错误',0,0,0),(17,'2026-03-28 21:36:18.177','2026-03-28 21:36:18.177',NULL,1,'admin','127.0.0.1',2,'密码错误',0,0,0),(18,'2026-03-28 21:36:29.172','2026-03-28 21:36:29.172',NULL,0,'admin','127.0.0.1',2,'验证码错误',0,0,0),(19,'2026-03-28 21:37:00.334','2026-03-28 21:37:00.334',NULL,1,'admin','127.0.0.1',1,'登录成功',0,0,0),(20,'2026-03-28 21:41:31.385','2026-03-28 21:41:31.385',NULL,1,'admin','127.0.0.1',1,'登录成功',0,0,0),(21,'2026-03-29 21:25:45.764','2026-03-29 21:25:45.764',NULL,2,'admin_test','127.0.0.1',1,'登录成功',0,0,0),(22,'2026-03-29 21:26:30.587','2026-03-29 21:26:30.587',NULL,1,'admin','127.0.0.1',1,'登录成功',0,0,0),(23,'2026-03-29 21:26:56.891','2026-03-29 21:26:56.891',NULL,2,'admin_test','127.0.0.1',1,'登录成功',0,0,0),(24,'2026-04-04 15:15:44.949','2026-04-04 15:15:44.949',NULL,0,'admin','127.0.0.1',2,'验证码错误',0,0,0),(25,'2026-04-04 15:15:55.430','2026-04-04 15:15:55.430',NULL,1,'admin','127.0.0.1',1,'登录成功',0,0,0),(26,'2026-04-05 15:24:19.787','2026-04-05 15:24:19.787',NULL,0,'admin_test','127.0.0.1',2,'验证码错误',0,0,0),(27,'2026-04-05 15:24:28.183','2026-04-05 15:24:28.183',NULL,2,'admin_test','127.0.0.1',1,'登录成功',0,0,0),(28,'2026-04-05 18:03:44.139','2026-04-05 18:03:44.139',NULL,1,'admin','127.0.0.1',1,'登录成功',0,0,0);
/*!40000 ALTER TABLE `sys_login_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_menu`
--

DROP TABLE IF EXISTS `sys_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_menu` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '父菜单ID',
  `type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '1' COMMENT '菜单类型（0:目录  1:菜单  2:按钮）',
  `path` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '路由路径',
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `component` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '路由组件路径',
  `redirect` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '重定向',
  `perms` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '授权标识',
  `title` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单标题',
  `title_i18n_key` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '国际化key',
  `icon` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '菜单标题',
  `fixed_in_tabs` tinyint(1) DEFAULT NULL COMMENT '固定在标签页',
  `hide_in_menu` tinyint(1) DEFAULT NULL COMMENT '不在菜单中显示',
  `keep_alive` tinyint(1) DEFAULT NULL COMMENT '缓存配置',
  `link` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '外链',
  `link_mode` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '链接模式',
  `nested_route_render_end` tinyint(1) DEFAULT NULL COMMENT '是否在当前路由层级结束渲染',
  `sort` bigint unsigned DEFAULT '0' COMMENT '排序',
  `status` bigint unsigned DEFAULT '1' COMMENT '状态（1:启用  2:禁用）',
  `create_by` bigint unsigned DEFAULT '0' COMMENT '创建人用户ID',
  `update_by` bigint unsigned DEFAULT '0' COMMENT '更新人用户ID',
  `delete_by` bigint unsigned DEFAULT '0' COMMENT '删除人用户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_sys_menu_name` (`name`),
  KEY `idx_sys_menu_delete_time` (`delete_time`)
) ENGINE=InnoDB AUTO_INCREMENT=110 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_menu`
--

LOCK TABLES `sys_menu` WRITE;
/*!40000 ALTER TABLE `sys_menu` DISABLE KEYS */;
INSERT INTO `sys_menu` VALUES (1,'2026-02-10 13:03:52.000','2026-04-05 19:36:12.851',NULL,0,'1','home','Home','/home/index.vue','',NULL,'首页','routes.home','material-symbols:dashboard-outline-rounded',1,0,0,'','',0,1,1,0,1,0),(2,'2026-02-02 13:06:25.000',NULL,NULL,0,'0','demos','Demos',NULL,NULL,NULL,'演示','routes.demo','hugeicons:codesandbox',NULL,NULL,NULL,NULL,NULL,NULL,2,1,0,0,0),(5,'2026-02-02 17:53:03.000',NULL,NULL,2,'0','fallback','Fallback',NULL,NULL,NULL,'异常页','routes.exception','solar:shield-warning-broken',NULL,NULL,NULL,NULL,NULL,NULL,2,1,0,0,0),(6,'2026-02-02 17:56:16.000',NULL,NULL,5,'1','403','403','/demos/fallback/403.vue',NULL,NULL,'403','routes.403','mdi:forbid',NULL,NULL,NULL,NULL,NULL,NULL,1,1,0,0,0),(7,'2026-02-02 18:02:57.000',NULL,NULL,5,'1','404','404','/demos/fallback/404.vue',NULL,NULL,'404','routes.404','ic:baseline-browser-not-supported',NULL,NULL,NULL,NULL,NULL,NULL,2,1,0,0,0),(8,'2026-02-02 18:04:15.000',NULL,NULL,5,'1','500','500','/demos/fallback/500.vue',NULL,NULL,'500','routes.500','streamline-flex:monitor-error',NULL,NULL,NULL,NULL,NULL,NULL,3,1,0,0,0),(9,'2026-02-02 18:06:13.000',NULL,NULL,2,'0','external','External',NULL,NULL,NULL,'外部页面','routes.externalPage','ant-design:link-outlined',NULL,NULL,NULL,NULL,NULL,NULL,3,1,0,0,0),(10,'2026-02-02 18:07:57.000',NULL,NULL,9,'0','iframe','Iframe',NULL,NULL,NULL,'内嵌','routes.embedded','material-symbols:iframe',NULL,NULL,NULL,'','',NULL,1,1,0,0,0),(11,'2026-02-02 18:11:56.000',NULL,NULL,10,'1','form','Form','{}',NULL,NULL,'复杂表单','routes.complexForm','lets-icons:form',NULL,NULL,NULL,'https://naive-ui.pro-components.cn/zh-CN/os-theme/components/form-list#list-nest.vue','iframe',NULL,1,1,0,0,0),(12,'2026-02-02 18:16:28.000',NULL,NULL,10,'1','edit-data-table','EditDataTable','{}',NULL,NULL,'编辑表格','routes.editTable','material-symbols:table-outline',NULL,NULL,NULL,'https://naive-ui.pro-components.cn/zh-CN/os-theme/components/edit-data-table#async.vue','iframe',NULL,2,1,0,0,0),(13,'2026-02-02 18:16:56.000',NULL,NULL,10,'1','baidu-iframe','BaiduIframe','{}',NULL,NULL,'百度','routes.baiduIframe','ri:baidu-fill',NULL,NULL,NULL,'https://www.baidu.com','iframe',NULL,3,1,0,0,0),(15,'2026-02-02 18:26:31.000',NULL,NULL,9,'0','link','Link',NULL,NULL,NULL,'外链','routes.externalLink','akar-icons:link-out',NULL,NULL,NULL,NULL,NULL,NULL,2,1,0,0,0),(16,'2026-02-02 18:27:44.000',NULL,NULL,15,'1','vite','Vite','',NULL,NULL,'Vite',NULL,'logos:vitejs',NULL,NULL,NULL,'https://vite.dev','newWindow',NULL,1,1,0,0,0),(17,'2026-02-02 18:29:27.000',NULL,NULL,15,'1','vue','Vue','',NULL,NULL,'Vue',NULL,'logos:vue',NULL,NULL,NULL,'https://vuejs.org/','newWindow',NULL,2,1,0,0,0),(20,'2026-02-02 19:08:32.000',NULL,NULL,2,'1','download','Download','/demos/download/index.vue',NULL,NULL,'文件下载','routes.fileDownload','material-symbols:download',NULL,NULL,NULL,NULL,NULL,NULL,4,1,0,0,0),(21,'2026-02-02 19:42:35.000',NULL,NULL,2,'0','nested-detail','NestedDetail','/demos/nested/index.vue',NULL,NULL,'嵌套详情页','routes.nestedDetail','bx:detail',NULL,NULL,NULL,NULL,NULL,NULL,5,1,0,0,0),(22,'2026-02-02 19:45:09.000',NULL,NULL,21,'1','detail','Detail','/demos/nested/detail.vue',NULL,NULL,'详情页','routes.detail',NULL,NULL,1,NULL,NULL,NULL,NULL,1,1,0,0,0),(23,'2026-02-02 19:48:24.000',NULL,NULL,2,'0','nested-detail2','NestedDetail2','/demos/nested/demo2/index.vue','nestedDetail2:detail1',NULL,'嵌套详情页（2）','routes.nestedDetail2','bx:detail',NULL,NULL,NULL,NULL,NULL,1,6,1,0,0,0),(24,'2026-02-02 19:51:07.000',NULL,NULL,23,'1','detail1','nestedDetail2:detail1','/demos/nested/demo2/detail1.vue',NULL,NULL,'详情页（1）','routes.detail1',NULL,NULL,1,NULL,NULL,NULL,NULL,1,1,0,0,0),(25,'2026-02-02 19:52:07.000',NULL,NULL,23,'1','detail2','nestedDetail2:detail2','/demos/nested/demo2/detail2.vue',NULL,NULL,'详情页（2）','routes.detail2',NULL,NULL,1,NULL,NULL,NULL,NULL,1,1,0,0,0),(26,'2026-02-02 23:12:53.000',NULL,NULL,2,'0','keep-alive','keep-alive',NULL,NULL,NULL,'缓存路由','routes.keepAlive','octicon:cache-16',NULL,NULL,NULL,NULL,NULL,NULL,1,1,0,0,0),(27,'2026-02-02 23:15:45.000',NULL,NULL,26,'1','demo1','demo1','/demos/keep-alive/demo1.vue',NULL,NULL,'基础缓存','routes.keepAliveDemo1',NULL,NULL,NULL,1,NULL,NULL,NULL,1,1,0,0,0),(28,'2026-02-02 23:17:23.000',NULL,NULL,26,'1','demo2','demo2','/demos/keep-alive/demo2.vue',NULL,NULL,'条件缓存','routes.keepAliveDemo2',NULL,NULL,NULL,1,NULL,NULL,NULL,2,1,0,0,0),(29,'2026-02-02 23:19:32.000',NULL,NULL,2,'1','tabs','Tabs','/demos/tabs/index.vue',NULL,NULL,'多标签','routes.tabs','mdi:tab',NULL,NULL,NULL,NULL,NULL,NULL,8,1,0,0,0),(30,'2026-02-02 23:20:59.000',NULL,NULL,2,'1','page-component','PageComponent','/demos/page-component/index.vue',NULL,NULL,'页面级组件','routes.pageComponent','material-symbols:pageview-outline',NULL,NULL,NULL,NULL,NULL,NULL,9,1,0,0,0),(31,'2026-02-02 23:22:31.000',NULL,NULL,2,'1','editor','Editor','/demos/wang-editor/index.vue',NULL,NULL,'富文本','routes.richText','material-symbols:edit-document-outline',NULL,NULL,NULL,NULL,NULL,NULL,10,1,0,0,0),(32,'2026-02-02 23:23:48.000',NULL,NULL,2,'1','complex-form','ComplexForm','/demos/complex-form/index.vue',NULL,NULL,'复杂表单','routes.complexForm','material-symbols:dynamic-form-outline',NULL,NULL,NULL,NULL,NULL,NULL,11,1,0,0,0),(33,'2026-02-02 23:24:45.000',NULL,NULL,2,'1','icon','Icon','/demos/icon/index.vue',NULL,NULL,'图标选择器','routes.iconSelector','mdi:image-outline',NULL,NULL,NULL,NULL,NULL,NULL,12,1,0,0,0),(34,'2026-02-02 23:25:53.000',NULL,NULL,2,'1','loading','Loading','/demos/loading/index.vue',NULL,NULL,'Loading 指令','routes.loading','line-md:loading-twotone-loop',NULL,NULL,NULL,NULL,NULL,NULL,13,1,0,0,0),(35,'2026-02-03 07:26:50.000','2026-02-05 16:48:16.161',NULL,0,'0','system','System','','',NULL,'系统管理','routes.system','ant-design:setting-outlined',0,0,0,'','',0,3,1,0,0,0),(36,'2026-02-02 23:28:23.000',NULL,NULL,35,'1','user','User','/system/user/index.vue',NULL,NULL,'用户管理','routes.userManagement','ant-design:user-outlined',NULL,NULL,NULL,NULL,NULL,NULL,1,1,0,0,0),(37,'2026-02-02 23:29:21.000',NULL,NULL,35,'1','role','Role','/system/role/index.vue',NULL,NULL,'角色管理','routes.roleManagement','carbon:user-role',NULL,NULL,NULL,NULL,NULL,NULL,2,1,0,0,0),(38,'2026-02-06 07:57:44.000','2026-03-20 22:56:21.277',NULL,35,'1','menu','Menu','/system/menu/index.vue','',NULL,'菜单管理','routes.menuManagement','ant-design:menu-outlined',0,0,0,'','',0,3,1,0,0,0),(39,'2026-03-21 15:14:50.000','2026-03-20 23:15:39.858',NULL,35,'1','/system/dict','dict','/system/dict/index.vue','',NULL,'字典管理','','material-symbols:batch-prediction',0,0,0,'','',0,4,1,0,0,0),(40,'2026-03-21 12:16:48.759','2026-03-21 12:16:48.759',NULL,35,'1','/system/department','department','/system/department/index.vue','',NULL,'部门管理','','mingcute:department-fill',0,0,0,'','',0,5,1,0,0,0),(41,'2026-03-24 19:49:04.708','2026-03-24 19:49:04.708',NULL,35,'1','/position','position','/system/position/index.vue','',NULL,'职务管理','','gis:position-man',0,0,0,'','',0,6,1,0,0,0),(44,'2026-03-25 13:22:28.000','2026-03-29 19:01:22.476',NULL,36,'2','','apiPerm_system_user_delete','','','system:user:delete','用户管理-删除用户','','',0,1,0,'','',0,2,1,0,0,0),(46,'2026-03-25 13:22:28.000','2026-03-29 19:01:22.446',NULL,37,'2','','apiPerm_system_role_power','','','system:role:power','角色管理-角色授权','','',0,1,0,'','',0,4,1,0,0,0),(48,'2026-03-25 13:22:28.000','2026-03-29 19:01:22.493',NULL,36,'2','','apiPerm_system_user_add','','','system:user:add','用户管理-添加用户','','',0,1,0,'','',0,1,1,0,0,0),(49,'2026-03-25 05:22:28.000','2026-03-29 19:01:22.503',NULL,36,'2','','apiPerm_system_user_edit','','','system:user:edit','用户管理-编辑用户','','',0,1,0,'','',0,3,1,0,0,0),(50,'2026-03-25 13:22:28.000','2026-03-29 19:01:22.392',NULL,37,'2','','apiPerm_system_role_add','','','system:role:add','角色管理-新增角色','','',0,1,0,'','',0,1,1,0,0,0),(51,'2026-03-25 05:22:28.000','2026-03-29 19:01:22.421',NULL,38,'2','','apiPerm_system_menu_edit','','','system:menu:edit','菜单管理-编辑菜单','','',0,1,0,'','',0,2,1,0,0,0),(64,'2026-03-25 05:22:28.000','2026-03-29 19:01:22.497',NULL,38,'2','','apiPerm_system_menu_delete','','','system:menu:delete','菜单管理-删除菜单','','',0,1,0,'','',0,3,1,0,0,0),(67,'2026-03-25 05:22:28.000','2026-03-29 19:01:22.431',NULL,37,'2','','apiPerm_system_role_edit','','','system:role:edit','角色管理-编辑角色','','',0,1,0,'','',0,2,1,0,0,0),(76,'2026-03-25 05:22:28.000','2026-03-29 19:01:22.438',NULL,81,'2','','apiPerm_system_config_edit','','','system:config:edit','参数配置-编辑参数','','',0,1,0,'','',0,2,1,0,0,0),(78,'2026-03-25 05:22:28.000','2026-03-29 19:01:22.415',NULL,37,'2','','apiPerm_system_role_delete','','','system:role:delete','角色管理-删除角色','','',0,1,0,'','',0,3,1,0,0,0),(79,NULL,'2026-03-26 20:15:44.643',NULL,82,'1','log-login','LogLogin','monitor/log-login/index.vue','',NULL,'登录日志','routes.loginLog','mdi:text-box-search-outline',0,0,1,'','',0,1,1,0,0,0),(80,NULL,'2026-03-26 20:15:47.942',NULL,82,'1','log-oper','LogOper','monitor/log-oper/index.vue','',NULL,'操作日志','routes.operLog','mdi:clipboard-text-outline',0,0,1,'','',0,2,1,0,0,0),(81,NULL,'2026-03-26 20:15:38.781',NULL,35,'1','sys-config','SysConfig','system/sys-config/index.vue','',NULL,'参数配置','routes.sysConfig','mdi:cog-outline',0,0,1,'','',0,7,1,0,0,0),(82,'2026-03-24 22:02:42.233','2026-03-24 22:02:42.233',NULL,0,'0','monitor','systemMonitor','','',NULL,'系统监控','','ph:monitor-play-fill',0,0,0,'','',0,4,1,0,0,0),(83,'2026-03-26 07:08:02.000','2026-03-25 23:08:14.997',NULL,38,'2','','','','','system:menu:add','菜单管理-新增菜单','','',0,0,0,'','',0,1,1,0,0,0),(88,'2026-03-27 05:28:43.000','2026-03-29 19:53:23.368',NULL,81,'2','','system_config_edit','','','system:config:add','参数配置-新增参数','','',0,0,0,'','',0,1,1,0,0,0),(89,'2026-03-29 19:07:23.103','2026-03-29 19:07:23.103',NULL,36,'2','','system_user_export','','','system:user:export','用户管理-导出用户','','',0,0,0,'','',0,4,1,0,0,0),(90,'2026-03-29 19:07:52.601','2026-03-29 19:07:52.601',NULL,36,'2','','system_user_import','','','system:user:import','用户管理-导入用户','','',0,0,0,'','',0,5,1,0,0,0),(91,'2026-03-29 19:09:02.419','2026-03-29 19:09:02.419',NULL,37,'2','','system_role_export','','','system:role:export','角色管理-导出角色','','',0,0,0,'','',0,5,1,0,0,0),(92,'2026-03-29 19:09:24.690','2026-03-29 19:09:24.690',NULL,37,'2','','system_role_import','','','system:role:import','角色管理-导入角色','','',0,0,0,'','',0,6,1,0,0,0),(93,'2026-03-30 03:11:56.000','2026-03-29 19:12:42.577',NULL,39,'2','','system_dict_add','','','system:dict:add','字典管理-新增字典','','',0,0,0,'','',0,1,1,0,0,0),(94,'2026-03-29 19:12:35.998','2026-03-29 19:12:35.998',NULL,39,'2','','system_dict_edit','','','system:dict:edit','字典管理-编辑字典','','',0,0,0,'','',0,2,1,0,0,0),(95,'2026-03-29 19:13:03.258','2026-03-29 19:13:03.258',NULL,39,'2','','system_dict_delete','','','system:dict:delete','字典管理-删除字典','','',0,0,0,'','',0,3,1,0,0,0),(96,'2026-03-30 03:14:09.000','2026-03-29 19:14:24.835',NULL,39,'2','','system_dict_power','','','system:dict:power','字典管理-字典配置','','',0,0,0,'','',0,4,1,0,0,0),(97,'2026-03-29 19:41:08.473','2026-03-29 19:41:08.473',NULL,40,'2','','system_department_add','','','system:department:add','部门管理-新增部门','','',0,0,0,'','',0,1,1,0,0,0),(98,'2026-03-29 19:42:25.746','2026-03-29 19:42:25.746',NULL,40,'2','','system_department_add_lowerLevel','','','system:department:add:lowerLevel','部门管理-添加下级','','',0,0,0,'','',0,2,1,0,0,0),(99,'2026-03-29 19:43:32.764','2026-03-29 19:43:32.764',NULL,40,'2','','system_department_export','','','system:department:export','部门管理-导出部门','','',0,0,0,'','',0,3,1,0,0,0),(100,'2026-03-29 19:44:14.215','2026-03-29 19:44:14.215',NULL,41,'2','','system_position_add','','','system:position:add','职务管理-新增职务','','',0,0,0,'','',0,1,1,0,0,0),(101,'2026-03-29 19:44:37.554','2026-03-29 19:44:37.554',NULL,41,'2','','system_position_edit','','','system:position:edit','职务管理-编辑职务','','',0,0,0,'','',0,2,1,0,0,0),(102,'2026-03-29 19:45:01.969','2026-03-29 19:45:01.969',NULL,41,'2','','system_position_delete','','','system:position:delete','职务管理-删除职务','','',0,0,0,'','',0,3,1,0,0,0),(103,'2026-03-29 19:45:23.706','2026-03-29 19:45:23.706',NULL,41,'2','','system_position_export','','','system:position:export','职务管理-导出职务','','',0,0,0,'','',0,4,1,0,0,0),(104,'2026-03-29 19:45:51.251','2026-03-29 19:45:51.251',NULL,41,'2','','system_position_import','','','system:position:import','职务管理-导入职务','','',0,0,0,'','',0,5,1,0,0,0),(105,'2026-03-29 19:46:55.609','2026-03-29 19:46:55.609',NULL,81,'2','','system_config_export','','','system:config:export','参数配置-导出参数','','',0,0,0,'','',0,3,1,0,0,0),(106,'2026-03-29 19:47:12.889','2026-03-29 19:47:12.889',NULL,81,'2','','system_config_import','','','system:config:import','参数配置-导入参数','','',0,0,0,'','',0,4,1,0,0,0),(107,'2026-04-05 18:05:43.716','2026-04-05 18:05:43.716',NULL,0,'0','online','Online','','',NULL,'在线开发','','ic:outline-cloud',0,0,0,'','',0,5,1,0,0,0),(108,'2026-04-06 02:07:05.000','2026-04-05 18:08:03.255',NULL,107,'1','devform','Devform','/develop/devform/index.vue','',NULL,'表单开发','','icon-park-outline:form',0,0,0,'','',0,0,1,0,0,0);
/*!40000 ALTER TABLE `sys_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_online_form`
--

DROP TABLE IF EXISTS `sys_online_form`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_online_form` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  `table_name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '物理表名',
  `entity_name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Go 实体名 PascalCase',
  `route_group` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Gin 路由组段',
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '说明',
  `sync_status` bigint DEFAULT '0' COMMENT '0未同步1已同步',
  `create_by` bigint unsigned DEFAULT '0' COMMENT '创建人用户ID',
  `update_by` bigint unsigned DEFAULT '0' COMMENT '更新人用户ID',
  `delete_by` bigint unsigned DEFAULT '0' COMMENT '删除人用户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_sys_online_form_phys_table_name` (`table_name`),
  KEY `idx_sys_online_form_delete_time` (`delete_time`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_online_form`
--

LOCK TABLES `sys_online_form` WRITE;
/*!40000 ALTER TABLE `sys_online_form` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_online_form` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_online_form_field`
--

DROP TABLE IF EXISTS `sys_online_form_field`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_online_form_field` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  `online_form_id` bigint unsigned NOT NULL COMMENT '表单ID',
  `sort` bigint DEFAULT '0' COMMENT '排序',
  `column_name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '列名snake',
  `db_type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'mysql类型',
  `length` bigint DEFAULT '255' COMMENT 'varchar长度',
  `decimal_scale` bigint DEFAULT '2' COMMENT 'decimal小数位',
  `nullable` tinyint(1) DEFAULT '0',
  `comment` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '列注释',
  `list_show` tinyint(1) DEFAULT '1' COMMENT '列表显示',
  `form_show` tinyint(1) DEFAULT '1' COMMENT '表单显示',
  `is_query` tinyint(1) DEFAULT '0' COMMENT '可查询',
  `query_type` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'eq' COMMENT 'eq或like',
  `create_by` bigint unsigned DEFAULT '0' COMMENT '创建人用户ID',
  `update_by` bigint unsigned DEFAULT '0' COMMENT '更新人用户ID',
  `delete_by` bigint unsigned DEFAULT '0' COMMENT '删除人用户ID',
  PRIMARY KEY (`id`),
  KEY `idx_sys_online_form_field_delete_time` (`delete_time`),
  KEY `idx_sys_online_form_field_online_form_id` (`online_form_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_online_form_field`
--

LOCK TABLES `sys_online_form_field` WRITE;
/*!40000 ALTER TABLE `sys_online_form_field` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_online_form_field` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_oper_log`
--

DROP TABLE IF EXISTS `sys_oper_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_oper_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  `title` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '模块标题',
  `method` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'HTTP方法',
  `path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '路径',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '用户ID',
  `account` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '账号',
  `ip` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'IP',
  `latency_ms` bigint DEFAULT NULL COMMENT '耗时毫秒',
  `status` bigint DEFAULT NULL COMMENT 'HTTP状态码',
  `create_by` bigint unsigned DEFAULT '0' COMMENT '创建人用户ID',
  `update_by` bigint unsigned DEFAULT '0' COMMENT '更新人用户ID',
  `delete_by` bigint unsigned DEFAULT '0' COMMENT '删除人用户ID',
  PRIMARY KEY (`id`),
  KEY `idx_sys_oper_log_delete_time` (`delete_time`),
  KEY `idx_sys_oper_log_method` (`method`),
  KEY `idx_sys_oper_log_user_id` (`user_id`),
  KEY `idx_sys_oper_log_account` (`account`)
) ENGINE=InnoDB AUTO_INCREMENT=193 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_oper_log`
--

LOCK TABLES `sys_oper_log` WRITE;
/*!40000 ALTER TABLE `sys_oper_log` DISABLE KEYS */;
INSERT INTO `sys_oper_log` VALUES (1,'2026-03-24 21:25:33.153','2026-03-24 21:25:33.153',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',10,200,0,0,0),(2,'2026-03-24 21:44:52.518','2026-03-24 21:44:52.518',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',10,200,0,0,0),(3,'2026-03-24 21:45:18.019','2026-03-24 21:45:18.019',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',6,200,0,0,0),(4,'2026-03-24 21:45:26.577','2026-03-24 21:45:26.577',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',7,200,0,0,0),(5,'2026-03-24 21:45:32.129','2026-03-24 21:45:32.129',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',4,200,0,0,0),(6,'2026-03-24 21:49:10.587','2026-03-24 21:49:10.587',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',5,200,0,0,0),(7,'2026-03-24 21:49:21.313','2026-03-24 21:49:21.313',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',7,200,0,0,0),(8,'2026-03-24 21:49:29.834','2026-03-24 21:49:29.834',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',4,200,0,0,0),(9,'2026-03-24 21:49:42.865','2026-03-24 21:49:42.865',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',5,200,0,0,0),(10,'2026-03-24 21:49:50.053','2026-03-24 21:49:50.053',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',6,200,0,0,0),(11,'2026-03-24 21:50:01.221','2026-03-24 21:50:01.221',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',7,200,0,0,0),(12,'2026-03-24 21:50:08.079','2026-03-24 21:50:08.079',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',5,200,0,0,0),(13,'2026-03-24 21:50:14.312','2026-03-24 21:50:14.312',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',8,200,0,0,0),(14,'2026-03-24 21:50:20.334','2026-03-24 21:50:20.334',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',4,200,0,0,0),(15,'2026-03-24 21:50:27.310','2026-03-24 21:50:27.310',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',5,200,0,0,0),(16,'2026-03-24 21:50:32.409','2026-03-24 21:50:32.409',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',5,200,0,0,0),(17,'2026-03-24 22:02:42.261','2026-03-24 22:02:42.261',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',9,200,0,0,0),(18,'2026-03-24 22:03:01.643','2026-03-24 22:03:01.643',NULL,'角色管理','POST','/api/role/setRolePower',1,'admin','127.0.0.1',7,200,0,0,0),(19,'2026-03-24 22:06:08.309','2026-03-24 22:06:08.309',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',7,200,0,0,0),(20,'2026-03-24 22:06:14.431','2026-03-24 22:06:14.431',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',6,200,0,0,0),(21,'2026-03-24 22:06:20.629','2026-03-24 22:06:20.629',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',7,200,0,0,0),(22,'2026-03-24 22:17:11.014','2026-03-24 22:17:11.014',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',10,200,0,0,0),(23,'2026-03-24 22:33:31.529','2026-03-24 22:33:31.529',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',8,200,0,0,0),(24,'2026-03-24 22:33:37.512','2026-03-24 22:33:37.512',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',5,200,0,0,0),(25,'2026-03-24 22:34:05.838','2026-03-24 22:34:05.838',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',5,200,0,0,0),(26,'2026-03-24 22:34:17.192','2026-03-24 22:34:17.192',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',5,200,0,0,0),(27,'2026-03-24 22:34:24.645','2026-03-24 22:34:24.645',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',7,200,0,0,0),(28,'2026-03-24 22:34:35.135','2026-03-24 22:34:35.135',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',5,200,0,0,0),(29,'2026-03-24 22:34:41.160','2026-03-24 22:34:41.160',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',5,200,0,0,0),(30,'2026-03-24 22:34:47.341','2026-03-24 22:34:47.341',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',5,200,0,0,0),(31,'2026-03-24 22:34:56.018','2026-03-24 22:34:56.018',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',5,200,0,0,0),(32,'2026-03-24 22:35:02.034','2026-03-24 22:35:02.034',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',8,200,0,0,0),(33,'2026-03-24 22:35:11.830','2026-03-24 22:35:11.830',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',7,200,0,0,0),(34,'2026-03-24 22:35:17.567','2026-03-24 22:35:17.567',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',8,200,0,0,0),(35,'2026-03-24 22:35:22.970','2026-03-24 22:35:22.970',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',6,200,0,0,0),(36,'2026-03-24 22:35:31.472','2026-03-24 22:35:31.472',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',7,200,0,0,0),(37,'2026-03-24 22:35:37.123','2026-03-24 22:35:37.123',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',5,200,0,0,0),(38,'2026-03-24 22:35:43.819','2026-03-24 22:35:43.819',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',4,200,0,0,0),(39,'2026-03-24 22:35:49.447','2026-03-24 22:35:49.447',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',7,200,0,0,0),(40,'2026-03-24 22:35:59.908','2026-03-24 22:35:59.908',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',11,200,0,0,0),(41,'2026-03-24 23:41:32.619','2026-03-24 23:41:32.619',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',9,200,0,0,0),(42,'2026-03-24 23:41:56.625','2026-03-24 23:41:56.625',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',7,200,0,0,0),(43,'2026-03-24 23:41:57.887','2026-03-24 23:41:57.887',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',7,200,0,0,0),(44,'2026-03-24 23:41:59.017','2026-03-24 23:41:59.017',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',6,200,0,0,0),(45,'2026-03-24 23:42:00.879','2026-03-24 23:42:00.879',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',6,200,0,0,0),(46,'2026-03-24 23:42:02.309','2026-03-24 23:42:02.309',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',5,200,0,0,0),(47,'2026-03-24 23:42:03.692','2026-03-24 23:42:03.692',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',8,200,0,0,0),(48,'2026-03-24 23:42:05.549','2026-03-24 23:42:05.549',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',8,200,0,0,0),(49,'2026-03-24 23:42:06.879','2026-03-24 23:42:06.879',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',3,200,0,0,0),(50,'2026-03-24 23:42:08.563','2026-03-24 23:42:08.563',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',8,200,0,0,0),(51,'2026-03-24 23:42:09.731','2026-03-24 23:42:09.731',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',5,200,0,0,0),(52,'2026-03-24 23:42:11.151','2026-03-24 23:42:11.151',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',7,200,0,0,0),(53,'2026-03-24 23:42:13.239','2026-03-24 23:42:13.239',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',7,200,0,0,0),(54,'2026-03-24 23:42:14.984','2026-03-24 23:42:14.984',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',4,200,0,0,0),(55,'2026-03-24 23:42:16.510','2026-03-24 23:42:16.510',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',12,200,0,0,0),(56,'2026-03-24 23:42:18.386','2026-03-24 23:42:18.386',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',10,200,0,0,0),(57,'2026-03-24 23:42:21.083','2026-03-24 23:42:21.083',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',13,200,0,0,0),(58,'2026-03-24 23:49:55.665','2026-03-24 23:49:55.665',NULL,'个人中心','PUT','/api/profile/info',1,'admin','127.0.0.1',11,200,0,0,0),(59,'2026-03-24 23:58:57.989','2026-03-24 23:58:57.989',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',10,200,0,0,0),(60,'2026-03-25 00:04:50.416','2026-03-25 00:04:50.416',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',8,200,0,0,0),(61,'2026-03-25 00:04:55.639','2026-03-25 00:04:55.639',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',5,200,0,0,0),(62,'2026-03-25 00:05:12.876','2026-03-25 00:05:12.876',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',5,200,0,0,0),(63,'2026-03-25 00:05:36.399','2026-03-25 00:05:36.399',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',8,200,0,0,0),(64,'2026-03-25 00:05:38.603','2026-03-25 00:05:38.603',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',10,200,0,0,0),(65,'2026-03-25 00:05:40.543','2026-03-25 00:05:40.543',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',9,200,0,0,0),(66,'2026-03-25 00:07:32.120','2026-03-25 00:07:32.120',NULL,'角色管理','POST','/api/role/setRolePower',1,'admin','127.0.0.1',8,200,0,0,0),(67,'2026-03-25 00:10:24.379','2026-03-25 00:10:24.379',NULL,'角色管理','POST','/api/role/setRolePower',1,'admin','127.0.0.1',9,200,0,0,0),(68,'2026-03-25 00:11:20.552','2026-03-25 00:11:20.552',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',38,200,0,0,0),(69,'2026-03-25 00:12:44.284','2026-03-25 00:12:44.284',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',6,200,0,0,0),(70,'2026-03-25 00:13:19.995','2026-03-25 00:13:19.995',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',6,200,0,0,0),(71,'2026-03-25 23:03:19.035','2026-03-25 23:03:19.035',NULL,'角色管理','POST','/api/role/setRolePower',2,'admin_test','127.0.0.1',14,200,0,0,0),(72,'2026-03-25 23:04:07.079','2026-03-25 23:04:07.079',NULL,'角色管理','POST','/api/role/setRolePower',2,'admin_test','127.0.0.1',12,200,0,0,0),(73,'2026-03-25 23:04:27.527','2026-03-25 23:04:27.527',NULL,'菜单管理','DELETE','/api/menu/delete/:id',2,'admin_test','127.0.0.1',9,200,0,0,0),(74,'2026-03-25 23:04:31.024','2026-03-25 23:04:31.024',NULL,'菜单管理','DELETE','/api/menu/delete/:id',2,'admin_test','127.0.0.1',8,200,0,0,0),(75,'2026-03-25 23:04:36.757','2026-03-25 23:04:36.757',NULL,'菜单管理','DELETE','/api/menu/delete/:id',2,'admin_test','127.0.0.1',9,200,0,0,0),(76,'2026-03-25 23:05:27.815','2026-03-25 23:05:27.815',NULL,'菜单管理','PUT','/api/menu/edit',2,'admin_test','127.0.0.1',7,200,0,0,0),(77,'2026-03-25 23:06:07.952','2026-03-25 23:06:07.952',NULL,'菜单管理','PUT','/api/menu/edit',2,'admin_test','127.0.0.1',8,200,0,0,0),(78,'2026-03-25 23:06:25.823','2026-03-25 23:06:25.823',NULL,'菜单管理','PUT','/api/menu/edit',2,'admin_test','127.0.0.1',7,200,0,0,0),(79,'2026-03-25 23:06:38.625','2026-03-25 23:06:38.625',NULL,'菜单管理','PUT','/api/menu/edit',2,'admin_test','127.0.0.1',10,200,0,0,0),(80,'2026-03-25 23:06:41.541','2026-03-25 23:06:41.541',NULL,'菜单管理','DELETE','/api/menu/delete/:id',2,'admin_test','127.0.0.1',16,200,0,0,0),(81,'2026-03-25 23:06:56.875','2026-03-25 23:06:56.875',NULL,'菜单管理','PUT','/api/menu/edit',2,'admin_test','127.0.0.1',10,200,0,0,0),(82,'2026-03-25 23:07:19.096','2026-03-25 23:07:19.096',NULL,'菜单管理','PUT','/api/menu/edit',2,'admin_test','127.0.0.1',15,200,0,0,0),(83,'2026-03-25 23:07:36.827','2026-03-25 23:07:36.827',NULL,'菜单管理','PUT','/api/menu/edit',2,'admin_test','127.0.0.1',6,200,0,0,0),(84,'2026-03-25 23:08:02.283','2026-03-25 23:08:02.283',NULL,'菜单管理','PUT','/api/menu/edit',2,'admin_test','127.0.0.1',6,200,0,0,0),(85,'2026-03-25 23:08:15.004','2026-03-25 23:08:15.004',NULL,'菜单管理','PUT','/api/menu/edit',2,'admin_test','127.0.0.1',7,200,0,0,0),(86,'2026-03-25 23:12:07.039','2026-03-25 23:12:07.039',NULL,'角色管理','POST','/api/role/setRolePower',2,'admin_test','127.0.0.1',13,200,0,0,0),(87,'2026-03-25 23:14:09.663','2026-03-25 23:14:09.663',NULL,'角色管理','POST','/api/role/setRolePower',2,'admin_test','127.0.0.1',11,200,0,0,0),(88,'2026-03-25 23:14:21.135','2026-03-25 23:14:21.135',NULL,'角色管理','POST','/api/role/setRolePower',2,'admin_test','127.0.0.1',8,200,0,0,0),(89,'2026-03-25 23:15:30.842','2026-03-25 23:15:30.842',NULL,'角色管理','POST','/api/role/setRolePower',2,'admin_test','127.0.0.1',14,200,0,0,0),(90,'2026-03-25 23:22:09.321','2026-03-25 23:22:09.321',NULL,'角色管理','POST','/api/role/setRolePower',2,'admin_test','127.0.0.1',9,200,0,0,0),(91,'2026-03-25 23:22:21.209','2026-03-25 23:22:21.209',NULL,'角色管理','POST','/api/role/setRolePower',2,'admin_test','127.0.0.1',12,200,0,0,0),(92,'2026-03-25 23:24:56.403','2026-03-25 23:24:56.403',NULL,'角色管理','PUT','/api/role/edit',1,'admin','127.0.0.1',21,200,0,0,0),(93,'2026-03-25 23:49:56.125','2026-03-25 23:49:56.125',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',27,200,0,0,0),(94,'2026-03-26 00:11:45.109','2026-03-26 00:11:45.109',NULL,'参数配置','PUT','/api/config/edit',1,'admin','127.0.0.1',0,200,0,0,0),(95,'2026-03-26 00:11:46.533','2026-03-26 00:11:46.533',NULL,'参数配置','PUT','/api/config/edit',1,'admin','127.0.0.1',0,200,0,0,0),(96,'2026-03-26 00:12:59.641','2026-03-26 00:12:59.641',NULL,'参数配置','PUT','/api/config/edit',1,'admin','127.0.0.1',24,200,0,0,0),(97,'2026-03-26 00:21:11.972','2026-03-26 00:21:11.972',NULL,'参数配置','PUT','/api/config/edit',1,'admin','127.0.0.1',12,200,0,0,0),(98,'2026-03-26 20:15:25.396','2026-03-26 20:15:25.396',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',16,200,0,0,0),(99,'2026-03-26 20:15:38.790','2026-03-26 20:15:38.790',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',8,200,0,0,0),(100,'2026-03-26 20:15:44.651','2026-03-26 20:15:44.651',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',8,200,0,0,0),(101,'2026-03-26 20:15:47.949','2026-03-26 20:15:47.949',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',8,200,0,0,0),(102,'2026-03-26 20:33:55.188','2026-03-26 20:33:55.188',NULL,'参数配置','PUT','/api/config/edit',1,'admin','127.0.0.1',19,200,0,0,0),(103,'2026-03-26 20:34:04.193','2026-03-26 20:34:04.193',NULL,'参数配置','PUT','/api/config/edit',1,'admin','127.0.0.1',9,200,0,0,0),(104,'2026-03-26 20:34:17.039','2026-03-26 20:34:17.039',NULL,'参数配置','PUT','/api/config/edit',1,'admin','127.0.0.1',10,200,0,0,0),(105,'2026-03-26 20:42:43.172','2026-03-26 20:42:43.172',NULL,'参数配置','PUT','/api/config/edit',1,'admin','127.0.0.1',18,200,0,0,0),(106,'2026-03-26 20:43:19.327','2026-03-26 20:43:19.327',NULL,'参数配置','PUT','/api/config/edit',1,'admin','127.0.0.1',9,200,0,0,0),(107,'2026-03-26 20:58:01.582','2026-03-26 20:58:01.582',NULL,'参数配置','POST','/api/config/add',1,'admin','127.0.0.1',35,200,0,0,0),(108,'2026-03-26 20:58:42.970','2026-03-26 20:58:42.970',NULL,'参数配置','PUT','/api/config/edit',2,'admin_test','127.0.0.1',10,200,0,0,0),(109,'2026-03-26 20:58:54.886','2026-03-26 20:58:54.886',NULL,'参数配置','PUT','/api/config/edit',2,'admin_test','127.0.0.1',7,200,0,0,0),(110,'2026-03-26 20:59:49.888','2026-03-26 20:59:49.888',NULL,'菜单管理','PUT','/api/menu/edit',2,'admin_test','127.0.0.1',7,200,0,0,0),(111,'2026-03-26 21:00:00.682','2026-03-26 21:00:00.682',NULL,'菜单管理','PUT','/api/menu/edit',2,'admin_test','127.0.0.1',6,200,0,0,0),(112,'2026-03-26 21:05:59.566','2026-03-26 21:05:59.566',NULL,'菜单管理','PUT','/api/menu/edit',2,'admin_test','127.0.0.1',10,200,0,0,0),(113,'2026-03-26 21:10:15.525','2026-03-26 21:10:15.525',NULL,'菜单管理','PUT','/api/menu/edit',2,'admin_test','127.0.0.1',7,200,0,0,0),(114,'2026-03-26 21:28:43.364','2026-03-26 21:28:43.364',NULL,'菜单管理','POST','/api/menu/add',2,'admin_test','127.0.0.1',8,200,0,0,0),(115,'2026-03-26 21:28:57.323','2026-03-26 21:28:57.323',NULL,'菜单管理','PUT','/api/menu/edit',2,'admin_test','127.0.0.1',8,200,0,0,0),(116,'2026-03-26 21:56:24.521','2026-03-26 21:56:24.521',NULL,'角色管理','POST','/api/role/setRolePower',2,'admin_test','127.0.0.1',9,200,0,0,0),(117,'2026-03-26 22:00:01.288','2026-03-26 22:00:01.288',NULL,'参数配置','PUT','/api/config/edit',2,'admin_test','127.0.0.1',8,200,0,0,0),(118,'2026-03-26 22:06:11.290','2026-03-26 22:06:11.290',NULL,'参数配置','PUT','/api/config/edit',2,'admin_test','127.0.0.1',9,200,0,0,0),(119,'2026-03-28 21:37:15.040','2026-03-28 21:37:15.040',NULL,'参数配置','PUT','/api/config/edit',1,'admin','127.0.0.1',10,200,0,0,0),(120,'2026-03-28 21:41:15.664','2026-03-28 21:41:15.664',NULL,'参数配置','PUT','/api/config/edit',1,'admin','127.0.0.1',16,200,0,0,0),(121,'2026-03-28 21:41:36.317','2026-03-28 21:41:36.317',NULL,'参数配置','PUT','/api/config/edit',1,'admin','127.0.0.1',8,200,0,0,0),(122,'2026-03-29 13:19:12.881','2026-03-29 13:19:12.881',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',78,200,0,0,0),(123,'2026-03-29 16:47:52.454','2026-03-29 16:47:52.454',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',40,200,0,0,0),(124,'2026-03-29 16:47:56.507','2026-03-29 16:47:56.507',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',43,200,0,0,0),(125,'2026-03-29 16:48:43.794','2026-03-29 16:48:43.794',NULL,'字典管理','POST','/api/dict/type/add',1,'admin','127.0.0.1',14,200,0,0,0),(126,'2026-03-29 16:48:54.186','2026-03-29 16:48:54.186',NULL,'字典管理','POST','/api/dict/data/add',1,'admin','127.0.0.1',9,200,0,0,0),(127,'2026-03-29 16:49:07.796','2026-03-29 16:49:07.796',NULL,'字典管理','POST','/api/dict/data/add',1,'admin','127.0.0.1',14,200,0,0,0),(128,'2026-03-29 16:52:29.400','2026-03-29 16:52:29.400',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',40,200,0,0,0),(129,'2026-03-29 16:52:32.899','2026-03-29 16:52:32.899',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',36,200,0,0,0),(130,'2026-03-29 16:53:05.833','2026-03-29 16:53:05.833',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',129,200,0,0,0),(131,'2026-03-29 16:59:58.167','2026-03-29 16:59:58.167',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',38,200,0,0,0),(132,'2026-03-29 17:02:40.624','2026-03-29 17:02:40.624',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',81,200,0,0,0),(133,'2026-03-29 17:02:45.684','2026-03-29 17:02:45.684',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',32,200,0,0,0),(134,'2026-03-29 17:58:23.348','2026-03-29 17:58:23.348',NULL,'用户管理','POST','/api/user/import',1,'admin','127.0.0.1',16,200,0,0,0),(135,'2026-03-29 18:05:05.119','2026-03-29 18:05:05.119',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',17,200,0,0,0),(136,'2026-03-29 18:08:08.739','2026-03-29 18:08:08.739',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',22,200,0,0,0),(137,'2026-03-29 18:08:12.414','2026-03-29 18:08:12.414',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',16,200,0,0,0),(138,'2026-03-29 18:08:17.153','2026-03-29 18:08:17.153',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',17,200,0,0,0),(139,'2026-03-29 18:08:23.976','2026-03-29 18:08:23.976',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',14,200,0,0,0),(140,'2026-03-29 19:07:23.129','2026-03-29 19:07:23.129',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',33,200,0,0,0),(141,'2026-03-29 19:07:52.607','2026-03-29 19:07:52.607',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',9,200,0,0,0),(142,'2026-03-29 19:09:02.434','2026-03-29 19:09:02.434',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',21,200,0,0,0),(143,'2026-03-29 19:09:24.699','2026-03-29 19:09:24.699',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',10,200,0,0,0),(144,'2026-03-29 19:11:56.133','2026-03-29 19:11:56.133',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',11,200,0,0,0),(145,'2026-03-29 19:12:36.010','2026-03-29 19:12:36.010',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',10,200,0,0,0),(146,'2026-03-29 19:12:42.585','2026-03-29 19:12:42.585',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',11,200,0,0,0),(147,'2026-03-29 19:13:03.266','2026-03-29 19:13:03.266',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',12,200,0,0,0),(148,'2026-03-29 19:14:09.265','2026-03-29 19:14:09.265',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',18,200,0,0,0),(149,'2026-03-29 19:14:24.844','2026-03-29 19:14:24.844',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',9,200,0,0,0),(150,'2026-03-29 19:41:08.488','2026-03-29 19:41:08.488',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',19,200,0,0,0),(151,'2026-03-29 19:42:25.755','2026-03-29 19:42:25.755',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',17,200,0,0,0),(152,'2026-03-29 19:43:32.775','2026-03-29 19:43:32.775',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',18,200,0,0,0),(153,'2026-03-29 19:44:14.230','2026-03-29 19:44:14.230',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',15,200,0,0,0),(154,'2026-03-29 19:44:37.563','2026-03-29 19:44:37.563',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',10,200,0,0,0),(155,'2026-03-29 19:45:01.984','2026-03-29 19:45:01.984',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',16,200,0,0,0),(156,'2026-03-29 19:45:23.712','2026-03-29 19:45:23.712',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',11,200,0,0,0),(157,'2026-03-29 19:45:51.298','2026-03-29 19:45:51.298',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',60,200,0,0,0),(158,'2026-03-29 19:46:55.619','2026-03-29 19:46:55.619',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',13,200,0,0,0),(159,'2026-03-29 19:47:12.895','2026-03-29 19:47:12.895',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',7,200,0,0,0),(160,'2026-03-29 19:53:23.407','2026-03-29 19:53:23.407',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',39,200,0,0,0),(161,'2026-03-29 21:25:29.363','2026-03-29 21:25:29.363',NULL,'角色管理','POST','/api/role/setRolePower',1,'admin','127.0.0.1',12,200,0,0,0),(162,'2026-03-29 21:26:42.956','2026-03-29 21:26:42.956',NULL,'角色管理','POST','/api/role/setRolePower',1,'admin','127.0.0.1',8,200,0,0,0),(163,'2026-03-29 21:27:10.284','2026-03-29 21:27:10.284',NULL,'角色管理','POST','/api/role/setRolePower',2,'admin_test','127.0.0.1',15,200,0,0,0),(164,'2026-03-29 21:27:37.677','2026-03-29 21:27:37.677',NULL,'角色管理','POST','/api/role/setRolePower',2,'admin_test','127.0.0.1',10,200,0,0,0),(165,'2026-03-29 21:29:43.504','2026-03-29 21:29:43.504',NULL,'角色管理','POST','/api/role/setRolePower',2,'admin_test','127.0.0.1',13,200,0,0,0),(166,'2026-03-29 21:30:16.078','2026-03-29 21:30:16.078',NULL,'角色管理','POST','/api/role/setRolePower',2,'admin_test','127.0.0.1',9,200,0,0,0),(167,'2026-04-05 15:25:09.775','2026-04-05 15:25:09.775',NULL,'角色管理','POST','/api/role/setRolePower',2,'admin_test','127.0.0.1',9,200,0,0,0),(168,'2026-04-05 18:05:43.729','2026-04-05 18:05:43.729',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',18,200,0,0,0),(169,'2026-04-05 18:07:05.304','2026-04-05 18:07:05.304',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',15,200,0,0,0),(170,'2026-04-05 18:08:03.267','2026-04-05 18:08:03.267',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',15,200,0,0,0),(171,'2026-04-05 18:11:39.229','2026-04-05 18:11:39.229',NULL,'系统','POST','/api/devform/add',1,'admin','127.0.0.1',11,200,0,0,0),(172,'2026-04-05 18:23:29.231','2026-04-05 18:23:29.231',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',12,200,1,1,0),(173,'2026-04-05 18:27:22.539','2026-04-05 18:27:22.539',NULL,'用户管理','PUT','/api/user/edit',1,'admin','127.0.0.1',22,200,1,1,0),(174,'2026-04-05 18:51:21.964','2026-04-05 18:51:21.964',NULL,'系统','DELETE','/api/devform/delete/:id',1,'admin','127.0.0.1',26,200,1,1,0),(175,'2026-04-05 18:51:33.529','2026-04-05 18:51:33.529',NULL,'系统','POST','/api/devform/add',1,'admin','127.0.0.1',13,200,1,1,0),(176,'2026-04-05 18:51:53.825','2026-04-05 18:51:53.825',NULL,'系统','POST','/api/devform/add',1,'admin','127.0.0.1',6,200,1,1,0),(177,'2026-04-05 18:52:37.299','2026-04-05 18:52:37.299',NULL,'系统','POST','/api/devform/add',1,'admin','127.0.0.1',14,200,1,1,0),(178,'2026-04-05 18:52:54.971','2026-04-05 18:52:54.971',NULL,'系统','POST','/api/devform/add',1,'admin','127.0.0.1',22,200,1,1,0),(179,'2026-04-05 18:54:50.860','2026-04-05 18:54:50.860',NULL,'系统','POST','/api/devform/add',1,'admin','127.0.0.1',57,200,1,1,0),(180,'2026-04-05 18:59:26.805','2026-04-05 18:59:26.805',NULL,'系统','POST','/api/devform/field/save',1,'admin','127.0.0.1',16,200,1,1,0),(181,'2026-04-05 19:06:38.896','2026-04-05 19:06:38.896',NULL,'系统','POST','/api/devform/field/save',1,'admin','127.0.0.1',11,200,1,1,0),(182,'2026-04-05 19:06:48.581','2026-04-05 19:06:48.581',NULL,'系统','POST','/api/devform/field/save',1,'admin','127.0.0.1',13,200,1,1,0),(183,'2026-04-05 19:06:57.233','2026-04-05 19:06:57.233',NULL,'系统','POST','/api/devform/sync/:id',1,'admin','127.0.0.1',70,200,1,1,0),(184,'2026-04-05 19:17:04.533','2026-04-05 19:17:04.533',NULL,'菜单管理','POST','/api/menu/add',1,'admin','127.0.0.1',24,200,1,1,0),(185,'2026-04-05 19:20:50.377','2026-04-05 19:20:50.377',NULL,'系统','POST','/api/demo_dict/add',1,'admin','127.0.0.1',16,200,1,1,0),(186,'2026-04-05 19:27:21.092','2026-04-05 19:27:21.092',NULL,'系统','POST','/api/devform/sync/:id',1,'admin','127.0.0.1',82,200,1,1,0),(187,'2026-04-05 19:33:19.638','2026-04-05 19:33:19.638',NULL,'菜单管理','DELETE','/api/menu/delete/:id',1,'admin','127.0.0.1',12,200,1,1,0),(188,'2026-04-05 19:36:12.859','2026-04-05 19:36:12.859',NULL,'菜单管理','PUT','/api/menu/edit',1,'admin','127.0.0.1',8,200,1,1,0),(189,'2026-04-05 19:37:01.955','2026-04-05 19:37:01.955',NULL,'系统','POST','/api/devform/field/save',1,'admin','127.0.0.1',8,200,1,1,0),(190,'2026-04-05 19:37:05.697','2026-04-05 19:37:05.697',NULL,'系统','POST','/api/devform/sync/:id',1,'admin','127.0.0.1',69,200,1,1,0),(191,'2026-04-05 19:40:50.942','2026-04-05 19:40:50.942',NULL,'系统','DELETE','/api/devform/delete/:id',1,'admin','127.0.0.1',16,200,1,1,0),(192,'2026-04-05 19:44:08.908','2026-04-05 19:44:08.908',NULL,'系统','POST','/api/devform/add',1,'admin','127.0.0.1',0,200,1,1,0);
/*!40000 ALTER TABLE `sys_oper_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_role`
--

DROP TABLE IF EXISTS `sys_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_role` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  `code` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色编码',
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称',
  `description` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '描述',
  `status` bigint unsigned DEFAULT '1' COMMENT '角色状态',
  `create_by` bigint unsigned DEFAULT '0' COMMENT '创建人用户ID',
  `update_by` bigint unsigned DEFAULT '0' COMMENT '更新人用户ID',
  `delete_by` bigint unsigned DEFAULT '0' COMMENT '删除人用户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_sys_role_code` (`code`),
  KEY `idx_sys_role_delete_time` (`delete_time`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_role`
--

LOCK TABLES `sys_role` WRITE;
/*!40000 ALTER TABLE `sys_role` DISABLE KEYS */;
INSERT INTO `sys_role` VALUES (1,'2026-02-05 20:34:48.631','2026-03-25 23:24:56.383',NULL,'admin','超级管理员','',1,0,0,0),(11,'2026-02-06 00:31:29.818','2026-02-07 16:03:56.253',NULL,'admin2','管理员2','123456',1,0,0,0);
/*!40000 ALTER TABLE `sys_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_role_menu`
--

DROP TABLE IF EXISTS `sys_role_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_role_menu` (
  `sys_role_id` bigint unsigned NOT NULL,
  `sys_menu_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`sys_role_id`,`sys_menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_role_menu`
--

LOCK TABLES `sys_role_menu` WRITE;
/*!40000 ALTER TABLE `sys_role_menu` DISABLE KEYS */;
INSERT INTO `sys_role_menu` VALUES (11,1),(11,2),(11,4),(11,5),(11,6),(11,7),(11,8),(11,9),(11,10),(11,11),(11,12),(11,13),(11,15),(11,16),(11,17),(11,20),(11,21),(11,22),(11,23),(11,24),(11,25),(11,26),(11,27),(11,28),(11,29),(11,30),(11,31),(11,32),(11,33),(11,34),(11,35),(11,36),(11,37),(11,38),(11,39),(11,40),(11,41),(11,79),(11,80),(11,81),(11,82);
/*!40000 ALTER TABLE `sys_role_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user`
--

DROP TABLE IF EXISTS `sys_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  `uuid` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'uuid',
  `account` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户账号',
  `password` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户密码',
  `u_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名称',
  `u_nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户昵称',
  `u_mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户手机号码',
  `u_email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户邮箱',
  `u_avatar` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户头像',
  `last_login_time` datetime(3) DEFAULT NULL COMMENT '账号最后一次登录时间',
  `gender` bigint unsigned DEFAULT NULL COMMENT '用户性别',
  `status` bigint unsigned DEFAULT NULL COMMENT '用户状态',
  `department_id` bigint unsigned DEFAULT '0' COMMENT '所属部门ID',
  `job_level_id` bigint unsigned DEFAULT '0' COMMENT '职务级别ID(冗余首项)',
  `create_by` bigint unsigned DEFAULT '0' COMMENT '创建人用户ID',
  `update_by` bigint unsigned DEFAULT '0' COMMENT '更新人用户ID',
  `delete_by` bigint unsigned DEFAULT '0' COMMENT '删除人用户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_sys_user_uuid` (`uuid`),
  UNIQUE KEY `uni_sys_user_account` (`account`),
  KEY `idx_sys_user_delete_time` (`delete_time`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user`
--

LOCK TABLES `sys_user` WRITE;
/*!40000 ALTER TABLE `sys_user` DISABLE KEYS */;
INSERT INTO `sys_user` VALUES (1,'2026-02-12 21:14:14.000','2026-04-05 18:03:44.132',NULL,'123456789','admin','69e62a233c1907e411ce0b9823e8b9f04f1511e8d02f59beaf9d9d4450a361d0','超级管理员','管理员','18000000000','','','2026-04-05 18:03:44.131',1,1,1,0,0,0,0),(2,'2026-03-20 21:43:29.832','2026-04-05 15:24:28.179',NULL,'6b1641e6-ecfc-4566-91d9-0645561fc900','admin_test','c4318372f98f4c46ed3a32c16ee4d7a76c832886d887631c0294b3314f34edf1','李四','测试管理员','18011111111','','','2026-04-05 15:24:28.178',1,1,2,1,0,0,0),(3,'2026-03-29 17:58:23.341','2026-04-05 18:27:22.516',NULL,'bfa76c93-5a2b-44f0-899f-9c6082be259d','admin_test1','8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92','测试管理员','测试管理员','18000000002','','',NULL,2,1,0,0,0,1,0);
/*!40000 ALTER TABLE `sys_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user_department`
--

DROP TABLE IF EXISTS `sys_user_department`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user_department` (
  `sys_user_id` bigint unsigned NOT NULL,
  `sys_department_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`sys_user_id`,`sys_department_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user_department`
--

LOCK TABLES `sys_user_department` WRITE;
/*!40000 ALTER TABLE `sys_user_department` DISABLE KEYS */;
INSERT INTO `sys_user_department` VALUES (1,1),(2,2),(2,3);
/*!40000 ALTER TABLE `sys_user_department` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user_job_level`
--

DROP TABLE IF EXISTS `sys_user_job_level`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user_job_level` (
  `sys_user_id` bigint unsigned NOT NULL,
  `sys_job_level_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`sys_user_id`,`sys_job_level_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user_job_level`
--

LOCK TABLES `sys_user_job_level` WRITE;
/*!40000 ALTER TABLE `sys_user_job_level` DISABLE KEYS */;
INSERT INTO `sys_user_job_level` VALUES (2,1),(2,3);
/*!40000 ALTER TABLE `sys_user_job_level` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user_role`
--

DROP TABLE IF EXISTS `sys_user_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user_role` (
  `sys_role_id` bigint unsigned NOT NULL,
  `sys_user_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`sys_role_id`,`sys_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user_role`
--

LOCK TABLES `sys_user_role` WRITE;
/*!40000 ALTER TABLE `sys_user_role` DISABLE KEYS */;
INSERT INTO `sys_user_role` VALUES (1,1),(11,1),(11,2),(11,3);
/*!40000 ALTER TABLE `sys_user_role` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-04-05 11:49:31
