create table if not exists users
(
    id           uuid PRIMARY KEY,
    email        VARCHAR(255) NOT NULL UNIQUE,
    name         VARCHAR(255) NOT NULL UNIQUE,
    password     VARCHAR(255) NOT NULL,
    refreshtoken VARCHAR(255)
);

create table if not exists friends
(
    id           uuid PRIMARY KEY,
    userSender   uuid NOT NULL,
    userReceiver uuid NOT NULL,
    status       INT  NOT NULL,
    FOREIGN KEY (userSender) REFERENCES users (id)
        ON DELETE CASCADE,
    FOREIGN KEY (userReceiver) REFERENCES users (id)
        ON DELETE CASCADE
)