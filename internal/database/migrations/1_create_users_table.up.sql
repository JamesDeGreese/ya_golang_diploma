CREATE TABLE users (
   id SERIAL PRIMARY KEY,
   login    VARCHAR(255) UNIQUE,
   password   VARCHAR(255),
   auth_token   VARCHAR(255),
   balance   BIGINT DEFAULT 0
);