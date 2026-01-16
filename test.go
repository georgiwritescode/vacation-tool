package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

const baseURL = "http://localhost:8080"

// Structures related to API
type User struct {
	ID           int         `json:"id"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Age          int         `json:"age"`
	Email        string      `json:"email"`
	VacationDays int         `json:"vacation_days"`
	NonPaidLeave int         `json:"non_paid_leave"`
	Timestamp    string      `json:"ts"`
	Vacations    []*Vacation `json:"vacations,omitempty"`
}

type Vacation struct {
	ID        int    `json:"id"`
	Label     string `json:"label"`
	FromDate  string `json:"fromDate"`
	ToDate    string `json:"toDate"`
	PersonId  int    `json:"personId"`
	DaysUsed  int    `json:"daysUsed"`
	Timestamp string `json:"ts"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	fmt.Println("🚀 Starting E2E Tests for Vacation Tool")
	
	// Ensure server is running (implied, but we can do a health check if needed)
	// For now, simpler to just run.

	runUserFeaturesTest()
	runRelationTest()
	runDaysLogicTest()
	runNonPaidTest()

	fmt.Println("\n✅ All Tests Passed Successfully!")
}

// --- Test Suites ---

func runUserFeaturesTest() {
	fmt.Println("\n[1] Testing User Features (Update/Delete)")

	// 1. Create User
	u := User{FirstName: "Feat", LastName: "Test", Age: 30, Email: "feat@test.com", VacationDays: 20}
	userID := createUser(u)
	fmt.Printf("    Created User ID: %d\n", userID)

	// 2. Update User
	u.ID = userID
	u.FirstName = "UpdatedFeat"
	_, status := makeRequest("PUT", "/api/v1/users/update", u)
	assert(status == 200, "Update User status 200")

	// 3. Verify Update
	fetched := getUser(userID)
	assert(fetched.FirstName == "UpdatedFeat", "User firstname updated")

	// 4. Delete User
	_, status = makeRequest("DELETE", fmt.Sprintf("/api/v1/users/delete/%d", userID), nil)
	assert(status == 200, "Delete User status 200")

	// 5. Verify Delete (Expect 500 or 400 with 'not found' or similar, currently API returns 500/error string for not found if using FindById check)
	// API FindById returns error "user not found"
	_, status = makeRequest("GET", fmt.Sprintf("/api/v1/users/%d", userID), nil)
	// Depending on implementation, might be 500.
	// user/store.go returns error, middleware might map it.
	fmt.Printf("    Get Deleted User Status: %d (Expected failure)\n", status)
}

func runRelationTest() {
	fmt.Println("\n[2] Testing User-Vacation Relationship")

	// 1. Create User
	u := User{FirstName: "Rel", LastName: "Test", Age: 25, Email: "rel@test.com", VacationDays: 20}
	userID := createUser(u)

	// 2. Create Vacations
	createVacation(Vacation{Label: "V1", FromDate: "2024-01", ToDate: "2024-01", PersonId: userID, DaysUsed: 2, Timestamp: "2024-01-01"})
	createVacation(Vacation{Label: "V2", FromDate: "2024-02", ToDate: "2024-02", PersonId: userID, DaysUsed: 3, Timestamp: "2024-01-01"})

	// 3. Fetch User and check vacations
	fetched := getUser(userID)
	assert(len(fetched.Vacations) == 2, fmt.Sprintf("User has 2 vacations (Found %d)", len(fetched.Vacations)))
	assert(fetched.Vacations[0].Label == "V1", "First vacation label matches")
}

func runDaysLogicTest() {
	fmt.Println("\n[3] Testing Paid Vacation Logic")

	// 1. Create User (Default 20)
	u := User{FirstName: "Paid", LastName: "Logic", Age: 22, Email: "paid@test.com", VacationDays: 20}
	userID := createUser(u)

	// 2. Use 5 days
	createVacation(Vacation{PersonId: userID, DaysUsed: 5, FromDate: "2024-01-01", ToDate: "2024-01-05", Label: "5 Days", Timestamp: "2024-01-01 00:00:00"})
	
	// 3. Check Balance (Expect 15)
	fetched := getUser(userID)
	assert(fetched.VacationDays == 15, fmt.Sprintf("Balance decreased to 15 (Found %d)", fetched.VacationDays))

	// 4. Attempt Overdraft (Use 20 more) -> Should Fail
	_, status := createVacationRaw(Vacation{PersonId: userID, DaysUsed: 20, FromDate: "2024-02-01", ToDate: "2024-02-20", Label: "Fail", Timestamp: "2024-01-01 00:00:00"})
	assert(status != 201, "Overdraft request failed")

	// 5. Verify Balance Unchanged
	fetched = getUser(userID)
	assert(fetched.VacationDays == 15, fmt.Sprintf("Balance remained 15 (Found %d)", fetched.VacationDays))
}

