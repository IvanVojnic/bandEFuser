create table if not exists users(
    id uuid primary key,
    email varchar(255) not null UNIQUE,
    name varchar(255) not null UNIQUE,
    password varchar(255) not null,
    refreshtoken varchar(255)
);

create table if not exists friends(
    id uuid primary key,
    userSender uuid not null,
    userReceiver uuid not null,
    status int not null,
    FOREIGN KEY (userSender) REFERENCES users(id)
      ON DELETE CASCADE,
    FOREIGN KEY (userReceiver) REFERENCES users(id)
      ON DELETE CASCADE
)