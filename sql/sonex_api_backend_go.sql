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

 Date: 18/12/2024 17:06:56
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
  "event_fee" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
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
  "event_amount1" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
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
  "event_amount1" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
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
  "event_tick" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
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
CREATE INDEX "factory_pool_created_index_chain_id_block_number" ON "public"."event_sonex_factory_pool_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "factory_pool_created_index_chain_id_event_pool" ON "public"."event_sonex_factory_pool_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_pool" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "factory_pool_created_index_chain_id_tx_hash" ON "public"."event_sonex_factory_pool_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
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
CREATE INDEX "pool_burn_index_chain_id_event_owner" ON "public"."event_sonex_pool_burn" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_owner" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "pool_burn_index_chain_id_tx_hash" ON "public"."event_sonex_pool_burn" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
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
CREATE INDEX "pool_mint_index_chain_id_tx_event_owner" ON "public"."event_sonex_pool_mint" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_owner" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "pool_mint_index_chain_id_tx_event_sender" ON "public"."event_sonex_pool_mint" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_sender" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "pool_mint_index_chain_id_tx_hash" ON "public"."event_sonex_pool_mint" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
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
CREATE INDEX "pool_swap_index_chain_id_tx_event_recipient" ON "public"."event_sonex_pool_swap" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_recipient" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "pool_swap_index_chain_id_tx_event_sender" ON "public"."event_sonex_pool_swap" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_sender" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "pool_swap_index_chain_id_tx_hash" ON "public"."event_sonex_pool_swap" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
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
CREATE UNIQUE INDEX "position_manager_collect_index_chain_id_tx_hash" ON "public"."event_sonex_position_manager_collect" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_position_manager_collect
-- ----------------------------
ALTER TABLE "public"."event_sonex_position_manager_collect" ADD CONSTRAINT "event_sonex_position_manager_increase_liquidity_copy1_pkey1" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_position_manager_decrease_liquidity
-- ----------------------------
CREATE INDEX "position_manager_decrease_liquidity_index_chain_id_block_number" ON "public"."event_sonex_position_manager_decrease_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_decrease_liquidity_index_chain_id_token_id" ON "public"."event_sonex_position_manager_decrease_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_token_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "position_manager_decrease_liquidity_index_chain_id_tx_hash" ON "public"."event_sonex_position_manager_decrease_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_position_manager_decrease_liquidity
-- ----------------------------
ALTER TABLE "public"."event_sonex_position_manager_decrease_liquidity" ADD CONSTRAINT "event_sonex_position_manager_increase_liquidity_copy1_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_position_manager_increase_liquidity
-- ----------------------------
CREATE INDEX "position_manager_increase_liquidity_index_chain_id_block_number" ON "public"."event_sonex_position_manager_increase_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "position_manager_increase_liquidity_index_chain_id_token_id" ON "public"."event_sonex_position_manager_increase_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_token_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "position_manager_increase_liquidity_index_chain_id_tx_hash" ON "public"."event_sonex_position_manager_increase_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_position_manager_increase_liquidity
-- ----------------------------
ALTER TABLE "public"."event_sonex_position_manager_increase_liquidity" ADD CONSTRAINT "event_sonex_position_manager_transfer_copy1_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_sonex_position_manager_transfer
-- ----------------------------
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
CREATE UNIQUE INDEX "position_manager_transfer_index_chain_id_tx_hash" ON "public"."event_sonex_position_manager_transfer" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_sonex_position_manager_transfer
-- ----------------------------
ALTER TABLE "public"."event_sonex_position_manager_transfer" ADD CONSTRAINT "event_sonex_pool_swap_copy1_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Initialize data
-- ----------------------------
INSERT INTO "public"."chain_scanner_prograss" ("id", "created_at", "updated_at", "chain_id", "height") VALUES (124502484323401728, '2024-12-10 08:32:30+00', '2024-12-13 02:40:50+00', 1946, 4800000);


CREATE TABLE "public"."sonex_token_price_record" (
                                                     "id" int8 NOT NULL DEFAULT nextval('sonex_token_price_record_id'::regclass),
                                                     "token_symbol" varchar(255) COLLATE "pg_catalog"."default",
                                                     "token_address" varchar(255) COLLATE "pg_catalog"."default",
                                                     "price" varchar COLLATE "pg_catalog"."default",
                                                     "created_at" timestamp(6),
                                                     CONSTRAINT "sonex_token_price_record_pkey" PRIMARY KEY ("id")
);

ALTER TABLE "public"."sonex_token_price_record"
    OWNER TO "postgres";

CREATE SEQUENCE "public"."sonex_token_price_record_id"
    INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

ALTER SEQUENCE "public"."sonex_token_price_record_id" OWNER TO "postgres";

CREATE INDEX "token_address_idx" ON "public"."sonex_token_price_record" USING btree (
    "token_address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
    );