
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table settings(
  id serial primary key,
  key varchar(255) not null,
  val bytea not null,
  flag boolean not null default false,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create unique index idx_settings_key on settings(key);

create table locales(
  id serial primary key,
  code varchar(255) not null,
  lang varchar(8) not null default 'en-US',
  message text not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create unique index idx_locales_code_lang on locales(code, lang);
create index idx_locales_code on locales(code);
create index idx_locales_lang on locales(lang);

create table notices(
  id serial primary key,
  lang varchar(8) not null default 'en-US',
  content text not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create index idx_notices_lang on notices(lang);

create table leave_words(
  id serial primary key,
  content text not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);

create table attachments(
  id serial primary key,
  title varchar(255) not null,
  name varchar(48) not null,
  media_type varchar(32) not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create unique index idx_attachments_uid on attachments(name);
create index idx_attachments_title on attachments(title);
create index idx_attachments_media_type on attachments(media_type);

create table pages(
  id serial primary key,
  name varchar(32) not null,
  title varchar(255) not null,
  lang varchar(8) not null default 'en-US',
  href varchar(255) not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create unique index idx_pages_lang_name on pages(lang, name);
create index idx_pages_name on pages(name);
create index idx_pages_lang on pages(lang);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table pages;
drop table attachments;
drop table leave_words;
drop table notices;
drop table locales;
drop table settings;
