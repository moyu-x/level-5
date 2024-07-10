package fake_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
)

func TestFake(t *testing.T) {
	// fmt.Println(gofakeit.IPv4Address())
	faker := gofakeit.New(0)
	t.Log(faker.IPv4Address())
	t.Log(gofakeit.IPv4Address())
}
