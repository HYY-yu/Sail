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

create
    unique index project_key_uindex
    on project (`key`);

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

create
    unique index namespace_pid_name_uindex
    on namespace (`project_group_id`, `name`);

CREATE TABLE `sail`.`staff`
(
    `id`            int PRIMARY KEY AUTO_INCREMENT,
    `name`          varchar(30)  NOT NULL,
    `password`      varchar(100) NOT NULL,
    `refresh_token` varchar(200) NOT NULL DEFAULT '',
    `create_time`   timestamp    NOT NULL,
    `create_by`     int          NOT NULL
);

create
    unique index staff_name_uindex
    on staff (`name`);

INSERT INTO staff (id, name, password, create_time, create_by) VALUE (1, 'Admin',
                                                                      '$2a$10$9QsXUNwjuYBdSlNA4zX/OucUcVJ/MdyqyOarzE/qdJRyw2qOjhFLS',
                                                                      NOW(), 1);

CREATE TABLE `sail`.`staff_group_rel`
(
    `id`               int PRIMARY KEY AUTO_INCREMENT,
    `project_group_id` int NOT NULL,
    `staff_id`         int NOT NULL,
    `role_type`        int NOT NULL COMMENT '权限角色'
);

INSERT INTO staff_group_rel (project_group_id, staff_id, role_type) VALUE (0, 1, 1);

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
    `config_type`      varchar(10) NOT NULL
);

create
    unique index config_big_key_uindex
    on config (`project_id`, `namespace_id`, `project_group_id`, `name`, `config_type`);

CREATE TABLE `sail`.`config_link`
(
    `id`               int PRIMARY KEY AUTO_INCREMENT,
    `config_id`        int NOT NULL,
    `public_config_id` int NOT NULL
);

create
    unique index config_link_config_id_uindex
    on config_link (`config_id`, `public_config_id`);


CREATE TABLE `sail`.`config_history`
(
    `id`          int PRIMARY KEY AUTO_INCREMENT,
    `config_id`   int       NOT NULL,
    `reversion`   int       NOT NULL,
    `op_type`     int       NOT NULL,
    `create_time` timestamp NOT NULL,
    `create_by`   int       NOT NULL
);


create
    unique index config_history_config_id_uindex
    on config_history (`config_id`, `reversion`);

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
