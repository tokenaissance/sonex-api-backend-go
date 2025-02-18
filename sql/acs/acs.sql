/*
 Navicat Premium Dump SQL

 Source Server         : postgres
 Source Server Type    : PostgreSQL
 Source Server Version : 160003 (160003)
 Source Host           : 10.0.0.150:5432
 Source Catalog        : sonex-api-backend-go
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 160003 (160003)
 File Encoding         : 65001

 Date: 18/02/2025 20:49:56
*/


-- ----------------------------
-- Table structure for acs_increase_liquidity_points_record
-- ----------------------------
DROP TABLE IF EXISTS "public"."acs_increase_liquidity_points_record";
CREATE TABLE "public"."acs_increase_liquidity_points_record" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "user_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "points" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_index" int8 NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "token_id" int8 NOT NULL,
  "pool" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "token0" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "token1" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "amount0" numeric NOT NULL,
  "amount1" numeric NOT NULL
)
;

-- ----------------------------
-- Table structure for acs_swap_acs_record
-- ----------------------------
DROP TABLE IF EXISTS "public"."acs_swap_acs_record";
CREATE TABLE "public"."acs_swap_acs_record" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "user_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "acs" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_index" int8 NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "pool" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "token0" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "token1" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "amount0" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "amount1" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Indexes structure for table acs_increase_liquidity_points_record
-- ----------------------------
CREATE UNIQUE INDEX "acs_liquidity_points_idx_chain_id_hash_index" ON "public"."acs_increase_liquidity_points_record" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "acs_liquidity_points_idx_tx_at" ON "public"."acs_increase_liquidity_points_record" USING btree (
  "tx_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table acs_increase_liquidity_points_record
-- ----------------------------
ALTER TABLE "public"."acs_increase_liquidity_points_record" ADD CONSTRAINT "acs_points_record_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table acs_swap_acs_record
-- ----------------------------
CREATE UNIQUE INDEX "acs_swap_acs_record_chain_id_hash_index" ON "public"."acs_swap_acs_record" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table acs_swap_acs_record
-- ----------------------------
ALTER TABLE "public"."acs_swap_acs_record" ADD CONSTRAINT "acs_swap_acs_record_pkey" PRIMARY KEY ("id");
