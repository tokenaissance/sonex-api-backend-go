
-- ----------------------------
-- Table structure for farm_pools
-- ----------------------------
DROP TABLE IF EXISTS "public"."farm_pools";
-- Farm质押池表
CREATE TABLE farm_pools (
    id SERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL DEFAULT 0, -- 链ID
    incentive_id VARCHAR(64) NOT NULL DEFAULT '', -- 质押池ID
    start_time TIMESTAMP, -- 质押池开放时间
    end_time TIMESTAMP, -- 质押池结束时间
    reward_token_address VARCHAR(42) NOT NULL DEFAULT '', -- 奖励token代币地址
    reward_token_decimals SMALLINT NOT NULL DEFAULT 0, -- 奖励token代币精度
    reward_token_amount VARCHAR(50) NOT NULL DEFAULT '', -- 奖励token数量
    reward_token_price VARCHAR(50) NOT NULL DEFAULT '', -- 奖励token价格
    refundee_address VARCHAR(42) NOT NULL DEFAULT '', -- 项目方回收奖励代币地址

    pool_address VARCHAR(42) NOT NULL DEFAULT '', -- v3池子合约地址
    fee_ratio VARCHAR(10) NOT NULL DEFAULT 0.0000, -- v3池子手续费率
    token0_address VARCHAR(42) NOT NULL DEFAULT '', -- token0代币地址
    token0_symbol VARCHAR(10) NOT NULL DEFAULT '', -- token0符号
    token0_decimals SMALLINT NOT NULL DEFAULT 0, -- token0精度
    token1_address VARCHAR(42) NOT NULL DEFAULT '', -- token1代币地址
    token1_symbol VARCHAR(10) NOT NULL DEFAULT '', -- token1符号
    token1_decimals SMALLINT NOT NULL DEFAULT 0, -- token1精度
    
    staked_liquidity_usd VARCHAR(50) NOT NULL DEFAULT '', -- 质押池质押的流动性的价值（USD）
    total_token0_amount VARCHAR(50) NOT NULL DEFAULT '', -- token0质押总量
    total_token1_amount VARCHAR(50) NOT NULL DEFAULT '', -- token1质押总量
    apr VARCHAR(50) NOT NULL DEFAULT '', -- 质押池年化收益率
    join_num INTEGER NOT NULL DEFAULT 0, -- 质押池参与质押人数
    tx_hash varchar(255) NOT NULL DEFAULT '', -- 交易哈希
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP NULL DEFAULT NULL, -- 更新时间
    UNIQUE ("incentive_id", "chain_id")
);

-- 表注释
COMMENT ON TABLE farm_pools IS 'Farm质押池表';

-- 字段注释
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
COMMENT ON COLUMN "farm_pools"."staked_liquidity_usd" IS '质押池质押的流动性的价值（USD）';
COMMENT ON COLUMN "farm_pools"."total_token0_amount" IS 'token0质押总量';
COMMENT ON COLUMN "farm_pools"."total_token1_amount" IS 'token1质押总量';
COMMENT ON COLUMN "farm_pools"."apr" IS '质押池年化收益率';
COMMENT ON COLUMN "farm_pools"."join_num" IS '质押池参与质押人数';
COMMENT ON COLUMN "farm_pools"."tx_hash" IS '质押池创建时的交易哈希';
COMMENT ON COLUMN "farm_pools"."created_at" IS '创建时间';
COMMENT ON COLUMN "farm_pools"."updated_at" IS '更新时间';

DROP TABLE IF EXISTS "public"."farm_user_pools";
CREATE TABLE "public"."farm_user_pools" (
  "id" SERIAL PRIMARY KEY,
  "incentive_id" VARCHAR(64) NOT NULL DEFAULT '',
  "token_id" BIGINT NOT NULL DEFAULT 0,
  "user_address" VARCHAR(42) NOT NULL DEFAULT '',
  "chain_id" BIGINT NOT NULL DEFAULT 0,
  "status" VARCHAR(10) NOT NULL CHECK (status IN ('stake', 'unstake')),
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NULL DEFAULT NULL
);

CREATE UNIQUE INDEX "idx_farm_user_pools_unique" ON "public"."farm_user_pools" ("user_address", "incentive_id", "token_id",  "chain_id", "status");

-- 表注释
COMMENT ON TABLE "public"."farm_user_pools" IS '用户Farm质押池关联表';

-- 字段注释
COMMENT ON COLUMN "public"."farm_user_pools"."incentive_id" IS '质押池ID';
COMMENT ON COLUMN "public"."farm_user_pools"."token_id" IS 'LP Token ID';
COMMENT ON COLUMN "public"."farm_user_pools"."user_address" IS '用户地址';
COMMENT ON COLUMN "public"."farm_user_pools"."chain_id" IS '链ID';
COMMENT ON COLUMN "public"."farm_user_pools"."status" IS '状态(stake:质押中 unstake:已解除)';
COMMENT ON COLUMN "public"."farm_user_pools"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."farm_user_pools"."updated_at" IS '更新时间';