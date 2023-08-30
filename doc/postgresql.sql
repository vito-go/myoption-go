-- 用户
CREATE TABLE IF NOT EXISTS public.user_info
(
    id              serial8               NOT NULL PRIMARY KEY,        -- 用户信息表主键ID
    user_id         varchar(32)           NOT NULL,                    -- 用户ID
    nick            character varying(16) NOT NULL,                    -- 用户昵称
    x25519_pub_key  character varying(64) NOT NULL,                    -- X25519公钥
    ed25519_pub_key character varying(64) NOT NULL,                    -- Ed25519公钥
    status          int2                     DEFAULT 1,                -- 用户状态，1表示激活，其他值表示禁用或其他状态
    create_time     timestamp with time zone,                          -- 用户信息创建时间
    update_time     timestamp with time zone DEFAULT current_timestamp -- 用户信息最后更新时间，默认为当前时间
);
CREATE UNIQUE INDEX ON public.user_info (user_id);

COMMENT ON TABLE public.user_info IS '用户信息表';
COMMENT ON COLUMN public.user_info.id IS '用户信息表主键ID';
COMMENT ON COLUMN public.user_info.user_id IS '用户ID';
COMMENT ON COLUMN public.user_info.nick IS '用户昵称';
COMMENT ON COLUMN public.user_info.status IS '用户状态，1表示激活，其他值表示禁用或其他状态';
COMMENT ON COLUMN public.user_info.create_time IS '用户信息创建时间';
COMMENT ON COLUMN public.user_info.update_time IS '用户信息最后更新时间，默认为当前时间';
COMMENT ON COLUMN public.user_info.x25519_pub_key IS 'X25519公钥';
COMMENT ON COLUMN public.user_info.ed25519_pub_key IS 'Ed25519公钥';


-- 用户
CREATE TABLE IF NOT EXISTS public.user_key
(
    id                  serial8                NOT NULL PRIMARY KEY,       -- 用户信息表主键ID
    user_id             varchar(32)            NOT NULL,                   -- 用户ID
    password            character varying(256) NOT NULL,                   -- 用户密码
    salt                character varying(32)  NOT NULL,                   -- 盐值 base64
    x25519_pri_enc_key  character varying(64)  NOT NULL,                   -- X25519私钥加密密钥
    ed25519_pri_enc_key character varying(64)  NOT NULL,                   -- Ed25519私钥加密密钥
    status              int2                     DEFAULT 1,                -- 用户状态，1表示激活，其他值表示禁用或其他状态
    create_time         timestamp with time zone,                          -- 用户信息创建时间
    update_time         timestamp with time zone DEFAULT current_timestamp -- 用户信息最后更新时间，默认为当前时间
);
CREATE UNIQUE INDEX ON public.user_key (user_id);

COMMENT ON TABLE public.user_key IS '用户密钥表';
COMMENT ON COLUMN public.user_key.id IS '用户信息表主键ID';
COMMENT ON COLUMN public.user_key.user_id IS '用户ID';
COMMENT ON COLUMN public.user_key.password IS '用户密码';
COMMENT ON COLUMN public.user_key.salt IS '盐';
COMMENT ON COLUMN public.user_key.x25519_pri_enc_key IS 'X25519私钥加密密钥';
COMMENT ON COLUMN public.user_key.ed25519_pri_enc_key IS 'Ed25519私钥加密密钥';
COMMENT ON COLUMN public.user_key.status IS '用户状态，1表示激活，其他值表示禁用或其他状态';
COMMENT ON COLUMN public.user_key.create_time IS '用户信息创建时间';
COMMENT ON COLUMN public.user_key.update_time IS '用户信息最后更新时间，默认为当前时间';


