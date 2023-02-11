use gousers;
create table users (
    id int primary key auto_increment,
    name varchar(30) not null,
    age int not null
);

create table friends (
    userid int,
    friendid int,
	primary key(userid, friendid),
    foreign key(userid) references users(id) on delete cascade on update cascade,
    foreign key(friendid) references users(id) on delete cascade on update cascade
);