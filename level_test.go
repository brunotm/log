package log

import "testing"

func TestParseLevel(t *testing.T) {
	l, err := ParseLevel("DEBUG")
	if err != nil {
		t.Error(err)
	}

	if l != DEBUG {
		t.Errorf("level %s is not debug", l.String())
	}

	if _, err := ParseLevel("___DEBUG__"); err == nil {
		t.Errorf("got no error in parsing an invalid level")
	}

}
