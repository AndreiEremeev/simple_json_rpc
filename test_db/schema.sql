create extension if not exists "uuid-ossp";

create table users(
    id uuid primary key,
    login text unique,
    created_at timestamp
);
