create database patient;
use patient;
create table abnormal(
	abnormal_id       int auto_increment,
	patient_id varchar(255), 
	value double,      
	timestamp  datetime,
	primary key (abnormal_id)
);
create table patient(
	patient_id varchar(255),
	edge_id  varchar(255),
	doctor_id  varchar(255),
	primary key(patient_id)
);
create table doctor(
	doctor_id varchar(255),
	primary key(doctor_id)
);
create table edge_node(
	edge_id varchar(255),
	primary key (edge_id)
);
alter table abnormal
	add constraint foreign key (patient_id) references patient(patient_id);
alter table patient
	add constraint foreign key (doctor_id) references doctor(doctor_id),
	add constraint foreign key (edge_id) references edge_node(edge_id);
insert into edge_node (edge_id) values ("edge");
insert into doctor (doctor_id) values ("doctor_admin");
insert into patient (patient_id, doctor_id, edge_id) values ("admin", "doctor_admin", "edge");