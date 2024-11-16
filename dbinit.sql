CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    balance INT8
);


CREATE TABLE users (
    id SERIAL PRIMARY KEY,                      
    name VARCHAR(100) NOT NULL,                
    email VARCHAR(150) UNIQUE NOT NULL,         
    password VARCHAR(255) NOT NULL,             
    wallet_pin VARCHAR(255) NOT NULL,           
    is_approved BOOLEAN DEFAULT FALSE,          
    wallet_balance NUMERIC(10, 2) DEFAULT 0.0   
);

CREATE TABLE vendors (
    id SERIAL PRIMARY KEY,                      
    user_name VARCHAR(100) UNIQUE NOT NULL,     
    stall_name VARCHAR(150) NOT NULL,           
    password VARCHAR(255) NOT NULL              
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,                      
    request_id UUID DEFAULT gen_random_uuid(),  
    transaction_id UUID DEFAULT gen_random_uuid(), 
    transaction_type INT NOT NULL,              
    amount NUMERIC(10, 2) NOT NULL,             
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    from_user_id INT REFERENCES users(id),      
    to_vendor_id INT REFERENCES vendors(id)    
);
