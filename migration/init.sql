CREATE DATABASE wb;
-- CREATE USER almaz WITH SUPERUSER CREATEDB PASSWORD 'almaz';
GRANT ALL PRIVILEGES ON DATABASE wb TO almaz;
\c wb
CREATE TABLE IF NOT EXISTS models
(
    order_uid text primary key,
    model json  not null
);