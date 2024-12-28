
-- ----------------------------
-- Table structure for xpool_pools
-- ----------------------------
DROP TABLE IF EXISTS "public"."xpool_pools";
-- 单币质押池表
CREATE TABLE xpool_pools (
    "id" int8 NOT NULL,
    chain_id BIGINT NOT NULL DEFAULT 0,
    pool_address VARCHAR(42) NOT NULL DEFAULT '',
    staking_token_name VARCHAR(50) NOT NULL DEFAULT '',
    staking_token_symbol VARCHAR(10) NOT NULL DEFAULT '',
    staking_token_address VARCHAR(42) NOT NULL DEFAULT '',
    staking_token_decimals SMALLINT NOT NULL DEFAULT 18,
    staking_token_price VARCHAR(50) NOT NULL DEFAULT 0, 
    reward_token_name VARCHAR(50) NOT NULL DEFAULT '',
    reward_token_symbol VARCHAR(10) NOT NULL DEFAULT '',
    reward_token_address VARCHAR(42) NOT NULL DEFAULT '',
    reward_token_decimals SMALLINT NOT NULL DEFAULT 18,
    reward_token_price VARCHAR(50) NOT NULL DEFAULT 0,
    reward_per_block VARCHAR(50) NOT NULL DEFAULT 0,
    total_staked VARCHAR(50) NOT NULL DEFAULT 0,
    total_claimed VARCHAR(50) NOT NULL DEFAULT 0,
    block_time SMALLINT NOT NULL DEFAULT 0,
    tvl VARCHAR(50) NOT NULL DEFAULT 0,
    apr VARCHAR(50) NOT NULL DEFAULT 0,
    join_num INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT check_status CHECK (status IN ('active', 'closed')),
    CONSTRAINT "xpool_pools_pkey" PRIMARY KEY ("id")
);
DROP SEQUENCE IF EXISTS "xpool_pools_id_seq";
CREATE SEQUENCE IF NOT EXISTS "xpool_pools_id_seq" START WITH 1;
ALTER TABLE "public"."xpool_pools" ALTER COLUMN id SET DEFAULT nextval('xpool_pools_id_seq');

ALTER TABLE "public"."xpool_pools" OWNER TO "postgres";

COMMENT ON TABLE xpool_pools IS '单币质押池表';
COMMENT ON COLUMN xpool_pools.staking_token_name IS '质押的代币名称';
COMMENT ON COLUMN xpool_pools.staking_token_symbol IS '质押的代币符号';
COMMENT ON COLUMN xpool_pools.staking_token_address IS '质押的代币合约地址';
COMMENT ON COLUMN xpool_pools.staking_token_decimals IS '质押的代币合约精度';
COMMENT ON COLUMN xpool_pools.staking_token_price IS '质押的代币的价格';
COMMENT ON COLUMN xpool_pools.reward_token_name IS '奖励的代币名称';
COMMENT ON COLUMN xpool_pools.reward_token_symbol IS '奖励的代币符号';
COMMENT ON COLUMN xpool_pools.reward_token_address IS '奖励的代币合约地址';
COMMENT ON COLUMN xpool_pools.reward_token_decimals IS '奖励的代币合约精度';
COMMENT ON COLUMN xpool_pools.reward_token_price IS '奖励的代币的价格';
COMMENT ON COLUMN xpool_pools.pool_address IS '质押池合约地址';
COMMENT ON COLUMN xpool_pools.reward_per_block IS '出块奖励';
COMMENT ON COLUMN xpool_pools.total_staked IS '累计池子的质押币数量';
COMMENT ON COLUMN xpool_pools.total_claimed IS '累计池子的奖励金额';
COMMENT ON COLUMN xpool_pools.block_time IS '出块速度（s）';
COMMENT ON COLUMN xpool_pools.tvl IS '池子的质押总价值（USD）';
COMMENT ON COLUMN xpool_pools.apr IS '年化收益率';
COMMENT ON COLUMN xpool_pools.join_num IS '参与人数';
COMMENT ON COLUMN xpool_pools.status IS '运行状态';

CREATE UNIQUE INDEX ak_xpool_pools_pool_address_chain_id ON xpool_pools(pool_address, chain_id);
CREATE INDEX idx_xpool_pools_staking_token_symbol_chain_id ON xpool_pools(staking_token_symbol, chain_id);
CREATE INDEX idx_xpool_pools_addresses_chain_id ON xpool_pools(pool_address, staking_token_address, reward_token_address, chain_id);
CREATE INDEX idx_xpool_pools_status_tvl ON xpool_pools(status, tvl DESC);

