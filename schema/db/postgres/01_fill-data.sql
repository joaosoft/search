-- migrate up

INSERT INTO search.address VALUES
(1, 'rua dos testes', 1, 'portugal')
;

INSERT INTO search.person VALUES
(1, 'a', 'aa', 30, true, 1),
(2, 'b', 'bb', 31, true, 1),
(3, 'c', 'cc', 32, true, 1),
(4, 'd', 'dd', 33, true, 1),
(5, 'e', 'ee', 34, true, 1)
;


-- migrate down
DELETE FROM search.address;
DELETE FROM search.person;