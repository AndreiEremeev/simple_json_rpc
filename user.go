package main

import (
	"database/sql"
	"fmt"
	//"time"
)

const (
	NotFound      = "user not found"
	InternalError = "internal error"
	AlreadyExists = "user already exists"
	ShortLogin    = "Login too short"
)

type User struct {
	Login     string  `json:"login"`
	Id        string  `json:"id"`
	CreatedAt float64 `json:"createdAt"`
}

type UsersLogin struct {
	Login string `json:login`
}

type UsersId struct {
	Id string `json:id`
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
	switch err := row.Scan(&result.Id, &result.Login, &result.CreatedAt); err {
	case sql.ErrNoRows:
		return fmt.Errorf(NotFound)
	case nil:
		*out = result
		return nil
	default:
		return fmt.Errorf(InternalError)
	}
}

func (um *UserManager) GetUserById(in *UsersId, out *User) error {
	var result User
	row := um.DB.QueryRow("select id::text, login, extract(epoch from created_at) from users where id::text = $1", in.Id)
	switch err := row.Scan(&result.Id, &result.Login, &result.CreatedAt); err {
	case sql.ErrNoRows:
		return fmt.Errorf(NotFound)
	case nil:
		*out = result
		return nil
	default:
		return fmt.Errorf(InternalError)
	}
}

func (um *UserManager) CreateUser(in *UsersLogin, out *User) error {
	if len(in.Login) < 4 {
		return fmt.Errorf(ShortLogin)
	}
	var result User
	row := um.DB.QueryRow("insert into users(id, login, created_at) values(uuid_generate_v4(), $1, now()) on conflict(login) do nothing returning id::text, login, extract(epoch from created_at)", in.Login)
	switch err := row.Scan(&result.Id, &result.Login, &result.CreatedAt); err {
	case sql.ErrNoRows:
		return fmt.Errorf(AlreadyExists)
	case nil:
		*out = result
		return nil
	default:
		return fmt.Errorf(InternalError)
	}
}

func (um *UserManager) EditUser(in *User, out *User) error {
	var result User
	err := um.GetUserByLogin(&UsersLogin{Login: in.Login}, &result)
	if err != nil {
		if err.Error() != NotFound {
			return err
		}
	} else {
		if result.Id != in.Id {
			return fmt.Errorf(AlreadyExists)
		}
	}
	row := um.DB.QueryRow("update users set login = $2, created_at = to_timestamp($3) where id::text = $1 returning id::text, login, extract(epoch from created_at)", in.Id, in.Login, in.CreatedAt)
	switch err = row.Scan(&result.Id, &result.Login, &result.CreatedAt); err {
	case sql.ErrNoRows:
		return fmt.Errorf(NotFound)
	case nil:
		*out = result
		return nil
	default:
		return fmt.Errorf(InternalError)
	}
}
