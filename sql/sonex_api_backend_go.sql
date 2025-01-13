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
-- Records of chain_scanner_prograss
-- ----------------------------
BEGIN;
INSERT INTO "public"."chain_scanner_prograss" ("id", "created_at", "updated_at", "chain_id", "height") VALUES (124502484323401728, '2024-12-09 00:00:00+00', '2024-12-09 00:00:00+00', 1946, 6500000);
COMMIT;

-- ----------------------------
-- Records of sonex_token
-- ----------------------------
BEGIN;
INSERT INTO "public"."sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (2, '2025-01-09 03:41:47.6758+00', '2025-01-09 03:41:47.6758+00', NULL, 'ETH', 'ETH', '0x0000000000000000000000000000000000000000', 18, 'https://img.sonex.so/eth-logo.png', 0, 0, 'f', 0.000000000000000000, '{"audits": "", "social": {"github": "https://github.com/ethereum/go-ethereum", "reddit": "https://reddit.com/r/ethereum", "discord": "", "twitter": "https://twitter.com/ethereum", "bitclout": "", "telegram": ""}, "website": "https://www.ethereum.org/", "whitepaper": "https://github.com/ethereum/wiki/wiki/White-Paper", "description": "Ethereum is a decentralized blockchain with smart contract functionality. Ether is the native cryptocurrency of the platform. Among cryptocurrencies, ether is second only to bitcoin in market capitalization. It is open-source software."}');
INSERT INTO "public"."sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (3, '2025-01-09 03:42:14.659032+00', '2025-01-09 03:42:14.659032+00', NULL, 'Wrapped Ethereum', 'WETH', '0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C', 18, 'https://img.sonex.so/weth-logo.png', 0, 0, 'f', 0.000000000000000000, '{"audits": "https://hacken.io/audits/singularitydao/", "social": {"github": "", "reddit": "", "discord": "", "twitter": "", "bitclout": "", "telegram": ""}, "website": "", "whitepaper": "", "description": "What is WETH (Wrapped ETH)? WETH is the tokenized/packaged form of ETH that you use to pay for items when you interact with Ethereum dApps. WETH follows the ERC-20 token standards, enabling it to achieve interoperability with other ERC-20 tokens.  This offers more utility to holders as they can use it across networks and dApps. You can stake, yield farm, lend, and provide liquidity to various liquidity pools with WETH.  Also, unlike ETH, which doesn’t conform to its own ERC-20 standard and thus has lower interoperability as it can’t be used on other chains besides Ethereum, WETH can be used on cheaper and high throughput alternatives like Binance, Polygon, Solana, and Cardano. The price of WETH will always be the same as ETH because it maintains a 1:1 wrapping ratio."}');
INSERT INTO "public"."sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (4, '2025-01-09 03:43:02.235669+00', '2025-01-09 03:43:02.235669+00', NULL, 'Bridged USDC', 'USDC.e', '0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391', 6, 'https://img.sonex.so/usdc-logo.png', 0, 3, 'f', 0.000000000000000000, '{"audits": "https://certificate.quantstamp.com/view/a-star-network-custom-signature", "social": {"github": "https://github.com/AstarNetwork/Astar", "reddit": "", "discord": "https://discord.gg/astarnetwork", "twitter": "https://twitter.com/AstarNetwork", "bitclout": "", "telegram": "https://t.me/PlasmAnnouncement"}, "website": "https://astar.network/", "whitepaper": "https://docs.astar.network/", "description": "Astar Network helps developers build decentralized applications (dApps) and layer 2 solutions by providing interoperable web 3.0 infrastructure, financial incentives, incubation programmes and technical support. It seeks to provide the best solution for all developers by supporting EVM and making a parachain where EVM and WASM smart contracts can co-exist and communicate with each other."}');
INSERT INTO "public"."sonex_token" ("id", "created_at", "updated_at", "deleted_at", "name", "symbol", "address", "decimals", "logo", "sort", "recommendOrder", "isHidden", "volume", "additional") VALUES (5, '2025-01-09 03:43:10.847477+00', '2025-01-09 03:43:10.847477+00', NULL, 'Astar Token', 'ASTR', '0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655', 18, 'https://img.sonex.so/astar-logo.png', 0, 1, 'f', 0.000000000000000000, '{"audits": "https://certificate.quantstamp.com/view/a-star-network-custom-signature", "social": {"github": "https://github.com/AstarNetwork/Astar", "reddit": "", "discord": "https://discord.gg/astarnetwork", "twitter": "https://twitter.com/AstarNetwork", "bitclout": "", "telegram": "https://t.me/PlasmAnnouncement"}, "website": "https://astar.network/", "whitepaper": "https://docs.astar.network/", "description": "Astar Network helps developers build decentralized applications (dApps) and layer 2 solutions by providing interoperable web 3.0 infrastructure, financial incentives, incubation programmes and technical support. It seeks to provide the best solution for all developers by supporting EVM and making a parachain where EVM and WASM smart contracts can co-exist and communicate with each other."}');
COMMIT;
-- ----------------------------
-- Records of token
-- ----------------------------
BEGIN;
INSERT INTO "public"."token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (3, 1946, 'Wrapped Ethereum', 'WETH', 18, '0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C', 'https://img.sonex.so/weth-logo.png', 999, 'f');
INSERT INTO "public"."token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (2, 1946, 'ETH', 'ETH', 18, '0x0000000000000000000000000000000000000000', 'https://img.sonex.so/eth-logo.png', 1000, 'f');
INSERT INTO "public"."token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (4, 1946, 'Bridged USDC', 'USDC.e', 6, '0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391', 'https://img.sonex.so/usdc-logo.png', 998, 'f');
INSERT INTO "public"."token" ("id", "chain_id", "name", "symbol", "decimals", "address", "logo", "sort", "hidden") VALUES (5, 1946, 'Astar Token', 'ASTR', 18, '0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655', 'https://img.sonex.so/astar-logo.png', 900, 'f');
COMMIT;