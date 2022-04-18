package uuid

import "github.com/google/uuid"

// Generate gives a random UUID
func Generate() string {
	return uuid.NewString()
}
