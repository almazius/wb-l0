CREATE TABLE IF NOT EXISTS models
(
    order_uid text primary key,
    model json  not null
);