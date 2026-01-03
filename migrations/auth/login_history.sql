CREATE TABLE `login_history` (
    `id` char(36) NOT NULL COMMENT '登录历史ID (UUID v7)',
    `user_id` char(36) NOT NULL COMMENT '用户ID (UUID v7)',
    `ip` varchar(45) NOT NULL COMMENT '登录IP地址',
    `device_type` varchar(20) DEFAULT NULL COMMENT '设备类型：Web/Android/iOS',
    `device_id` varchar(255) DEFAULT NULL COMMENT '设备标识',
    `user_agent` varchar(500) DEFAULT NULL COMMENT 'User-Agent',
    `login_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_login_at` (`login_at`),
    CONSTRAINT `fk_login_history_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录历史表';

