CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS role(
    id int generated always as identity,
    role varchar(100),
    primary key (id)
);

INSERT INTO role (role) values ('user');
INSERT INTO role (role) values ('admin');

CREATE TABLE IF NOT EXISTS "user"(
    id uuid DEFAULT uuid_generate_v4(),
    role_id int DEFAULT 1,
    password varchar(150) not null ,
    login varchar(50) unique not null ,
    primary key (id),
    foreign key (role_id)
        references role (id)
);

CREATE INDEX idx_login ON "user"(login);

CREATE TABLE IF NOT EXISTS session(
    id int generated always as identity,
    token uuid not null,
    user_id uuid,
    expires_at timestamp not null,
    primary key (id),
    foreign key (user_id)
        references "user" (id) on delete cascade
);

