package objects

type User struct {
	Id   int
	Role string
}

func NewUser(id int, role string) *User {
	return &User{Id: id, Role: role}
}
