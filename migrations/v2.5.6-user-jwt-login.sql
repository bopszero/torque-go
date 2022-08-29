CREATE TABLE `user_jwt_refresh`
(
    `id`          char(32) CHARACTER SET ascii COLLATE ascii_general_ci     NOT NULL,
    `uid`         bigint unsigned                                           NOT NULL,
    `device_uid`  varchar(255) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
    `expire_time` bigint unsigned                                           NOT NULL,
    `rotate_time` bigint unsigned                                           NOT NULL DEFAULT '0',
    `create_time` bigint                                                    NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_expire_time` (`expire_time`) USING BTREE,
    KEY `idx_create_time` (`create_time`) USING BTREE,
    KEY `idx_device_uid_expire_time` (`device_uid`, `expire_time`) USING BTREE,
    KEY `idx_uid_expire_time` (`uid`, `expire_time`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
;
