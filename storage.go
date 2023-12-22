package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(int, *Account) error
	GetAccount() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=gobank password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	// err := s.dropAccountTable()
	// if err != nil {
	// 	return err
	// }
	return s.createAccountTable()
}

func (s *PostgresStore) dropAccountTable() error {
	query := `DROP TABLE IF EXISTS accounts CASCADE`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		encrypted_password varchar(100),
		balance serial,
		created_at timestamp DEFAULT current_timestamp,
		updated_at timestamp DEFAULT current_timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(a *Account) error {
	query := `INSERT INTO accounts 
		(first_name, last_name, number, encrypted_password, balance) 
		VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.Query(
		query,
		a.FirstName,
		a.LastName,
		a.Number,
		a.EncryptedPassword,
		a.Balance,
	)

	return err
}

func (s *PostgresStore) DeleteAccount(id int) error {
	rows, err := s.db.Query(`DELETE FROM accounts WHERE id = $1 RETURNING id`, id)
	if err != nil {
		return err
	}

	return rowsEffected(rows, id)
}

func (s *PostgresStore) UpdateAccount(id int, a *Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM accounts WHERE number = $1`, number)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with number %d not found", number)
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM accounts WHERE id = $1`, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with id [%d] not found", id)
}

func (s *PostgresStore) GetAccount() ([]*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM accounts`)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)

	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	return account, err
}

func rowsEffected(rows *sql.Rows, id int) error {
	var rowsAffected int
	for rows.Next() {
		rowsAffected++
	}

	fmt.Println("rowsAffected: ", rowsAffected)

	if rowsAffected == 0 {
		return fmt.Errorf("no account with ID %d found", id)
	}

	return nil
}
