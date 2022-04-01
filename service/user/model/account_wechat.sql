CREATE TABLE `account_wechat`
(
    `account`      varchar(128) NOT NULL,
    `player_id`    int(11)      NOT NULL,
    `union_id`     varchar(128) NOT NULL DEFAULT '',
    `access_token` varchar(255) NOT NULL DEFAULT '',
    `user_info`    varchar(255) NOT NULL DEFAULT '',
    `create_time`  timestamp    NULL     DEFAULT CURRENT_TIMESTAMP,
    `last_login`   timestamp    NULL     DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`account`),
    UNIQUE KEY `idx_player_id` (`player_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;