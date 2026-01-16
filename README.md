# Vacation Tool

## Purpose
The **Vacation Tool** is a comprehensive application designed to streamline the management of employee leave. It features a **Web Dashboard** for easy management and a **robust JSON API** for integration. It handles user profiles and vacation requests with intelligent balance tracking, ensuring accurate accounting of paid and non-paid leave.

## Key Features

### 1. User Management
*   **CRUD Operations**: Create, Read, Update, and Delete employee profiles.
*   **Digital Profile**: Stores essential details like Name, Age, and Email.
*   **Leave Balances**: Tracks two types of leave:
    *   **Vacation Days (Paid)**: Standard paid time off (Default: 20 days).
    *   **Non-Paid Leave (Unpaid)**: Additional unpaid days off (Default: 0).

### 2. Web Interface
Use the browser to view and manage data:
*   **Dashboard** (`http://localhost:8080/`): Overview and quick links.
*   **Employees** (`http://localhost:8080/users`): List of all employees and their balances.
*   **Vacations** (`http://localhost:8080/vacations`): List of all scheduled vacations.

### 3. API & Vacation Tracking
*   **JSON API**: All data endpoints are namespaced under `/api/v1/`.
*   **Vacation Requests**: Users can request vacations with specific dates and duration.
*   **History**: Full history of vacations is linked to each user.

### 4. Smart Deduction Logic
The system enforces strict business rules for leave consumption:
1.  **Deduct Paid First**: Requests automatically consume `Vacation Days` first.
2.  **Fallback to Unpaid**: If `Vacation Days` are exhausted, the remaining duration is deducted from `Non-Paid Leave`.
3.  **Overdraft Protection**: If the User lacks sufficient *total* days (Paid + Unpaid) to cover the request, the vacation is rejected.

## Technology Stack
*   **Language**: Go (Golang)
*   **Database**: MariaDB
*   **Containerization**: Docker & Docker Compose

## Getting Started

### Prerequisites
*   Go 1.22+
*   Docker & Docker Compose

### Running the Application

1.  **Start the Database**:
    ```bash
    docker-compose up -d
    ```

2.  **Run the Server**:
    Using Makefile:
    ```bash
    make run
    ```
    Or manually:
    ```bash
    go run cmd/main.go
    ```
    The server will start on `http://localhost:8080`.

### Running Tests
A comprehensive E2E test suite verifies all logic.
```bash
go run test.go
```
This will run a series of tests to verify CRUD operations, database relationships, and the complex leave deduction logic.
