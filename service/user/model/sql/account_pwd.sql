CREATE TABLE `account_pwd`
(
    `account`     varchar(255) NOT NULL,
    `player_id`   int(11)      NOT NULL,
    `pwd`         varchar(255)      DEFAULT NULL,
    `create_time` timestamp    NULL DEFAULT CURRENT_TIMESTAMP,
    `last_login`  timestamp    NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`account`),
    UNIQUE KEY `idx_player_id` (`player_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;