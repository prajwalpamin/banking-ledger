package postgres

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/prajwalpamin/banking-ledger/internal/domain"
)

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *accountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(account *domain.Account) error {
	query := `
        INSERT INTO accounts (id, balance, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
    `

	_, err := r.db.Exec(query,
		account.ID,
		account.Balance,
		account.CreatedAt,
		account.UpdatedAt,
	)

	return err
}

func (r *accountRepository) GetByID(id uuid.UUID) (*domain.Account, error) {
	query := `
        SELECT id, balance, created_at, updated_at
        FROM accounts
        WHERE id = $1
    `

	account := &domain.Account{}
	err := r.db.QueryRow(query, id).Scan(
		&account.ID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return account, err
}

func (r *accountRepository) Update(account *domain.Account) error {
	query := `
        UPDATE accounts
        SET balance = $1, updated_at = $2
        WHERE id = $3
    `

	account.UpdatedAt = time.Now()
	_, err := r.db.Exec(query,
		account.Balance,
		account.UpdatedAt,
		account.ID,
	)

	return err
}
