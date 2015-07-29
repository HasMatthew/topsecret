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
    date_time datetime default NULL
);
