CREATE TABLE tbl_persons(  
    id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    age INT,
    email  VARCHAR(255),
    ts TIMESTAMP DEFAULT CURRENT_TIME ON UPDATE CURRENT_TIMESTAMP
);

CREATE Table tbl_vacations(
        id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(255),
        from_date  VARCHAR(255),
        to_date  VARCHAR(255),
        person_id INT,
        ts TIMESTAMP DEFAULT CURRENT_TIME ON UPDATE CURRENT_TIMESTAMP,
        Foreign Key (person_id) REFERENCES tbl_persons(id) 
);

INSERT INTO tbl_persons (ID, first_name, last_name, Age, Email) VALUES
(1, 'John', 'Doe', 30, 'john.doe@example.com'),
(2, 'Alice', 'Smith', 25, 'alice.smith@example.com'),
(3, 'Bob', 'Johnson', 40, 'bob.johnson@example.com'),
(4, 'Emily', 'Brown', 35, 'emily.brown@example.com'),
(5, 'Michael', 'Davis', 28, 'michael.davis@example.com');

INSERT INTO tbl_vacations (id, name, from_date, to_date, person_id) VALUES
(1, 'Vacation 1', '2024-04-06', '2024-04-10', 1),
(2, 'Vacation 2', '2024-05-15', '2024-05-20', 2),
(3, 'Vacation 3', '2024-06-01', '2024-06-05', 3),
(4, 'Vacation 4', '2024-07-10', '2024-07-15', 4),
(5, 'Vacation 5', '2024-08-20', '2024-08-25', 5),
(6, 'Vacation 6', '2024-09-06', '2024-09-10', 1),
(7, 'Vacation 7', '2024-10-15', '2024-10-20', 2),
(8, 'Vacation 8', '2024-11-01', '2024-11-05', 3),
(9, 'Vacation 9', '2024-12-10', '2024-12-15', 4),
(10, 'Vacation 10', '2025-01-20', '2025-01-25', 5),
(11, 'Vacation 11', '2025-02-06', '2025-02-10', 1),
(12, 'Vacation 12', '2025-03-15', '2025-03-20', 2),
(13, 'Vacation 13', '2025-04-01', '2025-04-05', 3),
(14, 'Vacation 14', '2025-05-10', '2025-05-15', 4),
(15, 'Vacation 15', '2025-06-20', '2025-06-25', 5),
(16, 'Vacation 16', '2025-07-06', '2025-07-10', 1),
(17, 'Vacation 17', '2025-08-15', '2025-08-20', 2),
(18, 'Vacation 18', '2025-09-01', '2025-09-05', 3),
(19, 'Vacation 19', '2025-10-10', '2025-10-15', 4),
(20, 'Vacation 20', '2025-11-20', '2025-11-25', 5);