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

 Date: 14/01/2025 21:00:39
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

-- ----------------------------
-- Table structure for event_sonex_factory_pool_created
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_sonex_factory_pool_created";
CREATE TABLE "public"."event_sonex_factory_pool_created" (
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
DROP TABLE IF EXISTS "public"."event_sonex_pool_burn";
CREATE TABLE "public"."event_sonex_pool_burn" (
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
DROP TABLE IF EXISTS "public"."event_sonex_pool_mint";
CREATE TABLE "public"."event_sonex_pool_mint" (
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
DROP TABLE IF EXISTS "public"."event_sonex_pool_swap";
CREATE TABLE "public"."event_sonex_pool_swap" (
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
DROP TABLE IF EXISTS "public"."event_sonex_position_manager_collect";
CREATE TABLE "public"."event_sonex_position_manager_collect" (
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
DROP TABLE IF EXISTS "public"."event_sonex_position_manager_decrease_liquidity";
CREATE TABLE "public"."event_sonex_position_manager_decrease_liquidity" (
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
DROP TABLE IF EXISTS "public"."event_sonex_position_manager_increase_liquidity";
CREATE TABLE "public"."event_sonex_position_manager_increase_liquidity" (
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
DROP TABLE IF EXISTS "public"."event_sonex_position_manager_transfer";
CREATE TABLE "public"."event_sonex_position_manager_transfer" (
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
-- Table structure for sonex_token
-- ----------------------------
DROP TABLE IF EXISTS "public"."sonex_token";
CREATE TABLE "public"."sonex_token" (
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
DROP TABLE IF EXISTS "public"."sonex_token_price_record";
CREATE TABLE "public"."sonex_token_price_record" (
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
-- Indexes structure for table chain_scanner_prograss
-- ----------------------------
CREATE UNIQUE INDEX "chain_scanner_prograss_index_chain_id" ON "public"."chain_scanner_prograss" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table chain_scanner_prograss
-- ----------------------------
ALTER TABLE "public"."chain_scanner_prograss" ADD CONSTRAINT "chain_scanner_prograss_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_factory_pool_created
-- ----------------------------
CREATE INDEX "sonex_factory_pool_created_index_chain_id_block_number" ON "public"."event_sonex_factory_pool_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "sonex_factory_pool_created_index_chain_id_tx_hash_tx_index" ON "public"."event_sonex_factory_pool_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "sonex_factory_pool_created_index_pool" ON "public"."event_sonex_factory_pool_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_pool" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_factory_pool_created
-- ----------------------------
ALTER TABLE "public"."event_sonex_factory_pool_created" ADD CONSTRAINT "event_sonex_pool_created_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_pool_burn
-- ----------------------------
CREATE INDEX "pool_burn_index_chain_id_block_number" ON "public"."event_sonex_pool_burn" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "pool_burn_index_chain_id_owner" ON "public"."event_sonex_pool_burn" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_owner" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "pool_burn_index_chain_id_pool" ON "public"."event_sonex_pool_burn" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "extra_pool" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "pool_burn_index_chain_id_tx_hash_tx_index" ON "public"."event_sonex_pool_burn" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_pool_burn
-- ----------------------------
ALTER TABLE "public"."event_sonex_pool_burn" ADD CONSTRAINT "event_sonex_pool_mint_copy1_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_pool_mint
-- ----------------------------
CREATE INDEX "pool_mint_index_chain_block_number" ON "public"."event_sonex_pool_mint" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "pool_mint_index_chain_id_owner" ON "public"."event_sonex_pool_mint" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_owner" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "pool_mint_index_chain_id_pool" ON "public"."event_sonex_pool_mint" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "extra_pool" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "pool_mint_index_chain_id_sender" ON "public"."event_sonex_pool_mint" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_sender" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "pool_mint_index_chain_id_tx_hash_tx_index" ON "public"."event_sonex_pool_mint" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_pool_mint
-- ----------------------------
ALTER TABLE "public"."event_sonex_pool_mint" ADD CONSTRAINT "event_sonex_factory_pool_created_copy1_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_pool_swap
-- ----------------------------
CREATE INDEX "pool_swap_index_chain_block_number" ON "public"."event_sonex_pool_swap" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "pool_swap_index_chain_id_pool" ON "public"."event_sonex_pool_swap" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "extra_pool" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "pool_swap_index_chain_id_recipient" ON "public"."event_sonex_pool_swap" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_recipient" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "pool_swap_index_chain_id_sender" ON "public"."event_sonex_pool_swap" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_sender" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "pool_swap_index_chain_id_tx_hash_tx_index" ON "public"."event_sonex_pool_swap" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_pool_swap
-- ----------------------------
ALTER TABLE "public"."event_sonex_pool_swap" ADD CONSTRAINT "event_sonex_pool_mint_copy1_pkey1" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_position_manager_collect
-- ----------------------------
CREATE INDEX "position_manager_collect_index_chain_id_block_number" ON "public"."event_sonex_position_manager_collect" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_collect_index_chain_id_event_token_id" ON "public"."event_sonex_position_manager_collect" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_token_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "position_manager_collect_index_chain_id_tx_hash_tx_index" ON "public"."event_sonex_position_manager_collect" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_position_manager_collect
-- ----------------------------
ALTER TABLE "public"."event_sonex_position_manager_collect" ADD CONSTRAINT "event_sonex_position_manager_increase_liquidity_copy1_pkey1" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_position_manager_decrease_liquidity
-- ----------------------------
CREATE UNIQUE INDEX "manager_decrease_liquidity_index_chain_id_tx_hash_tx_index" ON "public"."event_sonex_position_manager_decrease_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_decrease_liquidity_index_chain_id_block_number" ON "public"."event_sonex_position_manager_decrease_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_decrease_liquidity_index_chain_id_token_id" ON "public"."event_sonex_position_manager_decrease_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_token_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_position_manager_decrease_liquidity
-- ----------------------------
ALTER TABLE "public"."event_sonex_position_manager_decrease_liquidity" ADD CONSTRAINT "event_sonex_position_manager_increase_liquidity_copy1_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_position_manager_increase_liquidity
-- ----------------------------
CREATE UNIQUE INDEX "manager_increase_liquidity_index_chain_id_tx_hash_tx_index" ON "public"."event_sonex_position_manager_increase_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_increase_liquidity_index_chain_id_block_number" ON "public"."event_sonex_position_manager_increase_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_increase_liquidity_index_chain_id_token_id" ON "public"."event_sonex_position_manager_increase_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_token_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_position_manager_increase_liquidity
-- ----------------------------
ALTER TABLE "public"."event_sonex_position_manager_increase_liquidity" ADD CONSTRAINT "event_sonex_position_manager_transfer_copy1_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_position_manager_transfer
-- ----------------------------
CREATE UNIQUE INDEX "manager_transfer_index_chain_id_tx_hash_tx_index" ON "public"."event_sonex_position_manager_transfer" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_transfer_index_chain_id_block_number" ON "public"."event_sonex_position_manager_transfer" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_transfer_index_chain_id_event_from" ON "public"."event_sonex_position_manager_transfer" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_from" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_transfer_index_chain_id_event_to" ON "public"."event_sonex_position_manager_transfer" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_to" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_transfer_index_chain_id_token_id" ON "public"."event_sonex_position_manager_transfer" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_token_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_position_manager_transfer
-- ----------------------------
ALTER TABLE "public"."event_sonex_position_manager_transfer" ADD CONSTRAINT "event_sonex_pool_swap_copy1_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sonex_token
-- ----------------------------
CREATE INDEX "idx_token_address" ON "public"."sonex_token" USING btree (
  "address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_token_is_hidden" ON "public"."sonex_token" USING btree (
  "isHidden" "pg_catalog"."bool_ops" ASC NULLS LAST
);
CREATE INDEX "idx_token_symbol" ON "public"."sonex_token" USING btree (
  "symbol" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sonex_token
-- ----------------------------
ALTER TABLE "public"."sonex_token" ADD CONSTRAINT "sonex_token_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sonex_token_price_record
-- ----------------------------
CREATE INDEX "token_address_idx" ON "public"."sonex_token_price_record" USING btree (
  "token_address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sonex_token_price_record
-- ----------------------------
ALTER TABLE "public"."sonex_token_price_record" ADD CONSTRAINT "sonex_token_price_record_pkey" PRIMARY KEY ("id");

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


/*
 Navicat Premium Dump SQL

 Source Server         : test001
 Source Server Type    : PostgreSQL
 Source Server Version : 160004 (160004)
 Source Host           : 10.1.84.8:5432
 Source Catalog        : sonex-api-backend-go
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 160004 (160004)
 File Encoding         : 65001

 Date: 17/01/2025 21:18:04
*/


-- ----------------------------
-- Table structure for sonex_pool_tvl_record
-- ----------------------------
DROP TABLE IF EXISTS "public"."sonex_pool_tvl_record";
CREATE TABLE "public"."sonex_pool_tvl_record" (
  "id" int8 NOT NULL DEFAULT nextval('sonex_token_price_record_id'::regclass),
  "pool_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tvl" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamp(6) NOT NULL,
  "usd_decimals" int8 NOT NULL
)
;
ALTER TABLE "public"."sonex_pool_tvl_record" OWNER TO "apibackendgo";

-- ----------------------------
-- Indexes structure for table sonex_pool_tvl_record
-- ----------------------------
CREATE INDEX "token_address_idx_copy1" ON "public"."sonex_pool_tvl_record" USING btree (
  "pool_address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sonex_pool_tvl_record
-- ----------------------------
ALTER TABLE "public"."sonex_pool_tvl_record" ADD CONSTRAINT "sonex_token_price_record_copy1_pkey" PRIMARY KEY ("id");
