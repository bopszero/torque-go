ALTER TABLE withdraw
	ADD fee DECIMAL(20,8) DEFAULT 0 NOT NULL AFTER amount,
    ADD execute_status TINYINT NOT NULL DEFAULT 0 AFTER txn_hash;

ALTER TABLE torque
	ADD coin_fee DECIMAL(20,8) DEFAULT 0 NOT NULL AFTER coin_amount,
    ADD execute_status TINYINT NOT NULL DEFAULT 0 AFTER transactionhash;

CREATE TABLE `system_withdrawal_address` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `network` varchar(16) CHARACTER SET ascii NOT NULL,
  `currency` char(5) CHARACTER SET ascii NOT NULL,
  `address` varchar(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
  `key` varchar(255) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
  `status` bigint(20) NOT NULL,
  `create_time` bigint(20) NOT NULL,
  `update_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_network_address` (`network`,`address`),
  KEY `idx_status_network_currency` (`status`,`network`,`currency`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_withdrawal_request` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `address_id` int(10) unsigned NOT NULL,
  `currency` char(5) CHARACTER SET ascii NOT NULL,
  `status` tinyint(4) NOT NULL,
  `amount` decimal(32,18) NOT NULL,
  `amount_estimated_fee` decimal(32,18) NOT NULL,
  `combined_txn_hash` varchar(128) CHARACTER SET ascii DEFAULT NULL,
  `combined_signed_bytes` blob,
  `create_uid` bigint(20) unsigned NOT NULL,
  `create_time` bigint(20) NOT NULL,
  `update_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_currency_status_create_time` (`currency`,`create_time`) USING BTREE,
  KEY `idx_address_id` (`address_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_withdrawal_txn` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `request_id` bigint(20) unsigned NOT NULL,
  `ref_code` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `currency` char(5) CHARACTER SET ascii NOT NULL,
  `status` tinyint(4) NOT NULL,
  `hash` varchar(128) CHARACTER SET ascii DEFAULT NULL,
  `to_address` varchar(128) CHARACTER SET ascii COLLATE ascii_bin DEFAULT NULL,
  `output_index` bigint(20) unsigned NOT NULL,
  `fee_price` decimal(32,18) NOT NULL,
  `fee_max_quantity` int(10) unsigned NOT NULL,
  `signed_bytes` blob,
  `create_time` bigint(20) NOT NULL,
  `update_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_currency_hash_output_index` (`currency`,`hash`,`output_index`) USING BTREE,
  KEY `idx_request_id_status_to_address` (`request_id`,`status`,`to_address`) USING BTREE,
  KEY `idx_currency_to_address` (`currency`,`to_address`) USING BTREE,
  KEY `idx_ref_code` (`ref_code`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
