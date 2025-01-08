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

 Date: 09/01/2025 01:17:45
*/

-- ----------------------------
-- Sequence structure for sonex_token_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sonex_token_id_seq";
CREATE SEQUENCE "public"."sonex_token_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sonex_token_price_record_id
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sonex_token_price_record_id";
CREATE SEQUENCE "public"."sonex_token_price_record_id" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for token_id
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."token_id";
CREATE SEQUENCE "public"."token_id" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

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
-- Table structure for event_farm_incentive_created
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_incentive_created";
CREATE TABLE "public"."event_farm_incentive_created" (
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
  "event_reward_token" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_pool" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_incentive_id" int8 NOT NULL,
  "event_start_time" int8 NOT NULL,
  "event_weight" int8 NOT NULL,
  "event_timestamp" int8 NOT NULL
)
;

-- ----------------------------
-- Table structure for event_farm_incentive_ended
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_incentive_ended";
CREATE TABLE "public"."event_farm_incentive_ended" (
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
  "event_incentive_id" int8 NOT NULL,
  "event_weight" int8 NOT NULL,
  "event_timestamp" int8 NOT NULL
)
;

-- ----------------------------
-- Table structure for event_farm_incentive_updated
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_incentive_updated";
CREATE TABLE "public"."event_farm_incentive_updated" (
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
  "event_incentive_id" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_start_time" int8 NOT NULL,
  "event_weight" int8 NOT NULL,
  "event_total_liquidity" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_farm_paused
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_paused";
CREATE TABLE "public"."event_farm_paused" (
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
  "event_account" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_farm_reward_claimed
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_reward_claimed";
CREATE TABLE "public"."event_farm_reward_claimed" (
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
  "event_incentive_id" int8 NOT NULL,
  "event_amount" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_farm_reward_pool_updated
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_reward_pool_updated";
CREATE TABLE "public"."event_farm_reward_pool_updated" (
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
  "event_reward_token" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_total_weight" int8 NOT NULL,
  "event_timestamp" int8 NOT NULL
)
;

-- ----------------------------
-- Table structure for event_farm_reward_rate_updated
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_reward_rate_updated";
CREATE TABLE "public"."event_farm_reward_rate_updated" (
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
  "event_reward_token" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_old_rate_pre_second" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_new_rate_pre_second" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_timestamp" int8 NOT NULL
)
;

-- ----------------------------
-- Table structure for event_farm_token_deposited
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_token_deposited";
CREATE TABLE "public"."event_farm_token_deposited" (
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
  "event_token_id" int8 NOT NULL,
  "event_owner" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_farm_token_staked
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_token_staked";
CREATE TABLE "public"."event_farm_token_staked" (
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
  "event_token_id" int8 NOT NULL,
  "event_incentive_id" int8 NOT NULL,
  "event_liquidity" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_farm_token_unstaked
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_token_unstaked";
CREATE TABLE "public"."event_farm_token_unstaked" (
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
  "event_token_id" int8 NOT NULL,
  "event_incentive_id" int8 NOT NULL,
  "event_liquidity" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_farm_unpaused
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_unpaused";
CREATE TABLE "public"."event_farm_unpaused" (
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
  "event_account" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
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
-- Table structure for event_xpool_factory_pool_created
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_xpool_factory_pool_created";
CREATE TABLE "public"."event_xpool_factory_pool_created" (
  "id" int8 NOT NULL DEFAULT nextval('event_xpool_factory_pool_created_id_seq'::regclass),
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
  "event_pool" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_staking_token" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_reward_token" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_reward_pre_block" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_xpool_factory_pool_destroied
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_xpool_factory_pool_destroied";
CREATE TABLE "public"."event_xpool_factory_pool_destroied" (
  "id" int8 NOT NULL DEFAULT nextval('event_xpool_factory_pool_destroied_id_seq'::regclass),
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
  "event_pool" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_eth_balance" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_staking_balance" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_reward_balance" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_recipient" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_xpool_pool_claimed
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_xpool_pool_claimed";
CREATE TABLE "public"."event_xpool_pool_claimed" (
  "id" int8 NOT NULL DEFAULT nextval('event_xpool_pool_claimed_id_seq'::regclass),
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
  "event_user" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_height" int8 NOT NULL
)
;

-- ----------------------------
-- Table structure for event_xpool_pool_deposited
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_xpool_pool_deposited";
CREATE TABLE "public"."event_xpool_pool_deposited" (
  "id" int8 NOT NULL DEFAULT nextval('event_xpool_pool_deposited_id_seq'::regclass),
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
  "event_user" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_height" int8 NOT NULL
)
;

-- ----------------------------
-- Table structure for event_xpool_pool_withdrawed
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_xpool_pool_withdrawed";
CREATE TABLE "public"."event_xpool_pool_withdrawed" (
  "id" int8 NOT NULL DEFAULT nextval('event_xpool_pool_withdrawed_id_seq'::regclass),
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
  "event_user" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_height" int8 NOT NULL
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
  "token_symbol" varchar(255) COLLATE "pg_catalog"."default",
  "token_address" varchar(255) COLLATE "pg_catalog"."default",
  "price" varchar COLLATE "pg_catalog"."default",
  "created_at" timestamp(6)
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
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."sonex_token_id_seq"', 1, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."sonex_token_price_record_id"', 1, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."token_id"', 1, true);

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
-- Indexes structure for table event_farm_incentive_created
-- ----------------------------
CREATE UNIQUE INDEX "event_farm_incentive_created_index_chain_id_incentive_id" ON "public"."event_farm_incentive_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_incentive_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "event_farm_incentive_created_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_incentive_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table event_farm_incentive_created
-- ----------------------------
ALTER TABLE "public"."event_farm_incentive_created" ADD CONSTRAINT "event_farm_incentive_created_tx_hash_chain_id_key" UNIQUE ("tx_hash", "chain_id");

-- ----------------------------
-- Primary Key structure for table event_farm_incentive_created
-- ----------------------------
ALTER TABLE "public"."event_farm_incentive_created" ADD CONSTRAINT "event_farm_incentive_created_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_incentive_ended
-- ----------------------------
CREATE UNIQUE INDEX "event_farm_incentive_ended_index_chain_id_incentive_id" ON "public"."event_farm_incentive_ended" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_incentive_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "event_farm_incentive_ended_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_incentive_ended" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table event_farm_incentive_ended
-- ----------------------------
ALTER TABLE "public"."event_farm_incentive_ended" ADD CONSTRAINT "event_farm_incentive_ended_tx_hash_chain_id_key" UNIQUE ("tx_hash", "chain_id");

-- ----------------------------
-- Primary Key structure for table event_farm_incentive_ended
-- ----------------------------
ALTER TABLE "public"."event_farm_incentive_ended" ADD CONSTRAINT "event_farm_incentive_ended_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_incentive_updated
-- ----------------------------
CREATE UNIQUE INDEX "event_farm_incentive_updated_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_incentive_updated" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table event_farm_incentive_updated
-- ----------------------------
ALTER TABLE "public"."event_farm_incentive_updated" ADD CONSTRAINT "event_farm_token_staked_copy1_tx_hash_chain_id_key" UNIQUE ("tx_hash", "chain_id");

-- ----------------------------
-- Primary Key structure for table event_farm_incentive_updated
-- ----------------------------
ALTER TABLE "public"."event_farm_incentive_updated" ADD CONSTRAINT "event_farm_token_staked_copy1_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_paused
-- ----------------------------
CREATE UNIQUE INDEX "event_farm_paused_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_paused" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table event_farm_paused
-- ----------------------------
ALTER TABLE "public"."event_farm_paused" ADD CONSTRAINT "event_farm_unpaused_copy1_tx_hash_chain_id_key" UNIQUE ("tx_hash", "chain_id");

-- ----------------------------
-- Primary Key structure for table event_farm_paused
-- ----------------------------
ALTER TABLE "public"."event_farm_paused" ADD CONSTRAINT "event_farm_unpaused_copy1_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_reward_claimed
-- ----------------------------
CREATE INDEX "event_farm_reward_claimed_index_chain_id_incentive_id" ON "public"."event_farm_reward_claimed" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_incentive_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "event_farm_reward_claimed_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_reward_claimed" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "event_farm_reward_claimed_index_chain_id_user" ON "public"."event_farm_reward_claimed" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_user" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table event_farm_reward_claimed
-- ----------------------------
ALTER TABLE "public"."event_farm_reward_claimed" ADD CONSTRAINT "event_farm_reward_claimed_tx_hash_chain_id_key" UNIQUE ("tx_hash", "chain_id");

-- ----------------------------
-- Primary Key structure for table event_farm_reward_claimed
-- ----------------------------
ALTER TABLE "public"."event_farm_reward_claimed" ADD CONSTRAINT "event_farm_reward_claimed_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_reward_pool_updated
-- ----------------------------
CREATE UNIQUE INDEX "event_farm_reward_pool_updated_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_reward_pool_updated" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table event_farm_reward_pool_updated
-- ----------------------------
ALTER TABLE "public"."event_farm_reward_pool_updated" ADD CONSTRAINT "event_farm_reward_rate_updated_copy1_tx_hash_chain_id_key" UNIQUE ("tx_hash", "chain_id");

-- ----------------------------
-- Primary Key structure for table event_farm_reward_pool_updated
-- ----------------------------
ALTER TABLE "public"."event_farm_reward_pool_updated" ADD CONSTRAINT "event_farm_reward_rate_updated_copy1_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_reward_rate_updated
-- ----------------------------
CREATE UNIQUE INDEX "event_farm_reward_rate_updated_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_reward_rate_updated" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table event_farm_reward_rate_updated
-- ----------------------------
ALTER TABLE "public"."event_farm_reward_rate_updated" ADD CONSTRAINT "event_farm_reward_claimed_copy1_tx_hash_chain_id_key" UNIQUE ("tx_hash", "chain_id");

-- ----------------------------
-- Primary Key structure for table event_farm_reward_rate_updated
-- ----------------------------
ALTER TABLE "public"."event_farm_reward_rate_updated" ADD CONSTRAINT "event_farm_reward_claimed_copy1_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_token_deposited
-- ----------------------------
CREATE INDEX "event_farm_token_deposited_index_chain_id_owner" ON "public"."event_farm_token_deposited" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_owner" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "event_farm_token_deposited_index_chain_id_token_id" ON "public"."event_farm_token_deposited" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_token_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "event_farm_token_deposited_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_token_deposited" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table event_farm_token_deposited
-- ----------------------------
ALTER TABLE "public"."event_farm_token_deposited" ADD CONSTRAINT "event_farm_deposit_transferred_tx_hash_chain_id_key" UNIQUE ("tx_hash", "chain_id");

-- ----------------------------
-- Primary Key structure for table event_farm_token_deposited
-- ----------------------------
ALTER TABLE "public"."event_farm_token_deposited" ADD CONSTRAINT "event_farm_deposit_transferred_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_token_staked
-- ----------------------------
CREATE INDEX "event_farm_token_staked_index_chain_id_incentive_id" ON "public"."event_farm_token_staked" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_incentive_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "event_farm_token_staked_index_chain_id_token_id" ON "public"."event_farm_token_staked" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_token_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "event_farm_token_staked_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_token_staked" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table event_farm_token_staked
-- ----------------------------
ALTER TABLE "public"."event_farm_token_staked" ADD CONSTRAINT "event_farm_token_staked_tx_hash_chain_id_key" UNIQUE ("tx_hash", "chain_id");

-- ----------------------------
-- Primary Key structure for table event_farm_token_staked
-- ----------------------------
ALTER TABLE "public"."event_farm_token_staked" ADD CONSTRAINT "event_farm_token_staked_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_token_unstaked
-- ----------------------------
CREATE INDEX "event_farm_token_unstaked_index_chain_id_incentive_id" ON "public"."event_farm_token_unstaked" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_incentive_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "event_farm_token_unstaked_index_chain_id_token_id" ON "public"."event_farm_token_unstaked" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_token_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "event_farm_token_unstaked_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_token_unstaked" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table event_farm_token_unstaked
-- ----------------------------
ALTER TABLE "public"."event_farm_token_unstaked" ADD CONSTRAINT "event_farm_token_unstaked_tx_hash_chain_id_key" UNIQUE ("tx_hash", "chain_id");

-- ----------------------------
-- Primary Key structure for table event_farm_token_unstaked
-- ----------------------------
ALTER TABLE "public"."event_farm_token_unstaked" ADD CONSTRAINT "event_farm_token_unstaked_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_unpaused
-- ----------------------------
CREATE UNIQUE INDEX "event_farm_unpaused_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_unpaused" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table event_farm_unpaused
-- ----------------------------
ALTER TABLE "public"."event_farm_unpaused" ADD CONSTRAINT "event_farm_paused_copy1_tx_hash_chain_id_key" UNIQUE ("tx_hash", "chain_id");

-- ----------------------------
-- Primary Key structure for table event_farm_unpaused
-- ----------------------------
ALTER TABLE "public"."event_farm_unpaused" ADD CONSTRAINT "event_farm_paused_copy1_pkey" PRIMARY KEY ("id");

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
-- Indexes structure for table event_xpool_factory_pool_created
-- ----------------------------
CREATE INDEX "xpool_factory_pool_created_index_chain_id_block_number" ON "public"."event_xpool_factory_pool_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "xpool_factory_pool_created_index_chain_id_pool" ON "public"."event_xpool_factory_pool_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_pool" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "xpool_factory_pool_created_index_chain_id_reward_token" ON "public"."event_xpool_factory_pool_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_reward_token" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "xpool_factory_pool_created_index_chain_id_staking_token" ON "public"."event_xpool_factory_pool_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_staking_token" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "xpool_factory_pool_created_index_chain_id_tx_hash_tx_index" ON "public"."event_xpool_factory_pool_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_xpool_factory_pool_created
-- ----------------------------
ALTER TABLE "public"."event_xpool_factory_pool_created" ADD CONSTRAINT "event_sonex_position_manager_transfer_copy1_pkey1" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_xpool_factory_pool_destroied
-- ----------------------------
CREATE INDEX "xpool_factory_pool_destroied_index_chain_id_block_number" ON "public"."event_xpool_factory_pool_destroied" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "xpool_factory_pool_destroied_index_chain_id_pool" ON "public"."event_xpool_factory_pool_destroied" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_pool" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "xpool_factory_pool_destroied_index_chain_id_recipient" ON "public"."event_xpool_factory_pool_destroied" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_recipient" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "xpool_factory_pool_destroied_index_chain_id_tx_hash_tx_index" ON "public"."event_xpool_factory_pool_destroied" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_xpool_factory_pool_destroied
-- ----------------------------
ALTER TABLE "public"."event_xpool_factory_pool_destroied" ADD CONSTRAINT "event_xpool_factory_pool_created_copy1_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_xpool_pool_claimed
-- ----------------------------
CREATE INDEX "xpool_pool_claimed_index_chain_id_block_number" ON "public"."event_xpool_pool_claimed" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "xpool_pool_claimed_index_chain_id_tx_hash_tx_index" ON "public"."event_xpool_pool_claimed" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "xpool_pool_claimed_index_chain_id_user" ON "public"."event_xpool_pool_claimed" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_user" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_xpool_pool_claimed
-- ----------------------------
ALTER TABLE "public"."event_xpool_pool_claimed" ADD CONSTRAINT "event_xpool_pool_destroied_copy1_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_xpool_pool_deposited
-- ----------------------------
CREATE INDEX "xpool_pool_deposited_index_chain_id_block_number" ON "public"."event_xpool_pool_deposited" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "xpool_pool_deposited_index_chain_id_tx_hash_tx_index" ON "public"."event_xpool_pool_deposited" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "xpool_pool_deposited_index_chain_id_user" ON "public"."event_xpool_pool_deposited" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_user" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_xpool_pool_deposited
-- ----------------------------
ALTER TABLE "public"."event_xpool_pool_deposited" ADD CONSTRAINT "event_xpool_pool_deposited_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_xpool_pool_withdrawed
-- ----------------------------
CREATE INDEX "xpool_pool_withdrawed_index_chain_id_block_number" ON "public"."event_xpool_pool_withdrawed" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "xpool_pool_withdrawed_index_chain_id_tx_hash_tx_index" ON "public"."event_xpool_pool_withdrawed" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "xpool_pool_withdrawed_index_chain_id_user" ON "public"."event_xpool_pool_withdrawed" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_user" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_xpool_pool_withdrawed
-- ----------------------------
ALTER TABLE "public"."event_xpool_pool_withdrawed" ADD CONSTRAINT "event_xpool_pool_claimed_copy1_pkey" PRIMARY KEY ("id");

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

-- ----------------------------
-- Initialize date
-- ----------------------------
INSERT INTO "chain_scanner_prograss" ("id", "created_at", "updated_at", "chain_id", "height") VALUES (124502484323401728, '2024-12-10 08:32:30+00', '2024-12-13 02:40:50+00', 1946, 5000000);

INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (30, '2024-12-26 02:37:45.507939+00', '2024-12-26 02:37:45.507939+00', NULL, 'ACT', 'ACT', '0x0a73297F79858Ad1d58e68B45A806D4e6C6D5459', 18, 'https://img.sonex.so/act-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "", "reddit": "https://reddit.com/r/ACT1", "discord": "", "twitter": "https://twitter.com/ACTICOMMUNITY", "bitclout": "", "telegram": "https://t.me/actportal"}, "website": "https://actsol.xyz/", "whitepaper": "", "description": "Act I: The AI Prophecy Project (ACT) is a decentralized, open-source protocol designed to redefine AI interaction. This initiative fosters a diverse “ecosystem” where various AI systems, driven by the ACT token, can collaborate with each other and engage with human users, sparking innovative and unconventional AI interactions. Unlike mainstream chatbot approaches that prioritize helpfulness, harmlessness, and honesty—as championed by companies like Microsoft, OpenAI, and Anthropic, Act I encourages a more open-ended, creative approach to AI development."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (10, '2024-12-26 02:36:58.252013+00', '2024-12-26 02:36:58.252013+00', NULL, 'XRP', 'XRP', '0x795c68F07d3ced21fc357A76E077b871457A74A2', 18, 'https://img.sonex.so/xrp-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "https://cmc.certik-skynet.com/redirect?project=xrp-ledger", "social": {"github": "https://github.com/ripple/rippled", "reddit": "https://reddit.com/r/ripple", "discord": "", "twitter": "https://twitter.com/Ripple", "bitclout": "", "telegram": ""}, "website": "https://xrpl.org/", "whitepaper": "https://ripple.com/files/ripple_consensus_whitepaper.pdf", "description": "Launched in 2012, the XRP Ledger (XRPL) is an open-source, permissionless and decentralized technology. Benefits of the XRP Ledger include its low-cost ($0.0002 to transact), speed (settling transactions in 3-5 seconds), scalability (1,500 transactions per second) and inherently green attributes (carbon-neutral and energy-efficient). The XRP Ledger also features the first decentralized exchange (DEX) and custom tokenization capabilities built into the protocol. Since 2012, the XRP Ledger has been operating reliably, having closed 70 million ledgers."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (1, '2024-12-25 07:38:42.066938+00', '2024-12-25 07:38:42.066938+00', NULL, 'Wrapped Ethereum', 'WETH', '0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C', 18, 'https://img.sonex.so/weth-logo.png', 0, 0, 'f', '0.000000000000000000', '{"audits": "https://hacken.io/audits/singularitydao/", "social": {"github": "", "reddit": "", "discord": "", "twitter": "", "bitclout": "", "telegram": ""}, "website": "", "whitepaper": "", "description": "What is WETH (Wrapped ETH)? WETH is the tokenized/packaged form of ETH that you use to pay for items when you interact with Ethereum dApps. WETH follows the ERC-20 token standards, enabling it to achieve interoperability with other ERC-20 tokens.  This offers more utility to holders as they can use it across networks and dApps. You can stake, yield farm, lend, and provide liquidity to various liquidity pools with WETH.  Also, unlike ETH, which doesn’t conform to its own ERC-20 standard and thus has lower interoperability as it can’t be used on other chains besides Ethereum, WETH can be used on cheaper and high throughput alternatives like Binance, Polygon, Solana, and Cardano. The price of WETH will always be the same as ETH because it maintains a 1:1 wrapping ratio."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (2, '2024-12-25 07:38:45.262868+00', '2024-12-25 07:38:45.262868+00', NULL, 'Sonex Token', 'SONEX', '0x8879309d4002b8C98163Ece3f58653F774de8416', 18, 'https://img.sonex.so/sonex-logo.png', 0, 0, 'f', '0.000000000000000000', '{"audits": "", "social": {"github": "", "reddit": "", "discord": "https://discord.com/invite/sonex", "twitter": "https://x.com/sonex_so", "bitclout": "", "telegram": "https://t.me/SONEXexchange"}, "website": "https://sonex.so/", "whitepaper": "", "description": "SONEX is your gateway to the Soneium ecosystem, redefining decentralized finance by delivering a unified, all-in-one DeFi hub for retail users. We’re breaking away from the reliance on centralized platforms to bring you true decentralization.Launching with liquidity provision and swap functionality, SONEX will evolve into a fully integrated DeFi hub."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (3, '2024-12-25 07:38:48.009602+00', '2024-12-25 07:38:48.009602+00', NULL, 'Astar Minato Testnet Token', 'ASTR', '0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655', 18, 'https://img.sonex.so/astar-logo.png', 0, 1, 'f', '0.000000000000000000', '{"audits": "https://certificate.quantstamp.com/view/a-star-network-custom-signature", "social": {"github": "https://github.com/AstarNetwork/Astar", "reddit": "", "discord": "https://discord.gg/astarnetwork", "twitter": "https://twitter.com/AstarNetwork", "bitclout": "", "telegram": "https://t.me/PlasmAnnouncement"}, "website": "https://astar.network/", "whitepaper": "https://docs.astar.network/", "description": "Astar Network helps developers build decentralized applications (dApps) and layer 2 solutions by providing interoperable web 3.0 infrastructure, financial incentives, incubation programmes and technical support. It seeks to provide the best solution for all developers by supporting EVM and making a parachain where EVM and WASM smart contracts can co-exist and communicate with each other."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (11, '2024-12-26 02:37:01.021809+00', '2024-12-26 02:37:01.021809+00', NULL, 'BNB', 'BNB', '0x13aA369B93343e7b89cFCae9725A2f22895244c6', 18, 'https://img.sonex.so/bnb-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "https://audits.hacken.io/binance", "social": {"github": "https://github.com/bnb-chain", "reddit": "https://reddit.com/r/bnbchainofficial", "discord": "", "twitter": "https://twitter.com/bnbchain", "bitclout": "", "telegram": "https://t.me/BNBchaincommunity"}, "website": "https://bnbchain.org/en", "whitepaper": "", "description": "Launched in July 2017, Binance is the biggest cryptocurrency exchange globally based on daily trading volume. Binance aims to bring cryptocurrency exchanges to the forefront of financial activity globally. The idea behind Binance’s name is to show this new paradigm in global finance — Binary Finance, or Binance.Aside from being the largest cryptocurrency exchange globally, Binance has launched a whole ecosystem of functionalities for its users. The Binance network includes the Binance Chain, Binance Smart Chain, Binance Academy, Trust Wallet and Research projects, which all employ the powers of blockchain technology to bring new-age finance to the world. BNB is an integral part of the successful functioning of many of the Binance sub-projects."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (16, '2024-12-26 02:37:12.038345+00', '2024-12-26 02:37:12.038345+00', NULL, 'MOVE', 'MOVE', '0x3FCd768bAdFE17b8A5e878109657F0eaC0e4a168', 18, 'https://img.sonex.so/move-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "https://github.com/movementlabsxyz/", "reddit": "", "discord": "https://discord.com/invite/movementlabsxyz", "twitter": "https://twitter.com/movementfdn", "bitclout": "", "telegram": "https://t.me/movementlabsxyz"}, "website": "https://www.movementnetwork.xyz/", "whitepaper": "https://docs.movementnetwork.xyz/", "description": "Movement Network is an ecosystem of Modular Move-Based Blockchains that enables developers to build secure, performant, and interoperable blockchain applications, bridging the gap between Move and EVM ecosystems. Movement Network is first Move-EVM L2 for Ethereum, alongside open-source tooling and protocols to facilitate the adoption of the Move programming language across blockchain ecosystems.With Movement, developers can launch high-performance Move VM rollups with ease. Backed by Polychain Capital, Binance Labs, Hack VC, Placeholder, and Archetype, Movement Labs revolutionizes blockchain interoperability and advance Move-based technologies in the Web3 space."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (17, '2024-12-26 02:37:14.21648+00', '2024-12-26 02:37:14.21648+00', NULL, 'ME', 'ME', '0xF533784b6dD49C363cfF778067EC5983292B83e7', 18, 'https://img.sonex.so/me-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "https://github.com/me-foundation", "reddit": "", "discord": "https://discord.com/invite/magiceden", "twitter": "https://twitter.com/mefndn", "bitclout": "", "telegram": ""}, "website": "https://mefoundation.com/", "whitepaper": "https://blog.mefoundation.com/blog/me-tokenomics/", "description": "$ME is powering the onchain economy and represents the most powerful suite of cross chain trading, minting, and wallet protocols. Adopted by Magic Eden as their official token, $ME will represent the largest community of onchain users who will earn $ME for trading all assets on all chains.$ME represents the power of decentralized applications. The largest builder of the protocols is Magic Eden, which offers both NFT and token trading across all major ecosystems. Magic Eden represents the entrypoint for new retail users to enter crypto, complete with a retail-friendly mobile app where people can trade all assets on all chains."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (21, '2024-12-26 02:37:24.290876+00', '2024-12-26 02:37:24.290876+00', NULL, 'HYPE', 'HYPE', '0x1a971E327F094BbC86B0B47238d72b6E28988182', 18, 'https://img.sonex.so/hype-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "", "reddit": "", "discord": "https://discord.com/invite/hyperliquid", "twitter": "https://twitter.com/HyperliquidX", "bitclout": "", "telegram": "https://t.me/hyperliquid_announcements"}, "website": "https://hyperliquid.xyz/", "whitepaper": "https://hyperliquid.gitbook.io/hyperliquid-docs", "description": "Hyperliquid represents a cutting-edge blockchain platform, specifically designed to enhance the efficiency and performance of decentralized finance (DeFi) applications. At its core, Hyperliquid introduces a novel Layer 1 (L1) blockchain that is meticulously engineered from the ground up to optimize performance and scalability. This optimization is largely attributed to its proprietary consensus mechanism, HyperBFT, which plays a pivotal role in ensuring rapid transaction finality and robust security measures."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (26, '2024-12-26 02:37:35.792948+00', '2024-12-26 02:37:35.792948+00', NULL, 'SUI', 'SUI', '0x4ecfd87B5fe41007669eF599ddCd91Df5761A34B', 18, 'https://img.sonex.so/sui-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "https://github.com/MystenLabs", "reddit": "", "discord": "https://discord.gg/sui", "twitter": "https://twitter.com/SuiNetwork", "bitclout": "", "telegram": ""}, "website": "https://sui.io/#", "whitepaper": "https://docs.sui.io/", "description": "Sui is a groundbreaking layer-1 blockchain platform designed to support the needs of global adoption by offering a secure, powerful, and scalable development platform. At its core, Sui leverages a novel object-centric data model and the secure Move programming language to address inefficiencies prevalent in existing blockchain architectures."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (9, '2024-12-26 02:36:45.045007+00', '2024-12-26 02:36:45.045007+00', NULL, 'PNUT', 'PNUT', '0xd742c38D302595Ca5974da3C956e1337C128b199', 18, 'https://img.sonex.so/pnut-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "", "reddit": "", "discord": "", "twitter": "https://twitter.com/pnutsolana", "bitclout": "", "telegram": "https://t.me/pnutportal"}, "website": "https://www.pnutsol.com/", "whitepaper": "", "description": "PEANUT FOREVER. #JusticeForPeanut. Peanut the Squirrel is a meme coin on Solana."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (18, '2024-12-26 02:37:16.405965+00', '2024-12-26 02:37:16.405965+00', NULL, 'MOODENG', 'MOODENG', '0x681f1fA48dD3754CFebA39d648EAf0D9C789dDBf', 18, 'https://img.sonex.so/moodeng-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "", "reddit": "", "discord": "", "twitter": "https://twitter.com/MooDengSOL", "bitclout": "", "telegram": ""}, "website": "https://www.moodengsol.com/", "whitepaper": "", "description": "An unofficial fan page / memecoin of Moo Deng on Solana"}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (29, '2024-12-26 02:37:42.833237+00', '2024-12-26 02:37:42.833237+00', NULL, 'DOGE', 'DOGE', '0xb6Edf7e6A8bCBDe1285A54485448336cC703eDb4', 18, 'https://img.sonex.so/doge-logo.png', 0, 0, 't', '0.000000000000000000', '{}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (19, '2024-12-26 02:37:18.885127+00', '2024-12-26 02:37:18.885127+00', NULL, 'GOUT', 'GOUT', '0x2B8D4beA97d8DF2499142912f7F60D33c849A88b', 18, 'https://img.sonex.so/gout-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "", "reddit": "", "discord": "", "twitter": "https://twitter.com/GOUTCOIN", "bitclout": "", "telegram": "https://t.me/GOUTCOIN"}, "website": "https://goutchain.com/", "whitepaper": "", "description": "GOUT COIN is a cryptocurrency project that has emerged with the aim of creating a unique ecosystem within the crypto space. The project was born out of a series of iconic moments and is designed to celebrate those who overcome challenges and achieve success. It is built on the Binance Smart Chain (BSC) and has a total token supply of 210 billion. The token distribution includes 115 billion for private sale, 30 billion each for the BNB pool and USDT pool, and 35 billion set aside for burning (black hole). GOUT COIN aims to be more than just a meme token. It represents a symbol for individuals who strive to overcome adversity and pursue victory. The economic model of GOUT COIN incorporates features such as LP staking with high Annual Percentage Rates (APR) and relatively low risk. Through its yield farming mechanism on platforms like PancakeSwap, users can stake tokens to earn rewards. Farms require users to stake two tokens to provide liquidity and in return, they receive either a liquidity position NFT or LP Tokens, which can then be staked in the Farm to earn GOUT while maintaining a position in other tokens. The project has already gained some traction with over 600 groups in 30+ countries and billions of exposures. It has a vision and mission to create a global MEME ecosystem where users can freely experience the intertwining of fun and wealth. The community aspect is emphasized, inviting people to join and participate in various activities to earn GOUT through social interaction, work, and investment. Additionally, the project offers 24/7 support to assist users, especially those new to the world of cryptocurrency and facing difficulties in the buying process. Overall, GOUT COIN is an evolving project within the crypto space with its own unique characteristics and goals."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (15, '2024-12-26 02:37:09.784499+00', '2024-12-26 02:37:09.784499+00', NULL, 'THE', 'THE', '0xe1e6a75038C468fd487ae86F8Cf7f54b55204DaD', 18, 'https://img.sonex.so/the-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "https://github.com/ThenafiBNB/THENA-Contracts", "reddit": "", "discord": "https://discord.gg/thena", "twitter": "https://twitter.com/ThenaFi_", "bitclout": "", "telegram": "https://t.me/Thena_Fi"}, "website": "https://thena.fi/", "whitepaper": "https://thena.gitbook.io/thena/", "description": "THENA is an innovative decentralized exchange (DEX) and a liquidity layer built on BNB Chain. THENA offers spot and perpetual trading for a wide range of assets.At the core of THENA''s value proposition is THE Liquidity Marketplace, which uses a novel ve(3,3) model that allows DeFi protocols to incentivize liquidity by directly influencing veTHE holders'' votes. By adjusting the amount of $THE rewards distributed for specific pools, protocol teams can develop dynamic liquidity management strategies, attracting more liquidity providers and boosting their platforms'' Total Value Locked (TVL).THENA leverages its unique features and adaptive liquidity management capabilities to position itself as the go-to liquidity layer for protocols, while offering users an all-inclusive platform for their DeFi needs."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (24, '2024-12-26 02:37:31.130509+00', '2024-12-26 02:37:31.130509+00', NULL, 'PEPE', 'PEPE', '0x3cDc77726062841d0Db29074898Ba0662aeaf054', 18, 'https://img.sonex.so/pepe-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "https://fairyproof.com/report/Pepe", "social": {"github": "", "reddit": "", "discord": "", "twitter": "https://twitter.com/pepecoineth", "bitclout": "", "telegram": ""}, "website": "https://www.pepe.vip/", "whitepaper": "", "description": "PEPE is a deflationary memecoin launched on Ethereum. The cryptocurrency was created as a tribute to the Pepe the Frog internet meme, created by Matt Furie, which gained popularity in the early 2000s.The project aims to capitalize on the popularity of meme coins, like Shiba Inu and Dogecoin, and strives to establish itself as one of the top meme-based cryptocurrencies. PEPE appeals to the cryptocurrency community by instituting a no-tax policy and being up-front about its lack of utility, keeping things pure and simple as a memecoin."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (27, '2024-12-26 02:37:38.090249+00', '2024-12-26 02:37:38.090249+00', NULL, 'SOL', 'SOL', '0x87b3F4475Ecf85D89D6CEBe4224a0820C8F7583c', 18, 'https://img.sonex.so/sol-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "https://github.com/solana-labs", "reddit": "https://reddit.com/r/solana", "discord": "https://solana.com/discord", "twitter": "https://twitter.com/solana", "bitclout": "", "telegram": "https://t.me/solana"}, "website": "https://solana.com/", "whitepaper": "", "description": "Solana is a highly functional open source project that banks on blockchain technology’s permissionless nature to provide decentralized finance (DeFi) solutions.The Solana protocol is designed to facilitate decentralized app (DApp) creation. It aims to improve scalability by introducing a proof-of-history (PoH) consensus combined with the underlying proof-of-stake (PoS) consensus of the blockchain."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (12, '2024-12-26 02:37:03.374861+00', '2024-12-26 02:37:03.374861+00', NULL, 'SATS', 'SATS', '0xD6C0D4e730A67622CfE539818EB796014B006D31', 18, 'https://img.sonex.so/sats-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "", "reddit": "", "discord": "", "twitter": "https://twitter.com/satsvipordinal", "bitclout": "", "telegram": "https://t.me/satsbrc20_community"}, "website": "https://satscoin.vip/", "whitepaper": "", "description": "SATS (Ordinals) emerges as a fascinating player in the cryptocurrency landscape, intricately linked to Bitcoin. Each SATS represents the smallest unit of Bitcoin, equivalent to 0.00000001 BTC, highlighting its foundational connection to the pioneering cryptocurrency. This relationship underscores its significance in the broader blockchain ecosystem.Diving into its technical aspects, SATS operates on the Ordinals protocol, a unique layer that allows for the inscription of data directly onto the Bitcoin blockchain. This capability enhances the utility of SATS, enabling it to support a variety of applications beyond simple transactions. The Ordinals protocol facilitates the creation of digital artifacts, offering a new dimension to blockchain technology."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (20, '2024-12-26 02:37:21.538566+00', '2024-12-26 02:37:21.538566+00', NULL, 'NEIRO', 'NEIRO', '0xDeDd12fbd3CFC2604a4D7CA7dBE5182BEaEAc82E', 18, 'https://img.sonex.so/neiro-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "", "reddit": "", "discord": "", "twitter": "https://twitter.com/neiro", "bitclout": "", "telegram": "https://t.me/ETHNEIROCTO"}, "website": "https://www.neiroeth.io/", "whitepaper": "", "description": "Neiro is newly-adopted Shiba Inu dog, by the same woman who owned worldwide famous Kabosus, the doge behind the Doge meme. This project pays homage to Neiro, it is the first Neiro token deployed on Ethereum and it is 100% owned and managed with love by its community."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (25, '2024-12-26 02:37:33.373805+00', '2024-12-26 02:37:33.373805+00', NULL, 'BABYDOGE', 'BABYDOGE', '0xC21120d2e488c6627571F54e4d9b963737dcceC2', 18, 'https://img.sonex.so/babydoge-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "https://cmc.certik-skynet.com/redirect?project=baby-doge-coin", "social": {"github": "", "reddit": "", "discord": "https://discord.gg/babydogecoin1", "twitter": "https://twitter.com/babydogecoin", "bitclout": "", "telegram": "https://t.me/babydogearmy"}, "website": "https://www.babydoge.com/", "whitepaper": "", "description": "BabyDoge is the leading community meme coin on the BNB chain. Baby Doge was created initially as a joke with a mission to help spread awareness of animal adoption. Tweeted by Elon Musk BabyDoge was propelled into pop culture as an adorable iconThe primary utility for Baby Doge is building a fun meme community to spread awareness of animal welfare and pet adoption. The BabyDoge ai image generator empowers the BabyDoge community to create ai babydoge characters ai memes, and ai images on demand without any design experience. Additionally users can pay in babydoge to mint their images into NFTs.BabyDoge is the only meme with an official world record for helping save dogs donating over 81,000 pounds of dog food.BabyDoge has its own Swap BabyDogeSwap.com where it has its own AMM, core farms, farming as a service, free token locker, burn portal, card, Ai image generator, NFTs and more!"}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (4, '2024-12-25 07:38:50.773354+00', '2024-12-25 07:38:50.773354+00', NULL, 'DAI', 'DAI', '0x435d6EE7b2Cc7a1e5cF5c8dc29E651441fC07062', 18, 'https://img.sonex.so/dai-logo.png', 0, 2, 'f', '0.000000000000000000', '{}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (14, '2024-12-26 02:37:07.658454+00', '2024-12-26 02:37:07.658454+00', NULL, 'SHIB', 'SHIB', '0x94e9997a50a3F9C2Bc2922313aFa2d2CB4580ceB', 18, 'https://img.sonex.so/shib-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "https://cmc.certik-skynet.com/redirect?project=shiba-swap", "social": {"github": "", "reddit": "https://reddit.com/r/SHIBArmy", "discord": "https://discord.com/invite/shibatoken", "twitter": "https://twitter.com/shibtoken", "bitclout": "", "telegram": "https://t.me/ShibaInu_Dogecoinkiller"}, "website": "https://shibatoken.com/", "whitepaper": "https://github.com/shytoshikusama/woofwoofpaper/raw/main/SHIBA_INU_WOOF_WOOF.pdf", "description": "SHIB is a decentralized, community-led currency held by millions across the globe. Since its inception in late 2020, the Ethereum-based SHIB token has grown to become a worldwide phenomenon, and is now accepted as a form of payment at hundreds of locations, either directly or through third-party intermediaries.Shiba Inu coin was created anonymously in August 2020 under the pseudonym Ryoshi. Ryoshi says about himself that he is a nobody and not important and that the efforts to unmask his identity, even if successful, would be underwhelming.This meme coin quickly gained speed and value as a community of investors was drawn in by the cute charm of the coin paired with headlines and tweets from personalities like Elon Musk and Vitalik Buterin."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (13, '2024-12-26 02:37:05.581869+00', '2024-12-26 02:37:05.581869+00', NULL, 'USUAL', 'USUAL', '0xD31dbafb3246B6b73E6129FA018322416779abf6', 18, 'https://img.sonex.so/usual-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "", "reddit": "", "discord": "https://discord.usual.money/", "twitter": "https://twitter.com/usualmoney", "bitclout": "", "telegram": "https://telegram.usual.money/"}, "website": "https://usual.money/", "whitepaper": "https://docs.usual.money/start-here/why-usual", "description": "Usual is a decentralized Fiat Stablecoin issuer that redistributes ownership and governance through the $USUAL token.$USUAL is the governance token powering the Usual protocol. It gives users ownership and governance over the infrastructure and treasury. Uniquely designed with an intrinsic value tied directly to the protocol''s revenue model, $USUAL drives the adoption and use of USD0, aligning incentives for contributors and fueling protocol growth. Its innovative distribution model sets the stage for new DeFi possibilities, accelerating ecosystem expansion and sustainable decentralization."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (31, '2024-12-27 05:05:37.175378+00', '2024-12-27 05:05:37.175378+00', NULL, 'ENA', 'ENA', '0x800eB8C0Ac31AD57E1F73B6446B9E8Fd2cc53616', 18, 'https://img.sonex.so/ena-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "https://github.com/ethena-labs", "reddit": "", "discord": "https://discord.com/ethena", "twitter": "https://twitter.com/ethena_labs", "bitclout": "", "telegram": "https://t.me/ethena_labs"}, "website": "https://www.ethena.fi/", "whitepaper": "https://ethena-labs.gitbook.io/ethena-labs/", "description": "Ethena is a synthetic dollar protocol built on Ethereum that will provide a crypto-native solution for money that is not reliant on traditional banking system infrastructure, alongside a globally accessible dollar-denominated savings instrument — the ''Internet Bond''."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (7, '2024-12-25 07:39:42.074067+00', '2024-12-25 07:39:42.074067+00', NULL, 'ETH', 'ETH', '0x0000000000000000000000000000000000000000', 18, 'https://img.sonex.so/eth-logo.png', 0, 0, 'f', '0.000000000000000000', '{"audits": "", "social": {"github": "https://github.com/ethereum/go-ethereum", "reddit": "https://reddit.com/r/ethereum", "discord": "", "twitter": "https://twitter.com/ethereum", "bitclout": "", "telegram": ""}, "website": "https://www.ethereum.org/", "whitepaper": "https://github.com/ethereum/wiki/wiki/White-Paper", "description": "Ethereum is a decentralized blockchain with smart contract functionality. Ether is the native cryptocurrency of the platform. Among cryptocurrencies, ether is second only to bitcoin in market capitalization. It is open-source software."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (28, '2024-12-26 02:37:40.3002+00', '2024-12-26 02:37:40.3002+00', NULL, 'AAVE', 'AAVE', '0xEDebe43f08534A7ED3F2F8B87cC48bd901bae181', 18, 'https://img.sonex.so/aave-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "https://cmc.certik-skynet.com/redirect?project=aave", "social": {"github": "https://github.com/aave", "reddit": "https://reddit.com/r/Aave_Official", "discord": "https://aave.com/discord", "twitter": "https://twitter.com/aave", "bitclout": "", "telegram": "https://t.me/Aavesome"}, "website": "https://aave.com/", "whitepaper": "https://github.com/aave/aave-protocol/blob/master/docs/Aave_Protocol_Whitepaper_v1_0.pdf", "description": "Aave is a decentralized finance protocol that allows people to lend and borrow crypto.Lenders earn interest by depositing digital assets into specially created liquidity pools. Borrowers can then use their crypto as collateral to take out a flash loan using this liquidity.Aave (which means “ghost” in Finnish) was originally known as ETHLend when it launched in November 2017, but the rebranding to Aave happened in September 2018. (This helps explain why this token’s ticker is so different from its name!)AAVE provides holders with discounted fees on the platform, and it also serves as a governance token — giving owners a say in the future development of the protocol."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (23, '2024-12-26 02:37:28.806999+00', '2024-12-26 02:37:28.806999+00', NULL, 'ORDI', 'ORDI', '0x18357d16aaD531175ffF4159227dFe7E45aAeA21', 18, 'https://img.sonex.so/ordi-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "", "reddit": "", "discord": "", "twitter": "https://x.com/search?q=%24ORDI", "bitclout": "", "telegram": ""}, "website": "", "whitepaper": "", "description": "The total amount of bitcoins is 21 million, and one can be subdivided into 100 million satoshi (satoshi is sat, the smallest unit of bitcoin), and there are 2100 trillion satoshi in total. Here we must first understand the inscription, which is created by writing content to satoshi using the Ordinals protocol. The inscription does not need to use a separate token, side chain or change Bitcoin; What the Ordinals protocol does is to write information to each Satoshi, such as text, pictures, audio and video, etc. Due to the size limit of the Bitcoin block, the main information for inscribe (that is, mint mint) is mainly text and pictures , in the form of NFT and token."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (22, '2024-12-26 02:37:26.649226+00', '2024-12-26 02:37:26.649226+00', NULL, 'WLD', 'WLD', '0xa1A6ff7E3dC805f760B1bf0410eBF7a9c152803a', 18, 'https://img.sonex.so/wld-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "", "social": {"github": "https://github.com/worldcoin", "reddit": "", "discord": "https://discord.gg/worldcoin", "twitter": "https://twitter.com/worldcoin", "bitclout": "", "telegram": ""}, "website": "https://worldcoin.org/", "whitepaper": "https://whitepaper.worldcoin.org/", "description": "The mission of the Worldcoin project is to build the world’s largest identity and financial network as a public utility, giving ownership to everyone.The WLD token is designed as a utility token with governance properties, empowering users by giving them a say over the future of the protocol. Beyond conventional one-token-one-vote governance mechanisms, the introduction of World ID paves the way for one-person-one-vote mechanisms. These two mechanisms can also be combined to enable new ways of governance. After launch, the Worldcoin Foundation will solicit proposals and work with the community on how World ID and the WLD token should interact in Worldcoin''s governance model (see the project''s whitepaper for details on governance and decentralization)."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (5, '2024-12-25 07:38:53.368318+00', '2024-12-25 07:38:53.368318+00', NULL, 'Bridged USDC', 'USDC.e', '0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391', 6, 'https://img.sonex.so/usdc-logo.png', 0, 3, 'f', '0.000000000000000000', '{"audits": "", "social": {"github": "", "reddit": "", "discord": "https://discord.com/invite/sonex", "twitter": "https://x.com/sonex_so", "bitclout": "", "telegram": "https://t.me/SONEXexchange"}, "website": "https://sonex.so/", "whitepaper": "", "description": "SONEX is your gateway to the Soneium ecosystem, redefining decentralized finance by delivering a unified, all-in-one DeFi hub for retail users. We’re breaking away from the reliance on centralized platforms to bring you true decentralization.Launching with liquidity provision and swap functionality, SONEX will evolve into a fully integrated DeFi hub."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (33, '2024-12-30 02:48:56.676739+00', '2024-12-30 02:48:56.676739+00', NULL, 'LINE', 'LINE', '0xf5A42BA277A29F981eB9FE8D9A5b9f41C9F90866', 18, NULL, 0, 0, 't', '0.000000000000000000', '{}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (8, '2024-12-26 02:35:39.715137+00', '2024-12-26 02:35:39.715137+00', NULL, 'KOMA', 'KOMA', '0xd8cfd279a6db8AB2F41593e711A36224b4C84F6f', 18, 'https://img.sonex.so/koma-logo.png', 0, 0, 't', '0.000000000000000000', '{"audits": "https://cyberscope.io/audits/ff-koma", "social": {"github": "", "reddit": "", "discord": "", "twitter": "https://twitter.com/komabnb", "bitclout": "", "telegram": "https://t.me/komabnb"}, "website": "https://komabnb.com/", "whitepaper": "", "description": "Koma, son of Shib and protector of BNB, is a dog themed token built around community driven decentralization and charity."}');
INSERT INTO "sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (34, '2024-12-26 02:35:39.715137+00', '2024-12-26 02:35:39.715137+00', NULL, 'SONE', 'SONE', '0xA5D6513082EF1F157A33A066293309E74A8aF6Df', 18, NULL, 0, 0, 't', '0.000000000000000000', '{}');

INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (44, 1946, 'GOUT', 'GOUT', 18, '0x2B8D4beA97d8DF2499142912f7F60D33c849A88b', 'https://img.sonex.so/gout-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (45, 1946, 'THE', 'THE', 18, '0xe1e6a75038C468fd487ae86F8Cf7f54b55204DaD', 'https://img.sonex.so/the-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (33, 1946, 'ACT', 'ACT', 18, '0x0a73297F79858Ad1d58e68B45A806D4e6C6D5459', 'https://img.sonex.so/act-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (34, 1946, 'XRP', 'XRP', 18, '0x795c68F07d3ced21fc357A76E077b871457A74A2', 'https://img.sonex.so/xrp-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (35, 1946, 'Astar Minato Testnet Token', 'ASTR', 18, '0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655', 'https://img.sonex.so/astar-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (36, 1946, 'BNB', 'BNB', 18, '0x13aA369B93343e7b89cFCae9725A2f22895244c6', 'https://img.sonex.so/bnb-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (37, 1946, 'MOVE', 'MOVE', 18, '0x3FCd768bAdFE17b8A5e878109657F0eaC0e4a168', 'https://img.sonex.so/move-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (38, 1946, 'ME', 'ME', 18, '0xF533784b6dD49C363cfF778067EC5983292B83e7', 'https://img.sonex.so/me-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (39, 1946, 'HYPE', 'HYPE', 18, '0x1a971E327F094BbC86B0B47238d72b6E28988182', 'https://img.sonex.so/hype-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (40, 1946, 'SUI', 'SUI', 18, '0x4ecfd87B5fe41007669eF599ddCd91Df5761A34B', 'https://img.sonex.so/sui-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (41, 1946, 'PNUT', 'PNUT', 18, '0xd742c38D302595Ca5974da3C956e1337C128b199', 'https://img.sonex.so/pnut-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (42, 1946, 'MOODENG', 'MOODENG', 18, '0x681f1fA48dD3754CFebA39d648EAf0D9C789dDBf', 'https://img.sonex.so/moodeng-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (43, 1946, 'DOGE', 'DOGE', 18, '0xb6Edf7e6A8bCBDe1285A54485448336cC703eDb4', 'https://img.sonex.so/doge-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (47, 1946, 'SOL', 'SOL', 18, '0x87b3F4475Ecf85D89D6CEBe4224a0820C8F7583c', 'https://img.sonex.so/sol-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (48, 1946, 'SATS', 'SATS', 18, '0xD6C0D4e730A67622CfE539818EB796014B006D31', 'https://img.sonex.so/sats-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (49, 1946, 'NEIRO', 'NEIRO', 18, '0xDeDd12fbd3CFC2604a4D7CA7dBE5182BEaEAc82E', 'https://img.sonex.so/neiro-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (50, 1946, 'BABYDOGE', 'BABYDOGE', 18, '0xC21120d2e488c6627571F54e4d9b963737dcceC2', 'https://img.sonex.so/babydoge-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (51, 1946, 'DAI', 'DAI', 18, '0x435d6EE7b2Cc7a1e5cF5c8dc29E651441fC07062', 'https://img.sonex.so/dai-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (52, 1946, 'SHIB', 'SHIB', 18, '0x94e9997a50a3F9C2Bc2922313aFa2d2CB4580ceB', 'https://img.sonex.so/shib-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (53, 1946, 'USUAL', 'USUAL', 18, '0xD31dbafb3246B6b73E6129FA018322416779abf6', 'https://img.sonex.so/usual-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (54, 1946, 'ENA', 'ENA', 18, '0x800eB8C0Ac31AD57E1F73B6446B9E8Fd2cc53616', 'https://img.sonex.so/ena-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (55, 1946, 'ETH', 'ETH', 18, '0x0000000000000000000000000000000000000000', 'https://img.sonex.so/eth-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (56, 1946, 'AAVE', 'AAVE', 18, '0xEDebe43f08534A7ED3F2F8B87cC48bd901bae181', 'https://img.sonex.so/aave-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (57, 1946, 'ORDI', 'ORDI', 18, '0x18357d16aaD531175ffF4159227dFe7E45aAeA21', 'https://img.sonex.so/ordi-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (58, 1946, 'WLD', 'WLD', 18, '0xa1A6ff7E3dC805f760B1bf0410eBF7a9c152803a', 'https://img.sonex.so/wld-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (59, 1946, 'Bridged USDC', 'USDC.e', 6, '0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391', 'https://img.sonex.so/usdc-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (60, 1946, 'KOMA', 'KOMA', 18, '0xd8cfd279a6db8AB2F41593e711A36224b4C84F6f', 'https://img.sonex.so/koma-logo.png', 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (61, 1946, 'LINE', 'LINE', 18, '0xf5A42BA277A29F981eB9FE8D9A5b9f41C9F90866', NULL, 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (62, 1946, 'Wrapped Ethereum', 'WETH', 18, '0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C', 'https://img.sonex.so/weth-logo.png', 999, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (63, 1946, 'Sonex Token', 'SONEX', 18, '0x8879309d4002b8C98163Ece3f58653F774de8416', 'https://img.sonex.so/sonex-logo.png', 1000, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (64, 1946, 'SONE', 'SONE', 18, '0xA5D6513082EF1F157A33A066293309E74A8aF6Df', NULL, 0, 'f');
INSERT INTO "token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (46, 1946, 'PEPE', 'PEPE', 18, '0x3cDc77726062841d0Db29074898Ba0662aeaf054', 'https://img.sonex.so/pepe-logo.png', 0, 'f');

