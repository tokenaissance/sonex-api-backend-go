-- ----------------------------
-- Table structure for event_farm_incentive_created
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_incentive_created";
-- 创建质押池事件表
CREATE TABLE "public"."event_farm_incentive_created" (
  "id" int8 PRIMARY KEY,
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
  "event_pool" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_start_time" int8 NOT NULL,
  "event_end_time" int8 NOT NULL,
  "event_refundee" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_reward" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_incentive_id" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  UNIQUE ("tx_hash","chain_id")
);

DROP SEQUENCE IF EXISTS "event_farm_incentive_created_id_seq";
CREATE SEQUENCE IF NOT EXISTS "event_farm_incentive_created_id_seq" START WITH 1;
ALTER TABLE "public"."event_farm_incentive_created" ALTER COLUMN id SET DEFAULT nextval('event_farm_incentive_created_id_seq');

COMMENT ON TABLE "public"."event_farm_incentive_created" IS '创建质押池事件表';
COMMENT ON COLUMN "public"."event_farm_incentive_created"."event_reward_token" IS '奖励代币地址';
COMMENT ON COLUMN "public"."event_farm_incentive_created"."event_pool" IS 'V3池子地址';
COMMENT ON COLUMN "public"."event_farm_incentive_created"."event_start_time" IS '质押池开放时间';
COMMENT ON COLUMN "public"."event_farm_incentive_created"."event_end_time" IS '质押池结束时间';
COMMENT ON COLUMN "public"."event_farm_incentive_created"."event_refundee" IS '剩余奖励代币接收地址';
COMMENT ON COLUMN "public"."event_farm_incentive_created"."event_reward" IS '奖励代币数量';
COMMENT ON COLUMN "public"."event_farm_incentive_created"."event_incentive_id" IS '质押池ID';


-- ----------------------------
-- Table structure for event_farm_token_staked
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_token_staked";
-- 质押事件表
CREATE TABLE "public"."event_farm_token_staked" (
  "id" int8 PRIMARY KEY,
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
  "event_incentive_id" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_liquidity" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  UNIQUE ("tx_hash","chain_id")
);

DROP SEQUENCE IF EXISTS "event_farm_token_staked_id_seq";
CREATE SEQUENCE IF NOT EXISTS "event_farm_token_staked_id_seq" START WITH 1;
ALTER TABLE "public"."event_farm_token_staked" ALTER COLUMN id SET DEFAULT nextval('event_farm_token_staked_id_seq');

COMMENT ON TABLE "public"."event_farm_token_staked" IS '质押事件表';
COMMENT ON COLUMN "public"."event_farm_token_staked"."event_token_id" IS '质押的LP Token ID';
COMMENT ON COLUMN "public"."event_farm_token_staked"."event_liquidity" IS '质押的流动性数量';


-- ----------------------------
-- Table structure for event_farm_token_unstaked
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_token_unstaked";
-- 解除质押事件表
CREATE TABLE "public"."event_farm_token_unstaked" (
  "id" int8 PRIMARY KEY,
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
  "event_incentive_id" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  UNIQUE ("tx_hash","chain_id")
);

DROP SEQUENCE IF EXISTS "event_farm_token_unstaked_id_seq";
CREATE SEQUENCE IF NOT EXISTS "event_farm_token_unstaked_id_seq" START WITH 1;
ALTER TABLE "public"."event_farm_token_unstaked" ALTER COLUMN id SET DEFAULT nextval('event_farm_token_unstaked_id_seq');

COMMENT ON TABLE "public"."event_farm_token_unstaked" IS '解除质押事件表';
COMMENT ON COLUMN "public"."event_farm_token_unstaked"."event_token_id" IS '解除质押的LP Token ID';


-- ----------------------------
-- Table structure for event_farm_reward_claimed
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_reward_claimed";
-- 领取奖励事件表
CREATE TABLE "public"."event_farm_reward_claimed" (
  "id" int8 PRIMARY KEY,
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
  "event_to" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_reward" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  UNIQUE ("tx_hash","chain_id")
);

DROP SEQUENCE IF EXISTS "event_farm_reward_claimed_id_seq";
CREATE SEQUENCE IF NOT EXISTS "event_farm_reward_claimed_id_seq" START WITH 1;
ALTER TABLE "public"."event_farm_reward_claimed" ALTER COLUMN id SET DEFAULT nextval('event_farm_reward_claimed_id_seq');

COMMENT ON TABLE "public"."event_farm_reward_claimed" IS '领取奖励事件表';
COMMENT ON COLUMN "public"."event_farm_reward_claimed"."event_to" IS '接收奖励的地址';
COMMENT ON COLUMN "public"."event_farm_reward_claimed"."event_reward" IS '提取的奖励数量';


-- ----------------------------
-- Table structure for event_farm_deposit_transferred
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_deposit_transferred";
-- LP Token所有权转移事件表
CREATE TABLE "public"."event_farm_deposit_transferred" (
  "id" int8 PRIMARY KEY,
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
  "event_old_owner" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_new_owner" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  UNIQUE ("tx_hash","chain_id")
);

DROP SEQUENCE IF EXISTS "event_farm_deposit_transferred_id_seq";
CREATE SEQUENCE IF NOT EXISTS "event_farm_deposit_transferred_id_seq" START WITH 1;
ALTER TABLE "public"."event_farm_deposit_transferred" ALTER COLUMN id SET DEFAULT nextval('event_farm_deposit_transferred_id_seq');

COMMENT ON TABLE "public"."event_farm_deposit_transferred" IS 'LP Token所有权转移事件表';
COMMENT ON COLUMN "public"."event_farm_deposit_transferred"."event_token_id" IS 'LP Token ID';
COMMENT ON COLUMN "public"."event_farm_deposit_transferred"."event_old_owner" IS '原Owner地址';
COMMENT ON COLUMN "public"."event_farm_deposit_transferred"."event_new_owner" IS '新Owner地址';

-- ----------------------------
-- Table structure for event_farm_incentive_ended
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_farm_incentive_ended";
-- 质押池结束事件表
CREATE TABLE "public"."event_farm_incentive_ended" (
  "id" int8 PRIMARY KEY,
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
  "event_refund" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  UNIQUE ("tx_hash","chain_id")
);

DROP SEQUENCE IF EXISTS "event_farm_incentive_ended_id_seq";
CREATE SEQUENCE IF NOT EXISTS "event_farm_incentive_ended_id_seq" START WITH 1;
ALTER TABLE "public"."event_farm_incentive_ended" ALTER COLUMN id SET DEFAULT nextval('event_farm_incentive_ended_id_seq');

COMMENT ON TABLE "public"."event_farm_incentive_ended" IS '质押池结束事件表';
COMMENT ON COLUMN "public"."event_farm_incentive_ended"."event_incentive_id" IS '质押池ID';
COMMENT ON COLUMN "public"."event_farm_incentive_ended"."event_refund" IS '剩余奖励数量';





