package user

import (
	"database/sql"
	"fmt"

	"github.com/georgiwritescode/vacation-tool/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) FindById(id int) (*types.User, error) {

	rows, err := s.db.Query("SELECT id, first_name, last_name, age, email, vacation_days, non_paid_leave, ts from tbl_users where id = ?", id)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	vacations, err := s.getVacationsByUserId(user.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching vacations: %v", err)
	}
	user.Vacations = vacations

	return user, nil
}

func (s *Store) CreateUser(req *types.User) (int, error) {
	// Set default vacation days if 0 (or trust the DB default, but explicit is nicer if passed)
	if req.VacationDays == 0 {
		req.VacationDays = 20
	}
	res, err := s.db.Exec("INSERT INTO tbl_users (first_name, last_name, age, email, vacation_days, non_paid_leave) values (?, ?, ?, ?, ?, ?)", req.FirstName, req.LastName, req.Age, req.Email, req.VacationDays, req.NonPaidLeave)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func (s *Store) FetchAllUsers() ([]*types.User, error) {
	rows, err := s.db.Query("SELECT id, first_name, last_name, age, email, vacation_days, non_paid_leave, ts from tbl_users")
	if err != nil {
		return nil, err
	}

	var users []*types.User
	for rows.Next() {
		user, err := scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Store) UpdateUser(user *types.User) error {
	_, err := s.db.Exec("UPDATE tbl_users SET first_name=?, last_name=?, age=?, email=?, vacation_days=?, non_paid_leave=? WHERE id=?",
		user.FirstName, user.LastName, user.Age, user.Email, user.VacationDays, user.NonPaidLeave, user.ID)
	return err
}

func (s *Store) DeleteUser(id int) error {
	_, err := s.db.Exec("DELETE FROM tbl_users WHERE id=?", id)
	return err
}

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		&user.Email,
		&user.VacationDays,
		&user.NonPaidLeave,
		&user.Timestamp,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) getVacationsByUserId(userId int) ([]*types.Vacation, error) {
	rows, err := s.db.Query("SELECT * FROM tbl_vacations WHERE person_id = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vacations := make([]*types.Vacation, 0)
	for rows.Next() {
		v := new(types.Vacation)
		err := rows.Scan(
			&v.ID,
			&v.Label,
			&v.FromDate,
			&v.ToDate,
			&v.PersonId,
			&v.Timestamp,
			&v.DaysUsed,
		)
		if err != nil {
			return nil, err
		}
		vacations = append(vacations, v)
	}

	return vacations, nil
}
