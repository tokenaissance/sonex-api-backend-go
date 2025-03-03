/*
 Navicat Premium Dump SQL

 Source Server         : sonex-api-backend-go
 Source Server Type    : PostgreSQL
 Source Server Version : 160006 (160006)
 Source Host           : 10.1.84.8:5432
 Source Catalog        : sonex-api-backend-go
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 160006 (160006)
 File Encoding         : 65001

 Date: 03/03/2025 15:32:29
*/


-- ----------------------------
-- Table structure for puffer_increase_liquidity_points_record
-- ----------------------------
DROP TABLE IF EXISTS "public"."puffer_increase_liquidity_points_record";
CREATE TABLE "public"."puffer_increase_liquidity_points_record" (
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
  "amount0" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "amount1" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Indexes structure for table puffer_increase_liquidity_points_record
-- ----------------------------
CREATE UNIQUE INDEX "puffer_liquidity_points_idx_chain_id_hash_index" ON "public"."puffer_increase_liquidity_points_record" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "puffer_liquidity_points_idx_tx_at" ON "public"."puffer_increase_liquidity_points_record" USING btree (
  "tx_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "puffer_liquidity_points_idx_user" ON "public"."puffer_increase_liquidity_points_record" USING btree (
  "user_address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table puffer_increase_liquidity_points_record
-- ----------------------------
ALTER TABLE "public"."puffer_increase_liquidity_points_record" ADD CONSTRAINT "puffer_increase_liquidity_points_record_pkey" PRIMARY KEY ("id");
