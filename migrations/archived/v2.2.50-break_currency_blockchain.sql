ALTER TABLE currency
	CHANGE is_on_trading priority_trading tinyint NOT NULL,
	CHANGE is_on_wallet priority_wallet tinyint NOT NULL,
	DROP COLUMN network,
	DROP COLUMN latest_block_height,
	DROP INDEX uidx_code_network,
	ADD UNIQUE INDEX uidx_code(`code`)
;

CREATE TABLE `blockchain_network` (
  `id` smallint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `latest_block_height` bigint unsigned NOT NULL,
  `deposit_min_confirmations` smallint unsigned NOT NULL,
  `update_time` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
;
