CREATE TABLE IF NOT EXISTS price_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    phone_id INT NOT NULL,
    old_price DECIMAL(10, 2) NOT NULL,
    new_price DECIMAL(10, 2) NOT NULL,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    INDEX price_history_phone_id_index (phone_id),
    FOREIGN KEY (phone_id) REFERENCES phones(id) ON DELETE CASCADE
);