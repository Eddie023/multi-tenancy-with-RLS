/**
Create at least one role that has CREATE USER and CREATE DATABASE permissions
but is not a superuser such that you can still do administrative tasks but not bypass
other security checks.
**/
CREATE ROLE dev WITH LOGIN PASSWORD 'pass';

CREATE SCHEMA app AUTHORIZATION dev;

/**
Make the default schema as app when logged via dev user. 
**/
ALTER ROLE dev SET search_path TO app, public; 
