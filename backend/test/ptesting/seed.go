package ptesting

import (
	"math/rand/v2"
	"strconv"
	"testing"
)

type opt struct {
	seed  uint64
	tests int
}

type Option func(*opt)

func Seed(seed uint64) Option {
	return func(o *opt) {
		o.seed = seed
	}
}

func Count(count int) Option {
	return func(o *opt) {
		o.tests = count
	}
}

type Gen struct {
	seed *rand.PCG
	r    *rand.Rand
}

func ForAll(t *testing.T, opts ...Option) func(f func(*testing.T, *Gen)) {
	o := &opt{
		seed:  rand.Uint64(),
		tests: 10,
	}
	for _, opt := range opts {
		opt(o)
	}

	var seed = o.seed
	return func(f func(*testing.T, *Gen)) {
		for i := 0; i < o.tests; i++ {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				t.Logf("seed = %d", seed)
				pcg := rand.NewPCG(seed, seed+1)
				gen := &Gen{seed: pcg, r: rand.New(pcg)}

				t.Cleanup(func() {
					seed = pcg.Uint64()
				})

				f(t, gen)
			})
		}
	}
}
