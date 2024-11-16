package models

type User struct {
	USERNAME string
	PASSWORD string
}

func SetUsername(user *User, username string) {
	user.USERNAME = username
}
func SetPassword(user *User, password string) {
	user.PASSWORD = password
}