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
    domain_id BIGINT,
    list_id   TEXT,
    ranking   INT NOT NULL,
    PRIMARY KEY (domain_id, list_id)
);
