/*
 Navicat Premium Dump SQL

 Source Server         : sonex-api-backend-go
 Source Server Type    : PostgreSQL
 Source Server Version : 160004 (160004)
 Source Host           : 10.1.84.8:5432
 Source Catalog        : sonex-api-backend-go
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 160004 (160004)
 File Encoding         : 65001

 Date: 07/02/2025 23:26:14
*/


-- ----------------------------
-- Table structure for farm_points_record
-- ----------------------------
DROP TABLE IF EXISTS "public"."farm_points_record";
CREATE TABLE "public"."farm_points_record" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chian_id" int8 NOT NULL,
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
ALTER TABLE "public"."farm_points_record" OWNER TO "apibackendgo";

-- ----------------------------
-- Indexes structure for table farm_points_record
-- ----------------------------
CREATE UNIQUE INDEX "farm_points_record_index_chain_id_user_incentive_id" ON "public"."farm_points_record" USING btree (
  "chian_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "user_address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "incentive_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table farm_points_record
-- ----------------------------
ALTER TABLE "public"."farm_points_record" ADD CONSTRAINT "farm_points_record_pkey" PRIMARY KEY ("id");
