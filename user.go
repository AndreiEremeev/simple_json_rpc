package main

type User struct {
	Login     string `json:"login"`
	Uuid      string `json:"uuid"`
	CreatedAt string `json:"createdAt"`
}

type UsersLogin struct {
	Login string `json:login`
}

type UsersUuid struct {
	Uuid string `json:uuid`
}

type UserManager struct {
}

func NewUserManager() *UserManager {
	return &UserManager{}
}

func (um *UserManager) GetUserByLogin(in *UsersLogin, out *User) error {
	*out = User{
		Login: in.Login,
	}
	return nil
}
