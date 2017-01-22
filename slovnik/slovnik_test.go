package slovnik

import "testing"
import "os"

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

func TestParsePage(t *testing.T) {
	f, _ := os.Open("./sample.html")
	w := parsePage(f)

	const expectedWord = "hlavní"
	expectedTranslations := []string{
		"гла́вный",
		"основно́й",
		"центра́льный",
	}

	const expectedWordType = "přídavné jméno"

	if w.Word != expectedWord {
		t.Errorf("ParsePage word == %q, want %q", w.Word, expectedWord)
	}

	for i, trans := range w.Translations {
		if trans != expectedTranslations[i] {
			t.Errorf("ParsePage translation == %q, want %q", trans, expectedTranslations[i])
		}
	}

	if w.WordType != expectedWordType {
		t.Errorf("ParsePage wordType == %q, want %q", w.WordType, expectedWordType)
	}
}
