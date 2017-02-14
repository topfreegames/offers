-- offers api
-- https://github.com/topfreegames/offers
--
-- Licensed under the MIT license:
-- http://www.opensource.org/licenses/mit-license
-- Copyright © 2016 Top Free Games <backend@tfgco.com>

REVOKE ALL ON SCHEMA public FROM offers;
DROP DATABASE IF EXISTS offers;

DROP ROLE offers;

CREATE ROLE offers LOGIN
  SUPERUSER INHERIT CREATEDB CREATEROLE;

CREATE DATABASE offers
  WITH OWNER = offers
       ENCODING = 'UTF8'
       TABLESPACE = pg_default
       TEMPLATE = template0;

GRANT ALL ON SCHEMA public TO offers;
