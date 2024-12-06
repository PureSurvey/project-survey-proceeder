package objects

type User struct {
	Id           int
	Role         string
	IsSubscribed bool
}

func NewUser(id int, role string, isSubscribed bool) *User {
	return &User{Id: id, Role: role, IsSubscribed: isSubscribed}
}
