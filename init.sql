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

-- Insert sample data (optional - remove if not needed)
INSERT INTO tbl_users (first_name, last_name, age, email, vacation_days, non_paid_leave) VALUES
('John', 'Doe', 30, 'john.doe@example.com', 20, 0),
('Alice', 'Smith', 28, 'alice.smith@example.com', 20, 0),
('Bob', 'Johnson', 35, 'bob.johnson@example.com', 15, 5);

INSERT INTO tbl_vacations (label, from_date, to_date, person_id, days_used) VALUES
('Summer Vacation', '2026-07-01', '2026-07-10', 1, 10),
('Christmas Break', '2026-12-24', '2026-12-31', 2, 8);
