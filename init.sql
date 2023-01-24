create table user(
    id integer primary key autoincrement,
    name varchar not null,
    password_hash varchar not null
);

create table role(
    id integer primary key autoincrement,
    name varcharchar not null
);

create table user_roles(
    user_role_id integer not null primary key autoincrement,
    user_id integer not null,
    role_id integer not null,

    foreign key(user_id) references user(id),
    foreign key(role_id) references role(id),
    unique (user_id, role_id)
);

insert into role(id, name)
values
    (0, 'viewer'),
    (1, 'admin');