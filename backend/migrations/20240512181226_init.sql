-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION btree_gin;

CREATE TABLE projects
(
    id   BIGINT GENERATED ALWAYS AS IDENTITY,
    name VARCHAR NOT NULL
        CONSTRAINT check__projects__name CHECK ( name <> '' ),

    CONSTRAINT pk__projects PRIMARY KEY (id)
);

CREATE TYPE PLATFORM AS ENUM (
    'WEB',
    'IOS',
    'ANDROID',
    'OTHER'
    );

CREATE TABLE key_tags
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY,
    project_id BIGINT  NOT NULL,
    value      VARCHAR NOT NULL
        CONSTRAINT check__key_tags__name CHECK ( value <> '' ),

    CONSTRAINT pk__key_tags PRIMARY KEY (id)
);

ALTER TABLE key_tags
    ADD CONSTRAINT fk__key_tags__projects
        FOREIGN KEY (project_id)
            REFERENCES projects (id);

CREATE UNIQUE INDEX idx__key_tags__value ON key_tags (project_id, LOWER(value));

CREATE TABLE keys
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY,
    project_id BIGINT     NOT NULL,
    name       VARCHAR    NOT NULL
        CONSTRAINT check__keys__name CHECK ( name <> '' ),
    platforms  PLATFORM[] NOT NULL
        CONSTRAINT check__keys__platforms CHECK ( array_length(platforms, 1) > 0 ),
    tags       BIGINT[]   NOT NULL,
    CONSTRAINT pk__keys PRIMARY KEY (id)
);

ALTER TABLE keys
    ADD CONSTRAINT fk__keys__projects
        FOREIGN KEY (project_id)
            REFERENCES projects (id);

CREATE UNIQUE INDEX uidx__keys__name ON keys (project_id, name);

-- CREATE INDEX idx__keys__tags ON keys USING GIN (tags);
-- CREATE INDEX idx__keys__name ON keys USING GIN (project_id, (to_tsvector('english', name)));


-- +goose StatementEnd
