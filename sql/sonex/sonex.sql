DROP TABLE IF EXISTS public.sonex_token;
CREATE TABLE public.sonex_token (
    "id" int8 NOT NULL,
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