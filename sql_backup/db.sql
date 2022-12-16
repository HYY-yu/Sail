CREATE SCHEMA IF NOT EXISTS `sail` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `sail`.`project_group`
(
    `id`          int PRIMARY KEY AUTO_INCREMENT,
    `name`        varchar(50) UNIQUE NOT NULL,
    `create_time` timestamp          NOT NULL,
    `create_by`   int                NOT NULL,
    `delete_time` int                NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS `sail`.`project`
(
    `id`               int PRIMARY KEY AUTO_INCREMENT,
    `project_group_id` int                NOT NULL,
    `key`              varchar(50) UNIQUE NOT NULL,
    `name`             varchar(50)        NOT NULL,
    `create_time`      timestamp          NOT NULL,
    `create_by`        int                NOT NULL,
    `delete_time`      int                NOT NULL DEFAULT 0,
    INDEX (`project_group_id`)
);

CREATE TABLE IF NOT EXISTS `sail`.`namespace`
(
    `id`               int PRIMARY KEY AUTO_INCREMENT,
    `project_group_id` int          NOT NULL,
    `name`             varchar(50)  NOT NULL,
    `real_time`        bool         NOT NULL COMMENT '是否是实时发布',
    `secret_key`       varchar(100) NOT null,
    `create_time`      timestamp    NOT NULL,
    `create_by`        int          NOT NULL,
    `delete_time`      int          NOT NULL DEFAULT 0,
    UNIQUE KEY (`project_group_id`, `name`)
);

CREATE TABLE IF NOT EXISTS `sail`.`staff`
(
    `id`            int PRIMARY KEY AUTO_INCREMENT,
    `name`          varchar(30) UNIQUE NOT NULL,
    `password`      varchar(100)       NOT NULL,
    `refresh_token` varchar(200)       NOT NULL DEFAULT '',
    `create_time`   timestamp          NOT NULL,
    `create_by`     int                NOT NULL
);

CREATE TABLE IF NOT EXISTS `sail`.`staff_group_rel`
(
    `id`               int PRIMARY KEY AUTO_INCREMENT,
    `project_group_id` int     NOT NULL,
    `staff_id`         int     NOT NULL,
    `role_type`        tinyint NOT NULL COMMENT '权限角色',
    INDEX (`project_group_id`),
    INDEX (`staff_id`)
);

CREATE TABLE IF NOT EXISTS `sail`.`config`
(
    `id`               int PRIMARY KEY AUTO_INCREMENT,
    `name`             varchar(50) NOT NULL,
    `project_id`       int         NOT NULL,
    `project_group_id` int         NOT NULL COMMENT '公共配置只有project_group_id',
    `namespace_id`     int         NOT NULL,
    `is_public`        bool        NOT NULL,
    `is_link_public`   bool        NOT NULL,
    `is_encrypt`       bool        NOT NULL,
    `config_type`      varchar(10) NOT NULL,
    UNIQUE KEY (`project_id`, `namespace_id`, `project_group_id`, `name`, `config_type`)
);

CREATE TABLE IF NOT EXISTS `sail`.`config_link`
(
    `id`               int PRIMARY KEY AUTO_INCREMENT,
    `config_id`        int NOT NULL,
    `public_config_id` int NOT NULL,
    UNIQUE KEY (`config_id`, `public_config_id`),
    INDEX (`public_config_id`)
);

CREATE TABLE IF NOT EXISTS `sail`.`config_history`
(
    `id`          int PRIMARY KEY AUTO_INCREMENT,
    `config_id`   int       NOT NULL,
    `reversion`   int       NOT NULL,
    `op_type`     tinyint   NOT NULL,
    `create_time` timestamp NOT NULL,
    `create_by`   int       NOT NULL,
    UNIQUE KEY (`config_id`, `reversion`)
);

CREATE TABLE IF NOT EXISTS `sail`.`publish`
(
    `id`            int PRIMARY KEY AUTO_INCREMENT,
    `project_id`    int         NOT NULL,
    `namespace_id`  int         NOT NULL,
    `status`        tinyint     NOT NULL,
    `publish_token` varchar(50) NOT NULL,
    `create_time`   timestamp   NOT NULL Default CURRENT_TIMESTAMP,
    `update_time`   timestamp   NOT NULL Default CURRENT_TIMESTAMP,
    INDEX (`project_id`, `namespace_id`, `status`),
    UNIQUE INDEX (`publish_token`)
);

CREATE TABLE IF NOT EXISTS `sail`.`publish_config`
(
    `id`                   int PRIMARY KEY AUTO_INCREMENT,
    `publish_id`           int       NOT NULL,
    `config_id`            int       NOT NULL,
    `config_pre_reversion` int       NOT NULL,
    `status`               int       NOT NULL,
    `create_time`          timestamp NOT NULL Default CURRENT_TIMESTAMP,
    `update_time`          timestamp NOT NULL Default CURRENT_TIMESTAMP,
    INDEX (`publish_id`),
    INDEX (`config_id`)
);

CREATE TABLE IF NOT EXISTS `sail`.`publish_strategy`
(
    `id`          int PRIMARY KEY AUTO_INCREMENT,
    `publish_id`  int          NOT NULL,
    `type`        tinyint      NOT NULL COMMENT '发布类型',
    `data`        json         NOT NULL,
    `status`      tinyint      NOT NULL,
    `result`      varchar(500) NOT NULL,
    `create_time` timestamp    NOT NULL Default CURRENT_TIMESTAMP,
    `update_time` timestamp    NOT NULL Default CURRENT_TIMESTAMP,
    INDEX (`publish_id`)
);

INSERT INTO `sail`.staff (id, name, password, create_time, create_by) VALUE (1, 'Admin',
                                                                             '$2a$10$9QsXUNwjuYBdSlNA4zX/OucUcVJ/MdyqyOarzE/qdJRyw2qOjhFLS',
                                                                             NOW(), 1);

INSERT INTO `sail`.staff_group_rel (project_group_id, staff_id, role_type) VALUE (0, 1, 1);

