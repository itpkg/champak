
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table forum_articles(
  id serial primary key,
  user_id int not null,
  title varchar(255) not null,
  body text not null,
  type varchar(8) not null default 'markdown',
  vote int not null default 0,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create index idx_forum_articles on forum_articles(title);
create index idx_forum_type on forum_articles(type);

create table forum_tags(
  id serial primary key,
  name varchar(255) not null,
  vote int not null default 0,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create unique index idx_forum_tags on forum_tags(name);

create table forum_articles_tags(
  id serial primary key,
  article_id int not null,
  tag_id int not null
);
create unique index idx_forum_article_tag on forum_articles_tags(article_id, tag_id);

create table forum_comments(
  id serial primary key,
  user_id int not null,
  body text not null,
  type varchar(8) not null default 'markdown',
  vote int not null default 0,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create index idx_forum_comments on forum_comments(type);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table forum_comments;
drop table forum_articles_tags;
drop table forum_tags;
drop table forum_articles;
