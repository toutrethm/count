CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL COMMENT '真实姓名',
  `phone` varchar(32) NOT NULL COMMENT '登录账号',
  `password_hash` varchar(100) NOT NULL,
  `role` varchar(16) NOT NULL DEFAULT 'worker' COMMENT 'admin / worker',
  `station_role` varchar(255) NOT NULL DEFAULT '' COMMENT '能力码，逗号分隔',
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_phone` (`phone`),
  KEY `idx_users_role` (`role`),
  KEY `idx_users_station_role` (`station_role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `processes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(32) NOT NULL,
  `name` varchar(64) NOT NULL,
  `station_role` varchar(255) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_processes_code` (`code`),
  KEY `idx_processes_station_role` (`station_role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `worker_processes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `process_id` bigint unsigned NOT NULL,
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_worker_processes` (`user_id`, `process_id`),
  KEY `idx_worker_processes_process_id` (`process_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_no` varchar(64) NOT NULL,
  `qr_token` varchar(64) NOT NULL,
  `customer_name` varchar(64) NOT NULL DEFAULT '',
  `product_name` varchar(128) NOT NULL DEFAULT '',
  `spec` varchar(128) DEFAULT NULL,
  `quantity` decimal(18,3) NOT NULL DEFAULT 0,
  `delivery_date` date DEFAULT NULL,
  `status` varchar(32) NOT NULL DEFAULT 'draft',
  `created_by` bigint unsigned DEFAULT NULL,
  `remark` varchar(255) DEFAULT NULL,
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_orders_no` (`order_no`),
  UNIQUE KEY `uk_orders_qr_token` (`qr_token`),
  KEY `idx_orders_status` (`status`),
  KEY `idx_orders_created_by` (`created_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `order_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_id` bigint unsigned NOT NULL,
  `item_no` varchar(64) NOT NULL,
  `component_type` varchar(64) NOT NULL,
  `part_name` varchar(128) NOT NULL,
  `model` varchar(128) DEFAULT NULL,
  `spec` varchar(128) DEFAULT NULL,
  `dimensions` json DEFAULT NULL,
  `material` varchar(128) DEFAULT NULL,
  `quantity` decimal(18,3) NOT NULL DEFAULT 0,
  `unit` varchar(16) DEFAULT NULL,
  `remark` varchar(255) DEFAULT NULL,
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_order_items_order_id` (`order_id`),
  KEY `idx_order_items_component_type` (`component_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `order_processes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_id` bigint unsigned NOT NULL,
  `order_item_id` bigint unsigned NOT NULL,
  `process_id` bigint unsigned NOT NULL,
  `station_role` varchar(255) NOT NULL,
  `sort` int NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_processes` (`order_id`, `order_item_id`, `process_id`),
  KEY `idx_order_processes_process_id` (`process_id`),
  KEY `idx_order_processes_station_role` (`station_role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `wage_rules` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(64) NOT NULL,
  `component_type` varchar(64) NOT NULL,
  `min_diameter` decimal(10,2) NOT NULL DEFAULT 0,
  `max_diameter` decimal(10,2) NOT NULL DEFAULT 0,
  `min_length` decimal(10,2) NOT NULL DEFAULT 0,
  `max_length` decimal(10,2) NOT NULL DEFAULT 0,
  `base_unit_price` decimal(10,4) NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_wage_rules_code` (`code`),
  KEY `idx_wage_rules_component_type` (`component_type`),
  KEY `idx_wage_rules_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `scan_records` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_id` bigint unsigned NOT NULL,
  `order_item_id` bigint unsigned NOT NULL,
  `order_process_id` bigint unsigned NOT NULL,
  `process_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `station_role` varchar(255) NOT NULL,
  `wage_rule_id` bigint unsigned DEFAULT NULL,
  `wage_rule_code` varchar(64) NOT NULL DEFAULT '',
  `wage_unit_price` decimal(10,4) NOT NULL DEFAULT 0,
  `wage_amount` decimal(10,4) NOT NULL DEFAULT 0,
  `scanned_at` datetime(0) NOT NULL,
  `source` varchar(32) NOT NULL DEFAULT 'scan',
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_scan_records_order_process_id` (`order_process_id`),
  KEY `idx_scan_records_user_time` (`user_id`, `scanned_at`),
  KEY `idx_scan_records_order_time` (`order_id`, `scanned_at`),
  KEY `idx_scan_records_process_time` (`process_id`, `scanned_at`),
  KEY `idx_scan_records_station_role` (`station_role`),
  KEY `idx_scan_records_wage_rule_id` (`wage_rule_id`),
  KEY `idx_scan_records_wage_rule_code` (`wage_rule_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `processes` (`code`, `name`, `station_role`, `sort`, `status`)
VALUES
  ('blanking_center', '切料和打中心孔', 'blanking_center', 10, 1),
  ('turn_outer', '车外圆', 'turn_outer', 20, 1),
  ('turn_head', '车大头', 'turn_head', 30, 1),
  ('center_hole', '打中心孔', 'center_hole', 40, 1),
  ('turn_head_center', '车大头和打中心孔', 'turn_head_center', 50, 1),
  ('drill_tap_small', '钻孔攻牙', 'drill_tap_small', 60, 1),
  ('drill_tap_batch', '批量钻孔攻牙', 'drill_tap_batch', 70, 1),
  ('turn_sleeve', '车套', 'turn_sleeve_auto,turn_sleeve_manual', 80, 1)
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `station_role` = VALUES(`station_role`),
  `sort` = VALUES(`sort`),
  `status` = VALUES(`status`);

INSERT INTO `wage_rules` (`code`, `component_type`, `min_diameter`, `max_diameter`, `min_length`, `max_length`, `base_unit_price`, `status`)
VALUES
  ('guide_pillar_d0_25_l0_150', 'guide_pillar', 0, 25, 0, 150, 0.2000, 1),
  ('guide_pillar_d0_25_l151_250', 'guide_pillar', 0, 25, 151, 250, 0.2800, 1),
  ('guide_pillar_d0_25_l251_350', 'guide_pillar', 0, 25, 251, 350, 0.3500, 1),
  ('guide_pillar_d30_l0_150', 'guide_pillar', 30, 30, 0, 150, 0.2400, 1),
  ('guide_pillar_d30_l151_250', 'guide_pillar', 30, 30, 151, 250, 0.3500, 1),
  ('guide_pillar_d30_l251_350', 'guide_pillar', 30, 30, 251, 350, 0.4000, 1),
  ('guide_pillar_d35_l0_150', 'guide_pillar', 35, 35, 0, 150, 0.3000, 1),
  ('guide_pillar_d35_l151_250', 'guide_pillar', 35, 35, 151, 250, 0.3800, 1),
  ('guide_pillar_d35_l251_350', 'guide_pillar', 35, 35, 251, 350, 0.4500, 1),
  ('guide_pillar_d40_l0_150', 'guide_pillar', 40, 40, 0, 150, 0.3600, 1),
  ('guide_pillar_d40_l151_250', 'guide_pillar', 40, 40, 151, 250, 0.4200, 1),
  ('guide_pillar_d40_l251_350', 'guide_pillar', 40, 40, 251, 350, 0.5300, 1),
  ('top_pin_d0_30_l0_150', 'top_pin', 0, 30, 0, 150, 0.2000, 1),
  ('top_pin_d0_30_l151_250', 'top_pin', 0, 30, 151, 250, 0.2600, 1),
  ('top_pin_d0_30_l251_350', 'top_pin', 0, 30, 251, 350, 0.3200, 1),
  ('top_pin_d35_l0_150', 'top_pin', 35, 35, 0, 150, 0.2600, 1),
  ('top_pin_d35_l151_250', 'top_pin', 35, 35, 151, 250, 0.3200, 1),
  ('top_pin_d35_l251_350', 'top_pin', 35, 35, 251, 350, 0.3800, 1)
ON DUPLICATE KEY UPDATE
  `component_type` = VALUES(`component_type`),
  `min_diameter` = VALUES(`min_diameter`),
  `max_diameter` = VALUES(`max_diameter`),
  `min_length` = VALUES(`min_length`),
  `max_length` = VALUES(`max_length`),
  `base_unit_price` = VALUES(`base_unit_price`),
  `status` = VALUES(`status`);

INSERT INTO `users` (`username`, `phone`, `password_hash`, `role`, `station_role`, `status`)
VALUES
  ('管理员', '13800000000', '$2a$10$7aWlx1WSuW.bL/E5FyKnsehUSjfHC9bHtmcpSrPXPLjZipt32w8R6', 'admin', '', 1),
  ('赵吴洋', '17742523836', '$2a$10$cnVyQvDo7/8kFEJPEWdeB.DfvCry/ffuq6Um4BUYRkEO5xyNFH.ZK', 'worker', 'blanking_center,center_hole,drill_tap_batch', 1),
  ('刘小印', '13959578057', '$2a$10$4zNlnrOj8i3rapBbxt7.6eYMds1FDNBM4irXLK4ZW8wJ4wYazDBym', 'worker', 'turn_outer,turn_sleeve_manual', 1),
  ('王海龙', '18250505219', '$2a$10$INnGSmPN3s7R283R6Z.V5e/o8JUm/k3wQSWCgypbxg5VERX3voMWi', 'worker', 'turn_outer,turn_sleeve_auto', 1),
  ('陈家兴', '15106017652', '$2a$10$cFL8BRkdf9nrdTtMie43E.MWsal6f1NiiuBqxaZst2G/YBj/MB7KO', 'worker', 'turn_head,center_hole,drill_tap_small,turn_sleeve_auto', 1),
  ('郑志权', '15501148077', '$2a$10$klCnxwYQJJhnfqCDbLH3/uepZPfKF1ERI.V2cmuztLmfyXxZ4WPfi', 'worker', 'turn_sleeve_manual', 1),
  ('郑志勇', '18039005017', '$2a$10$2dEEeLgUQH2raZgoCsFqzOpwmI6D0f49OnJKMBoeCW.6Ef9XZ6ei.', 'worker', 'turn_sleeve_manual', 1),
  ('吴丽先', '13950776176', '$2a$10$6MD6FUanLdBS4HZH2KXOdO3pcpXq6oP5Hc7hrTkggvxQ6SsH5ingm', 'worker', 'turn_sleeve_manual', 1)
ON DUPLICATE KEY UPDATE
  `username` = VALUES(`username`),
  `password_hash` = VALUES(`password_hash`),
  `role` = VALUES(`role`),
  `station_role` = VALUES(`station_role`),
  `status` = VALUES(`status`);
