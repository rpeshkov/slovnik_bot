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

	if w.Word != expectedWord {
		t.Errorf("ParsePage word == %q, want %q", w.Word, expectedWord)
	}

	expectedTranslations := []string{
		"гла́вный",
		"основно́й",
		"центра́льный",
	}

	if len(w.Translations) != len(expectedTranslations) {
		t.Errorf("ParsePage len(translation) == %d, want %d", len(w.Translations), len(expectedTranslations))
	}

	for i, trans := range w.Translations {
		if trans != expectedTranslations[i] {
			t.Errorf("ParsePage translation == %q, want %q", trans, expectedTranslations[i])
		}
	}

	const expectedWordType = "přídavné jméno"
	if w.WordType != expectedWordType {
		t.Errorf("ParsePage wordType == %q, want %q", w.WordType, expectedWordType)
	}

	expectedSynonyms := []string{
		"ústřední",
		"podstatný",
		"základní",
		"zásadní",
	}

	if len(w.Synonyms) != len(expectedSynonyms) {
		t.Errorf("ParsePage len(synonyms) == %d, want %d", len(w.Synonyms), len(expectedSynonyms))

		for i, synonym := range w.Synonyms {
			if synonym != expectedSynonyms[i] {
				t.Errorf("ParsePage synonym == %q, want %q", synonym, expectedSynonyms[i])
			}
		}
	}

	expectedAntonyms := []string{
		"vedlejší",
		"podřadný",
		"podružný",
	}

	if len(w.Antonyms) != len(expectedAntonyms) {
		t.Errorf("ParsePage len(antonyms) == %d, want %d", len(w.Antonyms), len(expectedAntonyms))

		for i, antonym := range w.Antonyms {
			if antonym != expectedAntonyms[i] {
				t.Errorf("ParsePage antonym == %q, want %q", antonym, expectedAntonyms[i])
			}
		}
	}
}

func TestParseAltPage(t *testing.T) {
	f, _ := os.Open("./sample_issue8.html")
	w := parsePage(f)

	const expectedWord = "soutěživý"

	if w.Word != expectedWord {
		t.Errorf("ParsePage word == %q, want %q", w.Word, expectedWord)
	}

	expectedTranslations := []string{
		"состяза́тельный",
	}

	if len(w.Translations) != len(expectedTranslations) {
		t.Errorf("ParsePage len(translation) == %d, want %d", len(w.Translations), len(expectedTranslations))
	}

	for i, trans := range w.Translations {
		if trans != expectedTranslations[i] {
			t.Errorf("ParsePage translation == %q, want %q", trans, expectedTranslations[i])
		}
	}
}
