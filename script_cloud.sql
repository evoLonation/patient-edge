create database patient_cloud;
use patient_cloud;
create table abnormal(
	abnormal_id       int auto_increment,
	patient_id varchar(255),  # 传感器的名称
	value double,      
	timestamp  datetime,
	primary key (abnormal_id)
);
create table patient(
	patient_id varchar(255),
	edge_name  varchar(255),
	doctor_id  varchar(255),
	primary key(patient_id)
);
create table doctor(
	doctor_id varchar(255),
	primary key(doctor_id)
);
