CREATE DATABASE IF NOT EXISTS BillingEngine;

USE BillingEngine;

-- Create the loans table
CREATE TABLE loans
(
	id               BIGINT AUTO_INCREMENT PRIMARY KEY,
	reference_id     VARCHAR(255) NOT NULL,
	user_id          BIGINT       NOT NULL,
	amount           BIGINT       NOT NULL,
	rate_percentage  INT          NOT NULL,
	repayment_amount BIGINT       NOT NULL,
	status           INT          NOT NULL,
	created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at       TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
);

-- Create the repayments table
CREATE TABLE repayments
(
	id           BIGINT AUTO_INCREMENT PRIMARY KEY,
	loan_id      BIGINT       NOT NULL,
	reference_id VARCHAR(255) NOT NULL,
	amount       BIGINT       NOT NULL,
	created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at   TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
);

-- Add indexes for faster queries in descending order
CREATE INDEX idx_user_id ON loans (user_id DESC);
CREATE INDEX idx_reference_id ON loans (reference_id DESC);
CREATE INDEX idx_loan_id ON repayments (loan_id DESC);
CREATE INDEX idx_reference_id ON repayments (reference_id DESC);
