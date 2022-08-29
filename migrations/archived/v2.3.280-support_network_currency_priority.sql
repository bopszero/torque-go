ALTER TABLE network_currency
	ADD priority SMALLINT UNSIGNED NOT NULL AFTER network
;

ALTER TABLE blockchain_network
	ADD `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL AFTER `code`,
	ADD `token_transfer_code_name` varchar(16) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL AFTER `name`,
	ADD `currency` char(5) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL AFTER `token_transfer_code_name`
;
