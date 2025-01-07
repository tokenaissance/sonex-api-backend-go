/*
 Navicat Premium Dump SQL

 Source Server         : postgres
 Source Server Type    : PostgreSQL
 Source Server Version : 160003 (160003)
 Source Host           : 10.0.0.150:5432
 Source Catalog        : sonex_api_backend
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 160003 (160003)
 File Encoding         : 65001

 Date: 06/01/2025 21:46:02
*/

CREATE SEQUENCE "public"."token_id"
  INCREMENT 1
  MINVALUE  1
  MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for token
-- ----------------------------
DROP TABLE IF EXISTS "public"."token";
CREATE TABLE "public"."token" (
  "id" int8 NOT NULL DEFAULT nextval('token_id'::regclass),
  "chain_id" int8 NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "symbol" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "decimals" int8 NOT NULL,
  "address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "logo" varchar(255) COLLATE "pg_catalog"."default",
  "sort" int8 NOT NULL DEFAULT 0,
  "hidden" bool NOT NULL DEFAULT false
)
;

-- ----------------------------
-- Indexes structure for table token
-- ----------------------------
CREATE UNIQUE INDEX "token_index_chain_id_address" ON "public"."token" USING btree (
  "address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table token
-- ----------------------------
ALTER TABLE "public"."token" ADD CONSTRAINT "token_pkey" PRIMARY KEY ("id");