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

 Date: 09/03/2025 17:19:39
*/


-- ----------------------------
-- Table structure for xpool_points_record
-- ----------------------------
DROP TABLE IF EXISTS "public"."xpool_points_record";
CREATE TABLE "public"."xpool_points_record" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "user_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "pool_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "points" int8 NOT NULL,
  "amount" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "token_decimals" int8 NOT NULL,
  "price" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "usd_decimals" int8 NOT NULL
)
;

-- ----------------------------
-- Indexes structure for table xpool_points_record
-- ----------------------------
CREATE UNIQUE INDEX "xpool_points_record_index_chain_id_user_pool_address" ON "public"."xpool_points_record" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "user_address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "pool_address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table xpool_points_record
-- ----------------------------
ALTER TABLE "public"."xpool_points_record" ADD CONSTRAINT "xpool_points_record_pkey" PRIMARY KEY ("id");
