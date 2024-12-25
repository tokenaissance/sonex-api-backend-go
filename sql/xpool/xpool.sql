
-- ----------------------------
-- Table structure for xpool_pools
-- ----------------------------
DROP TABLE IF EXISTS "public"."xpool_pools";
-- 单币质押池表
CREATE TABLE xpool_pools (
    id BIGSERIAL PRIMARY KEY,
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
    CONSTRAINT check_status CHECK (status IN ('active', 'closed'))
);

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


