CREATE TABLE `users` (
    `id` char(36) NOT NULL COMMENT '用户ID (UUID v7)',
    `mobile` varchar(11) NOT NULL COMMENT '手机号',
    `password_hash` varchar(255) NOT NULL COMMENT '密码哈希',
    `status` tinyint NOT NULL DEFAULT 1 COMMENT '账户状态：1-启用，2-禁用，3-锁定',
    `locked_at` datetime DEFAULT NULL COMMENT '锁定时间',
    `lock_reason` varchar(255) DEFAULT NULL COMMENT '锁定原因',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_mobile` (`mobile`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

