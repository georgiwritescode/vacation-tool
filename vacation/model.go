package vacation

type Vacation struct {
	ID       int    `json:"id"`
	Label    string `json:"label"`
	FromDate string `json:"fromDate"`
	ToDate   string `json:"toDate"`
	PersonId int    `json:"personId"`
}
