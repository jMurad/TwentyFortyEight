CREATE TABLE IF NOT EXISTS players (
    id bigserial NOT NULL PRIMARY KEY,
    name character varying NOT NULL,
    score smallint NOT NULL,
    level smallint NOT NULL,
    sizebox character varying NOT NULL
);