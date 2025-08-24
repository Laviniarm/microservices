-- Databases
CREATE DATABASE IF NOT EXISTS `orders`;
CREATE DATABASE IF NOT EXISTS `payment`;

-- Estoque (orders)
USE `orders`;

CREATE TABLE IF NOT EXISTS inventory_items (
                                               id   VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL
    );

INSERT INTO inventory_items (id, name) VALUES
                                           ('SKU-1', 'Item 1'),
                                           ('SKU-2', 'Item 2'),
                                           ('SKU-3', 'Item 3')
    ON DUPLICATE KEY UPDATE name = VALUES(name);

-- Tabelas para persistir pedidos
CREATE TABLE IF NOT EXISTS orders (
                                      id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                      customer_id BIGINT NOT NULL,
                                      status VARCHAR(32) NOT NULL,
    created_at BIGINT NOT NULL
    );

CREATE TABLE IF NOT EXISTS order_items (
                                           id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                           order_id BIGINT NOT NULL,
                                           product_code VARCHAR(64) NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    quantity INT NOT NULL,
    CONSTRAINT fk_order_items_order
    FOREIGN KEY (order_id) REFERENCES orders(id)
    ON DELETE CASCADE
    );