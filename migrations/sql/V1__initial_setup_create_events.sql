create table if not exists users(
                                    id uuid primary key not null,
                                    email varchar(255) not null,
                                    name varchar(255) not null,
                                    password varchar(255) not null,
                                    refreshtoken varchar(255)
);

create table if not exists statuses(
                                       id int primary key not null,
                                       status varchar(255) not null
);

create table if not exists friends(
                                      id uuid primary key not null,
                                      userSender uuid not null,
                                      userReceiver uuid not null,
                                      status_id int not null,
                                      FOREIGN KEY (userSender) REFERENCES users(id)
                                          ON DELETE CASCADE,
                                      FOREIGN KEY (userReceiver) REFERENCES users(id)
                                          ON DELETE CASCADE,
                                      FOREIGN KEY (status_id) REFERENCES statuses(id)
                                          ON DELETE CASCADE
)