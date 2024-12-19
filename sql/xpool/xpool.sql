
-- ----------------------------
-- Table structure for xpool_pools
-- ----------------------------
DROP TABLE IF EXISTS "public"."xpool_pools";
-- 单币质押池表
CREATE TABLE xpool_pools (
    id BIGSERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL DEFAULT 0,
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
    pool_address VARCHAR(42) NOT NULL DEFAULT '',
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
COMMENT ON COLUMN xpool_pools.status IS '运行状态（1-正常，2-已撤销）';

CREATE UNIQUE INDEX ak_xpool_pools_pool_address_chain_id ON xpool_pools(pool_address, chain_id);
CREATE INDEX idx_xpool_pools_staking_token_symbol_chain_id ON xpool_pools(staking_token_symbol, chain_id);
CREATE INDEX idx_xpool_pools_addresses_chain_id ON xpool_pools(pool_address, staking_token_address, reward_token_address, chain_id);
CREATE INDEX idx_xpool_pools_status_tvl ON xpool_pools(status, tvl DESC);


CREATE TABLE xpool_user_action (
    id SERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL DEFAULT 0, -- 链ID
    block_number BIGINT NOT NULL DEFAULT 0, -- 区块高度
    tx_hash VARCHAR(66) NOT NULL DEFAULT '', -- 交易哈希
    user_address VARCHAR(42) NOT NULL DEFAULT '', -- 用户地址
    vault_id INT NOT NULL DEFAULT 0, -- 质押池ID
    vault_address VARCHAR(42) NOT NULL DEFAULT '', -- 质押池合约地址
    action TEXT CHECK (action IN ('create_pool', 'destroy_pool', 'stake', 'unstake', 'harvest', 'other')), -- 行为
    event_sign VARCHAR(100) NOT NULL DEFAULT '', -- 事件签名
    log_index INT NOT NULL DEFAULT 0, -- 事件顺序
    amount VARCHAR(50) NOT NULL DEFAULT '0.000000000000000000', -- 交易数量
    is_user BOOLEAN NOT NULL DEFAULT FALSE, -- 是否为用户操作（TRUE-是，FALSE-否）
    tx_time TIMESTAMP NULL DEFAULT NULL, -- 成交时间
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP, -- 更新时间
    UNIQUE (tx_hash, event_sign, log_index)
);

-- 索引
CREATE INDEX idx_xpool_user_action_user_address_log_index_tx_time 
ON xpool_user_action (user_address, log_index, tx_time);

CREATE INDEX idx_xpool_user_action_vault_id_is_user 
ON xpool_user_action (vault_id, is_user);

-- 添加列注释
COMMENT ON COLUMN xpool_user_action.tx_hash IS '交易哈希';
COMMENT ON COLUMN xpool_user_action.user_address IS '用户地址';
COMMENT ON COLUMN xpool_user_action.vault_id IS '质押池ID';
COMMENT ON COLUMN xpool_user_action.vault_address IS '质押池合约地址';
COMMENT ON COLUMN xpool_user_action.action IS '行为';
COMMENT ON COLUMN xpool_user_action.event_sign IS '事件签名';
COMMENT ON COLUMN xpool_user_action.amount IS '交易数量';
COMMENT ON COLUMN xpool_user_action.is_user IS '是否为用户操作（TRUE-是，FALSE-否）';
COMMENT ON COLUMN xpool_user_action.tx_time IS '成交时间';
COMMENT ON COLUMN xpool_user_action.created_at IS '创建时间';
COMMENT ON COLUMN xpool_user_action.updated_at IS '更新时间';

