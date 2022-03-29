CREATE TABLE `account_apple`
(
    `account`        varchar(128) NOT NULL,
    `player_id`      int(11)      NOT NULL,
    `email`          varchar(128)      DEFAULT NULL,
    `identity_token` varchar(255)      DEFAULT NULL,
    `access_token`   varchar(255)      DEFAULT NULL,
    `user_info`      varchar(255)      DEFAULT NULL,
    `create_time`    timestamp    NULL DEFAULT CURRENT_TIMESTAMP,
    `last_login`     timestamp    NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`account`),
    UNIQUE KEY `idx_player_id` (`player_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;