CREATE TABLE entity(
  Id SERIAL not null unique,
  EntityName varchar(255) not null,
  Price integer not null,
  IsActive bool not null 
);


insert into entity (name, Price, IsActive) values ('nssssame_1',1000,True);
insert into entity (EntityName, Price, IsActive) values ('nssssame_2',2000,False);
insert into entity (EntityName, Price, IsActive) values ('nssssame_3',3000,True);
insert into entity (EntityName, Price, IsActive) values ('nssssame_4',4000,False);
insert into entity (EntityName, Price, IsActive) values ('nssssame_5',5000,True);

select * from entity;