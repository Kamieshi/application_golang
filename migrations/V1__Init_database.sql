CREATE TABLE entity
(
    id          uuid         not null unique,
    entity_name varchar(255) not null,
    price       integer      not null,
    is_active   bool         not null,
    PRIMARY KEY (id)
);

CREATE TABLE users
(
    id            uuid         not null unique,
    username      varchar(255) not null unique,
    password_hash varchar(255) not null,
    is_admin      bool         not null,
    PRIMARY KEY (id)
);

CREATE TABLE sessions
(
    id            uuid         not null unique,
    user_id       uuid         not null,
    refresh_token varchar(255) not null,
    signature     varchar(255) not null,
    created_at    timestamp    not null,
    disabled      bool         not null,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE images
(
    id        uuid         not null unique,
    file_name varchar(400) not null,
    root_path varchar(500) not null,
    easy_link varchar(200) not null
);
