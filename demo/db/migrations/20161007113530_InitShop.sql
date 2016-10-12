
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table shop_products(
  id serial primary key,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);

create table shop_vendors(
  id serial primary key,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);

create table shop_prices(
  id serial primary key,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);

create table shop_fields(
  id serial primary key,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);

create table shop_deliverers(
  id serial primary key,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);

create table shop_payments(
  id serial primary key,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);

create table shop_orders(
  id serial primary key,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);

create table shop_bills(
  id serial primary key,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);

create table shop_returns(
  id serial primary key,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table shop_returns;
drop table shop_bills;
drop table shop_orders;
drop table shop_payments;
drop table shop_deliverers;
drop table shop_fields;
drop table shop_prices;
drop table shop_vendors;
drop table shop_products;
