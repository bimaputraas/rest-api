DROP DATABASE IF EXISTS example_data;

CREATE DATABASE example_data;

USE example_data;

CREATE TABLE users (
                       id            BIGINT(16) NOT NULL PRIMARY KEY AUTO_INCREMENT,
                       first_name    VARCHAR(100),
                       last_name     VARCHAR(100),
                       phone_number  VARCHAR(15) UNIQUE,
                       address       TEXT,
                       pin           VARCHAR(200),
                       created       DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE balances (
                          id              BIGINT(16) NOT NULL PRIMARY KEY AUTO_INCREMENT,
                          user_id         BIGINT(16) NOT NULL,
                          current_balance DECIMAL(15, 2) DEFAULT 0.00,
                          updated         VARCHAR(100) ,
                          FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE top_ups (
                         id             BIGINT(16) NOT NULL PRIMARY KEY AUTO_INCREMENT,
                         user_id        BIGINT(16) NOT NULL,
                         amount_top_up  DECIMAL(15, 2) ,
                         balance_before DECIMAL(15, 2) ,
                         balance_after  DECIMAL(15, 2) ,
                         created        DATETIME,
                         FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE payments (
                          id             BIGINT(16) NOT NULL PRIMARY KEY AUTO_INCREMENT,
                          user_id        BIGINT(16) NOT NULL,
                          amount         DECIMAL(15, 2) ,
                          remarks        VARCHAR(255),
                          balance_before DECIMAL(15, 2) ,
                          balance_after  DECIMAL(15, 2) ,
                          created        DATETIME,
                          FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE transfers (
                           id             BIGINT(16) NOT NULL PRIMARY KEY AUTO_INCREMENT,
                           user_id        BIGINT(16) NOT NULL,
                           amount         DECIMAL(15, 2) ,
                           remarks        VARCHAR(255),
                           balance_before DECIMAL(15, 2) ,
                           balance_after  DECIMAL(15, 2) ,
                           created        DATETIME,
                           FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE transactions (
                              id               BIGINT(16) NOT NULL PRIMARY KEY AUTO_INCREMENT,
                              user_id          BIGINT(16) NOT NULL,
                              top_up_id        BIGINT(16) ,
                              transfer_id      BIGINT(16) ,
                              payment_id       BIGINT(16) ,
                              transaction_type VARCHAR(10) ,
                              amount           DECIMAL(15, 2) ,
                              remarks          VARCHAR(255),
                              balance_before   DECIMAL(15, 2) ,
                              balance_after    DECIMAL(15, 2) ,
                              status           VARCHAR(10) ,
                              created          DATETIME,
                              FOREIGN KEY (user_id) REFERENCES users(id),
                              FOREIGN KEY (top_up_id) REFERENCES top_ups(id),
                              FOREIGN KEY (transfer_id) REFERENCES transfers(id),
                              FOREIGN KEY (payment_id) REFERENCES payments(id)
);
