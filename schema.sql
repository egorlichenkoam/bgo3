CREATE TABLE clients
(
    id        BIGSERIAL PRIMARY KEY,
    login     TEXT      NOT NULL UNIQUE,
    password  TEXT      NOT NULL,
    full_name TEXT      NOT NULL,
    passport  TEXT      NOT NULL,
    birthday  DATE      NOT NULL,
    status    TEXT      NOT NULL DEFAULT 'INACTIVE' CHECK (status in ('INACTIVE', 'ACTIVE')),
    created   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cards
(
    id       BIGSERIAL PRIMARY KEY,
    number   TEXT      NOT NULL,
    balance  BIGINT    NOT NULL DEFAULT 0,
    issuer   TEXT      NOT NULL CHECK (issuer in ('VISA', 'MasterCard', 'MIR')),
    holder   TEXT      NOT NULL,
    owner_id BIGINT    NOT NULL REFERENCES clients,
    status   TEXT      NOT NULL DEFAULT 'INACTIVE' CHECK (status in ('INACTIVE', 'ACTIVE')),
    created  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)

CREATE TABLE transactions
(
    id       BIGSERIAL PRIMARY KEY,
    card_id  BIGINT    NOT NULL REFERENCES cards,
    amount   BIGINT    NOT NULL DEFAULT 0,
    tx_type  TEXT      NOT NULL DEFAULT 'FROM' CHECK (tx_type in ('FROM', 'TO')),
    comments TEXT      NOT NULL,
    mcc      TEXT      NOT NULL,
    icon_id  BIGINT    NOT NULL REFERENCES icons,
    status   TEXT      NOT NULL DEFAULT 'PROCESS' CHECK (status in ('PROCESS', 'EXECUTED')),
    created  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)

CREATE TABLE icons
(
    id  BIGSERIAL PRIMARY KEY,
    url TEXT NOT NULL
)