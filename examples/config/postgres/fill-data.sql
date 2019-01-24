CREATE TABLE address (
  id_address    SERIAL PRIMARY KEY,
  street        TEXT,
  number        INTEGER,
  country       TEXT
);


CREATE TABLE person (
  id_person     SERIAL PRIMARY KEY,
  first_name    TEXT,
  last_name     TEXT,
  age           INTEGER,
  active        BOOLEAN,
  fk_address    INTEGER REFERENCES address (id_address)
);

INSERT INTO address VALUES
(1, 'rua dos testes', 1, 'portugal')
;

INSERT INTO person VALUES
(1, 'a', 'aa', 30, true, 1),
(2, 'b', 'bb', 31, true, 1),
(3, 'c', 'cc', 32, true, 1),
(4, 'd', 'dd', 33, true, 1),
(5, 'e', 'ee', 34, true, 1)
;