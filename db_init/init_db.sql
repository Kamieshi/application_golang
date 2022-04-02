CREATE TABLE entity(
  id SERIAL not null unique,
  entity_name varchar(255) not null,
  price integer not null,
  is_active bool not null 
);

CREATE TABLE user(
  id SERIAL not null unique,
  entity_name varchar(255) not null,
  price integer not null,
  is_active bool not null 
);



insert into entity (entity_name, price, is_active) values ('name_1',1000,True);
insert into entity (entity_name, price, is_active) values ('name_2',2000,False);
insert into entity (entity_name, price, is_active) values ('name_3',3000,True);
insert into entity (entity_name, price, is_active) values ('name_4',4000,False);
insert into entity (entity_name, price, is_active) values ('name_5',5000,True);

select * from entity;