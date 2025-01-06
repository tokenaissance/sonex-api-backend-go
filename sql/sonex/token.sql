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

 Date: 06/01/2025 20:13:22
*/


-- ----------------------------
-- Table structure for token
-- ----------------------------
DROP TABLE IF EXISTS "public"."token";
CREATE TABLE "public"."token" (
  "id" int8 NOT NULL,
  "chain_id" int8 NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "symbal" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "decimals" int8 NOT NULL,
  "address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "logo" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "sort" int8 NOT NULL,
  "hidden" bool NOT NULL DEFAULT false
)
;

-- ----------------------------
-- Records of token
-- ----------------------------
BEGIN;
INSERT INTO "public"."token" ("id", "chain_id", "name", "symbal", "decimals", "address", "logo", "sort", "hidden") VALUES (1, 1946, 'Wrapped Ethereum', 'WETH', 18, '0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C', 'https://img.sonex.so/weth-logo.png', 999, 'f');
COMMIT;

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
