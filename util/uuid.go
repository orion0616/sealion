package util

import (
	"github.com/google/uuid"
)

// CreateUUID returns uuid
func CreateUUID() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	uu := u.String()
	return uu, nil
}
