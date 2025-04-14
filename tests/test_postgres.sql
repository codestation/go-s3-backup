CREATE USER test WITH PASSWORD 'secret';

CREATE DATABASE test_schemas OWNER test;

\c test_schemas

CREATE SCHEMA schema1 AUTHORIZATION test;
CREATE SCHEMA schema2 AUTHORIZATION test;
CREATE SCHEMA schema3 AUTHORIZATION test;

CREATE TABLE schema1.table1 (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

ALTER TABLE schema1.table1 OWNER TO test;

INSERT INTO schema1.table1 (name) VALUES ('schema1_table1_row1'), ('schema1_table1_row2');

CREATE TABLE schema2.table1 (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

ALTER TABLE schema2.table1 OWNER TO test;

INSERT INTO schema2.table1 (name) VALUES ('schema2_table1_row1'), ('schema2_table1_row2');

CREATE TABLE schema3.table1 (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

ALTER TABLE schema3.table1 OWNER TO test;

INSERT INTO schema3.table1 (name) VALUES ('schema3_table1_row1'), ('schema3_table1_row2');

CREATE USER foo WITH PASSWORD 'secret';
CREATE USER bar WITH PASSWORD 'secret';
CREATE USER baz WITH PASSWORD 'secret';

CREATE DATABASE test_foo1 OWNER foo;
CREATE DATABASE test_foo2 OWNER foo;
CREATE DATABASE test_bar1 OWNER bar;
CREATE DATABASE test_bar2 OWNER bar;
CREATE DATABASE test_demo OWNER bar;
CREATE DATABASE test_baz1 OWNER baz;
CREATE DATABASE test_baz2 OWNER baz;

\c test_foo1
CREATE TABLE foo_table1 (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

ALTER TABLE foo_table1 OWNER TO foo;

INSERT INTO foo_table1 (name) VALUES ('foo1'), ('foo2'), ('foo3');

\c test_foo2
CREATE TABLE foo_table1 (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

ALTER TABLE foo_table1 OWNER TO foo;

INSERT INTO foo_table1 (name) VALUES ('foo1'), ('foo2'), ('foo3');

\c test_bar1
CREATE TABLE bar_table1 (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

ALTER TABLE bar_table1 OWNER TO bar;

INSERT INTO bar_table1 (name) VALUES ('bar1'), ('bar2'), ('bar3');

\c test_bar2
CREATE TABLE bar_table1 (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

ALTER TABLE bar_table1 OWNER TO bar;

INSERT INTO bar_table1 (name) VALUES ('bar1'), ('bar2'), ('bar3');

\c test_demo
CREATE TABLE bar_table1 (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

ALTER TABLE bar_table1 OWNER TO bar;

INSERT INTO bar_table1 (name) VALUES ('bar1'), ('bar2'), ('bar3');

\c test_baz1

CREATE TABLE baz_table1 (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

ALTER TABLE baz_table1 OWNER TO baz;

INSERT INTO baz_table1 (name) VALUES ('baz1'), ('baz2'), ('baz3');

\c test_baz2

CREATE TABLE baz_table1 (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

ALTER TABLE baz_table1 OWNER TO baz;

INSERT INTO baz_table1 (name) VALUES ('baz1'), ('baz2'), ('baz3');
