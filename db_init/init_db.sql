CREATE TABLE entity(
  id SERIAL not null unique,
  entity_name varchar(255) not null,
  price integer not null,
  is_active bool not null 
);

CREATE TABLE users(
  id SERIAL not null unique,
  username varchar(255) not null unique,
  password_hash varchar(255)  not null,
  is_admin bool not null
);

CREATE TABLE sessions(
    id SERIAL not null unique,
    user_id integer not null ,
    refresh_token varchar(255) not null ,
    signature varchar(255) not null ,
    created_at timestamp not null ,
    disabled bool not null
);

insert into entity (entity_name, price, is_active) values ('name_1',1000,True);
insert into entity (entity_name, price, is_active) values ('name_2',2000,False);
insert into entity (entity_name, price, is_active) values ('name_3',3000,True);
insert into entity (entity_name, price, is_active) values ('name_4',4000,False);
insert into entity (entity_name, price, is_active) values ('name_5',5000,True);

select * from entity;