CREATE TABLE tranco_lists
(
    id         TEXT PRIMARY KEY,
    created_on TIMESTAMP NOT NULL
);

CREATE TABLE tranco_domains
(
    id     SERIAL PRIMARY KEY,
    domain TEXT UNIQUE NOT NULL
);

CREATE TABLE tranco_rankings
(
    id        SERIAL PRIMARY KEY,
    domain_id BIGINT NOT NULL,
    list_id   TEXT   NOT NULL,
    ranking   INT    NOT NULL,
    UNIQUE (domain_id, list_id)
);
