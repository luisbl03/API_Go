package models

type User struct {
	username string
	password string
}

func GetUsername(u *User) string {
	return u.username
}
func GetPassword(u *User) string {
	return u.password
}
func SetUsername(u *User, username string) {
	u.username = username
}
func SetPassword(u *User, password string) {
	u.password = password
}