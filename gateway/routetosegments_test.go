package gateway

import (
	"testing"
)

func TestRouteToSegments(t *testing.T) {
	var actual, expected string

	expected = ""
	actual = RouteToSegments("/")
	if actual != expected {
		t.Errorf("expected %s, get %s", expected, actual)
	}

	expected = "foo.bar"
	actual = RouteToSegments("/foo/bar")
	if actual != expected {
		t.Errorf("expected %s, get %s", expected, actual)
	}

	expected = "foo.bar.egg"
	actual = RouteToSegments("/foo/bar egg")
	if actual != expected {
		t.Errorf("expected %s, get %s", expected, actual)
	}

}
