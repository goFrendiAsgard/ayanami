package service

import (
	"reflect"
	"testing"
)

func createTestIOList() IOList {
	var ioList IOList
	ioList = append(ioList, IO{VarName: "a", EventName: "consume.a"})
	ioList = append(ioList, IO{VarName: "alpha", EventName: "consume.a"})
	ioList = append(ioList, IO{VarName: "b", EventName: "publishServiceOutput.b"})
	ioList = append(ioList, IO{VarName: "b", EventName: "publishServiceOutput.any"})
	return ioList
}

func TestIoGetUniqueVarNames(t *testing.T) {
	ioList := createTestIOList()
	actual := ioList.GetUniqueVarNames()
	expected := []string{"a", "alpha", "b"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, get %#v", expected, actual)
	}
}

func TestIoGetUniqueEventNames(t *testing.T) {
	ioList := createTestIOList()
	actual := ioList.GetUniqueEventNames()
	expected := []string{"consume.a", "publishServiceOutput.b", "publishServiceOutput.any"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, get %#v", expected, actual)
	}
}

func TestIoGetVarEventNames(t *testing.T) {
	ioList := createTestIOList()

	actual := ioList.GetVarEventNames("a")
	expected := []string{"consume.a"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, get %#v", expected, actual)
	}

	actual = ioList.GetVarEventNames("alpha")
	expected = []string{"consume.a"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, get %#v", expected, actual)
	}

	actual = ioList.GetVarEventNames("b")
	expected = []string{"publishServiceOutput.b", "publishServiceOutput.any"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, get %#v", expected, actual)
	}

	var emptyList []string
	actual = ioList.GetVarEventNames("invalid")
	expected = emptyList
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, get %#v", expected, actual)
	}

}

func TestIoGetEventVarNames(t *testing.T) {
	ioList := createTestIOList()

	actual := ioList.GetEventVarNames("consume.a")
	expected := []string{"a", "alpha"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, get %#v", expected, actual)
	}

	actual = ioList.GetEventVarNames("publishServiceOutput.b")
	expected = []string{"b"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, get %#v", expected, actual)
	}

	actual = ioList.GetEventVarNames("publishServiceOutput.any")
	expected = []string{"b"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, get %#v", expected, actual)
	}

	var emptyList []string
	actual = ioList.GetEventVarNames("invalid")
	expected = emptyList
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, get %#v", expected, actual)
	}

}
