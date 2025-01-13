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

 Date: 13/01/2025 09:56:42
*/


-- ----------------------------
-- Table structure for event_pre_deposit_claimer_reward_claimed
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_pre_deposit_claimer_reward_claimed";
CREATE TABLE "public"."event_pre_deposit_claimer_reward_claimed" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_index" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "from_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "to_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "value" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_user" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_token" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_nonce" int8 NOT NULL
)
;

-- ----------------------------
-- Table structure for event_pre_deposit_staked
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_pre_deposit_staked";
CREATE TABLE "public"."event_pre_deposit_staked" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_index" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "from_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "to_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "value" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_user" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_token" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_timestamp" int8 NOT NULL
)
;

-- ----------------------------
-- Table structure for event_pre_deposit_withdrawn
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_pre_deposit_withdrawn";
CREATE TABLE "public"."event_pre_deposit_withdrawn" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_index" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "from_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "to_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "value" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_user" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_wastr_amount" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_vastr_amount" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_timestamp" int8 NOT NULL
)
;

-- ----------------------------
-- Indexes structure for table event_pre_deposit_claimer_reward_claimed
-- ----------------------------
CREATE UNIQUE INDEX "event_claimer_reward_claimed_index_chain_id_tx_hash_tx_index" ON "public"."event_pre_deposit_claimer_reward_claimed" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table event_pre_deposit_claimer_reward_claimed
-- ----------------------------
ALTER TABLE "public"."event_pre_deposit_claimer_reward_claimed" ADD CONSTRAINT "event_pre_deposit_staked_copy1_tx_hash_chain_id_key1" UNIQUE ("tx_hash", "chain_id");

-- ----------------------------
-- Primary Key structure for table event_pre_deposit_claimer_reward_claimed
-- ----------------------------
ALTER TABLE "public"."event_pre_deposit_claimer_reward_claimed" ADD CONSTRAINT "event_pre_deposit_staked_copy1_pkey1" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_pre_deposit_staked
-- ----------------------------
CREATE UNIQUE INDEX "event_pre_deposit_staked_index_chain_id_tx_hash_tx_index" ON "public"."event_pre_deposit_staked" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table event_pre_deposit_staked
-- ----------------------------
ALTER TABLE "public"."event_pre_deposit_staked" ADD CONSTRAINT "event_farm_token_staked_copy1_tx_hash_chain_id_key1" UNIQUE ("tx_hash", "chain_id");

-- ----------------------------
-- Primary Key structure for table event_pre_deposit_staked
-- ----------------------------
ALTER TABLE "public"."event_pre_deposit_staked" ADD CONSTRAINT "event_farm_token_staked_copy1_pkey1" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_pre_deposit_withdrawn
-- ----------------------------
CREATE UNIQUE INDEX "event_pre_deposit_withdrawn_index_chain_id_tx_hash_tx_index" ON "public"."event_pre_deposit_withdrawn" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table event_pre_deposit_withdrawn
-- ----------------------------
ALTER TABLE "public"."event_pre_deposit_withdrawn" ADD CONSTRAINT "event_pre_deposit_staked_copy1_tx_hash_chain_id_key" UNIQUE ("tx_hash", "chain_id");

-- ----------------------------
-- Primary Key structure for table event_pre_deposit_withdrawn
-- ----------------------------
ALTER TABLE "public"."event_pre_deposit_withdrawn" ADD CONSTRAINT "event_pre_deposit_staked_copy1_pkey" PRIMARY KEY ("id");
