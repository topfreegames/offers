-- offers api
-- https://github.com/topfreegames/offers
--
-- Licensed under the MIT license:
-- http://www.opensource.org/licenses/mit-license
-- Copyright © 2016 Top Free Games <backend@tfgco.com>

REVOKE ALL ON SCHEMA public FROM offers_test;
DROP DATABASE IF EXISTS offers_test;

DROP ROLE offers_test;

CREATE ROLE offers_test LOGIN
  SUPERUSER INHERIT CREATEDB CREATEROLE;

CREATE DATABASE offers_test
  WITH OWNER = offers_test
       ENCODING = 'UTF8'
       TABLESPACE = pg_default
       TEMPLATE = template0;

GRANT ALL ON SCHEMA public TO offers_test;
