package exists

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.riyazali.net/sqlite"
)

type exists struct {
	rdb *redis.Client
}

func (f *exists) Args() int           { return -1 }
func (f *exists) Deterministic() bool { return false }
func (f *exists) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	// var input []string
	// if len(values) >= 1 {
	// 	for _, v := range values {
	// 		input = append(input, v.Text())
	// 	}
	// } else {
	// 	ctx.ResultError(fmt.Errorf("must supply at least one argument to redis exists command"))
	// 	return
	// }

	result := f.rdb.Exists(context.TODO(), "akey")

	ctx.ResultInt(int(result.Val()))
}

// New returns a sqlite function for reading the contents of a file
func New(rdb *redis.Client) sqlite.Function {
	return &exists{rdb}
}
