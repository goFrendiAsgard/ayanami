package service

import (
	"testing"
)

func createTestDictionary() Dictionary {
	person := make(Dictionary)
	person["name"] = "Arya"
	person["surename"] = "Stark"
	person["affiliations"] = []interface{}{"faceless men", "winterfell"}
	dictionary := make(Dictionary)
	dictionary["person"] = person
	return dictionary
}

func TestDictionaryGet(t *testing.T) {
	dictionary := createTestDictionary()
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
	if actual != expected {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

	expected = nil
	actual = dictionary.Get("person.weapons")
	if actual != expected {
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
	if actual != expected {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

	expected = nil
	actual = dictionary.Get("person.affiliations.name")
	if actual != expected {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

}

func TestDictionaryContain(t *testing.T) {
	dictionary := createTestDictionary()
	var expected, actual bool

	expected = true
	actual = dictionary.Contain("person.name")
	if actual != expected {
		t.Errorf("Expected `%t`, get `%t`", expected, actual)
	}

	expected = true
	actual = dictionary.Contain("person.surename")
	if actual != expected {
		t.Errorf("Expected `%t`, get `%t`", expected, actual)
	}

	expected = false
	actual = dictionary.Contain("race")
	if actual != expected {
		t.Errorf("Expected `%t`, get `%t`", expected, actual)
	}

	expected = false
	actual = dictionary.Contain("person.weapons")
	if actual != expected {
		t.Errorf("Expected `%t`, get `%t`", expected, actual)
	}

	expected = true
	actual = dictionary.Contain("person.affiliations.0")
	if actual != expected {
		t.Errorf("Expected `%t`, get `%t`", expected, actual)
	}

	expected = true
	actual = dictionary.Contain("person.affiliations.1")
	if actual != expected {
		t.Errorf("Expected `%t`, get `%t`", expected, actual)
	}

	expected = false
	actual = dictionary.Contain("person.affiliations.2")
	if actual != expected {
		t.Errorf("Expected `%t`, get `%t`", expected, actual)
	}

	expected = false
	actual = dictionary.Contain("person.affiliations.name")
	if actual != expected {
		t.Errorf("Expected `%t`, get `%t`", expected, actual)
	}

}

func TestDictionarySet(t *testing.T) {
	dictionary := createTestDictionary()
	var err error
	var expected, actual interface{}

	err = dictionary.Set("person.affiliations.0", "the north")
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	expected = "the north"
	actual = dictionary.Get("person.affiliations.0")
	if actual != expected {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

	err = dictionary.Set("person.affiliations.2", "house of stark")
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	expected = "house of stark"
	actual = dictionary.Get("person.affiliations.2")
	if actual != expected {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

	err = dictionary.Set("person.name", "Sansa")
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	expected = "Sansa"
	actual = dictionary.Get("person.name")
	if actual != expected {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

	err = dictionary.Set("person.title", "queen in the north")
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	expected = "queen in the north"
	actual = dictionary.Get("person.title")
	if actual != expected {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

	err = dictionary.Set("location", "the north")
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	expected = "the north"
	actual = dictionary.Get("location")
	if actual != expected {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

	err = dictionary.Set("person.affiliations.-1", "the north")
	if err == nil {
		t.Errorf("Error expected")
	}
	expected = nil
	actual = dictionary.Get("person.affiliations.-1")
	if actual != expected {
		t.Errorf("Expected `%s`, get `%s`", expected, actual)
	}

}
