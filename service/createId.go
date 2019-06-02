package service

import (
	"fmt"
	"github.com/gofrs/uuid"
	"strings"
)

// CreateID create new UUID
func CreateID() (string, error) {
	// create ID
	UUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	ID := fmt.Sprintf("%s", UUID)
	// remove hyphens, since some message broker only support alpha-numeric
	ID = strings.Replace(ID, "-", "", -1)
	return ID, err
}
