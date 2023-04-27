package user

type UserModelInterface interface {
	CreateUser(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}
