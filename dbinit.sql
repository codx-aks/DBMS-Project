CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    balance INT8
);

CREATE TABLE users (
    name VARCHAR(100) NOT NULL,
    roll_no VARCHAR(9) PRIMARY KEY,
    pin INT NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    is_approved BOOLEAN DEFAULT FALSE,
    wallet_rem INT DEFAULT 0,
    other_details VARCHAR(255)
);

CREATE TABLE otp (
    roll_no VARCHAR(9) PRIMARY KEY REFERENCES users(roll_no),
    otp_last_generated INT NOT NULL,
    otp_last_generated_time DATETIME NOT NULL
);

CREATE TABLE vendors (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    image_url VARCHAR(255)
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,                       
    transaction_id UUID DEFAULT gen_random_uuid(),                         
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sender VARCHAR(9) NOT NULL,
    receiver INT NOT NULL,
    credit INT DEFAULT 0,
    debit INT DEFAULT 0,
    description VARCHAR(255),
    FOREIGN KEY (sender_roll_no) REFERENCES users(roll_no),
    FOREIGN KEY (receiver_vendor_id) REFERENCES vendors(id)
);

CREATE TABLE items (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    cost INT NOT NULL,
    image_url VARCHAR(255),
    description VARCHAR(255),
    vendor_id INT NOT NULL,
    FOREIGN KEY (vendor_id) REFERENCES vendors(id)
);