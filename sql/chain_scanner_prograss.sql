/*
 Navicat Premium Dump SQL

 Source Server         : postgres
 Source Server Type    : PostgreSQL
 Source Server Version : 160003 (160003)
 Source Host           : 10.0.0.150:5432
 Source Catalog        : points_service
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 160003 (160003)
 File Encoding         : 65001

 Date: 16/12/2024 18:21:18
*/


-- ----------------------------
-- Table structure for chain_scanner_prograss
-- ----------------------------
DROP TABLE IF EXISTS "public"."chain_scanner_prograss";
CREATE TABLE "public"."chain_scanner_prograss" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "height" int8 NOT NULL
)
;
ALTER TABLE "public"."chain_scanner_prograss" OWNER TO "postgres";

-- ----------------------------
-- Records of chain_scanner_prograss
-- ----------------------------
BEGIN;
INSERT INTO "public"."chain_scanner_prograss" ("id", "created_at", "updated_at", "chain_id", "height") VALUES (123418173985783808, '2024-12-07 08:43:50+00', '2024-12-16 09:21:57+00', 1946, 5573289);
COMMIT;

-- ----------------------------
-- Indexes structure for table chain_scanner_prograss
-- ----------------------------
CREATE UNIQUE INDEX "index_chain_scanner_prograss_chain_id" ON "public"."chain_scanner_prograss" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table chain_scanner_prograss
-- ----------------------------
ALTER TABLE "public"."chain_scanner_prograss" ADD CONSTRAINT "chain_scanner_prograss_pkey" PRIMARY KEY ("id");
