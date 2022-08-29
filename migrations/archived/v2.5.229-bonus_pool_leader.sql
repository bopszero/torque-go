CREATE TABLE `payout_stats`
(
    `id`                             int unsigned   NOT NULL AUTO_INCREMENT,
    `date`                           date           NOT NULL,
    `total_profit_amount`            decimal(32, 8) NOT NULL,
    `remaining_affiliate_amount`     decimal(32, 8) NOT NULL,
    `remaining_leader_reward_amount` decimal(32, 8) NOT NULL,
    `create_time`                    bigint         NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uidx_date` (`date`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE `leader_bonus_pool_execution`
(
    `id`                           int(10) unsigned             NOT NULL AUTO_INCREMENT,
    `hash`                         char(64) CHARACTER SET ascii NOT NULL,
    `from_date`                    date                         NOT NULL,
    `to_date`                      date                         NOT NULL,
    `total_amount`                 decimal(32, 8)               NOT NULL,
    `status`                       smallint(6)                  NOT NULL,
    `tier_senior_rate`             decimal(7, 4)                NOT NULL,
    `tier_senior_receiver_count`   smallint(5) unsigned         NOT NULL,
    `tier_regional_rate`           decimal(7, 4)                NOT NULL,
    `tier_regional_receiver_count` smallint(5) unsigned         NOT NULL,
    `tier_global_rate`             decimal(7, 4)                NOT NULL,
    `tier_global_receiver_count`   smallint(5) unsigned         NOT NULL,
    `create_time`                  bigint(20)                   NOT NULL,
    `update_time`                  bigint(20)                   NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uidx_hash` (`hash`),
    KEY `idx_status_to_date_from_date` (`status`, `to_date`, `from_date`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE `leader_bonus_pool_detail`
(
    `id`          bigint(20) unsigned                     NOT NULL AUTO_INCREMENT,
    `exec_id`     int(10) unsigned                        NOT NULL,
    `tier_type`   tinyint(3) unsigned                     NOT NULL,
    `uid`         bigint(20) unsigned                     NOT NULL,
    `order_id`    bigint(20) unsigned                     NOT NULL,
    `note`        varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `create_time` bigint(20)                              NOT NULL,
    `update_time` bigint(20)                              NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uidx_exec_id_tier_uid` (`exec_id`, `tier_type`, `uid`) USING BTREE,
    KEY `idx_order_id` (`order_id`) USING BTREE,
    KEY `idx_tier_uid_create_time` (`tier_type`, `uid`, `create_time`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
