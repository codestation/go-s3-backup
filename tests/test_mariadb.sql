CREATE USER test@'%' IDENTIFIED BY 'secret';

CREATE DATABASE test_foo;
CREATE DATABASE test_bar;
CREATE DATABASE test_baz;

USE test_foo;
CREATE TABLE foo_table1 (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO foo_table1 (name) VALUES ('foo1'), ('foo2'), ('foo3');

USE test_bar;
CREATE TABLE bar_table1 (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO bar_table1 (name) VALUES ('bar1'), ('bar2'), ('bar3');

USE test_baz;

CREATE TABLE baz_table1 (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO baz_table1 (name) VALUES ('baz1'), ('baz2'), ('baz3');
