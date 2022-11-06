DROP TABLE IF EXISTS configs CASCADE;
DROP TABLE IF EXISTS pairs CASCADE;
DROP TABLE IF EXISTS relevant_configs CASCADE;

CREATE TABLE configs
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    last_used  TIMESTAMP    NOT NULL DEFAULT NOW(),
    version    BIGINT NOT NULL,
    relevant BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE pairs (
    id         SERIAL PRIMARY KEY,
    config_id  INTEGER REFERENCES configs (id),
    key        VARCHAR(255) NOT NULL,
    value      VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);



