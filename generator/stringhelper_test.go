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
	// QuoteMap
	expectedQuoteMapResult := map[string]string{"00": `"rei"`, "01": `"shinji"`, "02": `"asuka"`, "03": `"\""`}
	actualQuoteMapResult := s.QuoteMap(map[string]string{"00": "rei", "01": "shinji", "02": "asuka", "03": `"`})
	if !reflect.DeepEqual(expectedQuoteMapResult, actualQuoteMapResult) {
		t.Errorf("expected %#v, get %#v", expectedQuoteMapResult, actualQuoteMapResult)
	}
	// IsMatch
	actualMatch := s.IsMatch("abc", "^a.+$")
	if true != actualMatch {
		t.Errorf("expected `%t`, get `%t`", true, actualMatch)
	}
	actualMatch = s.IsMatch("xabc", "^a.+$")
	if false != actualMatch {
		t.Errorf("expected `%t`, get `%t`", false, actualMatch)
	}
	// IsAlpha
	actualMatch = s.IsAlpha("abc")
	if true != actualMatch {
		t.Errorf("expected `%t`, get `%t`", true, actualMatch)
	}
	actualMatch = s.IsAlpha("123abc")
	if false != actualMatch {
		t.Errorf("expected `%t`, get `%t`", false, actualMatch)
	}
	// IsNumeric
	actualMatch = s.IsNumeric("123")
	if true != actualMatch {
		t.Errorf("expected `%t`, get `%t`", true, actualMatch)
	}
	actualMatch = s.IsNumeric("123abc")
	if false != actualMatch {
		t.Errorf("expected `%t`, get `%t`", false, actualMatch)
	}
	// IsAlphaNumeric
	actualMatch = s.IsAlphaNumeric("123abc")
	if true != actualMatch {
		t.Errorf("expected `%t`, get `%t`", true, actualMatch)
	}
	actualMatch = s.IsAlphaNumeric("123abc!")
	if false != actualMatch {
		t.Errorf("expected `%t`, get `%t`", false, actualMatch)
	}
}
