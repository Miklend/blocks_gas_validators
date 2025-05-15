CREATE TABLE IF NOT EXISTS ethereum_block_metrics (
    block_number BIGINT NOT NULL,
    block_time TIMESTAMPTZ NOT NULL,
    transactions_count INT,
    block_size_bytes BIGINT,
    gas_limit BIGINT,
    gas_used BIGINT,
    block_fullness DOUBLE PRECISION,
    block_author TEXT,
    gas_min DOUBLE PRECISION,
    gas_max DOUBLE PRECISION,
    gas_avg DOUBLE PRECISION,
    gas_stddev DOUBLE PRECISION,
    gas_all_prices JSONB,
    block_timestamp BIGINT NOT NULL,
    PRIMARY KEY (block_time, block_number)
) PARTITION BY RANGE (block_time);

CREATE TABLE IF NOT EXISTS polygon_block_metrics (
    block_number BIGINT NOT NULL,
    block_time TIMESTAMPTZ NOT NULL,
    transactions_count INT,
    block_size_bytes BIGINT,
    gas_limit BIGINT,
    gas_used BIGINT,
    block_fullness DOUBLE PRECISION,
    block_author TEXT,
    gas_min DOUBLE PRECISION,
    gas_max DOUBLE PRECISION,
    gas_avg DOUBLE PRECISION,
    gas_stddev DOUBLE PRECISION,
    gas_all_prices JSONB,
    block_timestamp BIGINT NOT NULL,
    PRIMARY KEY (block_time, block_number)
) PARTITION BY RANGE (block_time);

CREATE TABLE IF NOT EXISTS avalanche_block_metrics (
    block_number BIGINT NOT NULL,
    block_time TIMESTAMPTZ NOT NULL,
    transactions_count INT,
    block_size_bytes BIGINT,
    gas_limit BIGINT,
    gas_used BIGINT,
    block_fullness DOUBLE PRECISION,
    block_author TEXT,
    gas_min DOUBLE PRECISION,
    gas_max DOUBLE PRECISION,
    gas_avg DOUBLE PRECISION,
    gas_stddev DOUBLE PRECISION,
    gas_all_prices JSONB,
    block_timestamp BIGINT NOT NULL,
    PRIMARY KEY (block_time, block_number)
) PARTITION BY RANGE (block_time);

CREATE TABLE IF NOT EXISTS bnb_block_metrics (
    block_number BIGINT NOT NULL,
    block_time TIMESTAMPTZ NOT NULL,
    transactions_count INT,
    block_size_bytes BIGINT,
    gas_limit BIGINT,
    gas_used BIGINT,
    block_fullness DOUBLE PRECISION,
    block_author TEXT,
    gas_min DOUBLE PRECISION,
    gas_max DOUBLE PRECISION,
    gas_avg DOUBLE PRECISION,
    gas_stddev DOUBLE PRECISION,
    gas_all_prices JSONB,
    block_timestamp BIGINT NOT NULL,
    PRIMARY KEY (block_time, block_number)
) PARTITION BY RANGE (block_time);

CREATE TABLE IF NOT EXISTS base_block_metrics (
    block_number BIGINT NOT NULL,
    block_time TIMESTAMPTZ NOT NULL,
    transactions_count INT,
    block_size_bytes BIGINT,
    gas_limit BIGINT,
    gas_used BIGINT,
    block_fullness DOUBLE PRECISION,
    block_author TEXT,
    gas_min DOUBLE PRECISION,
    gas_max DOUBLE PRECISION,
    gas_avg DOUBLE PRECISION,
    gas_stddev DOUBLE PRECISION,
    gas_all_prices JSONB,
    block_timestamp BIGINT NOT NULL,
    PRIMARY KEY (block_time, block_number)
) PARTITION BY RANGE (block_time);

CREATE TABLE IF NOT EXISTS optimism_block_metrics (
    block_number BIGINT NOT NULL,
    block_time TIMESTAMPTZ NOT NULL,
    transactions_count INT,
    block_size_bytes BIGINT,
    gas_limit BIGINT,
    gas_used BIGINT,
    block_fullness DOUBLE PRECISION,
    block_author TEXT,
    gas_min DOUBLE PRECISION,
    gas_max DOUBLE PRECISION,
    gas_avg DOUBLE PRECISION,
    gas_stddev DOUBLE PRECISION,
    gas_all_prices JSONB,
    block_timestamp BIGINT NOT NULL,
    PRIMARY KEY (block_time, block_number)
) PARTITION BY RANGE (block_time);