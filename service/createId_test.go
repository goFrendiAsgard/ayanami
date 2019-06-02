package service

import (
	"regexp"
	"testing"
)

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{32}$")
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
