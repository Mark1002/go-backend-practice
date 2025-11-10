package mockgen

import "fmt"

type UserService struct {
	repo IUserRepo
}

var invalidUserIDError = fmt.Errorf("invalid user id")

func (u *UserService) Upsert(user User) error {
	if user.ID <= 0 {
		return invalidUserIDError
	}
	existingUser, err := u.repo.GetUserByID(user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return u.repo.Insert(user)
	}

	return u.repo.Update(user.ID, user)
}

func (u *UserService) GetUserByID(id int) (*User, error) {
	return u.repo.GetUserByID(id)
}

func (u *UserService) DeleteUserByID(id int) error {
	return u.repo.Delete(id)
}
