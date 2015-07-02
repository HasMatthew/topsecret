create database logs;
use logs;
create table clicks
(
    id varchar(255) NOT NULL PRIMARY KEY,
    advertiser_id int NOT NULL,
    site_id int NOT NULL,
    ip varchar(255) NOT NULL,
    ios_ifa varchar(255) default NULL,
    google_aid varchar(255) default NULL,
    windows_aid varchar(255) default NULL,
    created timestamp default current_timestamp
);
create table installs
(
    id varchar(255) NOT NULL,
    advertiser_id int NOT NULL,
    site_id int NOT NULL,
    ip varchar(255) NOT NULL,
    ios_ifa varchar(255) default NULL,
    google_aid varchar(255) default NULL,
    windows_aid varchar(255) default NULL,
    stat_click_id varchar(255) default NULL,
    created timestamp default current_timestamp
);
create table events
(
    id varchar(255) NOT NULL,
    advertiser_id int NOT NULL,
    site_id int NOT NULL,
    ip varchar(255) NOT NULL,
    ios_ifa varchar(255) default NULL,
    google_aid varchar(255) default NULL,
    windows_aid varchar(255) default NULL,
    stat_install_id varchar(255) default NULL,
    created timestamp default current_timestamp
);
create table persons
(
    first varchar(255) NOT NULL,
    last varchar(255) NOT NULL,
    age int NOT NULL,
    major varchar(255) default NULL,
    constellation varchar(255) default NULL,
    created timestamp default current_timestamp
);