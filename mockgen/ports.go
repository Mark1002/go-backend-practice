package mockgen

type IUserRepo interface {
	GetUserByID(id int) (*User, error)
	Insert(user User) error
	Update(id int, user User) error
}
