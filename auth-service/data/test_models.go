package data

import (
	"database/sql"
	"time"
)


type PostgresTestRepository struct {
	Conn *sql.DB
}

func NewPostgresTestRepository(db *sql.DB) *PostgresTestRepository {
	return &PostgresTestRepository{
		Conn: db,
	}
}


// GetAll returns a slice of all users, sorted by last name
func (r  *PostgresTestRepository) GetAll() ([]*User, error) {
	return []*User{}, nil
}

// GetByEmail returns one user by email
func (r  *PostgresTestRepository) GetByEmail(email string) (*User, error) {
	return &User{
		ID: 1,
		FirstName: "First",
		LastName: "Last",
		Email: "me@here.com",
		Password: "verysecret",
		Active: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// GetOne returns one user by id
func (r  *PostgresTestRepository) GetOne(id int) (*User, error) {
	return &User{
		ID: 1,
		FirstName: "First",
		LastName: "Last",
		Email: "me@here.com",
		Password: "",
		Active: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (r  *PostgresTestRepository) Update(user User) error {
	return nil
}

// Delete deletes one user from the database, by User.ID
// func (r  *PostgresTestRepository) Delete() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
// 	defer cancel()

// 	stmt := `delete from users where id = $1`

// 	_, err := db.ExecContext(ctx, stmt, u.ID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// DeleteByID deletes one user from the database, by ID
func (r  *PostgresTestRepository) DeleteByID(id int) error {
	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (r  *PostgresTestRepository) Insert(user User) (int, error) {
	return 2, nil
}

// ResetPassword is the method we will use to change a user's password.
func (r  *PostgresTestRepository) ResetPassword(password string, user User) error {
	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (r  *PostgresTestRepository) PasswordMatches(plainText string, user User) (bool, error) {
	return true, nil
}
