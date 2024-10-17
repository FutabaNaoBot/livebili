package gopool

import "github.com/panjf2000/ants/v2"

var pool, _ = ants.NewPool(-1)

func Go(f func()) {
	_ = pool.Submit(f)
}
