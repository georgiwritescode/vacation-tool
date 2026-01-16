package types

type Vacation struct {
	ID        int    `json:"id"`
	Label     string `json:"label"`
	FromDate  string `json:"fromDate"`
	ToDate    string `json:"toDate"`
	PersonId  int    `json:"personId"`
	Timestamp string `json:"ts"`
	DaysUsed  int    `json:"daysUsed"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Email     string      `json:"email"`
	VacationDays int      `json:"vacation_days"`
	NonPaidLeave int      `json:"non_paid_leave"`
	Timestamp string      `json:"ts"`
	Vacations []*Vacation `json:"vacations,omitempty"`
}

type UserStore interface {
	FindById(id int) (*User, error)
	CreateUser(user *User) (int, error)
	FetchAllUsers() ([]*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int) error
}

type VacationStore interface {
	FindById(int) (*Vacation, error)
	FindAll() ([]*Vacation, error)
	CreateVacation(*Vacation) (int, error)
	UpdateVacation(*Vacation) error
	DeleteVacation(int) error
	GetActiveVacations(string) ([]*Vacation, error)
}
