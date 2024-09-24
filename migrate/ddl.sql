DROP DATABASE IF EXISTS example_data;

CREATE DATABASE example_data;

USE example_data;

CREATE TABLE users (
    id            BINARY(16) NOT NULL PRIMARY KEY,
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL,
    phone_number  VARCHAR(15) NOT NULL UNIQUE,
    address       TEXT NOT NULL,
    pin           VARCHAR(6) NOT NULL,
    created       DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE balances (
    id              BINARY(16) NOT NULL PRIMARY KEY,
    user_id         BINARY(16) NOT NULL,
    current_balance DECIMAL(15, 2) DEFAULT 0.00,
    updated         DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE top_ups (
    id             BINARY(16) NOT NULL PRIMARY KEY,
    user_id        BINARY(16) NOT NULL,
    amount_top_up  DECIMAL(15, 2) ,
    balance_before DECIMAL(15, 2) ,
    balance_after  DECIMAL(15, 2) ,
    created        DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE payments (
    id             BINARY(16) NOT NULL PRIMARY KEY,
    user_id        BINARY(16) NOT NULL,
    amount         DECIMAL(15, 2) ,
    remarks        VARCHAR(255),
    balance_before DECIMAL(15, 2) ,
    balance_after  DECIMAL(15, 2) ,
    created        DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE transfers (
    id             BINARY(16) NOT NULL PRIMARY KEY,
    user_id        BINARY(16) NOT NULL,
    amount         DECIMAL(15, 2) ,
    remarks        VARCHAR(255),
    balance_before DECIMAL(15, 2) ,
    balance_after  DECIMAL(15, 2) ,
    created        DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE transactions (
    id               BINARY(16) NOT NULL PRIMARY KEY,
    user_id          BINARY(16) NOT NULL,
    top_up_id        BINARY(16),
    transfer_id      BINARY(16),
    payment_id       BINARY(16),
    transaction_type VARCHAR(10) ,
    amount           DECIMAL(15, 2) ,
    remarks          VARCHAR(255),
    balance_before   DECIMAL(15, 2) ,
    balance_after    DECIMAL(15, 2) ,
    status           VARCHAR(10) ,
    created          DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (top_up_id) REFERENCES top_ups(id),
    FOREIGN KEY (transfer_id) REFERENCES transfers(id),
    FOREIGN KEY (payment_id) REFERENCES payments(id)
);
