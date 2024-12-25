DROP TABLE IF EXISTS "public"."event_farm_incentive_created";
CREATE TABLE "public"."event_farm_incentive_created" (
  "id" SERIAL PRIMARY KEY,
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

COMMENT ON COLUMN "public"."event_farm_incentive_created"."event_reward_token" IS '奖励代币地址';
COMMENT ON COLUMN "public"."event_farm_incentive_created"."event_pool" IS 'Uniswap V3资金池地址';
COMMENT ON COLUMN "public"."event_farm_incentive_created"."event_start_time" IS '激励计划开始时间';
COMMENT ON COLUMN "public"."event_farm_incentive_created"."event_end_time" IS '激励计划结束时间';
COMMENT ON COLUMN "public"."event_farm_incentive_created"."event_refundee" IS '剩余奖励代币接收地址';
COMMENT ON COLUMN "public"."event_farm_incentive_created"."event_reward" IS '奖励代币数量';
COMMENT ON COLUMN "public"."event_farm_incentive_created"."event_incentive_id" IS '激励计划ID';


DROP TABLE IF EXISTS "public"."event_farm_token_staked";
CREATE TABLE "public"."event_farm_token_staked" (
  "id" SERIAL PRIMARY KEY,
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

COMMENT ON COLUMN "public"."event_farm_token_staked"."event_token_id" IS '质押的LP Token ID';
COMMENT ON COLUMN "public"."event_farm_token_staked"."event_liquidity" IS '质押的流动性数量';

DROP TABLE IF EXISTS "public"."event_farm_token_unstaked";
CREATE TABLE "public"."event_farm_token_unstaked" (
  "id" SERIAL PRIMARY KEY,
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

COMMENT ON COLUMN "public"."event_farm_token_unstaked"."event_token_id" IS '解除质押的LP Token ID';

DROP TABLE IF EXISTS "public"."event_farm_reward_claimed";
CREATE TABLE "public"."event_farm_reward_claimed" (
  "id" SERIAL PRIMARY KEY,
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

COMMENT ON COLUMN "public"."event_farm_reward_claimed"."event_to" IS '接收奖励的地址';
COMMENT ON COLUMN "public"."event_farm_reward_claimed"."event_reward" IS '提取的奖励数量';