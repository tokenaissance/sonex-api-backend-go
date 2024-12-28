
-- ----------------------------
-- Table structure for farm_pools
-- ----------------------------
DROP TABLE IF EXISTS "public"."farm_pools";
-- Farm质押池表
CREATE TABLE farm_pools (
    "id" int8 PRIMARY KEY,
    "chain_id" BIGINT NOT NULL DEFAULT 0,
    "incentive_id" VARCHAR(66) NOT NULL DEFAULT '',
    "start_time" TIMESTAMP, 
    "end_time" TIMESTAMP, 
    "reward_token_address" VARCHAR(42) NOT NULL DEFAULT '',
    "reward_token_decimals" SMALLINT NOT NULL DEFAULT 0, 
    "reward_token_amount" VARCHAR(50) NOT NULL DEFAULT '', 
    "reward_token_price" VARCHAR(50) NOT NULL DEFAULT '', 
    "refundee_address" VARCHAR(42) NOT NULL DEFAULT '', 
    "pool_address" VARCHAR(42) NOT NULL DEFAULT '', 
    "fee_ratio" VARCHAR(10) NOT NULL DEFAULT 0.0000, 
    "token0_address" VARCHAR(42) NOT NULL DEFAULT '',
    "token0_symbol" VARCHAR(10) NOT NULL DEFAULT '', 
    "token0_decimals" SMALLINT NOT NULL DEFAULT 0, 
    "token1_address" VARCHAR(42) NOT NULL DEFAULT '', 
    "token1_symbol" VARCHAR(10) NOT NULL DEFAULT '', 
    "token1_decimals" SMALLINT NOT NULL DEFAULT 0, 
    "staked_liquidity" VARCHAR(50) NOT NULL DEFAULT '',
    "staked_liquidity_usd" VARCHAR(50) NOT NULL DEFAULT '', 
    "total_token0_amount" VARCHAR(50) NOT NULL DEFAULT '', 
    "total_token1_amount" VARCHAR(50) NOT NULL DEFAULT '', 
    "apr" VARCHAR(50) NOT NULL DEFAULT '', 
    "join_num" INTEGER NOT NULL DEFAULT 0, 
    "tx_hash" varchar(255) NOT NULL DEFAULT '', 
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    "updated_at" TIMESTAMP NULL DEFAULT NULL,
    UNIQUE ("incentive_id", "chain_id")
);

DROP SEQUENCE IF EXISTS "farm_pools_id_seq";
CREATE SEQUENCE IF NOT EXISTS "farm_pools_id_seq" START WITH 1;
ALTER TABLE "public"."farm_pools" ALTER COLUMN id SET DEFAULT nextval('farm_pools_id_seq');


COMMENT ON TABLE farm_pools IS 'Farm质押池表';
COMMENT ON COLUMN "farm_pools"."chain_id" IS '链ID';
COMMENT ON COLUMN "farm_pools"."incentive_id" IS '质押池ID';
COMMENT ON COLUMN "farm_pools"."start_time" IS '质押池开放时间';
COMMENT ON COLUMN "farm_pools"."end_time" IS '质押池结束时间';
COMMENT ON COLUMN "farm_pools"."reward_token_address" IS '奖励token代币地址';
COMMENT ON COLUMN "farm_pools"."reward_token_decimals" IS '奖励token代币精度';
COMMENT ON COLUMN "farm_pools"."reward_token_amount" IS '奖励token数量';
COMMENT ON COLUMN "farm_pools"."reward_token_price" IS '奖励token价格';
COMMENT ON COLUMN "farm_pools"."refundee_address" IS '项目方回收奖励代币地址';
COMMENT ON COLUMN "farm_pools"."pool_address" IS 'v3池子合约地址';
COMMENT ON COLUMN "farm_pools"."fee_ratio" IS 'v3池子手续费率';
COMMENT ON COLUMN "farm_pools"."token0_address" IS 'token0代币地址';
COMMENT ON COLUMN "farm_pools"."token0_symbol" IS 'token0符号';
COMMENT ON COLUMN "farm_pools"."token0_decimals" IS 'token0精度';
COMMENT ON COLUMN "farm_pools"."token1_address" IS 'token1代币地址';
COMMENT ON COLUMN "farm_pools"."token1_symbol" IS 'token1符号';
COMMENT ON COLUMN "farm_pools"."token1_decimals" IS 'token1精度';
COMMENT ON COLUMN "farm_pools"."staked_liquidity" IS '质押池质押的流动性';
COMMENT ON COLUMN "farm_pools"."staked_liquidity_usd" IS '质押池质押的流动性的价值（USD）';
COMMENT ON COLUMN "farm_pools"."total_token0_amount" IS 'token0质押总量';
COMMENT ON COLUMN "farm_pools"."total_token1_amount" IS 'token1质押总量';
COMMENT ON COLUMN "farm_pools"."apr" IS '质押池年化收益率';
COMMENT ON COLUMN "farm_pools"."join_num" IS '质押池参与质押人数';
COMMENT ON COLUMN "farm_pools"."tx_hash" IS '质押池创建时的交易哈希';
COMMENT ON COLUMN "farm_pools"."created_at" IS '创建时间';
COMMENT ON COLUMN "farm_pools"."updated_at" IS '更新时间';