DROP TABLE IF EXISTS stock_price;
CREATE TABLE IF NOT EXISTS stock_price
(
    id           serial8               NOT NULL PRIMARY KEY,
    country_code character varying(16) NOT NULL,
    symbol_code  character varying(32) NOT NULL,
    today        int4                  NOT NULL,
    time_min     int4                  NOT NULL,
    price        double precision      NOT NULL,
    volume       bigint                NOT NULL,
    avg_price    double precision      NOT NULL,
    amount       bigint                NOT NULL,
    status       integer               NOT NULL DEFAULT 0,
    create_time  timestamp with time zone,
    update_time  timestamp with time zone       DEFAULT current_timestamp
);
CREATE UNIQUE INDEX ON stock_price (country_code, symbol_code, today, time_min);
COMMENT ON TABLE stock_price IS '股票价格表';
COMMENT ON COLUMN stock_price.id IS '用户信息表主键ID';
COMMENT ON COLUMN stock_price.country_code IS '国家代码';
COMMENT ON COLUMN stock_price.symbol_code IS '股票代码';
COMMENT ON COLUMN stock_price.today IS '日期';
COMMENT ON COLUMN stock_price.time_min IS '时间';
COMMENT ON COLUMN stock_price.price IS '价格';
COMMENT ON COLUMN stock_price.volume IS '成交量';
COMMENT ON COLUMN stock_price.avg_price IS '均价';
COMMENT ON COLUMN stock_price.amount IS '成交额';
COMMENT ON COLUMN stock_price.status IS '1:初始化，2:已更新';
COMMENT ON COLUMN stock_price.create_time IS '用户信息创建时间';
COMMENT ON COLUMN stock_price.update_time IS '用户信息最后更新时间，默认为当前时间';


--
-- 创建订单表 orders_binary_option
CREATE TABLE IF NOT EXISTS public.orders_binary_option
(
    id               SERIAL PRIMARY KEY,                          -- 订单编号，使用PostgreSQL的自增序列作为主键
    trans_id         VARCHAR(32)              not null,
    user_id          varchar(32)              NOT NULL,           -- 下单用户的唯一标识符
    country_code      VARCHAR(8)              NOT NULL,           -- 国家代码
    symbol_code      VARCHAR(10)              NOT NULL,           -- 股票代码
    strike_price     double precision         NOT NULL,           -- 价格，保留小数点后两位
    option           smallint                 NOT NULL,-- 1 CALL 2 put
    bet_money        integer                  NOT NULL,--下注金额
    session          INTEGER                  NOT NULL,           --0 2 3 5 10 15 20 30 60
    today            INTEGER                  NOT NULL,
    session_time_min INTEGER                  NOT NULL,           -- 200601021504
    order_time       TIMESTAMP with time zone NOT NULL,-- 下单时间
    settle_time      bigint                  NOT NULL DEFAULT 0,--
    settle_price     double precision         NOT NULL DEFAULT 0,-- 结算价格
    settle_result    SMALLINT                 NOT NULL DEFAULT 0, -- 1. user win 2 user lose
    order_status     SMALLINT                 NOT NULL DEFAULT 0, -- 订单状态，如未成交（pending）、已成交（filled）、部分成交（partially filled）等
    profit_loss      INTEGER                  NOT NULL default 0,
    create_time      timestamp with time zone,
    update_time      timestamp with time zone          DEFAULT current_timestamp
);
CREATE INDEX ON orders_binary_option (user_id, create_time);
CREATE INDEX ON orders_binary_option (today,order_status,session_time_min);
CREATE UNIQUE INDEX ON orders_binary_option (trans_id);



CREATE TABLE IF NOT EXISTS wallet_info
(
    id            BIGSERIAL PRIMARY KEY,                                       -- 钱包唯一标识 ID
    user_id       VARCHAR(64)              NOT NULL UNIQUE,                    -- 钱包所属用户 ID
    balance       BIGINT                   NOT NULL DEFAULT 0,                 -- 钱包可用余额，单位为分
    frozen_amount BIGINT                   NOT NULL DEFAULT 0,                 -- 钱包冻结金额，单位为分
    total_amount  BIGINT                   NOT NULL DEFAULT 0,                 -- 钱包总余额，单位为分
    create_time   TIMESTAMP with time zone NOT NULL DEFAULT current_timestamp, -- 数据创建时间
    update_time   TIMESTAMP with time zone NOT NULL DEFAULT current_timestamp, -- 数据更新时间
    CHECK (total_amount >= 0 and frozen_amount >= 0 and balance >= 0 and total_amount >= frozen_amount and
           total_amount >= balance)
);
COMMENT ON TABLE wallet_info IS '钱包';