-- ----------------------------
-- Table structure for xpool_user_action
-- ----------------------------
DROP TABLE IF EXISTS "public"."xpool_user_action";
-- 用户参与的质押池子表 
CREATE TABLE "public"."xpool_user_action" (
  "id" int8 NOT NULL,
  "chain_id" int8 NOT NULL DEFAULT 0,
  "block_number" int8 NOT NULL DEFAULT 0,
  "tx_hash" varchar(66) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "user_address" varchar(42) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "pool_id" int4 NOT NULL DEFAULT 0,
  "pool_address" varchar(42) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "action" text COLLATE "pg_catalog"."default",
  "event_sign" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "log_index" int4 NOT NULL DEFAULT 0,
  "amount" varchar(50) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '0.000000000000000000'::character varying,
  "is_user" bool NOT NULL DEFAULT false,
  "tx_time" timestamp(6),
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6),
  CONSTRAINT "xpool_user_action_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "xpool_user_action_tx_hash_event_sign_log_index_key" UNIQUE ("tx_hash", "event_sign", "log_index"),
  CONSTRAINT "xpool_user_action_action_check" CHECK (action = ANY (ARRAY['create_pool'::text, 'destroy_pool'::text, 'stake'::text, 'unstake'::text, 'harvest'::text, 'other'::text]))
);

DROP SEQUENCE IF EXISTS "xpool_user_action_id_seq";
CREATE SEQUENCE IF NOT EXISTS "xpool_user_action_id_seq" START WITH 1;
ALTER TABLE "public"."xpool_user_action" ALTER COLUMN id SET DEFAULT nextval('xpool_user_action_id_seq');

ALTER TABLE "public"."xpool_user_action" OWNER TO "postgres";

CREATE INDEX "idx_xpool_user_action_user_address_log_index_tx_time" ON "public"."xpool_user_action" USING btree (
  "user_address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "log_index" "pg_catalog"."int4_ops" ASC NULLS LAST,
  "tx_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);

CREATE INDEX "idx_xpool_user_action_pool_id_is_user" ON "public"."xpool_user_action" USING btree (
  "pool_id" "pg_catalog"."int4_ops" ASC NULLS LAST,
  "is_user" "pg_catalog"."bool_ops" ASC NULLS LAST
);

COMMENT ON TABLE "public"."xpool_user_action" IS 'Xpool用户行为流水表'; 
COMMENT ON COLUMN "public"."xpool_user_action"."tx_hash" IS '交易哈希';
COMMENT ON COLUMN "public"."xpool_user_action"."user_address" IS '用户地址';
COMMENT ON COLUMN "public"."xpool_user_action"."pool_id" IS '质押池ID';
COMMENT ON COLUMN "public"."xpool_user_action"."pool_address" IS '质押池合约地址';
COMMENT ON COLUMN "public"."xpool_user_action"."action" IS '行为';
COMMENT ON COLUMN "public"."xpool_user_action"."event_sign" IS '事件签名';
COMMENT ON COLUMN "public"."xpool_user_action"."amount" IS '交易数量';
COMMENT ON COLUMN "public"."xpool_user_action"."is_user" IS '是否为用户操作（TRUE-是，FALSE-否）';
COMMENT ON COLUMN "public"."xpool_user_action"."tx_time" IS '成交时间';
COMMENT ON COLUMN "public"."xpool_user_action"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."xpool_user_action"."updated_at" IS '更新时间';


-- ----------------------------
-- Table structure for xpool_user_pools
-- ----------------------------
DROP TABLE IF EXISTS "public"."xpool_user_pools";
-- 用户参与的质押池表 
CREATE TABLE "public"."xpool_user_pools" (
  "id" int8 NOT NULL,
  "chain_id" int8 NOT NULL DEFAULT 0, -- 链ID
  "pool_address" varchar(42) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying, -- 池地址
  "user_address" varchar(42) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying, -- 用户地址
  "staking_amount" varchar(50) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying, -- 质押数量
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间
  "updated_at" timestamp(6), -- 更新时间
  CONSTRAINT "xpool_user_pools_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "xpool_user_pools_pool_user_chain_idx" ON "public"."xpool_user_pools" ("pool_address", "user_address", "chain_id");

DROP SEQUENCE IF EXISTS "xpool_user_pools_id_seq";
CREATE SEQUENCE IF NOT EXISTS "xpool_user_pools_id_seq" START WITH 1;
ALTER TABLE "public"."xpool_user_pools" ALTER COLUMN id SET DEFAULT nextval('xpool_user_pools_id_seq');

ALTER TABLE "public"."xpool_user_pools" OWNER TO "postgres";

COMMENT ON TABLE "public"."xpool_user_pools" IS '用户参与的质押池表'; 
COMMENT ON COLUMN "public"."xpool_user_pools"."chain_id" IS '链ID';
COMMENT ON COLUMN "public"."xpool_user_pools"."pool_address" IS '质押池地址';
COMMENT ON COLUMN "public"."xpool_user_pools"."user_address" IS '用户地址';
COMMENT ON COLUMN "public"."xpool_user_pools"."staking_amount" IS '质押数量';
COMMENT ON COLUMN "public"."xpool_user_pools"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."xpool_user_pools"."updated_at" IS '更新时间';

