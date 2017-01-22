package slovnik

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"strings"

	"log"

	"golang.org/x/net/html"
)

// Word defines a structure with the word itself and possible translations of that word
type Word struct {
	Word         string
	Translations []string
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

	log.Println(query.String())

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
					}
				}
			}

			if t.Data == "a" && inTranslations {
				inTranslationLink = true
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

			break

		case tt == html.TextToken:
			t := z.Token()
			if inWord {
				w.Word = t.Data
			}
			if inTranslationLink {
				w.Translations = append(w.Translations, t.Data)
			}
			break
		}
	}
}
