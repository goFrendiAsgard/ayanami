package main

import (
	"fmt"
	"github.com/gofrs/uuid"
)

// SrvcCreateID create new UUID
func SrvcCreateID() (string, error) {
	// create ID
	UUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	ID := fmt.Sprintf("%s", UUID)
	return ID, err
}
