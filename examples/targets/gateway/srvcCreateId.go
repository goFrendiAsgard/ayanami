package main

import (
	"fmt"
	"github.com/gofrs/uuid"
)

// CreateID create new UUID
func CreateID() (string, error) {
	// create ID
	UUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	ID := fmt.Sprintf("%s", UUID)
	return ID, err
}
