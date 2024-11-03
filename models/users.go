package models

type User struct {
	username string
	password string
}

func (u *User) getUsername() string {
	return u.username
}

func (u *User) getPassword() string {
	return u.password
}
func (u *User) setUsername(username string) {
	u.username = username
}
func (u *User) setPassword(password string) {
	u.password = password
}