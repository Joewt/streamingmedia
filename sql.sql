
CREATE table users (
`id`  int(10) unsigned not null primary key auto_increment,
`login_name` varchar(64) unique key ,
`pwd` text not null
)

CREATE TABLE video_info (
`id` varchar(64) not null primary key ,
`author_id` int(10) unsigned ,
`name` text,
`display_ctime` text,
`create_time` datetime default current_timestamp
)

CREATE TABLE comments (
`id` varchar(64) not null primary key ,
`video_id` varchar(64)  ,
`author_id` int(10) unsigned ,
`content` text,
`time` datetime default current_timestamp
)