-- ----------------------------
-- Table structure for farm_user_pools
-- ----------------------------
DROP TABLE IF EXISTS "public"."farm_user_pools";
-- Farm用户质押池关联表
CREATE TABLE "public"."farm_user_pools" (
  "id" int8 PRIMARY KEY,
  "incentive_id" VARCHAR(66) NOT NULL DEFAULT '',
  "token_id" BIGINT NOT NULL DEFAULT 0,
  "user_address" VARCHAR(42) NOT NULL DEFAULT '',
  "chain_id" BIGINT NOT NULL DEFAULT 0,
  "status" VARCHAR(10) NOT NULL CHECK (status IN ('stake', 'withdraw')),
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NULL DEFAULT NULL,
  UNIQUE ("user_address", "incentive_id", "token_id",  "chain_id", "status")
);

DROP SEQUENCE IF EXISTS "farm_user_pools_id_seq";
CREATE SEQUENCE IF NOT EXISTS "farm_user_pools_id_seq" START WITH 1;
ALTER TABLE "public"."farm_user_pools" ALTER COLUMN id SET DEFAULT nextval('farm_user_pools_id_seq');

COMMENT ON TABLE "public"."farm_user_pools" IS 'Farm用户质押池关联表';
COMMENT ON COLUMN "public"."farm_user_pools"."incentive_id" IS '质押池ID';
COMMENT ON COLUMN "public"."farm_user_pools"."token_id" IS 'LP Token ID';
COMMENT ON COLUMN "public"."farm_user_pools"."user_address" IS '用户地址';
COMMENT ON COLUMN "public"."farm_user_pools"."chain_id" IS '链ID';
COMMENT ON COLUMN "public"."farm_user_pools"."status" IS '状态(stake:质押中 unstake:已解除)';
COMMENT ON COLUMN "public"."farm_user_pools"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."farm_user_pools"."updated_at" IS '更新时间';


-- ----------------------------
-- Table structure for farm_user_action
-- ----------------------------
DROP TABLE IF EXISTS "public"."farm_user_action";
-- Farm用户行为流水表 
CREATE TABLE "public"."farm_user_action" (
  id int8 PRIMARY KEY,
  "chain_id" int8 NOT NULL DEFAULT 0,
  "block_number" int8 NOT NULL DEFAULT 0,
  "tx_hash" varchar(66) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "user_address" varchar(42) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "to_address" varchar(42) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "incentive_id" VARCHAR(66) NOT NULL DEFAULT '',
  "token_id" int8 NOT NULL DEFAULT 0,
  "liquidity" varchar(50) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "pool_address" varchar(42) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "action" varchar(100) COLLATE "pg_catalog"."default",
  "event_sign" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "log_index" int4 NOT NULL DEFAULT 0,
  "amount" varchar(50) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "is_user" bool NOT NULL DEFAULT false,
  "tx_time" timestamp(6),
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6),
  CONSTRAINT "farm_user_action_tx_hash_event_sign_log_index_key" UNIQUE ("tx_hash", "event_sign", "log_index"),
  CONSTRAINT "farm_user_action_action_check" CHECK (action = ANY (ARRAY['create_pool'::text, 'destroy_pool'::text, 'stake'::text, 'unstake'::text, 'harvest'::text, 'transfer'::text, 'other'::text]))
);

DROP SEQUENCE IF EXISTS "farm_user_action_id_seq";
CREATE SEQUENCE IF NOT EXISTS "farm_user_action_id_seq" START WITH 1;
ALTER TABLE "public"."farm_user_action" ALTER COLUMN id SET DEFAULT nextval('farm_user_action_id_seq');

COMMENT ON TABLE "public"."farm_user_action" IS 'Farm用户行为流水表'; 
COMMENT ON COLUMN "public"."farm_user_action"."chain_id" IS '链ID';
COMMENT ON COLUMN "public"."farm_user_action"."block_number" IS '区块号';
COMMENT ON COLUMN "public"."farm_user_action"."tx_hash" IS '交易哈希';
COMMENT ON COLUMN "public"."farm_user_action"."user_address" IS '用户地址';
COMMENT ON COLUMN "public"."farm_user_action"."to_address" IS '接收地址';
COMMENT ON COLUMN "public"."farm_user_action"."incentive_id" IS '质押池ID';
COMMENT ON COLUMN "public"."farm_user_action"."token_id" IS 'LP Token ID';
COMMENT ON COLUMN "public"."farm_user_action"."liquidity" IS '质押的v3池子流动性'; 
COMMENT ON COLUMN "public"."farm_user_action"."pool_address" IS 'v3池子地址';
COMMENT ON COLUMN "public"."farm_user_action"."action" IS '操作类型';
COMMENT ON COLUMN "public"."farm_user_action"."event_sign" IS '事件签名';
COMMENT ON COLUMN "public"."farm_user_action"."log_index" IS '日志索引';
COMMENT ON COLUMN "public"."farm_user_action"."amount" IS '金额或数量';
COMMENT ON COLUMN "public"."farm_user_action"."is_user" IS '是否为用户操作';
COMMENT ON COLUMN "public"."farm_user_action"."tx_time" IS '交易时间';
COMMENT ON COLUMN "public"."farm_user_action"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."farm_user_action"."updated_at" IS '更新时间';


