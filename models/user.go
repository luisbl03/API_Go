package models

type User struct {
	username string
	password string
}

func getUsername(u *User) string {
	return u.username
}
func getPassword(u *User) string {
	return u.password
}