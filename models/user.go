package models

type User struct {
	Username string
	Password string
	Admin    bool // true if the user is an admin
}

var mockUsers = []User{
	{
		Username: "admin",
		Password: "admin123",
		Admin:    true,
	},
	{
		Username: "user",
		Password: "password123",
		Admin:    false,
	},
}

func UserMatchPassword(username string, password string) *User {
	for _, user := range mockUsers {
		if user.Username == username {
			if user.Password == password {
				return &user
			}
		}
	}
	return nil
}
