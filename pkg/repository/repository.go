package repository

import (
	"database/sql"
)

type (
	Repository interface {
		Save(query string, args []any) error
		FindAll(query string) (*sql.Rows, error)
		FindByID(query string, args []any) (*sql.Rows, error)
	}
	RepositoryDatabase struct {
		db *sql.DB
	}
)

var _ Repository = (*RepositoryDatabase)(nil)

func NewRepositoryDatabase(db *sql.DB) *RepositoryDatabase {
	return &RepositoryDatabase{
		db: db,
	}
}

func (r *RepositoryDatabase) Save(query string, args []any) error {
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepositoryDatabase) FindAll(query string) (*sql.Rows, error) {
	return r.db.Query(query)
}

func (r *RepositoryDatabase) FindByID(query string, args []any) (*sql.Rows, error) {
	return r.db.Query(query, args)
}
