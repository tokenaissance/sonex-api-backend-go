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

 Date: 21/01/2025 16:04:15
*/


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
  "event_incentive_id" int8 NOT NULL,
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
  "event_account" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_farm_reward_claimed
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_reward_claimed";
CREATE TABLE "public"."event_farm_reward_claimed" (
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
  "event_old_rate_per_second" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_new_rate_per_second" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_timestamp" int8 NOT NULL
)
;

-- ----------------------------
-- Table structure for event_farm_token_deposited
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_token_deposited";
CREATE TABLE "public"."event_farm_token_deposited" (
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
  "event_owner" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_farm_token_staked
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_token_staked";
CREATE TABLE "public"."event_farm_token_staked" (
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
  "event_token_id" int8 NOT NULL,
  "event_incentive_id" int8 NOT NULL,
  "event_total_liquidity" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_farm_token_unstaked
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_token_unstaked";
CREATE TABLE "public"."event_farm_token_unstaked" (
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
  "event_account" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Indexes structure for table event_farm_incentive_created
-- ----------------------------
CREATE UNIQUE INDEX "event_farm_incentive_created_chain_id_incentive_id" ON "public"."event_farm_incentive_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_incentive_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "farm_incentive_created_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_incentive_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

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
CREATE UNIQUE INDEX "farm_incentive_ended_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_incentive_ended" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_farm_incentive_ended
-- ----------------------------
ALTER TABLE "public"."event_farm_incentive_ended" ADD CONSTRAINT "event_farm_incentive_ended_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_incentive_updated
-- ----------------------------
CREATE INDEX "farm_incentive_updated_index_chain_id_incentive_id" ON "public"."event_farm_incentive_updated" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_incentive_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "farm_incentive_updated_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_incentive_updated" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_farm_incentive_updated
-- ----------------------------
ALTER TABLE "public"."event_farm_incentive_updated" ADD CONSTRAINT "event_farm_incentive_updated_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_paused
-- ----------------------------
CREATE UNIQUE INDEX "event_farm_paused_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_paused" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_farm_paused
-- ----------------------------
ALTER TABLE "public"."event_farm_paused" ADD CONSTRAINT "event_farm_paused_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_reward_claimed
-- ----------------------------
CREATE INDEX "event_farm_reward_claimed_index_chain_id_incentive_id" ON "public"."event_farm_reward_claimed" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_incentive_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "event_farm_reward_claimed_index_chain_id_user" ON "public"."event_farm_reward_claimed" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_user" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "farm_reward_claimed_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_reward_claimed" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_farm_reward_claimed
-- ----------------------------
ALTER TABLE "public"."event_farm_reward_claimed" ADD CONSTRAINT "event_farm_reward_claimed_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_reward_pool_updated
-- ----------------------------
CREATE UNIQUE INDEX "farm_reward_pool_updated_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_reward_pool_updated" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_farm_reward_pool_updated
-- ----------------------------
ALTER TABLE "public"."event_farm_reward_pool_updated" ADD CONSTRAINT "event_farm_reward_pool_updated_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_reward_rate_updated
-- ----------------------------
CREATE UNIQUE INDEX "farm_reward_rate_updated_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_reward_rate_updated" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_farm_reward_rate_updated
-- ----------------------------
ALTER TABLE "public"."event_farm_reward_rate_updated" ADD CONSTRAINT "event_farm_reward_rate_updated_pkey" PRIMARY KEY ("id");

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
CREATE UNIQUE INDEX "farm_token_deposited_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_token_deposited" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

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
CREATE INDEX "event_farm_token_staked_index_chain_id_user" ON "public"."event_farm_token_staked" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_user" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "farm_token_staked_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_token_staked" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

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
CREATE INDEX "event_farm_token_unstaked_index_chain_id_user" ON "public"."event_farm_token_unstaked" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "event_user" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "farm_token_unstaked_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_token_unstaked" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_farm_token_unstaked
-- ----------------------------
ALTER TABLE "public"."event_farm_token_unstaked" ADD CONSTRAINT "event_farm_token_unstaked_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_farm_unpaused
-- ----------------------------
CREATE UNIQUE INDEX "farm_unpaused_index_chain_id_tx_hash_tx_index" ON "public"."event_farm_unpaused" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_farm_unpaused
-- ----------------------------
ALTER TABLE "public"."event_farm_unpaused" ADD CONSTRAINT "event_farm_unpaused_pkey" PRIMARY KEY ("id");
