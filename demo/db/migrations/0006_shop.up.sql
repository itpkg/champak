create table mail_transport (
  id serial primary key,
  domain varchar(128) not null,
  transport varchar(128) not null default 'virtual:',
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create unique index idx_mail_transport_domain on mail_transport(domain);

create table mail_users (
  id serial primary key,
  email varchar(255) not null,
  password varchar(255) not null,
  name varchar(32),
  uid int not null default 1001,
  gid int not null default 1001,
  home varchar(128),
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create unique index idx_mail_users_email on mail_users(email);
create index idx_mail_users_name on mail_users(name);

create table mail_virtual (
  id serial primary key,
  source varchar(255) not null,
  destination varchar(255) not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create unique index idx_mail_virtual_source_destination on mail_virtual(source, destination);

create view postfix_mailboxes as
  select email, home||'/' as mailbox from mail_users
  union all
  select domain as email, 'dummy' as mailbox from mail_transport;
create view postfix_virtual as
  select email, email as address from mail_users
  union all
  select source, destination from mail_virtual;
