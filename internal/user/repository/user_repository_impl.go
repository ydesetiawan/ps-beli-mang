package repository

import (
	"github.com/jmoiron/sqlx"
	"ps-beli-mang/internal/user/model"
	"ps-beli-mang/pkg/errs"
	"strings"
)

type userRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepositoryImpl(db *sqlx.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

const (
	queryGetUserByIDAndRole       = `SELECT * FROM users WHERE id = $1 AND role = $2`
	queryGetUserByUsernameAndRole = `SELECT * FROM users WHERE username = $1 AND role = $2`
	queryInsertUser               = `
	WITH new_user AS (
		SELECT $1::char(26) AS id,
		       $2::varchar(50) AS username,
		       $3::varchar(100) AS password,
		       $4::varchar(255) AS email,
		       $5::varchar(10) AS role
	)
	INSERT INTO users (id, username, password, email, role)
	SELECT id, username, password, email, role
	FROM new_user
	WHERE NOT EXISTS (
		SELECT 1
		FROM users
		WHERE email = new_user.email
		  AND role = new_user.role
	)
	RETURNING id;
	`
)

func (r *userRepositoryImpl) GetUserByIDAndRole(id string, role string) (model.User, error) {
	var user model.User
	err := r.db.Get(&user, queryGetUserByIDAndRole, id, role)
	if err != nil {
		return model.User{}, errs.NewErrInternalServerErrors("execute query error [GetUserByIDAndRole]: ", err.Error())
	}
	return user, err
}

func (r *userRepositoryImpl) GetUserByUsernameAndRole(username string, role string) (model.User, error) {
	var user model.User
	err := r.db.Get(&user, queryGetUserByUsernameAndRole, username, role)
	if err != nil {
		return model.User{}, errs.NewErrInternalServerErrors("execute query error [GetUserByUsernameAndRole]: ", err.Error())
	}
	return user, err
}

func (r *userRepositoryImpl) Register(user *model.User) (string, error) {
	var lastInsertId = ""
	err := r.db.QueryRowx(queryInsertUser, user.ID, user.Username, user.Password, user.Email, user.Role).Scan(&lastInsertId)
	if err != nil {
		if strings.Contains(err.Error(), "users_username_key") {
			return lastInsertId, errs.NewErrDataConflict("username already exist", user.Username)
		}
		return lastInsertId, errs.NewErrDataConflict("execute query error [RegisterUser]: ", err.Error())
	}

	return lastInsertId, nil
}
