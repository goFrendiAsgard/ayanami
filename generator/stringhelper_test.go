package generator

import (
	"reflect"
	"testing"
)

func TestString(t *testing.T) {
	s := StringHelper{}
	// Quote
	expectedQuoteResult := `"echo \"hello\";"`
	actualQuoteResult := s.Quote(`echo "hello";`)
	if expectedQuoteResult != actualQuoteResult {
		t.Errorf("expected `%s`, get `%s`", expectedQuoteResult, actualQuoteResult)
	}
	// QuoteArray
	expectedQuoteArrayResult := []string{`"rei"`, `"shinji"`, `"asuka"`, `"\""`}
	actualQuoteArrayResult := s.QuoteArray([]string{"rei", "shinji", "asuka", `"`})
	if !reflect.DeepEqual(expectedQuoteArrayResult, actualQuoteArrayResult) {
		t.Errorf("expected %#v, get %#v", expectedQuoteArrayResult, actualQuoteArrayResult)
	}
	// QuoteArrayAndJoin
	expectedQuoteArrayAndJoinResult := `"rei", "shinji", "asuka", "\""`
	actualQuoteArrayAndJoinResult := s.QuoteArrayAndJoin([]string{"rei", "shinji", "asuka", `"`}, ", ")
	if expectedQuoteArrayAndJoinResult != actualQuoteArrayAndJoinResult {
		t.Errorf("expected `%s`, get `%s`", expectedQuoteArrayAndJoinResult, actualQuoteArrayAndJoinResult)
	}
	// IsMatch
	expectedMatch := true
	actualMatch := s.IsMatch("abc", "^a.+$")
	if expectedMatch != actualMatch {
		t.Errorf("expected `%t`, get `%t`", expectedMatch, actualMatch)
	}
	expectedMatch = false
	actualMatch = s.IsMatch("xabc", "^a.+$")
	if expectedMatch != actualMatch {
		t.Errorf("expected `%t`, get `%t`", expectedMatch, actualMatch)
	}
	// IsAlpha
	expectedMatch = true
	actualMatch = s.IsAlpha("abc")
	if expectedMatch != actualMatch {
		t.Errorf("expected `%t`, get `%t`", expectedMatch, actualMatch)
	}
	expectedMatch = false
	actualMatch = s.IsAlpha("123abc")
	if expectedMatch != actualMatch {
		t.Errorf("expected `%t`, get `%t`", expectedMatch, actualMatch)
	}
	// IsNumeric
	expectedMatch = true
	actualMatch = s.IsNumeric("123")
	if expectedMatch != actualMatch {
		t.Errorf("expected `%t`, get `%t`", expectedMatch, actualMatch)
	}
	expectedMatch = false
	actualMatch = s.IsNumeric("123abc")
	if expectedMatch != actualMatch {
		t.Errorf("expected `%t`, get `%t`", expectedMatch, actualMatch)
	}
	// IsAlphaNumeric
	expectedMatch = true
	actualMatch = s.IsAlphaNumeric("123abc")
	if expectedMatch != actualMatch {
		t.Errorf("expected `%t`, get `%t`", expectedMatch, actualMatch)
	}
	expectedMatch = false
	actualMatch = s.IsAlphaNumeric("123abc!")
	if expectedMatch != actualMatch {
		t.Errorf("expected `%t`, get `%t`", expectedMatch, actualMatch)
	}
}
