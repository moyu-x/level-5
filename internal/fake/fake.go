package fake

import "github.com/brianvoe/gofakeit/v7"

type Faker struct {
	*gofakeit.Faker
}

func New() *Faker {
	return &Faker{gofakeit.New(0)}
}
