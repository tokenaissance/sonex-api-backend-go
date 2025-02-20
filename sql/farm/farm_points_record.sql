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

 Date: 17/02/2025 17:19:39
*/


-- ----------------------------
-- Table structure for farm_points_record
-- ----------------------------
DROP TABLE IF EXISTS "public"."farm_points_record";
CREATE TABLE "public"."farm_points_record" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "user_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "points" int8 NOT NULL,
  "event_staked_id" int8 NOT NULL,
  "incentive_id" int8 NOT NULL,
  "amount0" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "amount1" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "price0" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "price1" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "usd_decimals" int8 NOT NULL,
  "staked_at" timestamptz(6) NOT NULL
)
;

-- ----------------------------
-- Indexes structure for table farm_points_record
-- ----------------------------
CREATE UNIQUE INDEX "farm_points_record_index_chain_id_user_incentive_id" ON "public"."farm_points_record" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "user_address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "incentive_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table farm_points_record
-- ----------------------------
ALTER TABLE "public"."farm_points_record" ADD CONSTRAINT "farm_points_record_pkey" PRIMARY KEY ("id");
