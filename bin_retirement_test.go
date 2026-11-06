package parseapi

import (
	"reflect"
	"testing"
)

func TestLegacyBINIsNotAPublicLookup(t *testing.T) {
	typ := reflect.TypeOf((*Client)(nil))
	if _, found := typ.MethodByName("BIN"); found {
		t.Fatal("legacy BIN method remains")
	}
	if _, found := typ.MethodByName("Card"); !found {
		t.Fatal("Card method is missing")
	}
}
