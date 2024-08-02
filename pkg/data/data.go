package data

import "context"

type Data interface {
	Reader(ctx context.Context)
	Writer(ctx context.Context)
}
