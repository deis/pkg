package time

import (
	"testing"
)

func TestUnMarshalText(t *testing.T) {
	dummyTime := Time{}

	goodTimeFormats := []string{
		"2006-01-02T15:04:05MST",
		"2006-01-02T15:04:05",
	}
	for _, goodTime := range goodTimeFormats {
		if dummyTime.UnmarshalText([]byte(goodTime)) != nil {
			t.Error("expected " + goodTime + " to be marshal-able.")
		}
	}

	badTime := "this is a bad time, isn't it?"
	if dummyTime.UnmarshalText([]byte(badTime)) == nil {
		t.Error("expected " + badTime + "to be unmarshal-able.")
	}
}
