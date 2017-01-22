package slovnik

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"strings"

	"golang.org/x/net/html"
)

// Word defines a structure with the word itself and possible translations of that word
type Word struct {
	Word         string
	Translations []string
	WordType     string
	Synonyms     []string
	Antonyms     []string
}

// Method for transforming Word struct to string
func (w Word) String() string {
	return fmt.Sprintf("*%s*\n%s", w.Word, strings.Join(w.Translations, ", "))
}

// GetTranslations from slovnik.seznam.cz for specified word
func GetTranslations(word string, langcode string) (Word, error) {
	urls := map[string]string{
		"cs": "https://slovnik.seznam.cz/cz-ru/",
		"ru": "https://slovnik.seznam.cz/ru/",
	}

	query, _ := url.Parse(urls[langcode])

	p := url.Values{}
	p.Add("q", word)

	query.RawQuery = p.Encode()

	resp, err := http.Get(query.String())
	if err != nil {
		return Word{}, err
	}
	return parsePage(resp.Body), nil
}

// DetectLanguage used to find out which language is used for the input string
func DetectLanguage(input string) string {
	const ru = "абвгдеёжзийклмнопрстуфхцчшщьыъэюя"
	for _, ch := range input {
		if strings.Contains(ru, strings.ToLower(string(ch))) {
			return "ru"
		}
	}
	return "cs"
}

func parsePage(pageBody io.Reader) Word {
	z := html.NewTokenizer(pageBody)

	inWord := false
	inTranslations := false
	inTranslationLink := false

	// class of span that contains wordtype
	inMorf := false

	foundSynonymsHeader := false
	inSynonymsBlock := false
	inSynonima := false

	foundAntonymsHeader := false
	inAntonymsBlock := false
	inAntonima := false

	w := Word{}
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return w

		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data == "h3" {
				for _, attr := range t.Attr {
					if attr.Key == "lang" && (attr.Val == "cs" || attr.Val == "ru") {
						inWord = true
					}
				}
			}

			if t.Data == "div" {
				for _, attr := range t.Attr {
					if attr.Key == "id" && attr.Val == "fastMeanings" {
						inTranslations = true
					} else if attr.Key == "class" && attr.Val == "other-meaning" && foundSynonymsHeader {
						inSynonymsBlock = true
					} else if attr.Key == "class" && attr.Val == "other-meaning" && foundAntonymsHeader {
						inAntonymsBlock = true
					}
				}
			}

			if t.Data == "a" && inTranslations {
				inTranslationLink = true
			}

			if t.Data == "a" && inSynonymsBlock {
				inSynonima = true
			}

			if t.Data == "a" && inAntonymsBlock {
				inAntonima = true
			}

			if t.Data == "span" {
				for _, attr := range t.Attr {
					if attr.Key == "class" && attr.Val == "morf" {
						inMorf = true
					}
				}
			}

			break

		case tt == html.EndTagToken:
			t := z.Token()
			if t.Data == "div" {
				inTranslations = false
			}

			if t.Data == "h3" {
				inWord = false
			}
			if t.Data == "a" && inTranslationLink {
				inTranslationLink = false
			}
			if t.Data == "a" && inSynonima {
				inSynonima = false
			}
			if t.Data == "a" && inAntonima {
				inAntonima = false
			}
			if t.Data == "span" && inMorf {
				inMorf = false
			}

			if t.Data == "div" && inSynonymsBlock {
				inSynonymsBlock = false
				foundSynonymsHeader = false
			}

			if t.Data == "div" && inAntonymsBlock {
				inAntonymsBlock = false
				foundAntonymsHeader = false
			}

			break

		case tt == html.TextToken:
			t := z.Token()
			if inWord {
				w.Word = t.Data
			}
			if inTranslationLink {
				w.Translations = append(w.Translations, t.Data)
			}

			if inMorf {
				w.WordType = t.Data
			}

			if t.Data == "Synonyma" {
				foundSynonymsHeader = true
			}

			if t.Data == "Antonyma" {
				foundAntonymsHeader = true
			}

			if inSynonima {
				w.Synonyms = append(w.Synonyms, t.Data)
			}

			if inAntonima {
				w.Antonyms = append(w.Antonyms, t.Data)
			}
			break
		}
	}
}
