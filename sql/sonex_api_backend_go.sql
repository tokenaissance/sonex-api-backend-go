/*
 Navicat Premium Data Transfer

 Source Server         : main_apibackendgo
 Source Server Type    : PostgreSQL
 Source Server Version : 160006 (160006)
 Source Host           : 10.1.84.16:5432
 Source Catalog        : sonex-api-backend-go
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 160006 (160006)
 File Encoding         : 65001

 Date: 16/01/2025 00:20:48
*/


-- ----------------------------
-- Sequence structure for sonex_token_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "sonex_token_id_seq";
CREATE SEQUENCE "sonex_token_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sonex_token_price_record_id
-- ----------------------------
DROP SEQUENCE IF EXISTS "sonex_token_price_record_id";
CREATE SEQUENCE "sonex_token_price_record_id" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for token_id
-- ----------------------------
DROP SEQUENCE IF EXISTS "token_id";
CREATE SEQUENCE "token_id" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for chain_scanner_prograss
-- ----------------------------
DROP TABLE IF EXISTS "chain_scanner_prograss";
CREATE TABLE "chain_scanner_prograss" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "height" int8 NOT NULL
)
;

-- ----------------------------
-- Table structure for event_pre_deposit_claimer_reward_claimed
-- ----------------------------
DROP TABLE IF EXISTS "event_pre_deposit_claimer_reward_claimed";
CREATE TABLE "event_pre_deposit_claimer_reward_claimed" (
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
DROP TABLE IF EXISTS "event_pre_deposit_staked";
CREATE TABLE "event_pre_deposit_staked" (
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
DROP TABLE IF EXISTS "event_pre_deposit_withdrawn";
CREATE TABLE "event_pre_deposit_withdrawn" (
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
-- Table structure for event_sonex_factory_pool_created
-- ----------------------------
DROP TABLE IF EXISTS "event_sonex_factory_pool_created";
CREATE TABLE "event_sonex_factory_pool_created" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_index" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "from_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "to_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "value" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_token0" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_token1" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_fee" int8 NOT NULL,
  "event_tick_spacing" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_pool" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_sonex_pool_burn
-- ----------------------------
DROP TABLE IF EXISTS "event_sonex_pool_burn";
CREATE TABLE "event_sonex_pool_burn" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_index" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "from_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "to_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "value" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_owner" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_tick_lower" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_tick_upper" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount0" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount1" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "extra_pool" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_sonex_pool_mint
-- ----------------------------
DROP TABLE IF EXISTS "event_sonex_pool_mint";
CREATE TABLE "event_sonex_pool_mint" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_index" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "from_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "to_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "value" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_sender" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_owner" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_tick_lower" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_tick_upper" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount0" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount1" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "extra_pool" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_sonex_pool_swap
-- ----------------------------
DROP TABLE IF EXISTS "event_sonex_pool_swap";
CREATE TABLE "event_sonex_pool_swap" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_index" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "from_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "to_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "value" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_sender" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_recipient" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount0" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount1" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_sqrt_price_x96" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_tick" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "extra_pool" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_sonex_position_manager_collect
-- ----------------------------
DROP TABLE IF EXISTS "event_sonex_position_manager_collect";
CREATE TABLE "event_sonex_position_manager_collect" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_index" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "from_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "to_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "value" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_token_id" int8 NOT NULL,
  "event_recipient" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount0" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount1" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_sonex_position_manager_decrease_liquidity
-- ----------------------------
DROP TABLE IF EXISTS "event_sonex_position_manager_decrease_liquidity";
CREATE TABLE "event_sonex_position_manager_decrease_liquidity" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_index" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "from_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "to_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "value" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_token_id" int8 NOT NULL,
  "event_liquidity" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount0" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount1" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_sonex_position_manager_increase_liquidity
-- ----------------------------
DROP TABLE IF EXISTS "event_sonex_position_manager_increase_liquidity";
CREATE TABLE "event_sonex_position_manager_increase_liquidity" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_index" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "from_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "to_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "value" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_token_id" int8 NOT NULL,
  "event_liquidity" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount0" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount1" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_sonex_position_manager_transfer
-- ----------------------------
DROP TABLE IF EXISTS "event_sonex_position_manager_transfer";
CREATE TABLE "event_sonex_position_manager_transfer" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_index" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "from_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "to_" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "value" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_from" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_to" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_token_id" int8 NOT NULL
)
;

-- ----------------------------
-- Table structure for pre_deposit_points_record
-- ----------------------------
DROP TABLE IF EXISTS "pre_deposit_points_record";
CREATE TABLE "pre_deposit_points_record" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "user_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "points" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_index" int8 NOT NULL,
  "tx_at" timestamptz(6) NOT NULL
)
;

-- ----------------------------
-- Table structure for pre_deposit_reward_claim_record
-- ----------------------------
DROP TABLE IF EXISTS "pre_deposit_reward_claim_record";
CREATE TABLE "pre_deposit_reward_claim_record" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "user_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "token" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "amount" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "nonce" int8 NOT NULL,
  "signature" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default",
  "tx_index" int8,
  "tx_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for pre_deposit_reward_snapshot
-- ----------------------------
DROP TABLE IF EXISTS "pre_deposit_reward_snapshot";
CREATE TABLE "pre_deposit_reward_snapshot" (
  "id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "user_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "wastr_amount" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "vastr_amount" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "reward_token" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "reward_amount" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for sonex_token
-- ----------------------------
DROP TABLE IF EXISTS "sonex_token";
CREATE TABLE "sonex_token" (
  "id" int8 NOT NULL DEFAULT nextval('sonex_token_id_seq'::regclass),
  "created_at" timestamptz(6) NOT NULL DEFAULT ('now'::text)::timestamp(6) with time zone,
  "updated_at" timestamptz(6) NOT NULL DEFAULT ('now'::text)::timestamp(6) with time zone,
  "deleted_at" timestamp(6),
  "name" varchar(120) COLLATE "pg_catalog"."default" NOT NULL,
  "symbol" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "address" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
  "decimals" int4 NOT NULL DEFAULT 18,
  "logo" text COLLATE "pg_catalog"."default",
  "sort" int4 NOT NULL DEFAULT 0,
  "recommendOrder" int4 NOT NULL DEFAULT 0,
  "isHidden" bool NOT NULL DEFAULT false,
  "volume" numeric(32,18) NOT NULL DEFAULT '0'::numeric,
  "additional" jsonb DEFAULT '{}'::jsonb
)
;

-- ----------------------------
-- Table structure for sonex_token_price_record
-- ----------------------------
DROP TABLE IF EXISTS "sonex_token_price_record";
CREATE TABLE "sonex_token_price_record" (
  "id" int8 NOT NULL DEFAULT nextval('sonex_token_price_record_id'::regclass),
  "token_symbol" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "token_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "price" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamp(6) NOT NULL,
  "usd_decimals" int8 NOT NULL
)
;

-- ----------------------------
-- Table structure for token
-- ----------------------------
DROP TABLE IF EXISTS "token";
CREATE TABLE "token" (
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
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"sonex_token_id_seq"', 1, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"sonex_token_price_record_id"', 1, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"token_id"', 1, true);

-- ----------------------------
-- Indexes structure for table chain_scanner_prograss
-- ----------------------------
CREATE UNIQUE INDEX "chain_scanner_prograss_index_chain_id" ON "chain_scanner_prograss" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table chain_scanner_prograss
-- ----------------------------
ALTER TABLE "chain_scanner_prograss" ADD CONSTRAINT "chain_scanner_prograss_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_pre_deposit_claimer_reward_claimed
-- ----------------------------
CREATE UNIQUE INDEX "event_claimer_reward_claimed_index_chain_id_tx_hash_tx_index" ON "event_pre_deposit_claimer_reward_claimed" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_pre_deposit_claimer_reward_claimed
-- ----------------------------
ALTER TABLE "event_pre_deposit_claimer_reward_claimed" ADD CONSTRAINT "event_pre_deposit_claimer_reward_claimed_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_pre_deposit_staked
-- ----------------------------
CREATE UNIQUE INDEX "event_pre_deposit_staked_index_chain_id_tx_hash_tx_index" ON "event_pre_deposit_staked" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_pre_deposit_staked
-- ----------------------------
ALTER TABLE "event_pre_deposit_staked" ADD CONSTRAINT "event_pre_deposit_staked_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_pre_deposit_withdrawn
-- ----------------------------
CREATE UNIQUE INDEX "event_pre_deposit_withdrawn_index_chain_id_tx_hash_tx_index" ON "event_pre_deposit_withdrawn" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_pre_deposit_withdrawn
-- ----------------------------
ALTER TABLE "event_pre_deposit_withdrawn" ADD CONSTRAINT "event_pre_deposit_withdrawn_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_factory_pool_created
-- ----------------------------
CREATE INDEX "sonex_factory_pool_created_index_chain_id_block_number" ON "event_sonex_factory_pool_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "sonex_factory_pool_created_index_chain_id_tx_hash_tx_index" ON "event_sonex_factory_pool_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "sonex_factory_pool_created_index_pool" ON "event_sonex_factory_pool_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_pool" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_factory_pool_created
-- ----------------------------
ALTER TABLE "event_sonex_factory_pool_created" ADD CONSTRAINT "event_sonex_pool_created_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_pool_burn
-- ----------------------------
CREATE INDEX "pool_burn_index_chain_id_block_number" ON "event_sonex_pool_burn" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "pool_burn_index_chain_id_owner" ON "event_sonex_pool_burn" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_owner" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "pool_burn_index_chain_id_pool" ON "event_sonex_pool_burn" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "extra_pool" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "pool_burn_index_chain_id_tx_hash_tx_index" ON "event_sonex_pool_burn" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_pool_burn
-- ----------------------------
ALTER TABLE "event_sonex_pool_burn" ADD CONSTRAINT "event_sonex_pool_burn_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_pool_mint
-- ----------------------------
CREATE INDEX "pool_mint_index_chain_block_number" ON "event_sonex_pool_mint" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "pool_mint_index_chain_id_owner" ON "event_sonex_pool_mint" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_owner" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "pool_mint_index_chain_id_pool" ON "event_sonex_pool_mint" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "extra_pool" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "pool_mint_index_chain_id_sender" ON "event_sonex_pool_mint" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_sender" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "pool_mint_index_chain_id_tx_hash_tx_index" ON "event_sonex_pool_mint" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_pool_mint
-- ----------------------------
ALTER TABLE "event_sonex_pool_mint" ADD CONSTRAINT "event_sonex_pool_mint_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_pool_swap
-- ----------------------------
CREATE INDEX "pool_swap_index_chain_block_number" ON "event_sonex_pool_swap" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "pool_swap_index_chain_id_pool" ON "event_sonex_pool_swap" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "extra_pool" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "pool_swap_index_chain_id_recipient" ON "event_sonex_pool_swap" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_recipient" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "pool_swap_index_chain_id_sender" ON "event_sonex_pool_swap" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_sender" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "pool_swap_index_chain_id_tx_hash_tx_index" ON "event_sonex_pool_swap" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_pool_swap
-- ----------------------------
ALTER TABLE "event_sonex_pool_swap" ADD CONSTRAINT "event_sonex_pool_swap_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_position_manager_collect
-- ----------------------------
CREATE INDEX "position_manager_collect_index_chain_id_block_number" ON "event_sonex_position_manager_collect" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_collect_index_chain_id_event_token_id" ON "event_sonex_position_manager_collect" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_token_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "position_manager_collect_index_chain_id_tx_hash_tx_index" ON "event_sonex_position_manager_collect" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_position_manager_collect
-- ----------------------------
ALTER TABLE "event_sonex_position_manager_collect" ADD CONSTRAINT "event_sonex_position_manager_collect_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_position_manager_decrease_liquidity
-- ----------------------------
CREATE UNIQUE INDEX "manager_decrease_liquidity_index_chain_id_tx_hash_tx_index" ON "event_sonex_position_manager_decrease_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_decrease_liquidity_index_chain_id_block_number" ON "event_sonex_position_manager_decrease_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_decrease_liquidity_index_chain_id_token_id" ON "event_sonex_position_manager_decrease_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_token_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_position_manager_decrease_liquidity
-- ----------------------------
ALTER TABLE "event_sonex_position_manager_decrease_liquidity" ADD CONSTRAINT "event_sonex_position_manager_decrease_liquidity_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_position_manager_increase_liquidity
-- ----------------------------
CREATE UNIQUE INDEX "manager_increase_liquidity_index_chain_id_tx_hash_tx_index" ON "event_sonex_position_manager_increase_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_increase_liquidity_index_chain_id_block_number" ON "event_sonex_position_manager_increase_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_increase_liquidity_index_chain_id_token_id" ON "event_sonex_position_manager_increase_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_token_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_position_manager_increase_liquidity
-- ----------------------------
ALTER TABLE "event_sonex_position_manager_increase_liquidity" ADD CONSTRAINT "event_sonex_position_manager_increase_liquidity_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_position_manager_transfer
-- ----------------------------
CREATE UNIQUE INDEX "manager_transfer_index_chain_id_tx_hash_tx_index" ON "event_sonex_position_manager_transfer" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_transfer_index_chain_id_block_number" ON "event_sonex_position_manager_transfer" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_transfer_index_chain_id_event_from" ON "event_sonex_position_manager_transfer" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_from" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_transfer_index_chain_id_event_to" ON "event_sonex_position_manager_transfer" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_to" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_transfer_index_chain_id_token_id" ON "event_sonex_position_manager_transfer" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_token_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_position_manager_transfer
-- ----------------------------
ALTER TABLE "event_sonex_position_manager_transfer" ADD CONSTRAINT "event_sonex_position_manager_transfer_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table pre_deposit_points_record
-- ----------------------------
CREATE INDEX "pre_deposit_points_record_chain_user" ON "pre_deposit_points_record" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "user_address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "pre_deposit_points_record_index_chain_id_tx_hash_tx_index" ON "pre_deposit_points_record" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table pre_deposit_points_record
-- ----------------------------
ALTER TABLE "pre_deposit_points_record" ADD CONSTRAINT "pre_deposit_points_record_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table pre_deposit_reward_claim_record
-- ----------------------------
CREATE UNIQUE INDEX "pre_deposit_reward_claim_record_index_chain_id_tx_hash_tx_index" ON "pre_deposit_reward_claim_record" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "pre_deposit_reward_claim_record_index_chain_id_user_nonce" ON "pre_deposit_reward_claim_record" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "user_address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "nonce" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table pre_deposit_reward_claim_record
-- ----------------------------
ALTER TABLE "pre_deposit_reward_claim_record" ADD CONSTRAINT "astar_claim_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table pre_deposit_reward_snapshot
-- ----------------------------
CREATE INDEX "pre_deposit_reward_snapshot_index_chain_id_user_token" ON "pre_deposit_reward_snapshot" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "user_address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "reward_token" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table pre_deposit_reward_snapshot
-- ----------------------------
ALTER TABLE "pre_deposit_reward_snapshot" ADD CONSTRAINT "pre_deposit_reward_snapshot_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sonex_token
-- ----------------------------
CREATE INDEX "idx_token_address" ON "sonex_token" USING btree (
  "address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_token_is_hidden" ON "sonex_token" USING btree (
  "isHidden" "pg_catalog"."bool_ops" ASC NULLS LAST
);
CREATE INDEX "idx_token_symbol" ON "sonex_token" USING btree (
  "symbol" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sonex_token
-- ----------------------------
ALTER TABLE "sonex_token" ADD CONSTRAINT "sonex_token_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sonex_token_price_record
-- ----------------------------
CREATE INDEX "token_address_idx" ON "sonex_token_price_record" USING btree (
  "token_address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sonex_token_price_record
-- ----------------------------
ALTER TABLE "sonex_token_price_record" ADD CONSTRAINT "sonex_token_price_record_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table token
-- ----------------------------
CREATE UNIQUE INDEX "token_index_chain_id_address" ON "token" USING btree (
  "address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table token
-- ----------------------------
ALTER TABLE "token" ADD CONSTRAINT "token_pkey" PRIMARY KEY ("id");
