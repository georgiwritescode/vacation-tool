-- Vacation Tool Database Initialization Script
-- This script creates the necessary tables for the vacation management system

-- Create users table
CREATE TABLE IF NOT EXISTS tbl_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    age INT NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    vacation_days INT NOT NULL DEFAULT 20,
    non_paid_leave INT NOT NULL DEFAULT 0,
    ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create vacations table
CREATE TABLE IF NOT EXISTS tbl_vacations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    label VARCHAR(255) NOT NULL,
    from_date DATE NOT NULL,
    to_date DATE NOT NULL,
    person_id INT NOT NULL,
    ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    days_used INT NOT NULL DEFAULT 0,
    FOREIGN KEY (person_id) REFERENCES tbl_users(id) ON DELETE CASCADE,
    INDEX idx_person_id (person_id),
    INDEX idx_dates (from_date, to_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;