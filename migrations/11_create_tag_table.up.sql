CREATE TABLE tags (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE phone_tags (
    phone_id INT,
    tag_id INT,
    FOREIGN KEY (phone_id) REFERENCES phones(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id),
    PRIMARY KEY (phone_id, tag_id)
);
