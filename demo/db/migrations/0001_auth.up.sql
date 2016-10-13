create table users(
  id serial primary key,
  name varchar(32) not null,
  email varchar(255) not null,
  uid varchar(36) not null,
  password varchar(255),
  provider_id varchar(255) not null,
  provider_type varchar(32) not null,
  sign_in_count int not null default 0,
  current_sign_in_at timestamp without time zone,
  current_sign_in_ip inet,
  last_sign_in_at timestamp without time zone,
  last_sign_in_ip inet,
  confirmed_at timestamp without time zone,
  locked_at timestamp without time zone,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create unique index idx_users_uid on users(uid);
create unique index idx_users_email on users(email);
create unique index idx_users_provider_id_type on users(provider_id, provider_type);
create index idx_users_name on users(name);
create index idx_users_provider_id on users(provider_id);
create index idx_users_provider_type on users(provider_type);

create table logs(
  id serial primary key,
  user_id int not null,
  type varchar(8) not null default 'info',
  message varchar(255) not null,
  created_at timestamp without time zone not null default now()
);
create index idx_logs_type on logs(type);

create table roles(
  id serial primary key,
  name varchar(32) not null,
  resource_id int not null default 0,
  resource_type varchar(255) not null default '-',
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create unique index idx_roles_name_resource_type_id on roles(name, resource_type, resource_id);
create index idx_roles_name on roles(name);
create index idx_roles_resource_type on roles(resource_type);

create table policies(
  id serial primary key,
  user_id int not null,
  role_id int not null,
  start_up date not null default current_date,
  shut_down date not null default '2016-10-07',
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null
);
create unique index idx_policies_user_role on policies(user_id, role_id);
