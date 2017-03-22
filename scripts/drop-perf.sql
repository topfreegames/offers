-- offers api
-- https://github.com/topfreegames/offers
--
-- Licensed under the MIT license:
-- http://www.opensource.org/licenses/mit-license
-- Copyright © 2017 Top Free Games <backend@tfgco.com>

REVOKE ALL ON SCHEMA public FROM offers_perf;
DROP DATABASE IF EXISTS offers_perf;

DROP ROLE offers_perf;

CREATE ROLE offers_perf LOGIN
  SUPERUSER INHERIT CREATEDB CREATEROLE;

CREATE DATABASE offers_perf
  WITH OWNER = offers_perf
       ENCODING = 'UTF8'
       TABLESPACE = pg_default
       TEMPLATE = template0;

GRANT ALL ON SCHEMA public TO offers_perf;

