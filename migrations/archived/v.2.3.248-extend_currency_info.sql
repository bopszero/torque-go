ALTER TABLE currency
	ADD COLUMN `icon_url` varchar(255) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL AFTER priority_wallet,
	ADD COLUMN `is_fiat` tinyint(1) NOT NULL AFTER `icon_url`,
	ADD COLUMN `symbol` varchar(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL AFTER `is_fiat`;

ALTER TABLE currency
	MODIFY COLUMN priority_trading TINYINT UNSIGNED NOT NULL,
	MODIFY COLUMN priority_wallet TINYINT UNSIGNED NOT NULL,
	ADD priority_display SMALLINT UNSIGNED NOT NULL AFTER price_usdt;

ALTER TABLE currency
	ADD decimal_places SMALLINT UNSIGNED DEFAULT 8 NOT NULL AFTER is_fiat
;
