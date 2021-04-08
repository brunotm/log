package log

import "testing"

func TestParseFormat(t *testing.T) {
	f, err := ParseFormat("JSON")
	if err != nil {
		t.Error(err)
	}

	if f != FormatJSON {
		t.Errorf("format %s is not json", f.String())
	}

	if _, err := ParseFormat("TEXTY"); err == nil {
		t.Errorf("got no error in parsing an invalid format")
	}

}
