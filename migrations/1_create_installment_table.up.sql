CREATE TABLE IF NOT EXISTS installments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    phone_id INT NOT NULL,
    three_months DECIMAL(10, 2) NOT NULL,
    six_months DECIMAL(10, 2) NOT NULL,
    twelve_months DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX installments_phone_id_index (phone_id),
    FOREIGN KEY (phone_id) REFERENCES phones(id) ON DELETE CASCADE
);