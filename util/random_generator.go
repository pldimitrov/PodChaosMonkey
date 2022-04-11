package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// IRandomGenerator interface used to mock rand for testing
type IRandomGenerator interface {
	Int() int
}

// RandomGenerator used to generate random integers
type RandomGenerator struct{}

// NewRandomGenerator constructor
func NewRandomGenerator() RandomGenerator {
	return RandomGenerator{}
}

// Int returns a random integer
func (r RandomGenerator) Int() int {
	return rand.Int()
}

// FakeRandomGenerator used to mock the behaviour or rand for unit tests
type FakeRandomGenerator struct {
	n int
}

// NewFakeRandomGenerator constructor
func NewFakeRandomGenerator(n int) FakeRandomGenerator {
	return FakeRandomGenerator{n: n}
}

// Int returns the predefined value set at construction time
func (f FakeRandomGenerator) Int() int {
	return f.n
}
