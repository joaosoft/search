-- migrate up

CREATE SCHEMA IF NOT EXISTS search;

CREATE TABLE search.address (
  id_address    SERIAL PRIMARY KEY,
  street        TEXT,
  number        INTEGER,
  country       TEXT
);


CREATE TABLE search.person (
  id_person     SERIAL PRIMARY KEY,
  first_name    TEXT,
  last_name     TEXT,
  age           INTEGER,
  active        BOOLEAN,
  fk_address    INTEGER REFERENCES address (id_address)
);


-- migrate down
DROP TABLE search.address;
DROP TABLE search.person;