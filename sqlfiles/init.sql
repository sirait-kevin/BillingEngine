CREATE DATABASE IF NOT EXISTS BillingEngine;

USE BillingEngine;

-- Create the 'users' table with necessary fields and indexes
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    encrypted_name BLOB NOT NULL,
    encrypted_email BLOB NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    );
