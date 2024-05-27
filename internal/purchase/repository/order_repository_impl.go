package repository

import "github.com/jmoiron/sqlx"

type orderRepositoryImpl struct {
	db *sqlx.DB
}

func NewOrderRepositoryImpl(db *sqlx.DB) OrderRepository {
	return &orderRepositoryImpl{db: db}
}
