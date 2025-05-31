package logic

import (
	"math/rand"
	"time"
)

// SeedManager handles the generation and management of random seeds.
type SeedManager struct {
	seed int64
	rng  *rand.Rand
}

// NewSeedManager creates a new SeedManager with a random seed.
func NewSeedManager() *SeedManager {
	seed := time.Now().UnixNano()
	return &SeedManager{
		seed: seed,
		rng:  rand.New(rand.NewSource(seed)),
	}
}

// NewSeedManagerWithFixedValue creates a new SeedManager with a fixed seed for testing.
func NewSeedManagerWithFixedValue(seed int64) *SeedManager {
	return &SeedManager{
		seed: seed,
		rng:  rand.New(rand.NewSource(seed)),
	}
}

// GetSeed returns the current seed value.
func (sm *SeedManager) GetSeed() int64 {
	return sm.seed
}

// SetSeed sets a new seed value.
func (sm *SeedManager) SetSeed(seed int64) {
	sm.seed = seed
	sm.rng = rand.New(rand.NewSource(seed))
}

// RandomFloat64 generates a random float64 value between 0 and 1.
func (sm *SeedManager) RandomFloat64() float64 {
	return sm.rng.Float64()
}

// RandomInt generates a random integer between min and max.
func (sm *SeedManager) RandomInt(min, max int) int {
	return sm.rng.Intn(max-min) + min
}
