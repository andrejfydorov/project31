use gousers;

insert into users (id, name, age) 
values (1, 'andre', 24),
(2, 'vova', 34),
(3, 'anton', 28);

insert into friends (userid, friendid) 
values (1, 2),
(1, 3),
(2, 3),
(3, 1);