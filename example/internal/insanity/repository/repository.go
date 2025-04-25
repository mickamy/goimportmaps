package repository

import (
	"github.com/mickamy/goimportmaps-example/internal/sanity/model"
)

type User struct {
}

func New() *User {
	return &User{}
}

func (r *User) Find(id string) (*model.User, error) {
	return &model.User{
		ID:   id,
		Name: "John Doe",
	}, nil
}
