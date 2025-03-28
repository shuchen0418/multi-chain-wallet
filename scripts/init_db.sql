-- 创建数据库
CREATE DATABASE IF NOT EXISTS multi_chain_wallet CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE multi_chain_wallet;

-- 创建钱包表
CREATE TABLE IF NOT EXISTS wallets (
    id VARCHAR(36) PRIMARY KEY,
    address VARCHAR(42) NOT NULL UNIQUE,
    priv_key_enc TEXT NOT NULL,
    mnemonic_enc TEXT,
    chain_type VARCHAR(20) NOT NULL,
    create_time BIGINT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_address (address),
    INDEX idx_chain_type (chain_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建交易表
CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(36) PRIMARY KEY,
    wallet_id VARCHAR(36) NOT NULL,
    tx_hash VARCHAR(66) NOT NULL UNIQUE,
    `from` VARCHAR(42) NOT NULL,
    `to` VARCHAR(42) NOT NULL,
    amount VARCHAR(78) NOT NULL,
    status VARCHAR(20) NOT NULL,
    chain_type VARCHAR(20) NOT NULL,
    create_time BIGINT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_wallet_id (wallet_id),
    INDEX idx_tx_hash (tx_hash),
    INDEX idx_from (`from`),
    INDEX idx_to (`to`),
    INDEX idx_chain_type (chain_type),
    FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci; 