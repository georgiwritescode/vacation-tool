package types

type Vacation struct {
	ID       int    `json:"id"`
	Label    string `json:"label"`
	FromDate string `json:"fromDate"`
	ToDate   string `json:"toDate"`
	PersonId int    `json:"personId"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
	Timestamp string `json:"ts"`
}

type UserStore interface {
	FindById(id int) (*User, error)
}

type VacationStore interface {
	FindById(int) (*Vacation, error)
	FindAll() ([]*Vacation, error)
	CreateVacation(*Vacation) (*Vacation, error)
	UpdateVacation(*Vacation) (*Vacation, error)
	DeleteVacation(int) error
}
