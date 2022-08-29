RENAME TABLE wallet_address TO `deposit_address_stock`;
ALTER TABLE deposit_address_stock
    CHANGE wallet_address_id id BIGINT UNSIGNED auto_increment NOT NULL,
    CHANGE date_created create_date datetime DEFAULT CURRENT_TIMESTAMP NOT NULL,
    MODIFY COLUMN address varchar(128) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    MODIFY COLUMN is_used BOOL DEFAULT 0 NOT NULL,
    DROP `deleted`,
	ADD currency char(5) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL AFTER coin_id,
	ADD network varchar(16) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL AFTER currency,
    ADD UNIQUE INDEX `uidx_currency_network_address` (`currency`,`network`,`address`)
;
UPDATE
	deposit_address_stock
SET
	currency = 'LTC',
	network = 'LTC'
WHERE
	coin_id = 1;
UPDATE
	deposit_address_stock
SET
	currency = 'BTC',
	network = 'BTC'
WHERE
	coin_id = 2;
UPDATE
	deposit_address_stock
SET
	currency = 'BCH',
	network = 'BCH'
WHERE
	coin_id = 20;
UPDATE
	deposit_address_stock
SET
	currency = 'ETH',
	network = 'ETH'
WHERE
	coin_id = 3;
UPDATE
	deposit_address_stock
SET
	currency = 'USDT',
	network = 'ETH'
WHERE
	coin_id = 4;



UPDATE
	deposit
SET
	txn_hash = uuid()
WHERE
	txn_hash = ''
;
ALTER TABLE deposit
	ADD currency char(5) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL AFTER coin_id,
	ADD network varchar(16) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL AFTER currency,
	ADD close_time bigint NOT NULL DEFAULT 0 AFTER status,
    DROP wallet_address_id,
    ADD UNIQUE INDEX `uidx_currency_network_txn_hash_txn_index` (`currency`,`network`,`txn_hash`,`txn_to_index`)
;
UPDATE
	deposit
SET
	currency = 'LTC',
	network = 'LTC'
WHERE
	coin_id = 1;
UPDATE
	deposit
SET
	currency = 'BTC',
	network = 'BTC'
WHERE
	coin_id = 2;
UPDATE
	deposit
SET
	currency = 'BCH',
	network = 'BCH'
WHERE
	coin_id = 20;
UPDATE
	deposit
SET
	currency = 'ETH',
	network = 'ETH'
WHERE
	coin_id = 3;
UPDATE
	deposit
SET
	currency = 'USDT',
	network = 'ETH'
WHERE
	coin_id = 4;



UPDATE
	torque
SET
	code = CONCAT('R-', SUBSTRING(REPLACE(uuid(), '-', ''), 1, 13))
WHERE
	code IS NULL
;
ALTER TABLE torque
	MODIFY COLUMN code char(15) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
	ADD currency char(5) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL AFTER coin_id,
	ADD network varchar(16) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL AFTER currency,
	ADD close_time bigint NOT NULL DEFAULT 0 AFTER status,
	ADD UNIQUE INDEX idx_withdraw_token(withdraw_confirm_token),
	ADD INDEX idx_uid_currencycreate_time(is_reinvest,user_id,currency,date_created)
;
UPDATE
	torque
SET
	currency = 'LTC',
	network = 'LTC'
WHERE
	coin_id = 1;
UPDATE
	torque
SET
	currency = 'BTC',
	network = 'BTC'
WHERE
	coin_id = 2;
UPDATE
	torque
SET
	currency = 'BCH',
	network = 'BCH'
WHERE
	coin_id = 20;
UPDATE
	torque
SET
	currency = 'ETH',
	network = 'ETH'
WHERE
	coin_id = 3;
UPDATE
	torque
SET
	currency = 'USDT',
	network = 'ETH'
WHERE
	coin_id = 4;



UPDATE
	withdraw
SET
	withdraw_confirm_token = NULL
WHERE
	withdraw_confirm_token = ''
;
ALTER TABLE withdraw
	MODIFY COLUMN code char(15) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
	ADD currency char(5) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL AFTER coin_id,
	ADD network varchar(16) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL AFTER currency,
	ADD close_time bigint NOT NULL DEFAULT 0 AFTER status,
	ADD UNIQUE INDEX uidx_code(code),
	ADD UNIQUE INDEX idx_withdraw_token(withdraw_confirm_token),
	ADD INDEX idx_uid_currency_create_time(user_id,currency,date_created)
;
UPDATE
	withdraw
SET
	currency = 'LTC',
	network = 'LTC'
WHERE
	coin_id = 1;
UPDATE
	withdraw
SET
	currency = 'BTC',
	network = 'BTC'
WHERE
	coin_id = 2;
UPDATE
	withdraw
SET
	currency = 'BCH',
	network = 'BCH'
WHERE
	coin_id = 20;
UPDATE
	withdraw
SET
	currency = 'ETH',
	network = 'ETH'
WHERE
	coin_id = 3;
UPDATE
	withdraw
SET
	currency = 'USDT',
	network = 'ETH'
WHERE
	coin_id = 4;



RENAME TABLE crypto_txns TO `deposit_crypto_txn`;
ALTER TABLE deposit_crypto_txn
	CHANGE coin currency char(5) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
	ADD network varchar(16) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL AFTER currency,
	CHANGE fromAddress from_address varchar(128) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
	CHANGE toAddress to_address varchar(128) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
	CHANGE txHash hash varchar(128) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
	CHANGE txAmt amount decimal(32,18) NOT NULL,
	CHANGE confirmation confirmations INT UNSIGNED NOT NULL,
	CHANGE `time` block_time BIGINT NOT NULL,
	CHANGE status is_accepted BOOL DEFAULT FALSE NOT NULL,
	DROP inserted_at,
	DROP updated_at,
	ADD create_time BIGINT NOT NULL,
	ADD update_time BIGINT NOT NULL,
    ADD UNIQUE INDEX `uidx_currency_network_hash_to_index` (`currency`,`network`,`hash`,`to_index`)
;
UPDATE
	deposit_crypto_txn
SET
	network = 'LTC'
WHERE
	currency = 'LTC';
UPDATE
	deposit_crypto_txn
SET
	network = 'BTC'
WHERE
	currency = 'BTC';
UPDATE
	deposit_crypto_txn
SET
	network = 'BCH'
WHERE
	currency = 'BCH';
UPDATE
	deposit_crypto_txn
SET
	network = 'ETH'
WHERE
	currency = 'ETH';
UPDATE
	deposit_crypto_txn
SET
	network = 'ETH'
WHERE
	currency = 'USDT';



ALTER TABLE
	system_withdrawal_request ADD network varchar(16) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL AFTER currency
;



CREATE TABLE `network_currency` (
  `id` smallint(5) unsigned NOT NULL AUTO_INCREMENT,
  `currency` char(5) CHARACTER SET ascii NOT NULL,
  `network` varchar(16) CHARACTER SET ascii NOT NULL,
  `withdrawal_fee` decimal(32,18) NOT NULL,
  `create_time` bigint(20) NOT NULL,
  `update_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_currency_network` (`currency`,`network`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci
