-- ----------------------------
-- Indexes structure for table event_farm_incentive_created
-- ----------------------------
CREATE UNIQUE INDEX "event_farm_incentive_created_index_incentive_id_chain_id" ON "public"."event_farm_incentive_created" USING btree (
  "event_incentive_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);