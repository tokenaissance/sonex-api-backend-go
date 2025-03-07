create table sonex_badge_staked_position (
    position_id BIGINT NOT NULL,
    position_value VARCHAR(255) NOT NULL, 
    user_address VARCHAR(255) NOT NULL, 
    PRIMARY KEY (position_id)
);


-- function for trigger of after_insert_event_farm_token_unstaked
CREATE OR REPLACE FUNCTION delete_sonex_badge_staked_position()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM sonex_badge_staked_position WHERE position_id=NEW.event_token_id;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION delete_sonex_badge_staked_position () IS 'function for trigger of after_insert_event_farm_token_unstaked';

CREATE TRIGGER after_insert_event_farm_token_unstaked
AFTER INSERT
ON event_farm_token_unstaked
FOR EACH ROW
EXECUTE FUNCTION delete_sonex_badge_staked_position();


create table sonex_badge_rewards (
    id SERIAL NOT NULL,
    user_address VARCHAR(255) NOT NULL, 
    points INTEGER NOT NULL, 
    nft_id INTEGER NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (user_address)
);