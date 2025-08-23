CREATE DATABASE IF NOT EXISTS `order`;
CREATE DATABASE IF NOT EXISTS `payment`;

CREATE TABLE IF NOT EXISTS inventory_items (
    id   VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL
    );

INSERT INTO inventory_items (id, name) VALUES
        ('SKU-1', 'Item 1'),
        ('SKU-2', 'Item 2'),
        ('SKU-3', 'Item 3')
    ON DUPLICATE KEY UPDATE name = VALUES(name);