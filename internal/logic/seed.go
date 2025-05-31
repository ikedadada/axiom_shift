package logic

import (
	"math/rand"
	"time"
)

// SeedManager handles the generation and management of random seeds.
type SeedManager struct {
	seed int64
}

// NewSeedManager creates a new SeedManager with a random seed.
func NewSeedManager() *SeedManager {
	return &SeedManager{
		seed: time.Now().UnixNano(),
	}
}

// NewSeedManagerWithFixedValue creates a new SeedManager with a fixed seed for testing.
func NewSeedManagerWithFixedValue(seed int64) *SeedManager {
	return &SeedManager{
		seed: seed,
	}
}

// GetSeed returns the current seed value.
func (sm *SeedManager) GetSeed() int64 {
	return sm.seed
}

// SetSeed sets a new seed value.
func (sm *SeedManager) SetSeed(seed int64) {
	sm.seed = seed
	rand.Seed(sm.seed)
}

// RandomFloat64 generates a random float64 value between 0 and 1.
func (sm *SeedManager) RandomFloat64() float64 {
	return rand.Float64()
}

// RandomInt generates a random integer between min and max.
func (sm *SeedManager) RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}