-- 添加表字段注释
CREATE UNIQUE INDEX ON wallet_info (user_id);
COMMENT ON COLUMN wallet_info.id IS '钱包唯一标识 ID';
COMMENT ON COLUMN wallet_info.user_id IS '钱包所属用户 ID';
COMMENT ON COLUMN wallet_info.balance IS '钱包可用余额，单位为分';
COMMENT ON COLUMN wallet_info.frozen_amount IS '钱包冻结金额，单位为分';
COMMENT ON COLUMN wallet_info.total_amount IS '钱包总余额，单位为分';
COMMENT ON COLUMN wallet_info.create_time IS '数据创建时间';
COMMENT ON COLUMN wallet_info.update_time IS '数据更新时间';


CREATE TABLE IF NOT EXISTS wallet_detail
(
    id              BIGSERIAL PRIMARY KEY,                                       -- 明细记录唯一标识 ID
    trans_id        VARCHAR(32),                                                 -- 支付平台交易 ID( 如：微信支付订单号，支付宝订单号等）
    trans_type      SMALLINT                 NOT NULL,                           -- 交易类型( 充值deposit，下单，提现 withdraw，盈利）
    user_id         VARCHAR(64)              NOT NULL,                           -- 钱包所属用户 ID
    amount          BIGINT                   NOT NULL,                           -- 交易金额，
    status          SMALLINT,                                                    -- 交易状态(0:初始化;1:WAITING;2 PENDING,100 CLOSE 200 SUCCESS)
    remark          VARCHAR(256),                                                -- 交易备注
    source_kind     SMALLINT,                                                    -- 充值/提现来源类型(银行卡、支付宝、微信等）
    source_trans_id VARCHAR(64),                                                 -- 充值/提现/  id
    from_account    VARCHAR(256),                                                -- 充值/提现　资金来源账户
    to_account      VARCHAR(256),                                                -- 充值/提现　资金去向账户
    balance         BIGINT                   NOT NULL DEFAULT 0,                 -- 钱包可用余额，单位为分
    create_time     TIMESTAMP with time zone NOT NULL DEFAULT current_timestamp, -- 数据创建时间
    update_time     TIMESTAMP with time zone NOT NULL DEFAULT current_timestamp  -- 数据更新时间
);

-- 添加表字段注释
CREATE INDEX ON wallet_detail (user_id);
CREATE UNIQUE INDEX ON wallet_detail (trans_id);
COMMENT ON TABLE wallet_detail IS '钱包明细表';

COMMENT ON COLUMN wallet_detail.id IS '明细记录唯一标识ID';
COMMENT ON COLUMN wallet_detail.trans_id IS '支付平台交易ID(如：微信支付订单号，支付宝订单号等）';
COMMENT ON COLUMN wallet_detail.trans_type IS '交易类型（充值deposit，下单，提现withdraw，盈利）';
COMMENT ON COLUMN wallet_detail.user_id IS '钱包所属用户ID';
COMMENT ON COLUMN wallet_detail.amount IS '交易金额，单位为分';
COMMENT ON COLUMN wallet_detail.status IS '交易状态（0:初始化;1:WAITING;2:PENDING;100:CLOSE;200:SUCCESS）';
COMMENT ON COLUMN wallet_detail.remark IS '交易备注';
COMMENT ON COLUMN wallet_detail.source_kind IS '充值/提现来源类型(银行卡、支付宝、微信等）';
COMMENT ON COLUMN wallet_detail.source_trans_id IS '充值/提现/下单/盈利 ID';
COMMENT ON COLUMN wallet_detail.from_account IS '充值/提现资金来源账户';
COMMENT ON COLUMN wallet_detail.to_account IS '充值/提现资金去向账户';
COMMENT ON COLUMN wallet_detail.balance IS '钱包可用余额，单位为分';
COMMENT ON COLUMN wallet_detail.create_time IS '数据创建时间';
COMMENT ON COLUMN wallet_detail.update_time IS '数据更新时间';

