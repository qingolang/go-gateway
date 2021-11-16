-- MySQL dump 10.13  Distrib 8.0.26, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: go_gateway_test2
-- ------------------------------------------------------
-- Server version	8.0.26

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `area`
--

DROP TABLE IF EXISTS `area`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `area` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `area_name` varchar(255) NOT NULL,
  `city_id` int NOT NULL,
  `user_id` int NOT NULL,
  `update_at` datetime NOT NULL,
  `create_at` datetime NOT NULL,
  `delete_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb3 COMMENT='area';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `area`
--

LOCK TABLES `area` WRITE;
/*!40000 ALTER TABLE `area` DISABLE KEYS */;
INSERT INTO `area` VALUES (2,'area_name',1,2,'2019-06-15 00:00:00','2019-06-15 00:00:00','2019-06-15 00:00:00');
/*!40000 ALTER TABLE `area` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gateway_admin`
--

DROP TABLE IF EXISTS `gateway_admin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gateway_admin` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `user_name` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名',
  `salt` varchar(50) NOT NULL DEFAULT '' COMMENT '盐',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
  `create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `update_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '更新时间',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb3 COMMENT='管理员表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gateway_admin`
--

LOCK TABLES `gateway_admin` WRITE;
/*!40000 ALTER TABLE `gateway_admin` DISABLE KEYS */;
INSERT INTO `gateway_admin` VALUES (1,'admin','admin','2823d896e9822c0833d41d4904f0c00756d718570fce49b9a379a62c804689d3','2020-04-10 16:42:05','2020-04-21 06:35:08',0);
/*!40000 ALTER TABLE `gateway_admin` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gateway_app`
--

DROP TABLE IF EXISTS `gateway_app`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gateway_app` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `app_id` varchar(255) NOT NULL DEFAULT '' COMMENT '租户id',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '租户名称',
  `secret` varchar(255) NOT NULL DEFAULT '' COMMENT '密钥',
  `white_ips` varchar(1000) NOT NULL DEFAULT '' COMMENT 'ip白名单，支持前缀匹配',
  `qpd` bigint NOT NULL DEFAULT '0' COMMENT '日请求量限制',
  `qps` bigint NOT NULL DEFAULT '0' COMMENT '每秒请求量限制',
  `create_at` datetime NOT NULL COMMENT '添加时间',
  `update_at` datetime NOT NULL COMMENT '更新时间',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除 1=删除',
  `status` int DEFAULT '1' COMMENT '状态',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='网关租户表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gateway_app`
--

LOCK TABLES `gateway_app` WRITE;
/*!40000 ALTER TABLE `gateway_app` DISABLE KEYS */;
/*!40000 ALTER TABLE `gateway_app` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gateway_service_access_control`
--

DROP TABLE IF EXISTS `gateway_service_access_control`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gateway_service_access_control` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_id` bigint NOT NULL DEFAULT '0' COMMENT '服务id',
  `open_auth` tinyint NOT NULL DEFAULT '0' COMMENT '是否开启权限 1=开启',
  `black_list` varchar(1000) NOT NULL DEFAULT '' COMMENT '黑名单ip',
  `white_list` varchar(1000) NOT NULL DEFAULT '' COMMENT '白名单ip',
  `clientip_flow_limit` int NOT NULL DEFAULT '0' COMMENT '客户端ip限流',
  `service_flow_limit` int NOT NULL DEFAULT '0' COMMENT '服务端限流',
  `open_api_white_list` int DEFAULT '0' COMMENT '是否开启api白名单 它依赖于open_auth是否开启JTW校验',
  `open_white_list` int DEFAULT '0' COMMENT '是否开启IP白名单',
  `open_black_list` int DEFAULT '0' COMMENT ' 是否开启IP黑名单',
  `api_white_list` varchar(1000) DEFAULT NULL COMMENT 'api白名单',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='网关权限控制表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gateway_service_access_control`
--

LOCK TABLES `gateway_service_access_control` WRITE;
/*!40000 ALTER TABLE `gateway_service_access_control` DISABLE KEYS */;
/*!40000 ALTER TABLE `gateway_service_access_control` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gateway_service_grpc_rule`
--

DROP TABLE IF EXISTS `gateway_service_grpc_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gateway_service_grpc_rule` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_id` bigint NOT NULL DEFAULT '0' COMMENT '服务id',
  `port` int NOT NULL DEFAULT '0' COMMENT '端口',
  `header_transfor` varchar(5000) NOT NULL DEFAULT '' COMMENT 'header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='网关路由匹配表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gateway_service_grpc_rule`
--

LOCK TABLES `gateway_service_grpc_rule` WRITE;
/*!40000 ALTER TABLE `gateway_service_grpc_rule` DISABLE KEYS */;
/*!40000 ALTER TABLE `gateway_service_grpc_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gateway_service_http_rule`
--

DROP TABLE IF EXISTS `gateway_service_http_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gateway_service_http_rule` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_id` bigint NOT NULL COMMENT '服务id',
  `rule_type` tinyint NOT NULL DEFAULT '0' COMMENT '匹配类型 0=url前缀url_prefix 1=域名domain ',
  `rule` varchar(255) NOT NULL DEFAULT '' COMMENT 'type=domain表示域名，type=url_prefix时表示url前缀',
  `need_https` tinyint NOT NULL DEFAULT '0' COMMENT '支持https 1=支持',
  `need_strip_uri` tinyint NOT NULL DEFAULT '0' COMMENT '启用strip_uri 1=启用',
  `need_websocket` tinyint NOT NULL DEFAULT '0' COMMENT '是否支持websocket 1=支持',
  `url_rewrite` varchar(5000) NOT NULL DEFAULT '' COMMENT 'url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔',
  `header_transfor` varchar(5000) NOT NULL DEFAULT '' COMMENT 'header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='网关路由匹配表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gateway_service_http_rule`
--

LOCK TABLES `gateway_service_http_rule` WRITE;
/*!40000 ALTER TABLE `gateway_service_http_rule` DISABLE KEYS */;
/*!40000 ALTER TABLE `gateway_service_http_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gateway_service_info`
--

DROP TABLE IF EXISTS `gateway_service_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gateway_service_info` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `load_type` tinyint NOT NULL DEFAULT '0' COMMENT '负载类型 0=http 1=tcp 2=grpc',
  `service_name` varchar(255) NOT NULL DEFAULT '' COMMENT '服务名称 6-128 数字字母下划线',
  `service_desc` varchar(255) NOT NULL DEFAULT '' COMMENT '服务描述',
  `create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '添加时间',
  `update_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '更新时间',
  `is_delete` tinyint DEFAULT '0' COMMENT '是否删除 1=删除',
  `status` int DEFAULT '1' COMMENT '服务状态',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='网关基本信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gateway_service_info`
--

LOCK TABLES `gateway_service_info` WRITE;
/*!40000 ALTER TABLE `gateway_service_info` DISABLE KEYS */;
/*!40000 ALTER TABLE `gateway_service_info` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gateway_service_load_balance`
--

DROP TABLE IF EXISTS `gateway_service_load_balance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gateway_service_load_balance` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_id` bigint NOT NULL DEFAULT '0' COMMENT '服务id',
  `check_method` tinyint NOT NULL DEFAULT '0' COMMENT '检查方法 0=tcpchk,检测端口是否握手成功',
  `check_timeout` int NOT NULL DEFAULT '0' COMMENT 'check超时时间,单位s',
  `check_interval` int NOT NULL DEFAULT '0' COMMENT '检查间隔, 单位s',
  `round_type` tinyint NOT NULL DEFAULT '2' COMMENT '轮询方式 0=random 1=round-robin 2=weight_round-robin 3=ip_hash',
  `ip_list` varchar(2000) NOT NULL DEFAULT '' COMMENT 'ip列表',
  `weight_list` varchar(2000) NOT NULL DEFAULT '' COMMENT '权重列表',
  `forbid_list` varchar(2000) NOT NULL DEFAULT '' COMMENT '禁用ip列表',
  `upstream_connect_timeout` int NOT NULL DEFAULT '0' COMMENT '建立连接超时, 单位s',
  `upstream_header_timeout` int NOT NULL DEFAULT '0' COMMENT '获取header超时, 单位s',
  `upstream_idle_timeout` int NOT NULL DEFAULT '0' COMMENT '链接最大空闲时间, 单位s',
  `upstream_max_idle` int NOT NULL DEFAULT '0' COMMENT '最大空闲链接数',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='网关负载表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gateway_service_load_balance`
--

LOCK TABLES `gateway_service_load_balance` WRITE;
/*!40000 ALTER TABLE `gateway_service_load_balance` DISABLE KEYS */;
/*!40000 ALTER TABLE `gateway_service_load_balance` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gateway_service_tcp_rule`
--

DROP TABLE IF EXISTS `gateway_service_tcp_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gateway_service_tcp_rule` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_id` bigint NOT NULL COMMENT '服务id',
  `port` int NOT NULL DEFAULT '0' COMMENT '端口号',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='网关路由匹配表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gateway_service_tcp_rule`
--

LOCK TABLES `gateway_service_tcp_rule` WRITE;
/*!40000 ALTER TABLE `gateway_service_tcp_rule` DISABLE KEYS */;
/*!40000 ALTER TABLE `gateway_service_tcp_rule` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-11-16 15:56:49
