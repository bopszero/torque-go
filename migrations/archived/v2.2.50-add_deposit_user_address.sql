CREATE TABLE `deposit_user_address` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uid` bigint unsigned NOT NULL,
  `currency` char(5) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
  `network` varchar(16) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
  `address` varchar(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
  `create_time` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_currency_network_address` (`currency`,`network`,`address`) USING BTREE,
  UNIQUE KEY `uidx_uid_currency_network` (`uid`,`currency`,`network`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
;

INSERT
	INTO
	deposit_user_address(
		uid, currency, network, address, create_time
	)
SELECT
	user_id AS uid
	, currency
	, network
	, min(address)
	, unix_timestamp(min(date_created)) AS create_time
FROM
	deposit
WHERE
	status = 'Rejected'
	AND amount = 0
GROUP BY
	user_id
	, currency
	, network
;
