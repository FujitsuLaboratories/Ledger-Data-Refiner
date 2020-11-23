CREATE USER :user WITH PASSWORD :passwd;
DROP DATABASE IF EXISTS :dbname;
CREATE DATABASE :dbname owner :user;
\c :dbname;