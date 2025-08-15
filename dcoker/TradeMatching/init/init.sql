CREATE DATABASE trade;

USE trade;

-- 撮合訂單表
CREATE TABLE `exchange_order` (
                                  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'id',
                                  `order_id` VARCHAR(64) NOT NULL COMMENT '撮合訂單號',
                                  `member_id` BIGINT NOT NULL COMMENT '會員號',
                                  `type` VARCHAR(32) NOT NULL COMMENT '掛單類型',
                                  `amount` DECIMAL(32,16) NOT NULL COMMENT '买入或卖出量，对于市价买入单表',
                                  `symbol` VARCHAR(64) NOT NULL COMMENT '交易符號',
                                  `traded_amount` DECIMAL(32,16) NOT NULL COMM[]ENT '成交量',
                                  `turnover` DECIMAL(32,16) NOT NULL COMMENT '成交額 對市價買賣有用',
                                  `coin_symbol` VARCHAR(32) NOT NULL COMMENT '币单位',
                                  `bases_symbol` VARCHAR(32) NOT NULL COMMENT '结算单位',
                                  `status` VARCHAR(32) NOT NULL COMMENT '订单状态',
                                  `direction` VARCHAR(16) NOT NULL COMMENT '订单方向',
                                  `price` DECIMAL(32,16) NOT NULL COMMENT '挂单价格',
                                  `time` DATETIME NOT NULL COMMENT '挂单时间',
                                  `completed_time` DATETIME DEFAULT NULL COMMENT '交易完成时间',
                                  `canceled_time` DATETIME DEFAULT NULL COMMENT '取消时间',
                                  `use_discount` CHAR(1) NOT NULL DEFAULT '0' COMMENT '是否使用折扣 0 不使用 1使用',
                                  INDEX `idx_order_id` (`order_id`),
                                  INDEX `idx_member_id` (`member_id`),
                                  INDEX `idx_symbol` (`symbol`),
                                  INDEX `idx_status` (`status`),
                                  INDEX `idx_time` (`time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='撮合訂單表';

-- 撮合交易明細表
CREATE TABLE `trade_detail` (
                                `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'id',
                                `buy_order_id` VARCHAR(64) NOT NULL COMMENT '買方訂單 ID',
                                `sell_order_id` VARCHAR(64) NOT NULL COMMENT '賣方訂單 ID',
                                `price` DECIMAL(32,16) NOT NULL COMMENT '成交價格',
                                `amount` DECIMAL(32,16) NOT NULL COMMENT '成交數量',
                                `symbol` VARCHAR(64) NOT NULL COMMENT '交易符號',
                                `remark` VARCHAR(255) DEFAULT NULL COMMENT '備注',
                                `trade_time` DATETIME NOT NULL COMMENT '审核时间',
                                `created_time` DATETIME NOT NULL COMMENT '创建时间',
                                INDEX `idx_buy_order_id` (`buy_order_id`),
                                INDEX `idx_sell_order_id` (`sell_order_id`),
                                INDEX `idx_symbol_trade_time` (`symbol`, `trade_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='撮合交易明細表';

-- 訂單交易明細表
CREATE TABLE `order_detail` (
                                `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'ID',
                                `order_id` VARCHAR(64) NOT NULL COMMENT '訂單號',
                                `member_id` BIGINT NOT NULL COMMENT '會員ID',
                                `event_type` VARCHAR(32) NOT NULL COMMENT '事件類型(CREATE, TRADE, CANCEL, AMEND, EXPIRE)',
                                `symbol` VARCHAR(64) NOT NULL COMMENT '交易符號',
                                `direction` VARCHAR(16) NOT NULL COMMENT '訂單方向 BUY/SELL',
                                `order_type` VARCHAR(16) NOT NULL COMMENT '訂單類型 LIMIT/MARKET',
                                `price` DECIMAL(32,16) DEFAULT NULL COMMENT '下單價',
                                `amount` DECIMAL(32,16) DEFAULT NULL COMMENT '掛單數量',
                                `traded_amount` DECIMAL(32,16) DEFAULT NULL COMMENT '已成交數量',
                                `untraded_amount` DECIMAL(32,16) DEFAULT NULL COMMENT '未成交數量',
                                `turnover` DECIMAL(32,16) DEFAULT NULL COMMENT '成交金額(成交價×成交量)',
                                `reason` VARCHAR(255) DEFAULT NULL COMMENT '事件原因或備註',
                                `ip_address` VARCHAR(64) DEFAULT NULL COMMENT '操作IP',
                                `device_info` VARCHAR(255) DEFAULT NULL COMMENT '設備信息',
                                `api_key_id` VARCHAR(64) DEFAULT NULL COMMENT 'API Key ID',
                                `created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '事件時間',
                                INDEX `idx_order_id` (`order_id`),
                                INDEX `idx_member_id` (`member_id`),
                                INDEX `idx_symbol_created_time` (`symbol`, `created_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='訂單交易明細表';
