package graph

import "github.com/kaenova/go-playground/gqlgen/model"

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userIdIncrement  int
	todosIdIncrement int

	users []*model.User
	todos []*model.Todo
}

func (r *Resolver) searchUserById(id int) model.User {
	var user model.User
	for _, v := range r.users {
		if v.ID == id {
			user = *v
		}
	}
	return user
}
