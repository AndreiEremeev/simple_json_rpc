package main

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func initDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	database, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err.Error())
	}
	return database
}

func Equal(a *User, b *User) bool {
	return (a.ID == b.ID && a.Login == b.Login && a.CreatedAt == b.CreatedAt)
}

func TestValidCreate(t *testing.T) {
	db := initDB()
	defer db.Close()
	um := NewUserManager(db)
	in := &UsersLogin{Login: "test"}
	test := &User{}
	err := um.CreateUser(in, test)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if test.Login != in.Login {
		t.Errorf("created User with Login != in.Login")
		return
	}
	testByLogin := &User{}
	err = um.GetUserByLogin(&UsersLogin{test.Login}, testByLogin)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	testByID := &User{}
	err = um.GetUserByID(&UsersID{test.ID}, testByID)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if !Equal(testByID, testByLogin) {
		t.Errorf("user selected by ID != user selected by Login")
		return
	}
	_, err = db.Exec("TRUNCATE TABLE users")
	if err != nil {
		panic(err.Error())
	}
}

func TestValidEdit(t *testing.T) {
	db := initDB()
	defer db.Close()
	um := NewUserManager(db)
	in := &UsersLogin{Login: "test"}
	test := &User{}
	err := um.CreateUser(in, test)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	test.Login = "another_test"
	test.CreatedAt = float64(time.Now().Unix())
	testEdited := &User{}
	err = um.EditUser(test, testEdited)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if !Equal(test, testEdited) {
		t.Errorf("Edit in  != Edit out")
		return
	}
	_, err = db.Exec("TRUNCATE TABLE users")
	if err != nil {
		panic(err.Error())
	}
}

func TestShortLoginCreate(t *testing.T) {
	db := initDB()
	defer db.Close()
	um := NewUserManager(db)
	in := &UsersLogin{Login: "tes"}
	test := &User{}
	err := um.CreateUser(in, test)
	if err == nil || err != ErrShortLogin {
		t.Errorf(err.Error())
		return
	}
	_, err = db.Exec("TRUNCATE TABLE users")
	if err != nil {
		panic(err.Error())
	}
}
func TestUserByLoginNotFound(t *testing.T) {
	db := initDB()
	defer db.Close()
	um := NewUserManager(db)
	test := &User{}
	err := um.GetUserByLogin(&UsersLogin{Login: "test"}, test)
	if err == nil || err != ErrNotFound {
		t.Errorf("unexisting user found or error occured")
		return
	}
	_, err = db.Exec("TRUNCATE TABLE users")
	if err != nil {
		panic(err.Error())
	}
}
func TestUserByIDNotFound(t *testing.T) {
	db := initDB()
	defer db.Close()
	um := NewUserManager(db)
	test := &User{}
	err := um.GetUserByID(&UsersID{ID: "550e8400-e29b-41d4-a716-446655440000"}, test)
	if err == nil || err != ErrNotFound {
		t.Errorf("unexisting user found or db error occured")
		return
	}
	_, err = db.Exec("TRUNCATE TABLE users")
	if err != nil {
		panic(err.Error())
	}
}
func TestUserAlreadyExists(t *testing.T) {
	db := initDB()
	defer db.Close()
	um := NewUserManager(db)
	in := &UsersLogin{Login: "test"}
	test := &User{}
	err := um.CreateUser(in, test)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	err = um.CreateUser(in, test)
	if err == nil || err != ErrAlreadyExists {
		t.Errorf("user created twice")
		return
	}
	_, err = db.Exec("TRUNCATE TABLE users")
	if err != nil {
		panic(err.Error())
	}
}

func TestUserEditAlreadyExists(t *testing.T) {
	db := initDB()
	defer db.Close()
	um := NewUserManager(db)
	in := &UsersLogin{Login: "test"}
	test := &User{}
	err := um.CreateUser(in, test)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	in.Login = "another_test"
	anotherTest := &User{}
	err = um.CreateUser(in, anotherTest)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	dummy := &User{}
	anotherTest.Login = "test"
	err = um.EditUser(anotherTest, dummy)
	if err == nil || err != ErrAlreadyExists {
		t.Errorf("user login changed to login of another user or db error occured")
	}
	_, err = db.Exec("TRUNCATE TABLE users")
	if err != nil {
		panic(err.Error())
	}
}

func TestUnexistingEdit(t *testing.T) {
	db := initDB()
	defer db.Close()
	um := NewUserManager(db)
	in := &User{
		Login:     "test",
		ID:        "550e8400-e29b-41d4-a716-446655440000",
		CreatedAt: float64(time.Now().Unix()),
	}
	test := &User{}
	err := um.EditUser(in, test)
	if err == nil || err != ErrNotFound {
		t.Errorf("edited unexisting user")
		return
	}
	_, err = db.Exec("TRUNCATE TABLE users")
	if err != nil {
		panic(err.Error())
	}
}
