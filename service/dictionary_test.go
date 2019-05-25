package service

import (
	"testing"
)

func TestDictionary(t *testing.T) {
	person := make(Dictionary)
	person["name"] = "Arya"
	person["surename"] = "Stark"
	person["affiliations"] = []interface{}{"faceless men", "winterfell"}
	dictionary := make(Dictionary)
	dictionary["person"] = person

	var expected, actual interface{}

	expected = "Arya"
	actual = dictionary.Get("person.name")
	if actual != expected {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

	expected = "Stark"
	actual = dictionary.Get("person.surename")
	if actual != expected {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

	expected = nil
	actual = dictionary.Get("race")
	if actual != nil {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

	expected = nil
	actual = dictionary.Get("person.weapons")
	if actual != nil {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

	expected = "faceless men"
	actual = dictionary.Get("person.affiliations.0")
	if actual != expected {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

	expected = "winterfell"
	actual = dictionary.Get("person.affiliations.1")
	if actual != expected {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

	expected = nil
	actual = dictionary.Get("person.affiliations.2")
	if actual != nil {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

	expected = nil
	actual = dictionary.Get("person.affiliations.name")
	if actual != nil {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

}
