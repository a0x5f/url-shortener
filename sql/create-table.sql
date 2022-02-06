CREATE DATABASE "url_shortener";

\c url_shortener;

CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE
);