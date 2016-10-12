
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table reading_books(
  id serial primary key,
  author varchar(36) not null,
  publisher varchar(36) not null,
  title varchar(255) not null,
  type varchar(8) not null default 'epub3',
  lang varchar(8) not null default 'en-US',
  file varchar(255) not null,
  home varchar(255) not null,
  vote int not null default 0,
  published_at date not null default current_date,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create unique index idx_reading_books_file on reading_books(file);
create index idx_reading_books_author on reading_books(author);
create index idx_reading_books_publisher on reading_books(publisher);
create index idx_reading_books_type on reading_books(type);
create index idx_reading_books_lang on reading_books(lang);

create table reading_notes(
  id serial primary key,
  user_id int not null,
  book_id int not null,
  body text not null,
  type varchar(8) not null default 'markdown',
  vote int not null default 0,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create index idx_reading_notes_type on reading_notes(type);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table reading_notes;
drop table reading_books;
