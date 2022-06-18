CREATE TABLE users (
    id          UUID,
    name        VARCHAR(50),
    email       VARCHAR(50),
    password    VARCHAR(256),
    created_at  TIMESTAMP,
    updated_at  TIMESTAMP
);
