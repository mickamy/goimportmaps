package usecase

import (
	"github.com/mickamy/goimportmaps-example/internal/sanity/model"
	"github.com/mickamy/goimportmaps-example/internal/sanity/repository"
)

type FindUserUseCase struct {
	repo repository.User
}

func (u *FindUserUseCase) Do(id string) (*model.User, error) {
	user, err := u.repo.Find(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
