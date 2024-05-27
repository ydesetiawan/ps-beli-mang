package repository

import "github.com/jmoiron/sqlx"

type merchantRepositoryImpl struct {
	db *sqlx.DB
}

func NewMerchantRepositoryImpl(db *sqlx.DB) MerchantRepository {
	return &merchantRepositoryImpl{db: db}
}