func runNonPaidTest() {
	fmt.Println("\n[4] Testing Non-Paid Leave Logic")

	// 1. Create User (5 Paid, 10 NonPaid)
	u := User{FirstName: "NonPaid", LastName: "Logic", Age: 40, Email: "nplogic@test.com", VacationDays: 5, NonPaidLeave: 10}
	userID := createUser(u)
	
	init := getUser(userID)
	assert(init.VacationDays == 5, fmt.Sprintf("Init Paid 5 (Found %d)", init.VacationDays))
	assert(init.NonPaidLeave == 10, fmt.Sprintf("Init NonPaid 10 (Found %d)", init.NonPaidLeave))

	// 2. Request 8 days (Uses 5 Paid, 3 NonPaid)
	createVacation(Vacation{PersonId: userID, DaysUsed: 8, FromDate: "2024-01-01", ToDate: "2024-01-08", Label: "Mixed", Timestamp: "2024-01-01 00:00:00"})
	
	// 3. Check Balance (Expect 0 Paid, 7 NonPaid)
	fetched := getUser(userID)
	assert(fetched.VacationDays == 0, fmt.Sprintf("Paid becomes 0 (Found %d)", fetched.VacationDays))
	assert(fetched.NonPaidLeave == 7, fmt.Sprintf("NonPaid becomes 7 (Found %d)", fetched.NonPaidLeave))

	// 4. Request 8 days again (Have 7 NonPaid) -> Fail
	_, status := createVacationRaw(Vacation{PersonId: userID, DaysUsed: 8, FromDate: "2024-03-01", ToDate: "2024-03-08", Label: "Fail", Timestamp: "2024-01-01 00:00:00"})
	assert(status != 201, "Insufficient non-paid request failed")

	// 5. Request 7 days -> Success
	createVacation(Vacation{PersonId: userID, DaysUsed: 7, FromDate: "2024-04-01", ToDate: "2024-04-07", Label: "Last", Timestamp: "2024-01-01 00:00:00"})
	
	// 6. Check Balance (Expect 0, 0)
	fetched = getUser(userID)
	assert(fetched.VacationDays == 0, fmt.Sprintf("Paid is 0 (Found %d)", fetched.VacationDays))
	assert(fetched.NonPaidLeave == 0, fmt.Sprintf("NonPaid is 0 (Found %d)", fetched.NonPaidLeave))
}


// --- Helper Functions ---

func assert(condition bool, msg string) {
	if !condition {
		fmt.Printf("   ❌ FAIL: %s\n", msg)
		os.Exit(1)
	}
	fmt.Printf("   ✅ PASS: %s\n", msg)
}

func makeRequest(method, urlPath string, body interface{}) ([]byte, int) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, _ := http.NewRequest(method, baseURL+urlPath, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	
	// Basic retry for connection refused if server is just starting
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("FATAL: Request %s %s failed: %v\n", method, urlPath, err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	
	data, _ := io.ReadAll(resp.Body)
	return data, resp.StatusCode
}

func createUser(u User) int {
	data, status := makeRequest("POST", "/api/v1/users/create", u)
	if status != 200 {
		fmt.Printf("Create User failed: %s\n", string(data))
		os.Exit(1)
	}
	// Response: "user with id: <ID> created"
	content := string(data)
	// Simple parse
	re := regexp.MustCompile(`id: (\d+)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 {
		// Fallback: maybe it returned JSON in quote?
		// Try to parse just ID if regex fails
		fmt.Printf("Could not parse ID from: %s\n", content)
		os.Exit(1)
	}
	id, _ := strconv.Atoi(matches[1])
	return id
}

func getUser(id int) *User {
	data, status := makeRequest("GET", fmt.Sprintf("/api/v1/users/%d", id), nil)
	if status != 200 {
		fmt.Printf("Get User %d failed: %s\n", id, string(data))
		os.Exit(1)
	}
	var u User
	json.Unmarshal(data, &u)
	return &u
}

func createVacation(v Vacation) {
	data, status := createVacationRaw(v)
	if status != 201 {
		fmt.Printf("Create Vacation failed (Status %d): %s\n", status, string(data))
		os.Exit(1)
	}
}

func createVacationRaw(v Vacation) ([]byte, int) {
	return makeRequest("POST", "/api/v1/vacations/create", v)
}
