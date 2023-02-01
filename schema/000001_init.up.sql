CREATE TABLE users
(
    id              serial       not null unique,
    name            varchar(255) not null,
    username        varchar(255) not null unique,
    password_hash   varchar(255) not null,
    public_api_key  varchar(255) not null,
    private_api_key varchar(255) not null
);

CREATE TABLE orders
(
    order_id              varchar(255) not null unique,
    user_id               integer      not null,
    cli_order_id          varchar(255) not null,
    type                  varchar(255) not null,
    symbol                varchar(255) not null,
    quantity              float8       not null,
    side                  varchar(255) not null,
    filled                float8       not null,
    timestamp             varchar(255) not null,
    last_update_timestamp varchar(255) not null,
    price                 float8
);

CREATE TABLE users_orders
(
    id       serial                                                      not null unique,
    user_id  int references users (id) on delete cascade                 not null,
    order_id varchar(255) references orders (order_id) on delete cascade not null
)
