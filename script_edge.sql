create database patient_edge;
use patient_edge;
create table temperature(
	temperature_id       int auto_increment,
	patient_id varchar(255),  # 传感器的名称
	value double,      
	timestamp  datetime,
	primary key (temperature_id)
);
create table patient(
	patient_id varchar(255),
	primary key(patient_id)
);
insert into patient (patient_id) values ("sensor-temperature-01");
