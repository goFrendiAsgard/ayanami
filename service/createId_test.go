package service

import (
	"regexp"
	"testing"
)

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func TestCreateID(t *testing.T) {
	ID, err := CreateID()
	if err != nil {
		t.Errorf("Getting error: %s", err)
	}
	if !isValidUUID(ID) {
		t.Errorf("Expected valid UUID, get %s", ID)
	}
}
