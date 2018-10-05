package transferrer

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// NewStore returns a store that implements the interface Store
func NewStore(db *sqlx.DB) Store {
	return store{db}
}

// Store is the interface that defines the methods to acces the databse
type Store interface {
	Account(string) (Account, error)
	Move(Transfer) error
}

// User encapsulates all information related with a user
type User struct {
	ID        int        `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// Account encapsulates all information related with an account
type Account struct {
	ID        int        `db:"id"`
	Number    string     `db:"number"`
	Balance   float64    `db:"balance"`
	Currency  string     `db:"currency"`
	Owner     string     `db:"owner"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// Movement encapsulates all information related with a movement of money
type Movement struct {
	ID          int        `db:"id"`
	Origin      string     `db:"origin"`
	Destination string     `db:"destination"`
	Amount      float64    `db:"amount"`
	CreatedAt   *time.Time `db:"created_at"`
}

type store struct {
	db *sqlx.DB
}

// Account returns the account of a given user
func (s store) Account(email string) (Account, error) {
	query := `SELECT a.number, a.balance, a.currency, a.owner
				FROM account a
				JOIN "user" u ON a.owner = u.email
				WHERE u.email = $1`

	a := Account{}
	rows, err := s.db.Query(query, email)
	if err != nil {
		log.Println("Error query: ", err)
		return a, err
	}

	for rows.Next() {
		if err := rows.Scan(&a.Number, &a.Balance, &a.Currency, &a.Owner); err != nil {
			log.Println("Error scan: ", err)
			return a, err
		}
		break
	}

	return a, nil
}

// Move makes a transfer of money between two users in a transactional way. Updates the balance of each user and then inserts new register to movements table.
func (s store) Move(t Transfer) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	{
		stmt, err := tx.Prepare(`UPDATE account set balance = balance + $1 where owner = $2`)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()

		if _, err := stmt.Exec(-t.Amount, t.OriginUser); err != nil {
			tx.Rollback() // return an error too, we may want to wrap them
			return err
		}
	}

	{
		stmt, err := tx.Prepare(`UPDATE account set balance = balance + $1 where owner = $2`)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()

		if _, err := stmt.Exec(t.Amount, t.DestinationUser); err != nil {
			tx.Rollback() // return an error too, we may want to wrap them
			return err
		}
	}

	{
		stmt, err := tx.Prepare(`INSERT INTO movement(origin, destination, amount) VALUES ($1, $2, $3)`)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()

		if _, err := stmt.Exec(t.OriginNumber, t.DestinationNumber, t.Amount); err != nil {
			tx.Rollback() // return an error too, we may want to wrap them
			return err
		}
	}

	return tx.Commit()
}
