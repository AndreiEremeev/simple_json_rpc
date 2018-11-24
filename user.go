package main

import (
	"database/sql"
	"fmt"
	//"time"
)

var (
	ErrNotFound      = fmt.Errorf("user not found")
	ErrInternalError = fmt.Errorf("internal error")
	ErrAlreadyExists = fmt.Errorf("user already exists")
	ErrShortLogin    = fmt.Errorf("login too short")
)

type User struct {
	Login     string  `json:"login"`
	ID        string  `json:"id"`
	CreatedAt float64 `json:"createdAt"`
}

type UsersLogin struct {
	Login string `json:login`
}

type UsersID struct {
	ID string `json:id`
}

type UserManager struct {
	DB *sql.DB
}

func NewUserManager(db *sql.DB) *UserManager {
	return &UserManager{DB: db}
}

func (um *UserManager) GetUserByLogin(in *UsersLogin, out *User) error {
	var result User
	row := um.DB.QueryRow("select id::text, login, extract(epoch from created_at) from users where login = $1", in.Login)
	switch err := row.Scan(&result.ID, &result.Login, &result.CreatedAt); err {
	case sql.ErrNoRows:
		return ErrNotFound
	case nil:
		*out = result
		return nil
	default:
		return ErrInternalError
	}
}

func (um *UserManager) GetUserByID(in *UsersID, out *User) error {
	var result User
	row := um.DB.QueryRow("select id::text, login, extract(epoch from created_at) from users where id::text = $1", in.ID)
	switch err := row.Scan(&result.ID, &result.Login, &result.CreatedAt); err {
	case sql.ErrNoRows:
		return ErrNotFound
	case nil:
		*out = result
		return nil
	default:
		return ErrInternalError
	}
}

func (um *UserManager) CreateUser(in *UsersLogin, out *User) error {
	if len(in.Login) < 4 {
		return ErrShortLogin
	}
	var result User
	row := um.DB.QueryRow("insert into users(id, login, created_at) values(uuid_generate_v4(), $1, now()) on conflict(login) do nothing returning id::text, login, extract(epoch from created_at)", in.Login)
	switch err := row.Scan(&result.ID, &result.Login, &result.CreatedAt); err {
	case sql.ErrNoRows:
		return ErrAlreadyExists
	case nil:
		*out = result
		return nil
	default:
		return ErrInternalError
	}
}

func (um *UserManager) EditUser(in *User, out *User) error {
	var result User
	err := um.GetUserByLogin(&UsersLogin{Login: in.Login}, &result)
	if err != nil && err != ErrNotFound {
		return err
	} else {
		if err == nil && result.ID != in.ID {
			return ErrAlreadyExists
		}
	}
	row := um.DB.QueryRow("update users set login = $2, created_at = to_timestamp($3) where id::text = $1 returning id::text, login, extract(epoch from created_at)", in.ID, in.Login, in.CreatedAt)
	switch err = row.Scan(&result.ID, &result.Login, &result.CreatedAt); err {
	case sql.ErrNoRows:
		return ErrNotFound
	case nil:
		*out = result
		return nil
	default:
		return ErrInternalError
	}
}
