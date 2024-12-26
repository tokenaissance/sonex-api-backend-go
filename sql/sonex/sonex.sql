DROP TABLE IF EXISTS public.sonex_token;
CREATE TABLE public.sonex_token (
    id                            BIGSERIAL      PRIMARY KEY,
    created_at timestamp with time zone DEFAULT ('now'::text)::timestamp(6) with time zone NOT NULL,
    updated_at timestamp with time zone DEFAULT ('now'::text)::timestamp(6) with time zone NOT NULL,
    deleted_at timestamp without time zone,
    name character varying(120) NOT NULL,
    symbol character varying,
    address character varying(128),
    decimals integer DEFAULT 18 NOT NULL,
    logo text,
    sort integer DEFAULT 0 NOT NULL,
    "recommendOrder" integer DEFAULT 0 NOT NULL,
    "isHidden" boolean DEFAULT false NOT NULL,
    volume numeric(32,18) DEFAULT '0'::numeric NOT NULL,
    additional jsonb DEFAULT '{}'::jsonb NOT NULL
);

-- 为 address 列创建索引
CREATE INDEX idx_token_address ON public.sonex_token (address);

-- 为 symbol 列创建索引
CREATE INDEX idx_token_symbol ON public.sonex_token (symbol);

-- 为 isHidden 列创建索引
CREATE INDEX idx_token_is_hidden ON public.sonex_token ("isHidden");


CREATE TABLE "public"."sonex_token_price_record" (
     "id" int8 NOT NULL DEFAULT nextval('sonex_token_price_record_id'::regclass),
     "token_symbol" varchar(255) COLLATE "pg_catalog"."default",
     "token_address" varchar(255) COLLATE "pg_catalog"."default",
     "price" varchar COLLATE "pg_catalog"."default",
     "created_at" timestamp(6),
     CONSTRAINT "sonex_token_price_record_pkey" PRIMARY KEY ("id")
);

ALTER TABLE "public"."sonex_token_price_record"
    OWNER TO "postgres";

CREATE SEQUENCE "public"."sonex_token_price_record_id"
    INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

ALTER SEQUENCE "public"."sonex_token_price_record_id" OWNER TO "postgres";

CREATE INDEX "token_address_idx" ON "public"."sonex_token_price_record" USING btree (
    "token_address" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
    );