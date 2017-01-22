package slovnik

import "testing"

func TestDetectLanguage(t *testing.T) {
	cases := []struct {
		in, lang string
	}{
		{"hlavní", "cs"},
		{"привет", "ru"},
		{"sиniy", "ru"},
	}

	for _, c := range cases {
		got := DetectLanguage(c.in)
		if got != c.lang {
			t.Errorf("DetectLanguage(%q) == %q, want %q", c.in, got, c.lang)
		}
	}

}
