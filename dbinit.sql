CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    balance INT8
);

CREATE TABLE users (
    name VARCHAR(100) NOT NULL,
    roll_no VARCHAR(9) PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    pin VARCHAR(255) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    is_approved BOOLEAN DEFAULT FALSE,
    wallet_balance INT DEFAULT 0
);

CREATE TABLE otp (
    roll_no VARCHAR(9) PRIMARY KEY REFERENCES users(roll_no),
    otp_last_generated INT NOT NULL,
    otp_last_generated_time TIMESTAMP NOT NULL
);

CREATE TABLE vendors (
    id SERIAL PRIMARY KEY,
    is_active BOOLEAN DEFAULT FALSE,
    password VARCHAR(255) NOT NULL,
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
    amount INT DEFAULT 0,
    description VARCHAR(255),
    FOREIGN KEY (sender) REFERENCES users(roll_no),
    FOREIGN KEY (receiver) REFERENCES vendors(id)
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    is_available BOOLEAN DEFAULT FALSE,
    cost INT NOT NULL,
    image_url VARCHAR(255),
    description VARCHAR(255),
    vendor_id INT NOT NULL,
    FOREIGN KEY (vendor_id) REFERENCES vendors(id)
);