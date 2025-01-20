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

 Date: 20/01/2025 17:06:10
*/


-- ----------------------------
-- Table structure for event_xpool_factory_operator_added
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_xpool_factory_operator_added";
CREATE TABLE "public"."event_xpool_factory_operator_added" (
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
  "event_operator" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for event_xpool_factory_operator_removed
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_xpool_factory_operator_removed";
CREATE TABLE "public"."event_xpool_factory_operator_removed" (
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
  "event_operator" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
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
ALTER TABLE "public"."event_xpool_factory_pool_destroied" OWNER TO "postgres";

-- ----------------------------
-- Table structure for event_xpool_factory_reward_per_block_set
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_xpool_factory_reward_per_block_set";
CREATE TABLE "public"."event_xpool_factory_reward_per_block_set" (
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
  "event_old" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_new" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
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
-- Indexes structure for table event_xpool_factory_operator_added
-- ----------------------------
CREATE UNIQUE INDEX "xpool_factory_operator_added_index_chain_id_tx_hash_tx_index" ON "public"."event_xpool_factory_operator_added" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_xpool_factory_operator_added
-- ----------------------------
ALTER TABLE "public"."event_xpool_factory_operator_added" ADD CONSTRAINT "event_xpool_factory_operator_added_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_xpool_factory_operator_removed
-- ----------------------------
CREATE UNIQUE INDEX "xpool_factory_operator_removed_index_chain_id_tx_hash_tx_index" ON "public"."event_xpool_factory_operator_removed" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_xpool_factory_operator_removed
-- ----------------------------
ALTER TABLE "public"."event_xpool_factory_operator_removed" ADD CONSTRAINT "event_xpool_factory_operator_removed_pkey" PRIMARY KEY ("id");

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
ALTER TABLE "public"."event_xpool_factory_pool_created" ADD CONSTRAINT "event_xpool_factory_pool_created_pkey" PRIMARY KEY ("id");

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
ALTER TABLE "public"."event_xpool_factory_pool_destroied" ADD CONSTRAINT "event_xpool_factory_pool_destroied_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table event_xpool_factory_reward_per_block_set
-- ----------------------------
CREATE UNIQUE INDEX "factory_reward_per_block_set_index_chain_id_tx_hash_tx_index" ON "public"."event_xpool_factory_reward_per_block_set" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tx_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table event_xpool_factory_reward_per_block_set
-- ----------------------------
ALTER TABLE "public"."event_xpool_factory_reward_per_block_set" ADD CONSTRAINT "event_xpool_pool_reward_per_block_set_pkey" PRIMARY KEY ("id");

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
ALTER TABLE "public"."event_xpool_pool_claimed" ADD CONSTRAINT "event_xpool_pool_claimed_pkey" PRIMARY KEY ("id");

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
ALTER TABLE "public"."event_xpool_pool_withdrawed" ADD CONSTRAINT "event_xpool_pool_withdrawed_pkey" PRIMARY KEY ("id");
