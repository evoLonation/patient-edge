create database patient;
use patient;
create table temperature(
	temperature_id       int auto_increment,
	patient_id varchar(255), 
	value double,      
	timestamp  datetime,
	primary key (temperature_id)
);
create table patient(
	patient_id varchar(255),
	primary key(patient_id)
);
insert into patient (patient_id) values ("admin");
alter table temperature 
	add constraint foreign key (patient_id) references patient(patient_id);