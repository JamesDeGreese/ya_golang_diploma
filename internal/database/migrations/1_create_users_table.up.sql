CREATE TABLE users (
   id SERIAL PRIMARY KEY,
   login    VARCHAR(255),
   password   VARCHAR(255),
   auth_token   VARCHAR(255)
);