CREATE TABLE `player`
(
    `player_id`        int(11)      NOT NULL NULL DEFAULT 0,
    `name`             varchar(255) NOT NULL      DEFAULT '',
    `gender`           int(11)      NOT NULL      DEFAULT '1' COMMENT '性别',
    `avator_url`       varchar(255) NOT NULL      DEFAULT '' COMMENT '头像地址',
    `invitation_id`    int(11)      NOT NULL      DEFAULT '0' COMMENT '邀请玩家id',
    `channel`          varchar(255) NOT NULL      DEFAULT '' COMMENT '渠道来源',
    `vip_lv`           int(11)      NOT NULL      DEFAULT '0' COMMENT 'vip等级',
    `vip_exp`          int(20)      NOT NULL      DEFAULT '0' COMMENT 'vip升级经验',
    `level`            int(11)      NOT NULL      DEFAULT '1' COMMENT '等级',
    `level_exp`        int(20)      NOT NULL      DEFAULT '0' COMMENT '等级经验',
    `ori_account_type` int(11)      NOT NULL COMMENT '玩家原始注册方式',
    `ori_account`      varchar(255) NOT NULL COMMENT '玩家原始注册账户名',
    `create_time`      timestamp    NULL          DEFAULT CURRENT_TIMESTAMP,
    `update_time`      timestamp    NULL          DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`player_id`),
    KEY `idx_account` (`ori_account`) USING HASH COMMENT '账户名索引'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4