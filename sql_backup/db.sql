CREATE SCHEMA `sail`;

CREATE TABLE `sail`.`project_group`
(
    `id`          int PRIMARY KEY AUTO_INCREMENT,
    `name`        varchar(50) NOT NULL,
    `create_time` timestamp   NOT NULL,
    `create_by`   int         NOT NULL,
    `delete_time` int         NOT NULL DEFAULT 0
);

create
unique index project_group_name_uindex
	on project_group (name);

CREATE TABLE `sail`.`project`
(
    `id`               int PRIMARY KEY AUTO_INCREMENT,
    `project_group_id` int         NOT NULL,
    `key`              varchar(50) NOT NULL,
    `name`             varchar(50) NOT NULL,
    `create_time`      timestamp   NOT NULL,
    `create_by`        int         NOT NULL,
    `delete_time`      int         NOT NULL DEFAULT 0
);

CREATE TABLE `sail`.`namespace`
(
    `id`               int PRIMARY KEY AUTO_INCREMENT,
    `project_group_id` int          NOT NULL,
    `name`             varchar(50)  NOT NULL,
    `real_time`        bool         NOT NULL COMMENT '是否是实时发布',
    `secret_key`       varchar(100) NOT null,
    `create_time`      timestamp    NOT NULL,
    `create_by`        int          NOT NULL,
    `delete_time`      int          NOT NULL DEFAULT 0
);

CREATE TABLE `sail`.`staff`
(
    `id`          int PRIMARY KEY AUTO_INCREMENT,
    `name`        varchar(30)  NOT NULL,
    `password`    varchar(100) NOT NULL,
    `create_time` timestamp    NOT NULL,
    `create_by`   int          NOT NULL,
    `delete_time` int          NOT NULL DEFAULT 0
);

CREATE TABLE `sail`.`staff_group_rel`
(
    `id`               int PRIMARY KEY AUTO_INCREMENT,
    `project_group_id` int NOT NULL,
    `staff_id`         int NOT NULL,
    `role_type`        int NOT NULL COMMENT '权限角色'
);

CREATE TABLE `sail`.`config`
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
    `config_key`       varchar(50) NOT NULL
);

CREATE TABLE `sail`.`config_link`
(
    `id`               int PRIMARY KEY AUTO_INCREMENT,
    `config_id`        int NOT NULL,
    `public_config_id` int NOT NULL
);

CREATE TABLE `sail`.`config_history`
(
    `id`          int PRIMARY KEY AUTO_INCREMENT,
    `config_id`   int       NOT NULL,
    `reversion`   int       NOT NULL,
    `create_time` timestamp NOT NULL,
    `create_by`   int       NOT NULL
);

CREATE TABLE `sail`.`publish_config`
(
    `id`                 int PRIMARY KEY AUTO_INCREMENT,
    `project_id`         int          NOT NULL,
    `namespace_id`       int          NOT NULL,
    `publish_type`       int          NOT NULL COMMENT '发布方式',
    `publish_data`       varchar(20)  NOT NULL COMMENT '发布数据',
    `publish_config_ids` varchar(100) NOT NULL,
    `status`             int          NOT NULL,
    `create_time`        timestamp    NOT NULL Default CURRENT_TIMESTAMP,
    `update_time`        timestamp    NOT NULL Default CURRENT_TIMESTAMP
);
