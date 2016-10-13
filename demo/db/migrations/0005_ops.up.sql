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
